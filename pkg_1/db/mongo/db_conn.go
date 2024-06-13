package mongo

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	Pages    Pages
	Notes    Notes
	Category Category
	Theme    Theme
}

type Pages struct {
	*mongo.Collection
}

type Notes struct {
	*mongo.Collection
}

type Category struct {
	*mongo.Collection
}

type Theme struct {
	*mongo.Collection
}

func New(connectString string, connectTimeout time.Duration) Storage {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectString).SetTLSConfig(&tls.Config{
		InsecureSkipVerify: true,
	}))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	pages := Pages{
		Collection: client.Database("memory-base").Collection("pages"),
	}

	notes := Notes{
		Collection: client.Database("memory-base").Collection("notes"),
	}

	category := Category{
		Collection: client.Database("memory-base").Collection("category"),
	}

	theme := Theme{
		Collection: client.Database("memory-base").Collection("theme"),
	}

	return Storage{
		Pages:    pages,
		Notes:    notes,
		Category: category,
		Theme:    theme,
	}
}

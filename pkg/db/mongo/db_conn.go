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
	pages Pages
	notes Notes
}

type Pages struct {
	*mongo.Collection
}

type Notes struct {
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
		Collection: client.Database("read-adviser").Collection("pages"),
	}

	notes := Notes{
		Collection: client.Database("read-adviser").Collection("notes"),
	}

	return Storage{
		pages: pages,
		notes: notes,
	}
}

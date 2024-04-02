package mongo

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"read-adviser-bot/lib/e"
	"read-adviser-bot/storage"
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

type Page struct {
	URL      string `bson:"url"`
	UserName string `bson:"username"`
}

type Note struct {
	Title    string `bson:"title"`
	Sphere   string `bson:"sphere"`
	Category string `bson:"category"`
	Content  string `bson:"content"`
}

func New(connectString string, connectTimeout time.Duration) Storage {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectString).SetTLSConfig(&tls.Config{
		InsecureSkipVerify: true, // Используйте только для отладки, не для продакшена
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

func (s Storage) Save(ctx context.Context, page *storage.Page) error {
	_, err := s.pages.InsertOne(ctx, Page{
		URL:      page.URL,
		UserName: page.UserName,
	})
	if err != nil {
		return e.Wrap("can't save page", err)
	}

	return nil
}

func (s Storage) SaveNote(ctx context.Context, note *storage.Note) error {
	_, err := s.notes.InsertOne(ctx, Note{
		Title:    note.Title,
		Sphere:   note.Sphere,
		Category: note.Category,
		Content:  note.Content,
	})
	if err != nil {
		return e.Wrap("can't save note", err)
	}

	return nil
}

func (s Storage) GetNote(ctx context.Context, conditionField string, conditionValue string) (note *storage.Note, err error) {
	defer func() { err = e.WrapIfErr("can't get note", err) }()

	//var n Note

	// Создаем агрегационный пайплайн с операторами $match и $sample
	pipeline := mongo.Pipeline{
		//{{"$match", bson.D{{conditionField, conditionValue}}}},
		{{"$sample", bson.D{{"size", 1}}}},
	}

	// Выполняем агрегацию
	cursor, err := s.notes.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregation: %w", err)
	}
	defer cursor.Close(context.Background())

	// Читаем результат
	var result bson.M
	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode result: %w", err)
		}
	}

	// Преобразование result (bson.M) в структуру Note
	var noteResult storage.Note
	if title, ok := result["title"].(string); ok {
		noteResult.Title = title
	}
	if sphere, ok := result["sphere"].(string); ok {
		noteResult.Sphere = sphere
	}
	if category, ok := result["category"].(string); ok {
		noteResult.Category = category
	}
	if content, ok := result["content"].(string); ok {
		noteResult.Content = content
	}

	return &noteResult, nil
}

func (s Storage) PickRandom(ctx context.Context, userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random page", err) }()

	pipe := bson.A{
		bson.M{"$sample": bson.M{"size": 1}},
	}

	cursor, err := s.pages.Aggregate(ctx, pipe)
	if err != nil {
		return nil, err
	}

	var p Page

	cursor.Next(ctx)

	err = cursor.Decode(&p)
	switch {
	case errors.Is(err, io.EOF):
		return nil, storage.ErrNoSavedPages
	case err != nil:
		return nil, err
	}

	return &storage.Page{
		URL:      p.URL,
		UserName: p.UserName,
	}, nil
}

func (s Storage) Remove(ctx context.Context, storagePage *storage.Page) error {
	_, err := s.pages.DeleteOne(ctx, toPage(storagePage).Filter())
	if err != nil {
		return e.Wrap("can't remove page", err)
	}

	return nil
}

func (s Storage) IsExists(ctx context.Context, storagePage *storage.Page) (bool, error) {
	count, err := s.pages.CountDocuments(ctx, toPage(storagePage).Filter())
	if err != nil {
		return false, e.Wrap("can't check if page exists", err)
	}

	return count > 0, nil
}

func toPage(p *storage.Page) Page {
	return Page{
		URL:      p.URL,
		UserName: p.UserName,
	}
}

func (p Page) Filter() bson.M {
	return bson.M{
		"url":      p.URL,
		"username": p.UserName,
	}
}

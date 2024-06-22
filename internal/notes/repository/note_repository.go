package repository

import (
	"context"
	"fmt"
	"github.com/MangoSociety/know_api/internal/notes/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"time"
)

type NoteRepository interface {
	GetAll() ([]domain.Note, error)
	GetById() (*domain.Note, error)
	Create(note domain.Note) error
	Update(note domain.Note) error
	Delete(id primitive.ObjectID) error
	GetByTitle(title string) (*domain.Note, error)
	GetRandomNoteByCategory(categoryIDs []primitive.ObjectID) (*domain.Note, error)
	GetByHex(hexID string) (*domain.Note, error)
}

type noteRepository struct {
	collection *mongo.Collection
}

func NewNoteRepository(dbClient *mongo.Client) NoteRepository {
	collection := dbClient.Database("note_taking").Collection("notes")
	return &noteRepository{collection: collection}
}

func (r *noteRepository) GetAll() ([]domain.Note, error) {
	var notes []domain.Note
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var note domain.Note
		cursor.Decode(&note)
		notes = append(notes, note)
	}
	return notes, nil
}

func (r *noteRepository) GetByHex(hexID string) (*domain.Note, error) {
	objectID, err := primitive.ObjectIDFromHex(hexID)
	var note domain.Note
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&note)
	if err != nil {
		return nil, fmt.Errorf("error finding document: %w", err)
	}

	return &note, nil
}

func (r *noteRepository) GetById() (*domain.Note, error) {
	panic("implement	me")
}

func (r *noteRepository) Create(note domain.Note) error {
	_, err := r.collection.InsertOne(context.Background(), note)
	return err
}

func (r *noteRepository) Update(note domain.Note) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": note.ID},
		bson.M{"$set": note},
		options.Update().SetUpsert(true),
	)
	return err
}

func (r *noteRepository) Delete(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *noteRepository) GetByTitle(title string) (*domain.Note, error) {
	var note domain.Note
	err := r.collection.FindOne(context.Background(), bson.M{"title": title}).Decode(&note)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &note, nil
}

func (r *noteRepository) GetRandomNoteByCategory(categoryIDs []primitive.ObjectID) (*domain.Note, error) {
	filter := bson.M{"category_ids": bson.M{"$in": categoryIDs}}

	count, err := r.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, nil
	}

	// Генерация случайного смещения
	rand.Seed(time.Now().UnixNano())
	skip := rand.Int63n(count)

	var note domain.Note
	err = r.collection.FindOne(context.Background(), filter, options.FindOne().SetSkip(skip)).Decode(&note)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

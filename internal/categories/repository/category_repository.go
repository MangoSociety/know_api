package repository

import (
	"context"
	"github.com/MangoSociety/know_api/internal/categories/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoryRepository interface {
	GetAll() ([]domain.Category, error)
	Create(category domain.Category) error
	Update(category domain.Category) error
	Delete(id primitive.ObjectID) error
	GetByName(name string) (*domain.Category, error)
}

type categoryRepository struct {
	collection *mongo.Collection
}

func NewCategoryRepository(dbClient *mongo.Client) CategoryRepository {
	collection := dbClient.Database("note_taking").Collection("categories")
	return &categoryRepository{collection: collection}
}

func (r *categoryRepository) GetAll() ([]domain.Category, error) {
	var categories []domain.Category
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var category domain.Category
		cursor.Decode(&category)
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *categoryRepository) Create(category domain.Category) error {
	_, err := r.collection.InsertOne(context.Background(), category)
	return err
}

func (r *categoryRepository) Update(category domain.Category) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": category.ID},
		bson.M{"$set": category},
		options.Update().SetUpsert(true),
	)
	return err
}

func (r *categoryRepository) Delete(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *categoryRepository) GetByName(name string) (*domain.Category, error) {
	var category domain.Category
	err := r.collection.FindOne(context.Background(), bson.M{"name": name}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

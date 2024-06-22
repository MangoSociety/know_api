package repository

import (
	"context"
	"github.com/MangoSociety/know_api/internal/spheres/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

type SphereRepository interface {
	GetAll() ([]domain.Sphere, error)
	Create(sphere domain.Sphere) error
	Update(sphere domain.Sphere) error
	Delete(id primitive.ObjectID) error
	GetByName(name string) (*domain.Sphere, error)
	GetById(id string) (*domain.Sphere, error)
}

type sphereRepository struct {
	collection *mongo.Collection
}

func NewSphereRepository(dbClient *mongo.Client) SphereRepository {
	collection := dbClient.Database("note_taking").Collection("spheres")
	return &sphereRepository{collection: collection}
}

func (r *sphereRepository) GetAll() ([]domain.Sphere, error) {
	var spheres []domain.Sphere
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var sphere domain.Sphere
		cursor.Decode(&sphere)
		spheres = append(spheres, sphere)
	}
	return spheres, nil
}

func (r *sphereRepository) Create(sphere domain.Sphere) error {
	_, err := r.collection.InsertOne(context.Background(), sphere)
	return err
}

func (r *sphereRepository) Update(sphere domain.Sphere) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": sphere.ID},
		bson.M{"$set": sphere},
		options.Update().SetUpsert(true),
	)
	return err
}

func (r *sphereRepository) Delete(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *sphereRepository) GetByName(name string) (*domain.Sphere, error) {
	var sphere domain.Sphere
	err := r.collection.FindOne(context.Background(), bson.M{"name": name}).Decode(&sphere)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &sphere, nil
}

func (r *sphereRepository) GetById(id string) (*domain.Sphere, error) {
	objectIDStr := strings.TrimPrefix(id, `ObjectID("`)
	objectIDStr = strings.TrimSuffix(objectIDStr, `")`)
	var sphere domain.Sphere
	objID, err := primitive.ObjectIDFromHex(objectIDStr)
	if err != nil {
		return nil, err
	}
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&sphere)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &sphere, nil
}

package repository

import (
	"context"
	"github.com/MangoSociety/know_api/internal/statistics/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatisticsRepository interface {
	GetAll() ([]domain.Statistics, error)
	Create(statistics domain.Statistics) error
}

type statisticsRepository struct {
	collection *mongo.Collection
}

func NewStatisticsRepository(dbClient *mongo.Client) StatisticsRepository {
	collection := dbClient.Database("note_taking").Collection("statistics")
	return &statisticsRepository{collection: collection}
}

func (r *statisticsRepository) GetAll() ([]domain.Statistics, error) {
	var statistics []domain.Statistics
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var stat domain.Statistics
		cursor.Decode(&stat)
		statistics = append(statistics, stat)
	}
	return statistics, nil
}

func (r *statisticsRepository) Create(statistics domain.Statistics) error {
	_, err := r.collection.InsertOne(context.Background(), statistics)
	return err
}

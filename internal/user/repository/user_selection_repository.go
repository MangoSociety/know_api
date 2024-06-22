package repository

import (
	"context"

	"github.com/MangoSociety/know_api/internal/user/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserSelectionRepository struct {
	collection *mongo.Collection
}

func NewUserSelectionRepository(dbClient *mongo.Client) *UserSelectionRepository {
	collection := dbClient.Database("note_taking").Collection("user_selections")
	return &UserSelectionRepository{
		collection: collection,
	}
}

func (r *UserSelectionRepository) GetByChatID(chatID int) (*domain.UserSelection, error) {
	filter := bson.M{"chat_id": chatID}
	var selection domain.UserSelection
	err := r.collection.FindOne(context.Background(), filter).Decode(&selection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &selection, nil
}

func (r *UserSelectionRepository) CreateOrUpdate(userSelection *domain.UserSelection) error {
	filter := bson.M{"chat_id": userSelection.ChatID}
	update := bson.M{
		"$set": userSelection,
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(context.Background(), filter, update, opts)
	return err
}

package repos

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"golang.org/x/net/context"
	"know_api/internal_1/models"
	"know_api/internal_1/repository/category"
	inst "know_api/pkg_1/db/mongo"
)

type categoryRepository struct {
	db *inst.Category
}

func (c categoryRepository) Create(ctx context.Context, category *models.Category) (*uuid.UUID, error) {
	category.ID, _ = uuid.New()
	_, err := c.db.InsertOne(ctx, category)
	if err != nil {
		return nil, err
	}
	return &category.ID, err
}

func (c categoryRepository) IsExists(ctx context.Context, title string) (bool, error) {
	filter := bson.M{"title": title}
	var result bson.M
	err := c.db.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func NewCategoryMGRepository(db *inst.Storage) category.MGRepository {
	return &categoryRepository{
		db: &db.Category,
	}
}

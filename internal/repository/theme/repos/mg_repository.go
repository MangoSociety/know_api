package repos

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"know_api/internal/models"
	"know_api/internal/repository/theme"
	mg "know_api/pkg/db/mongo"
)

type themeRepository struct {
	db *mg.Theme
}

func (t themeRepository) Create(ctx context.Context, theme *models.Theme) (*uuid.UUID, error) {
	theme.ID, _ = uuid.New()
	_, err := t.db.InsertOne(ctx, theme)
	if err != nil {
		return nil, err
	}
	return &theme.ID, err
}

func (t themeRepository) All(ctx context.Context) ([]models.Theme, error) {
	result, err := t.db.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var data []models.Theme
	if err := result.All(context.Background(), &data); err != nil {
		return nil, err
	}
	return data, nil
}

func NewThemeMGRepository(db *mg.Theme) theme.MgRepository {
	return &themeRepository{
		db: db,
	}
}

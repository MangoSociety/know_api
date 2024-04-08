package repos

import (
	"context"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"know_api/internal/models"
	"know_api/internal/questions"
	"know_api/pkg/db/mongo"
)

type questionsRepository struct {
	db *mongo.Storage
}

func (q questionsRepository) Create(ctx context.Context, question *models.Question) error {
	//TODO implement me
	panic("implement me")
}

func (q questionsRepository) Update(ctx context.Context, question *models.Question) error {
	//TODO implement me
	panic("implement me")
}

func (q questionsRepository) GetByID(ctx context.Context, questionId uuid.UUID) (*models.Question, error) {
	//TODO implement me
	panic("implement me")
}

func (q questionsRepository) GetByTitle(ctx context.Context, title string) (*models.Question, error) {
	//TODO implement me
	panic("implement me")
}

func (q questionsRepository) GetAllByCategoryID(ctx context.Context, categoryId uuid.UUID) ([]*models.Question, error) {
	//TODO implement me
	panic("implement me")
}

func NewQuestionsMGRepository(db *mongo.Storage) questions.MGRepository {
	return &questionsRepository{
		db: db,
	}
}

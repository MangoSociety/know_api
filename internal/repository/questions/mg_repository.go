package questions

import (
	"context"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"know_api/internal/models"
)

type MGRepository interface {
	Create(ctx context.Context, question *models.Question) (*uuid.UUID, error)
	//Update(ctx context.Context, question *models.Question) error
	//GetByID(ctx context.Context, questionId uuid.UUID) (*models.Question, error)
	//GetByTitle(ctx context.Context, title string) (*models.Question, error)
	//GetAllByCategoryID(ctx context.Context, categoryId uuid.UUID) ([]*models.Question, error)
}

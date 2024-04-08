package questions

import (
	"context"
	"know_api/internal/models"
)

type UseCase interface {
	AutoMigration(ctx context.Context) error
	GetQuestionById() (*models.Question, error)
	GetQuestionRandomFromCategory(idCategory int) (*models.Question, error)
	GetQuestionsByCategory(idCategory int) ([]*models.Question, error)
}

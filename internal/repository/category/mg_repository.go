package category

import (
	"context"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"know_api/internal/models"
)

type MGRepository interface {
	Create(ctx context.Context, category *models.Category) (*uuid.UUID, error)
	IsExists(ctx context.Context, title string) (bool, error)
	//All()
}

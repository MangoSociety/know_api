package theme

import (
	"context"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"know_api/internal_1/models"
)

type MgRepository interface {
	Create(ctx context.Context, theme *models.Theme) (*uuid.UUID, error)
	All(ctx context.Context) ([]models.Theme, error)
}

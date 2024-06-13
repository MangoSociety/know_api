package questions

import (
	"context"
	"github.com/google/go-github/github"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"know_api/internal_1/models"
)

type Repository interface {
	GetTree(ctx context.Context) ([]github.TreeEntry, error)
	GetFileContent(ctx context.Context, path string) (string, error)
	IsExistsTheme(ctx context.Context, theme string) (bool, error)
	CreateTheme(ctx context.Context, theme *models.Theme) (*uuid.UUID, error)
	IsExistsCategory(ctx context.Context, category string) (bool, error)
	CreateCategory(ctx context.Context, category *models.Category) (*uuid.UUID, error)
}

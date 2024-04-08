package questions

import (
	"context"
	"github.com/google/go-github/github"
)

// TODO: возможно в будущем стоит предусмотреть выбирать список репозиториев для миграции данных
type GHRepository interface {
	GetTree(ctx context.Context) ([]github.TreeEntry, error)
	GetFileContent(ctx context.Context, path string) (string, error)
}

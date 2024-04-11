package repository

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/go-github/github"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	gh "know_api/clients/github"
	"know_api/internal/models"
	"know_api/internal/questions"
	mg "know_api/pkg/db/mongo"
	"log"
)

type QuestionsRepository struct {
	input  gh.Client
	output mg.Storage
}

func (q QuestionsRepository) GetTree(ctx context.Context) ([]github.TreeEntry, error) {
	tree, response, err := q.input.Client.Git.GetTree(ctx, q.input.Owner, q.input.Repo, q.input.Sha, true)
	if err != nil {
		log.Println("Error get tree from github")
		log.Println(response)
		return nil, err
	}
	return tree.Entries, nil
}

func (q QuestionsRepository) GetFileContent(ctx context.Context, path string) (string, error) {
	// Получаем содержимое файла по пути
	fileContent, _, _, err := q.input.Client.Repositories.GetContents(ctx, q.input.Owner, q.input.Repo, path, nil)
	if err != nil {
		log.Println("Error get file content from github")
		return "", err
	}

	// Проверяем, является ли содержимое файлом
	if fileContent.GetType() == "file" {
		// Декодируем содержимое файла из формата base64
		content, err := base64.StdEncoding.DecodeString(*fileContent.Content)
		if err != nil {
			log.Println("Error decoding file content")
			return "", err
		}

		return string(content), nil
	}

	return "", fmt.Errorf("path is not a file: %s", path)
}

func (q QuestionsRepository) IsExistsTheme(ctx context.Context, theme string) (bool, error) {
	count, err := q.output.Theme.CountDocuments(ctx, bson.D{{"title", theme}})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (q QuestionsRepository) CreateTheme(ctx context.Context, theme *models.Theme) (*uuid.UUID, error) {
	theme.ID, _ = uuid.New()
	_, err := q.output.Theme.InsertOne(ctx, theme)
	if err != nil {
		return nil, err
	}
	return &theme.ID, err
}

func (q QuestionsRepository) IsExistsCategory(ctx context.Context, category string) (bool, error) {
	count, err := q.output.Category.CountDocuments(ctx, bson.D{{"title", category}})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (q QuestionsRepository) CreateCategory(ctx context.Context, category *models.Category) (*uuid.UUID, error) {
	category.ID, _ = uuid.New()
	_, err := q.output.Category.InsertOne(ctx, category)
	if err != nil {
		return nil, err
	}
	return &category.ID, err
}

func NewQuestionsRepository(input gh.Client, output mg.Storage) questions.Repository {
	return QuestionsRepository{
		input:  input,
		output: output,
	}
}

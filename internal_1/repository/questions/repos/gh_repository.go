package repos

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/go-github/github"
	gh "know_api/clients/github"
	"know_api/internal_1/repository/questions"
	"log"
)

type questionsGHRepository struct {
	ghClient *gh.Client
}

func (q questionsGHRepository) GetTree(ctx context.Context) ([]github.TreeEntry, error) {
	tree, response, err := q.ghClient.Client.Git.GetTree(ctx, q.ghClient.Owner, q.ghClient.Repo, q.ghClient.Sha, true)
	if err != nil {
		log.Println("Error get tree from github")
		log.Println(response)
		return nil, err
	}
	return tree.Entries, nil
}

func (q questionsGHRepository) GetFileContent(ctx context.Context, path string) (string, error) {
	// Получаем содержимое файла по пути
	fileContent, _, _, err := q.ghClient.Client.Repositories.GetContents(ctx, q.ghClient.Owner, q.ghClient.Repo, path, nil)
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

func NewQuestionsGHRepository(ghClient gh.Client) questions.GHRepository {
	return &questionsGHRepository{
		ghClient: &ghClient,
	}
}

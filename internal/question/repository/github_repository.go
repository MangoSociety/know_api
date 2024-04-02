package repository

import (
	"encoding/base64"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"log"
	gh_client "read-adviser-bot/clients/github"
	"read-adviser-bot/pkg/utils"
	"read-adviser-bot/storage"
	"read-adviser-bot/storage/mongo"
	"strings"
	"sync"
)

type QuestionsRepo struct {
	client  *gh_client.Client
	storage mongo.Storage
}

type FileResult struct {
	Path    string
	Success bool
}

func NewQuestionsRepo(client *gh_client.Client, storage mongo.Storage) *QuestionsRepo {
	return &QuestionsRepo{
		client:  client,
		storage: storage,
	}
}

func (q *QuestionsRepo) GetFileTree(context context.Context, path string) {
	result, _, err := q.client.Client.Git.GetTree(context, q.client.Owner, q.client.Repo, q.client.Sha, true)
	if err != nil {
		fmt.Printf("Problem in getting repository information %v\n", err)
	}

	data := filter(result.Entries, path)

	var wg sync.WaitGroup
	resultsChan := make(chan FileResult, len(data))

	for _, item := range data {
		wg.Add(1)
		go q.processFile(context, item, resultsChan, &wg)
	}

	// Ждем завершения всех горутин
	wg.Wait()
	close(resultsChan)

	// Собираем статистику
	successCount := 0
	var failedFiles []string
	for result := range resultsChan {
		if result.Success {
			successCount++
		} else {
			failedFiles = append(failedFiles, result.Path)
		}
	}

	fmt.Printf("Успешно обработано файлов: %d\n", successCount)
	if len(failedFiles) > 0 {
		fmt.Println("Не удалось обработать файлы:")
		for _, file := range failedFiles {
			fmt.Println(file)
		}
	}
}

func (q *QuestionsRepo) processFile(ctx context.Context, path string, resultsChan chan<- FileResult, wg *sync.WaitGroup) {
	defer wg.Done()

	currentFile, _, _, _ := q.client.Client.Repositories.GetContents(
		context.Background(),
		q.client.Owner,
		q.client.Repo,
		path,
		&github.RepositoryContentGetOptions{
			Ref: q.client.Sha,
		},
	)

	if currentFile == nil {
		resultsChan <- FileResult{Path: path, Success: false}
		return
	}

	contentBytes, err := base64.StdEncoding.DecodeString(*currentFile.Content)
	if err != nil {
		fmt.Printf("Problem in processing md file with %v\n", err)
	}

	contentStr := string(contentBytes)

	result, err := utils.ProcessingQuestion(contentStr)
	if err != nil {
		resultsChan <- FileResult{Path: path, Success: false}
		return
	}

	//if result.Title {
	//
	//}

	note := storage.Note{
		Title:    result.Title,
		Sphere:   result.Sphere,
		Category: result.Theme,
		Content:  result.Content,
	}

	if err := q.storage.SaveNote(ctx, &note); err != nil {
		resultsChan <- FileResult{Path: path, Success: false}
	} else {
		resultsChan <- FileResult{Path: path, Success: true}
	}

	// Здесь логика обработки файла и сохранения его содержимого в БД
	// Если все успешно, отправляем в канал результат с Success: true
	// Если есть ошибка, отправляем в канал результат с Success: false
	// Пример:
	// if err := saveToDB(path); err != nil {
	//	resultsChan <- FileResult{Path: path, Success: false}
	// } else {
	//	resultsChan <- FileResult{Path: path, Success: true}
	// }
}

func (q *QuestionsRepo) GetRandomQuestion(path string) string {
	currentFile, _, _, _ := q.client.Client.Repositories.GetContents(
		context.Background(),
		q.client.Owner,
		q.client.Repo,
		path,
		&github.RepositoryContentGetOptions{
			Ref: q.client.Sha,
		},
	)

	if currentFile == nil {
		return "null file"
	}

	contentBytes, err := base64.StdEncoding.DecodeString(*currentFile.Content)
	if err != nil {
		fmt.Printf("Problem in processing md file with %v\n", err)
	}

	contentStr := string(contentBytes)

	res1, _ := utils.ProcessingQuestion(contentStr)
	res2 := "Название вопроса: " + res1.Title + "\n"

	return res2 + res1.Content
}

func (q *QuestionsRepo) GetRandomQuestionStruct(path string) storage.Note {
	currentFile, _, _, err := q.client.Client.Repositories.GetContents(
		context.Background(),
		q.client.Owner,
		q.client.Repo,
		path,
		&github.RepositoryContentGetOptions{
			Ref: q.client.Sha,
		},
	)

	if err != nil {
		log.Println("Error in getting content from github", err)
	}

	if currentFile == nil {
		panic("null file")
	}

	contentBytes, err := base64.StdEncoding.DecodeString(*currentFile.Content)
	if err != nil {
		fmt.Printf("Problem in processing md file with %v\n", err)
	}

	contentStr := string(contentBytes)

	result, _ := utils.ProcessingQuestion(contentStr)
	return storage.Note{
		Title:    result.Title,
		Sphere:   result.Sphere,
		Category: result.Theme,
		Content:  result.Content,
	}
}

func (q *QuestionsRepo) GetInterviewQuestions(topic string) ([]string, error) {
	result, _, err := q.client.Client.Git.GetTree(context.Background(), q.client.Owner, q.client.Repo, q.client.Sha, true)
	if err != nil {
		fmt.Printf("Problem in getting repository information %v\n", err)
		return nil, err
	}
	return filter(result.Entries, topic), nil

}

// filter get tree file as string's list by topic's name,
// example topic: "android/interview_questions/"
func filter(data []github.TreeEntry, topic string) []string {
	var result []string
	for _, item := range data {
		if strings.Contains(*item.Path, topic) {
			result = append(result, *item.Path)
		}
	}
	return result
}

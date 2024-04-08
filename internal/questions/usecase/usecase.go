package usecase

import (
	"context"
	"github.com/google/go-github/github"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"know_api/internal/models"
	"know_api/internal/questions"
	"regexp"
	"strings"
	"sync"
)

type questionsUseCase struct {
	ghRepo questions.GHRepository
	mgRepo questions.MGRepository
}

type UniqueNote struct {
	Category UniqueCategory
	Theme    string
	Title    string
	Content  string
	External []string
	Internal []string
}

type UniqueCategory struct {
	Id    uuid.UUID
	Title string
}

func (q questionsUseCase) AutoMigration(ctx context.Context) error {
	tree, err := q.ghRepo.GetTree(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	successCount := 0
	failureCount := 0
	successChan := make(chan bool)
	failureChan := make(chan bool)

	for _, entry := range tree {
		if *entry.Type == "blob" && strings.Contains(*entry.Path, "android/interview_questions") {
			wg.Add(1)
			go func(entry github.TreeEntry) {
				defer wg.Done()
				data, err := q.ghRepo.GetFileContent(ctx, *entry.Path)
				if err != nil {
					failureChan <- true
					return
				}
				content := strings.ReplaceAll(data, "\n\n", "\n")
				note, err := parseUniqueNote(content)
				if note.Category.Title != "" && note.Theme != "" && note.Title != "" && note.Content != "" {
					// TODO: Проверка категории и ее актуализация + сохранение в базу
					successChan <- true
				} else {
					failureChan <- true
				}
			}(entry)
		}
	}

	go func() {
		wg.Wait()
		close(successChan)
		close(failureChan)
	}()

	for range successChan {
		successCount++
	}

	for range failureChan {
		failureCount++
	}

	return nil
}

func parseUniqueNote(text string) (*UniqueNote, error) {
	note := &UniqueNote{}

	// Ищем тему
	themeRegex := regexp.MustCompile(`Theme : (.*)`)
	themeMatch := themeRegex.FindStringSubmatch(text)
	if len(themeMatch) > 1 {
		note.Theme = strings.TrimSpace(themeMatch[1])
	}

	// Ищем заголовок
	titleRegex := regexp.MustCompile(`Title: (.*)`)
	titleMatch := titleRegex.FindStringSubmatch(text)
	if len(titleMatch) > 1 {
		note.Title = strings.TrimSpace(titleMatch[1])
	}

	// Ищем содержимое между маркерами ### Content и ### External Link
	contentRegex := regexp.MustCompile(`(?s)### Content\n(.*?)\n### External Link`)
	contentMatch := contentRegex.FindStringSubmatch(text)
	if len(contentMatch) > 1 {
		// Удаляем пустые строки и блоки кода, которые не содержат текста
		content := strings.TrimSpace(contentMatch[1])
		content = regexp.MustCompile(`(?m)^\s*$\n?`).ReplaceAllString(content, "")
		note.Content = content
	} else {
		// Если содержимое отсутствует, присваиваем пустую строку
		note.Content = ""
	}

	// Ищем внешние ссылки
	externalRegex := regexp.MustCompile(`### External Link\n\n(.*?)\n\n### Internal Link`)
	externalMatch := externalRegex.FindStringSubmatch(text)
	if len(externalMatch) > 1 {
		note.External = strings.Split(strings.TrimSpace(externalMatch[1]), "\n")
	}

	// Ищем внутренние ссылки
	internalRegex := regexp.MustCompile(`### Internal Link\n\n(.*?)\n\n`)
	internalMatch := internalRegex.FindStringSubmatch(text)
	if len(internalMatch) > 1 {
		note.Internal = strings.Split(strings.TrimSpace(internalMatch[1]), "\n")
	}

	return note, nil
}

func (q questionsUseCase) GetQuestionById() (*models.Question, error) {
	//TODO implement me
	panic("implement me")
}

func (q questionsUseCase) GetQuestionRandomFromCategory(idCategory int) (*models.Question, error) {
	//TODO implement me
	panic("implement me")
}

func (q questionsUseCase) GetQuestionsByCategory(idCategory int) ([]*models.Question, error) {
	//TODO implement me
	panic("implement me")
}

func NewQuestionsUseCase(ghRepo questions.GHRepository, mgRepo questions.MGRepository) questions.UseCase {
	return &questionsUseCase{
		ghRepo: ghRepo,
		mgRepo: mgRepo,
	}
}

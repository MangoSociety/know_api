package usecase

import (
	"context"
	"github.com/google/go-github/github"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"know_api/internal_1/models"
	"know_api/internal_1/questions"
	"log"
	"regexp"
	"strings"
	"sync"
)

type questionsUseCase struct {
	repository questions.Repository
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
	tree, err := q.repository.GetTree(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var processWg sync.WaitGroup // Для ожидания обработки категорий
	successCount := 0
	failureCount := 0
	successChan := make(chan bool)
	failureChan := make(chan bool)

	categoryChan := make(chan string)
	uniqueCategories := make(map[string]struct{})

	processWg.Add(1)
	go func() {
		defer processWg.Done()
		for category := range categoryChan {
			uniqueCategories[category] = struct{}{}
			log.Println("Added category:", category)
		}
	}()

	for _, entry := range tree {
		if *entry.Type == "blob" && strings.Contains(*entry.Path, "android/interview_questions") {
			wg.Add(1)
			go func(entry github.TreeEntry) {
				defer wg.Done()
				data, err := q.repository.GetFileContent(ctx, *entry.Path)
				if err != nil {
					log.Println("Error get file content from github for path:", *entry.Path)
					failureChan <- true
					return
				}
				content := strings.ReplaceAll(data, "\n\n", "\n")
				note, err := parseUniqueNote(content)
				if err != nil {
					log.Println("Error parse file content:", err)
					failureChan <- true
					return
				}
				//log.Println("parse note:", note)
				if note.Category.Title != "" && note.Theme.Title != "" && note.Title != "" && note.Content != "" {
					//isExistTheme, err := q.repository.IsExistsTheme(ctx, note.Theme.Title)
					//if err != nil {
					//	failureChan <- true
					//	return
					//}
					//if !isExistTheme {
					//	_, err := q.repository.CreateTheme(ctx, &note.Theme)
					//	if err != nil {
					//		failureChan <- true
					//		return
					//	}
					//}

					categoryChan <- note.Category.Title

					//isExistCategory, err := q.repository.IsExistsCategory(ctx, note.Category.Title)
					//if err != nil {
					//	failureChan <- true
					//	return
					//}
					//if !isExistCategory {
					//	_, err := q.repository.CreateCategory(ctx, &note.Category)
					//	if err != nil {
					//		failureChan <- true
					//		return
					//	}
					//}

					successChan <- true
				} else {
					log.Println("Error parse file content: not all fields are filled")
					failureChan <- true
				}
			}(entry)
		}
	}

	for range successChan {
		log.Println("Success ", successCount)
		successCount++
	}

	for range failureChan {
		log.Println("Failure ", failureCount)
		failureCount++
	}

	wg.Wait()           // Ожидаем обработку всех записей
	close(categoryChan) // Закрываем канал категорий после завершения всех горутин
	processWg.Wait()    // Ожидаем завершения обработки уникальных категорий

	// Выводим уникальные категории
	log.Println("Уникальные категории:")
	for item := range uniqueCategories {
		log.Println("Category:", item)
	}

	close(successChan)
	close(failureChan)
	log.Println("Migration finished\nSuccess:", successCount, "\nFailure:", failureCount)

	// Ваш код для закрытия остальных каналов и вывода итогов...
	return nil
}

//func (q questionsUseCase) AutoMigration(ctx context.Context) error {
//	tree, err := q.repository.GetTree(ctx)
//	if err != nil {
//		return err
//	}
//
//	var wg sync.WaitGroup
//	successCount := 0
//	failureCount := 0
//	successChan := make(chan bool)
//	failureChan := make(chan bool)
//
//	categoryChan := make(chan string)
//
//	uniqueCategories := make(map[string]struct{})
//
//	go func() {
//		for category := range categoryChan {
//			uniqueCategories[category] = struct{}{}
//			log.Println("Added category:", category)
//		}
//	}()
//
//	for _, entry := range tree {
//		if *entry.Type == "blob" && strings.Contains(*entry.Path, "android/interview_questions") {
//			wg.Add(1)
//			go func(entry github.TreeEntry) {
//				defer wg.Done()
//				data, err := q.repository.GetFileContent(ctx, *entry.Path)
//				if err != nil {
//					log.Println("Error get file content from github for path:", *entry.Path)
//					failureChan <- true
//					return
//				}
//				content := strings.ReplaceAll(data, "\n\n", "\n")
//				note, err := parseUniqueNote(content)
//				if err != nil {
//					log.Println("Error parse file content:", err)
//					failureChan <- true
//					return
//				}
//				//log.Println("parse note:", note)
//				if note.Category.Title != "" && note.Theme.Title != "" && note.Title != "" && note.Content != "" {
//					//isExistTheme, err := q.repository.IsExistsTheme(ctx, note.Theme.Title)
//					//if err != nil {
//					//	failureChan <- true
//					//	return
//					//}
//					//if !isExistTheme {
//					//	_, err := q.repository.CreateTheme(ctx, &note.Theme)
//					//	if err != nil {
//					//		failureChan <- true
//					//		return
//					//	}
//					//}
//
//					categoryChan <- note.Category.Title
//
//					//isExistCategory, err := q.repository.IsExistsCategory(ctx, note.Category.Title)
//					//if err != nil {
//					//	failureChan <- true
//					//	return
//					//}
//					//if !isExistCategory {
//					//	_, err := q.repository.CreateCategory(ctx, &note.Category)
//					//	if err != nil {
//					//		failureChan <- true
//					//		return
//					//	}
//					//}
//
//					successChan <- true
//				} else {
//					log.Println("Error parse file content: not all fields are filled")
//					failureChan <- true
//				}
//			}(entry)
//		}
//	}
//
//	go func() {
//
//	}()
//
//	for range successChan {
//		log.Println("Success ", successCount)
//		successCount++
//	}
//
//	for range failureChan {
//		log.Println("Failure ", failureCount)
//		failureCount++
//	}
//
//	wg.Wait()
//	close(categoryChan)
//
//	// Выводим уникальные категории
//	fmt.Println("Уникальные категории:")
//	for item := range uniqueCategories {
//		log.Println("Category:", item)
//	}
//	close(successChan)
//	close(failureChan)
//	log.Println("Migration finished\nSuccess:", successCount, "\nFailure:", failureCount)
//
//	return nil
//}

func parseUniqueNote(text string) (*models.Question, error) {
	note := &models.Question{}

	// Ищем тему
	categoryRegex := regexp.MustCompile(`Theme : (.*)`)
	categoryMatch := categoryRegex.FindStringSubmatch(text)
	if len(categoryMatch) > 1 {
		note.Category = models.Category{
			Title: strings.TrimSpace(categoryMatch[1]),
		}
	}

	// Ищем тему
	//themeRegex := regexp.MustCompile(`Sphere : (.*)`)
	//themeMatch := themeRegex.FindStringSubmatch(text)
	//if len(themeMatch) > 1 {
	//	note.Theme = models.Theme{
	//		Title: strings.TrimSpace(themeMatch[1]),
	//	}
	//}

	note.Theme = models.Theme{
		Title: "android",
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

func NewQuestionsUseCase(repository questions.Repository) questions.UseCase {
	return &questionsUseCase{
		repository: repository,
	}
}

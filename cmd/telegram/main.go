package main

import (
	"context"
	"fmt"
	tgClient "github.com/MangoSociety/know_api/clients/telegram"
	categoriesRepo "github.com/MangoSociety/know_api/internal/categories/repository"
	categoriesService "github.com/MangoSociety/know_api/internal/categories/service"
	configReal "github.com/MangoSociety/know_api/internal/config"
	repository2 "github.com/MangoSociety/know_api/internal/notes/repository"
	service2 "github.com/MangoSociety/know_api/internal/notes/service"
	spheresRepo "github.com/MangoSociety/know_api/internal/spheres/repository"
	spheresService "github.com/MangoSociety/know_api/internal/spheres/service"
	statitsticsRepository "github.com/MangoSociety/know_api/internal/statistics/repository"
	statisticsService "github.com/MangoSociety/know_api/internal/statistics/service"
	"github.com/MangoSociety/know_api/internal/telegram/consumer/event_consumer"
	"github.com/MangoSociety/know_api/internal/telegram/events/telegram"
	"github.com/MangoSociety/know_api/internal/user/repository"
	"github.com/MangoSociety/know_api/internal/user/service"
	"github.com/MangoSociety/know_api/pkg/mongodb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/labstack/gommon/log"
	"github.com/russross/blackfriday/v2"
	"strings"
	"time"
)

/*
Структура работы бота
- авто синк с гитхабом
*
*/

//cfg := config.MustLoad()
//
//bot, err := tg_bot_api.NewBotAPI(cfg.TgToken)
//if err != nil {
//	log.Panic(err)
//}
//bot.Debug = true
//
//storage := mongo.New(cfg.MongoConnect, 10*time.Second)
//gh_client := gh_instance.NewGithubClient(cfg.Github.Token, cfg.Github.Owner, cfg.Github.Repo, cfg.Github.Sha)
//
////ghRepo := repos2.NewQuestionsGHRepository(*gh_client)
////mgRepo := repos2.NewQuestionsMGRepository(&storage)
//repository := repository2.NewQuestionsRepository(*gh_client, storage)
//questionsUseCase := usecase.NewQuestionsUseCase(repository)
//
//eventsProcessor := telegram.NewProcessor(
//	tg_client.NewTelegramClient("api.telegram.org", cfg.TgToken, bot),
//	questionsUseCase,
//)
//
//log.Print("telegram processor started")
//
//consumer := event_consumer.NewConsumer(eventsProcessor, eventsProcessor, 100) //  NewConsumer(eventsProcessor, eventsProcessor, batchSize)
//
//ctx, cancel := context.WithCancel(context.Background())
//defer cancel()
//if err := consumer.Start(ctx); err != nil {
//	log.Fatal("service is stopped", err)
//}

func main() {
	cfg := configReal.MustLoadConfig()
	log.Printf("Config: %s", cfg.Mongo.Connect)
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Panic(err)
	}

	// Инициализация MongoDB
	dbClient, err := mongodb.NewClient(cfg.Mongo.Connect)
	if err != nil {
		log.Fatal(err)
	}

	// Инициализация репозиториев
	noteRepo := repository2.NewNoteRepository(dbClient)
	categoryRepo := categoriesRepo.NewCategoryRepository(dbClient)
	sphereRepo := spheresRepo.NewSphereRepository(dbClient)
	userSelectionRepo := repository.NewUserSelectionRepository(dbClient)
	statisticsRepo := statitsticsRepository.NewStatisticsRepository(dbClient)

	// Инициализация сервисов
	noteService := service2.NewNoteService(noteRepo)
	categoryService := categoriesService.NewCategoryService(categoryRepo)
	sphereService := spheresService.NewSphereService(sphereRepo)
	userSelectionService := service.NewUserSelectionService(userSelectionRepo)
	statService := statisticsService.NewStatisticsService(statisticsRepo)

	eventsProcessor := telegram.NewProcessor(
		tgClient.NewTelegramClient("api.telegram.org", cfg.Telegram.Token, bot),
		noteService,
		categoryService,
		sphereService,
		userSelectionService,
		statService,
	)
	consumer := event_consumer.NewConsumer(eventsProcessor, eventsProcessor, 100) //  NewConsumer(eventsProcessor, eventsProcessor, batchSize)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := consumer.Start(ctx); err != nil {
		log.Fatal("service is stopped", err)
	}
}

const testString = `#готово

Theme : #common
Title: Расскажи что такое SharedPreferences Какие данные можно хранить Какие плюсы и минусы
Sphere: #android

### Content

### Какие данные можно хранить

С SharedPreferences, вы можете хранить следующие типы данных:

- boolean
- float
- int
- long
- String
- Set<String> (используется для хранения коллекций строк, например, списка значений)

Эти типы данных покрывают большинство нужд приложения в сохранении простых конфигурационных параметров и пользовательских предпочтений.

### Плюсы SharedPreferences

1. **Простота использования**: Интерфейс SharedPreferences прост и интуитивно понятен, что делает его легким в использовании для хранения и извлечения простых данных.
2. **Легкость интеграции**: Он хорошо интегрирован в Android SDK, предоставляя прямой и эффективный способ сохранения легковесных данных.
3. **Быстрый доступ**: Доступ к данным, хранящимся в SharedPreferences, осуществляется быстро, что делает его подходящим для хранения данных, необходимых при старте приложения.
4. **Поддержка асинхронного сохранения**: С API 9 (Android 2.3, Gingerbread) и выше, SharedPreferences предлагает метод apply(), который асинхронно сохраняет изменения, минимизируя задержки UI.

### Минусы SharedPreferences

1. **Ограниченный объем и типы данных**: SharedPreferences подходит только для примитивных типов данных и не предназначен для хранения сложных объектов или больших объемов данных.
2. **Безопасность**: Данные, сохраненные в SharedPreferences, хранятся в виде обычных файлов XML без шифрования, что делает их уязвимыми для атак, если устройство скомпрометировано.
3. **Отсутствие поддержки структурированных данных**: SharedPreferences не подходит для хранения структурированных данных или реализации сложных иерархий настроек.
4. **Проблемы с многопоточностью**: При неправильном использовании может возникнуть состояние гонки или другие проблемы с многопоточностью, особенно если одновременно производится чтение и запись из разных потоков.

SharedPreferences является удобным инструментом для хранения небольших объемов данных, таких как настройки пользователя или простые флаги состояния. Однако для более сложных или объемных данных следует рассмотреть другие варианты хранения данных, такие как базы данных SQLite или хранилище на основе файлов с использованием внутренней или внешней памяти.

### External Link

-

### Internal Link

- ....
`

// Section представляет раздел внутри контента Markdown
type Section struct {
	Title   string `bson:"title"`
	Content string `bson:"content"`
}

// MarkdownData структура для хранения данных из Markdown
type MarkdownData struct {
	Theme        string    `bson:"theme"`
	Title        string    `bson:"title"`
	Sphere       string    `bson:"sphere"`
	Sections     []Section `bson:"sections"`
	ExternalLink string    `bson:"external_link"`
	InternalLink string    `bson:"internal_link"`
	Date         time.Time `bson:"date"`
}

//func main() {
//	data, _ := parseString(testString)
//	fmt.Println(data)
//}

// parseString разбирает строку Markdown и возвращает структуру MarkdownData
func parseString(content string) (*MarkdownData, error) {
	output := blackfriday.Run([]byte(content))
	sections, theme, title, sphere, externalLink, internalLink := parseHTMLContent(string(output))

	return &MarkdownData{
		Theme:        theme,
		Title:        title,
		Sphere:       sphere,
		Sections:     sections,
		ExternalLink: externalLink,
		InternalLink: internalLink,
		Date:         time.Now(),
	}, nil
}

// parseHTMLContent разбирает HTML-содержимое, сгенерированное из Markdown, и возвращает секции
func parseHTMLContent(htmlContent string) ([]Section, string, string, string, string, string) {
	var sections []Section
	var currentSection *Section
	var theme, title, sphere, externalLink, internalLink string

	lines := strings.Split(htmlContent, "\n")
	for _, line := range lines {
		if strings.Contains(line, "<h1>") && strings.Contains(line, "</h1>") {
			theme = extractContent(line, "h1")
		} else if strings.Contains(line, "<h2>") && strings.Contains(line, "</h2>") {
			title = extractContent(line, "h2")
		} else if strings.Contains(line, "<h3>") && strings.Contains(line, "</h3>") {
			sphere = extractContent(line, "h3")
		} else if strings.Contains(line, "<h4>") && strings.Contains(line, "</h4>") {
			if currentSection != nil {
				sections = append(sections, *currentSection)
			}
			currentSection = &Section{
				Title:   extractContent(line, "h4"),
				Content: "",
			}
		} else if strings.Contains(line, "<ul>") && strings.Contains(line, "</ul>") {
			if currentSection != nil {
				currentSection.Content += formatList(line, "ul") + "\n"
			}
		} else if strings.Contains(line, "<ol>") && strings.Contains(line, "</ol>") {
			if currentSection != nil {
				currentSection.Content += formatList(line, "ol") + "\n"
			}
		} else if strings.Contains(line, "<strong>") && strings.Contains(line, "</strong>") {
			if currentSection != nil {
				currentSection.Content += formatBoldText(line) + "\n"
			}
		} else if strings.Contains(line, "<em>") && strings.Contains(line, "</em>") {
			if currentSection != nil {
				currentSection.Content += formatItalicText(line) + "\n"
			}
		} else if strings.Contains(line, "<a href=") {
			if strings.Contains(line, "External Link") {
				externalLink = extractLink(line)
			} else if strings.Contains(line, "Internal Link") {
				internalLink = extractLink(line)
			}
		} else {
			if currentSection != nil {
				currentSection.Content += line + "\n"
			}
		}
	}

	if currentSection != nil {
		sections = append(sections, *currentSection)
	}

	return sections, theme, title, sphere, externalLink, internalLink
}

// extractContent извлекает содержимое из HTML-тега
func extractContent(line, tag string) string {
	openTag := fmt.Sprintf("<%s>", tag)
	closeTag := fmt.Sprintf("</%s>", tag)
	return strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(line, openTag, ""), closeTag, ""))
}

// formatList форматирует список в читабельный вид
func formatList(line, tag string) string {
	openTag := fmt.Sprintf("<%s>", tag)
	closeTag := fmt.Sprintf("</%s>", tag)
	itemOpenTag := "<li>"
	itemCloseTag := "</li>"
	formatted := strings.ReplaceAll(line, openTag, "")
	formatted = strings.ReplaceAll(formatted, closeTag, "")
	formatted = strings.ReplaceAll(formatted, itemOpenTag, "- ")
	formatted = strings.ReplaceAll(formatted, itemCloseTag, "\n")
	return formatted
}

// formatBoldText форматирует жирный текст
func formatBoldText(line string) string {
	return strings.ReplaceAll(strings.ReplaceAll(line, "<strong>", "**"), "</strong>", "**")
}

// formatItalicText форматирует курсивный текст
func formatItalicText(line string) string {
	return strings.ReplaceAll(strings.ReplaceAll(line, "<em>", "*"), "</em>", "*")
}

// extractLink извлекает ссылку из HTML-тега <a>
func extractLink(line string) string {
	start := strings.Index(line, "href=") + len("href=") + 1
	end := strings.Index(line[start:], "\"") + start
	return strings.TrimSpace(line[start:end])
}

package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/go-github/github"
	"go/format"
	"log"
	ghClient "read-adviser-bot/clients/github"
	telegram2 "read-adviser-bot/clients/telegram"
	"read-adviser-bot/config"
	event_consumer "read-adviser-bot/consumer/event-consumer"
	"read-adviser-bot/events/telegram"
	quest_repo "read-adviser-bot/internal/question/repository"
	quest_usecase "read-adviser-bot/internal/question/usecase"
	"read-adviser-bot/storage/mongo"
	"regexp"
	"strings"
	"time"
)

//func main() {
//	fmt.Println("start app")
//}

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

type Package struct {
	FullName      string
	Description   string
	StarsCount    int
	ForksCount    int
	LastUpdatedBy string
}

type createTree struct {
	BaseTree string        `json:"base_tree,omitempty"`
	Entries  []interface{} `json:"tree"`
}

type OneFile struct {
	file []string
}

func main() {

	var owner = "MangoSociety"
	var repo = "golang_road"
	var sha = "1fbd0e9ed5ee1239996720d79435b7a2feb5d507"

	cfg := config.MustLoad()

	bot, err := tgbotapi.NewBotAPI(cfg.TgBotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	storage := mongo.New(cfg.MongoConnectionString, 10*time.Second)
	gh_client := ghClient.NewGithubClient("whoiswho", owner, repo, sha)
	var repository = quest_repo.NewQuestionsRepo(gh_client, storage)
	var useCase = quest_usecase.NewQuestionsUseCase(repository)

	//repository.GetFileTree(context.Background(), "android/interview_questions/")

	log.Println("storage created", storage)

	eventsProcessor := telegram.NewProcessor(
		telegram2.NewTelegramClient(tgBotHost, cfg.TgBotToken, bot),
		*useCase,
		storage,
	)

	log.Print("service started")

	consumer := event_consumer.NewConsumer(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func filter(data []github.TreeEntry) []string {
	var result []string
	for _, item := range data {
		if strings.Contains(*item.Path, "android/interview_questions/") {
			result = append(result, *item.Path+"\n")
		}
	}
	return result
}

func processingContentMdFile(data string) (string, error) {
	return extractSubstring(data, "Sphere", "Theme")
}

func getSphere(data string) string {
	sphereFull, err := extractSubstring(data, "Sphere", "Theme")
	if err != nil {
		return ""
	}
	return extractTextBetweenBrackets(sphereFull)
}

func getTheme(data string) string {
	sphereFull, err := extractSubstring(data, "Theme", "Title")
	if err != nil {
		return ""
	}
	return extractTextBetweenBrackets(sphereFull)
}

func getTitle(data string) (string, error) {
	substring, err := extractSubstring(data, "Title:", "### Content")
	if err != nil {
		return "", fmt.Errorf("Title not found")
	}
	endIndex := strings.Index(substring, "\n")
	return substring[len("Title: "):endIndex], nil
}

func getContent(data string) (string, error) {
	startText := "### Content"
	endText := "### External"
	substring, err := extractSubstring(data, startText, endText)
	content := substring[len(startText):(len(substring) - len(endText))]
	if err != nil {
		return "", fmt.Errorf("Content not found")
	}
	return content, err
}

func contentOutput(data []Block) {
	for _, item := range data {
		if item.Type == "text" {
			//fmt.Printf("{\nKey: %v,\nValue: %v\n}\n,", item.Type, item.Text)
		} else {
			//var result = strings.ReplaceAll(item.Text, "\t", "\n")
			fmt.Println(item.Text)
		}

	}
}

func extractSubstring(mdContent, start, end string) (string, error) {
	startIndex := strings.Index(mdContent, start)
	if startIndex == -1 {
		return "", fmt.Errorf("Start string not found")
	}

	endIndex := strings.Index(mdContent[startIndex:], end)
	if endIndex == -1 {
		return "", fmt.Errorf("End string not found")
	}

	substring := mdContent[startIndex : startIndex+endIndex+len(end)]
	//fmt.Println(strings.Split(substring, "\n"))
	return substring, nil
}

func extractTextBetweenBrackets(input string) string {
	re := regexp.MustCompile(`\[\[([^]]+)\]\]`)

	match := re.FindStringSubmatch(input)

	if len(match) < 2 {
		return ""
	}

	return match[1]
}

type Block struct {
	Type string
	Text string
}

func analyzeBlocks(input string) []Block {
	var blocks []Block
	blockType := "text"
	var currentBlockText strings.Builder

	lines := strings.Split(input, "\n")
	isInCodeBlock := false

	for _, line := range lines {

		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(line, "```kotlin") {
			// Начало блока кода
			isInCodeBlock = true
			if currentBlockText.Len() > 0 {
				blocks = appendIfNotEmpty(blocks, Block{Type: blockType, Text: currentBlockText.String()})
				currentBlockText.Reset()
			}
			blockType = "code"
		} else if strings.HasPrefix(trimmedLine, "```") {
			// Конец блока кода
			isInCodeBlock = false
			if currentBlockText.Len() > 0 {
				blocks = appendIfNotEmpty(blocks, Block{Type: blockType, Text: cleanText(currentBlockText.String())})
				currentBlockText.Reset()
			}
			blockType = "text"
		} else {
			// Обычная строка
			if isInCodeBlock {
				currentBlockText.WriteString(line + "\n")
			} else {
				if currentBlockText.Len() > 0 {
					resultText := ""
					if blockType == "text" {
						resultText = cleanText(currentBlockText.String())
					} else {
						resultText = currentBlockText.String()
					}
					blocks = appendIfNotEmpty(blocks, Block{Type: blockType, Text: resultText})
					currentBlockText.Reset()
				}
				blockType = "text"
				currentBlockText.WriteString(line)
			}
		}
	}

	// Добавить последний блок
	if currentBlockText.Len() > 0 {
		blocks = appendIfNotEmpty(blocks, Block{Type: blockType, Text: cleanText(currentBlockText.String())})
	}

	return blocks
}

func appendIfNotEmpty(blocks []Block, block Block) []Block {
	// Проверить, содержит ли значение символы
	if strings.TrimSpace(block.Text) != "" {
		return append(blocks, block)
	}
	return blocks
}

func cleanText(text string) string {
	// Удалить пробелы с конца и лишние переносы строк
	//return strings.TrimSpace(strings.ReplaceAll(text, "\n", " "))
	return text
}

func formatKotlinCode(code string) (string, error) {
	formattedCode, err := format.Source([]byte(code))
	if err != nil {
		return "", fmt.Errorf("ошибка при форматировании кода: %w", err)
	}

	return string(formattedCode), nil
}
func replaceMultipleSpacesWithNewline(input string) string {
	// Создаем регулярное выражение для поиска 3 и более подряд идущих пробелов
	//re := regexp.MustCompile(` {3,}`)

	// Заменяем найденные подряд идущие пробелы на символ новой строки
	//result := re.ReplaceAllString(input, "\n")

	return input
}

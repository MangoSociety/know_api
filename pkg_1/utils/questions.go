package utils

import (
	"fmt"
	"know_api/internal_1/models"
	"regexp"
	"strings"
)

type Question struct {
	Title   string
	Sphere  string
	Theme   string
	Content string
}

func ProcessingQuestion(data string) (Question, error) {
	sphere, _ := getSphereFromQuestionMd(data)
	theme, _ := getThemeFromQuestionMd(data)
	title, _ := getTitleFromQuestionMd(data)
	content, _ := getContentFromQuestionMd(data)

	if title == "" || sphere == "" || content == "" {
		return Question{}, fmt.Errorf("Title not found")
	}

	return Question{
		Title:   title,
		Sphere:  sphere,
		Theme:   theme,
		Content: content,
	}, nil
}

func getTitleFromQuestionMd(data string) (string, error) {
	substring, err := ExtractSubstring(data, "Title:", "### Content")
	if err != nil {
		return "", fmt.Errorf("Title not found")
	}
	endIndex := strings.Index(substring, "\n")
	return substring[len("Title: "):endIndex], nil
}

func getSphereFromQuestionMd(data string) (string, error) {
	sphereFull, err := ExtractSubstring(data, "Sphere", "Theme")
	if err != nil {
		return "", err
	}
	return ExtractTextBetweenBrackets(sphereFull)
}

func getThemeFromQuestionMd(data string) (string, error) {
	sphereFull, err := ExtractSubstring(data, "Theme", "Title")
	if err != nil {
		return "", err
	}
	return ExtractTextBetweenBrackets(sphereFull)
}

func getContentFromQuestionMd(data string) (string, error) {
	startText := "### Content"
	endText := "### External"
	substring, err := ExtractSubstring(data, startText, endText)
	endIndex := len(substring) - len(endText)
	if endIndex < 0 {
		return "error with getContentFromQuestion", nil
	}
	content := substring[len(startText):(len(substring) - len(endText))]
	if err != nil {
		return "", fmt.Errorf("Content not found")
	}
	return content, err
}

func parseUniqueNote(text string) (*models.Question, error) {
	note := &models.Question{}

	// Ищем категорию
	categoryRegex := regexp.MustCompile(`Theme : (.*)`)
	categoryMatch := categoryRegex.FindStringSubmatch(text)
	if len(categoryMatch) > 1 {
		note.Theme = models.Theme{
			Title: strings.TrimSpace(categoryMatch[1]),
		}
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

package utils

import (
	"fmt"
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

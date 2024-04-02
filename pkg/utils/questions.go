package utils

import (
	"fmt"
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

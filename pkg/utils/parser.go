package utils

import (
	"bufio"
	"github.com/MangoSociety/know_api/internal/notes/models"
	"strings"
)

// ParseMarkdownFile парсит содержимое MD файла и возвращает Note
func ParseMarkdownFile(content string) models.Note {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var note models.Note

	var externalLink, internalLink []string

	for scanner.Scan() {
		line := scanner.Text()

		// Парсинг заголовка
		if strings.HasPrefix(line, "Title:") {
			note.Title = strings.TrimSpace(strings.TrimPrefix(line, "Title:"))
		}

		// Парсинг темы
		if strings.HasPrefix(line, "Theme :") {
			themes := strings.Fields(strings.TrimPrefix(line, "Theme :"))
			for _, theme := range themes {
				note.Theme = append(note.Theme, strings.TrimPrefix(theme, "#"))
			}
		}

		// Парсинг сферы
		if strings.HasPrefix(line, "Sphere:") {
			note.Sphere = strings.TrimSpace(strings.TrimPrefix(line, "Sphere:"))
			note.Sphere = strings.TrimPrefix(note.Sphere, "#")
		}

		// Парсинг контента
		if strings.HasPrefix(line, "### Content") {
			var contentBuilder strings.Builder
			for scanner.Scan() {
				contentLine := scanner.Text()
				if strings.HasPrefix(contentLine, "### External Link") || strings.HasPrefix(contentLine, "### Internal Link") {
					break
				}
				contentBuilder.WriteString(contentLine + "\n")
			}
			note.Content = contentBuilder.String()
		}

		// Парсинг внешних ссылок
		if strings.HasPrefix(line, "### External Link") {
			for scanner.Scan() {
				linkLine := scanner.Text()
				if strings.HasPrefix(linkLine, "### Internal Link") {
					break
				}
				externalLink = append(externalLink, linkLine)
			}
		}

		// Парсинг внутренних ссылок
		if strings.HasPrefix(line, "### Internal Link") {
			for scanner.Scan() {
				linkLine := scanner.Text()
				internalLink = append(internalLink, linkLine)
			}
		}
	}

	if len(externalLink) > 0 {
		note.ExternalLink = strings.Join(externalLink, "\n")
	}

	if len(internalLink) > 0 {
		note.InternalLink = strings.Join(internalLink, "\n")
	}

	return note
}

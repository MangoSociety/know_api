package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func ExtractTextBetweenBrackets(input string) (string, error) {
	re := regexp.MustCompile(`\[\[([^]]+)\]\]`)

	match := re.FindStringSubmatch(input)

	if len(match) < 2 {
		return "", errors.New("String is the smallest")
	}

	return match[1], nil
}

func ExtractSubstring(mdContent, start, end string) (string, error) {
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

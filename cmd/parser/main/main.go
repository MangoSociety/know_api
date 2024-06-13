package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Question структура для хранения заголовка и ответа
type Question struct {
	Title  string
	Answer string
}

func main() {
	// Отправка запроса на сайт
	res, err := http.Get("https://easyoffer.ru/question/2621")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Проверка на успешность запроса
	if res.StatusCode != 200 {
		log.Fatalf("Failed to fetch the page: %d %s", res.StatusCode, res.Status)
	}

	// Загрузка HTML-документа
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Создание структуры для хранения данных
	var question Question

	// Парсинг заголовка вопроса
	question.Title = doc.Find(".mb-5").Text()

	// Парсинг ответа
	question.Answer = doc.Find(".card-body").Text()

	fmt.Printf("Title: %s\n", question.Title)
	fmt.Printf("Answer: %s\n", question.Answer)
}

package usecase

import (
	"fmt"
	"math/rand"
	"read-adviser-bot/internal/question/repository"
	"read-adviser-bot/storage"
	"strings"
	"time"
)

type QuestionsUseCase struct {
	repo *repository.QuestionsRepo
}

func NewQuestionsUseCase(repo *repository.QuestionsRepo) *QuestionsUseCase {
	return &QuestionsUseCase{
		repo: repo,
	}
}

func (q *QuestionsUseCase) GetRandomQuestionGolang() string {
	random_quest, err := q.repo.GetInterviewQuestions("golang/interview_qustions/")
	if err != nil {
		fmt.Printf("get questions list with error %s", err.Error())
	}

	var result = "Твои 5 вопросов и ответов \n ----------------- \n"
	var quests = getRandomString(random_quest, 5)
	fmt.Println(quests)
	for index := 0; index < 5; index++ {
		var answer = q.repo.GetRandomQuestion(quests[index])
		result = result + answer + "\n----NEXT----\n"
	}

	//result := q.repo.GetRandomQuestion(getRandomString(random_quest))
	return result
}

func (q *QuestionsUseCase) GetRandomQuestionAndroid() string {
	random_quest, err := q.repo.GetInterviewQuestions("android/interview_questions/")
	if err != nil {
		fmt.Printf("get questions list with error %s", err.Error())
	}
	fmt.Println(random_quest)

	var result = "Твои 5 вопросов и ответов \n ----------------- \n"
	var quests = getRandomString(random_quest, 5)
	for index := 0; index < 5; index++ {
		var answer = q.repo.GetRandomQuestion(quests[index])
		result = result + answer + "\n----NEXT----\n"
	}

	//result := q.repo.GetRandomQuestion(getRandomString(random_quest))
	return result
}

func (q *QuestionsUseCase) GetRandomQuestionAndroidStruct() storage.Note {
	random_quest := q.repo.GetRandomQuestionStruct("android/interview_questions/")
	return random_quest
}

func getRandomString(numbers []string, count int) []string {
	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Создаем копию исходного списка
	shuffledNumbers := make([]string, len(numbers))
	copy(shuffledNumbers, numbers)

	// Используем алгоритм случайной перестановки (Fisher-Yates shuffle)
	for i := len(shuffledNumbers) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffledNumbers[i], shuffledNumbers[j] = shuffledNumbers[j], shuffledNumbers[i]
	}

	result := shuffledNumbers[:5]
	for i := 0; i < 5; i++ {
		result[i] = strings.ReplaceAll(shuffledNumbers[i], "\n", "")
	}

	// Возвращаем первые count элементов
	return result
}

package main

import (
	"context"
	tg_bot_api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gh_instance "know_api/clients/github"
	tg_client "know_api/clients/telegram"
	"know_api/config"
	"know_api/internal/instance/telegram/consumer/event_consumer"
	"know_api/internal/instance/telegram/events/telegram"
	"know_api/internal/questions/repos"
	"know_api/internal/questions/usecase"
	"know_api/pkg/db/mongo"
	"log"
	"time"
)

/*
Структура работы бота
- авто синк с гитхабом
*
*/
func main() {
	cfg := config.MustLoad()

	bot, err := tg_bot_api.NewBotAPI(cfg.TgToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	storage := mongo.New(cfg.MongoConnect, 10*time.Second)
	gh_client := gh_instance.NewGithubClient(cfg.Github.Token, cfg.Github.Owner, cfg.Github.Repo, cfg.Github.Sha)

	ghRepo := repos.NewQuestionsGHRepository(*gh_client)
	mgRepo := repos.NewQuestionsMGRepository(&storage)
	questionsUseCase := usecase.NewQuestionsUseCase(ghRepo, mgRepo)

	eventsProcessor := telegram.NewProcessor(
		tg_client.NewTelegramClient("api.telegram.org", cfg.TgToken, bot),
		questionsUseCase,
	)

	log.Print("telegram processor started")

	consumer := event_consumer.NewConsumer(eventsProcessor, eventsProcessor, 100) //  NewConsumer(eventsProcessor, eventsProcessor, batchSize)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := consumer.Start(ctx); err != nil {
		log.Fatal("service is stopped", err)
	}
}

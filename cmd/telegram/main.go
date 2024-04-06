package main

import (
	tg_bot_api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tg_client "know_api/clients/telegram"
	"know_api/config"
	"know_api/internal/instance/telegram/consumer/event_consumer"
	"know_api/internal/instance/telegram/events/telegram"
	"log"
)

func main() {
	cfg := config.MustLoad()

	bot, err := tg_bot_api.NewBotAPI(cfg.TgToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	//storage := mongo.New(cfg.MongoConnect, 10*time.Second)
	//gh_client := gh_client.NewGithubClient(cfg.Github.Token, cfg.Github.Owner, cfg.Github.Repo, cfg.Github.Sha)

	eventsProcessor := telegram.NewProcessor(
		tg_client.NewTelegramClient("api.telegram.org", cfg.TgToken, bot),
	)

	log.Print("telegram processor started")

	consumer := event_consumer.NewConsumer(eventsProcessor, eventsProcessor, 100) //  NewConsumer(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

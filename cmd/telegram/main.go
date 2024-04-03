package main

import (
	tg_bot_api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gh_client "know_api/clients/github"
	"know_api/config"
	"know_api/storage/mongo"
	"log"
	"time"
)

func main() {
	cfg := config.MustLoad()

	bot, err := tg_bot_api.NewBotAPI(cfg.TgBotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	storage := mongo.New(cfg.MongoConnectionString, 10*time.Second)
	gh_client := gh_client.NewGithubClient(token, owner, repo, sha)
	var repository = quest_repo.NewQuestionsRepo(gh_client, storage)
	var useCase = quest_usecase.NewQuestionsUseCase(repository)
}

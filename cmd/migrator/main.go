package main

import (
	categoriesRepo "github.com/MangoSociety/know_api/internal/categories/repository"
	"github.com/MangoSociety/know_api/internal/migrator/repository"
	"github.com/MangoSociety/know_api/internal/migrator/service"
	notesRepo "github.com/MangoSociety/know_api/internal/notes/repository"
	spheresRepo "github.com/MangoSociety/know_api/internal/spheres/repository"
	"github.com/MangoSociety/know_api/pkg/mongodb"
	"github.com/robfig/cron/v3"
	"log"
)

func main() {
	dbClient, err := mongodb.NewClient("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}

	gitHubRepo := repository.NewGitHubRepository("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", "MangoSociety", "repo_base", "refs/heads/main", "android_base")
	noteRepo := notesRepo.NewNoteRepository(dbClient)
	categoryRepo := categoriesRepo.NewCategoryRepository(dbClient)
	sphereRepo := spheresRepo.NewSphereRepository(dbClient)

	migratorService := service.NewMigratorService(gitHubRepo, noteRepo, categoryRepo, sphereRepo)

	c := cron.New()
	c.AddFunc("@every 6h", func() {
		if err := migratorService.Migrate(); err != nil {
			log.Println("Error during migration:", err)
		}
	})
	c.Start()

	// Запуск миграции при старте приложения
	if err := migratorService.Migrate(); err != nil {
		log.Println("Error during initial migration:", err)
	}

	// Ожидание для завершения cron-задач
	select {}
}

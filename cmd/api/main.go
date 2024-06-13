package main

import (
	"log"

	notesHttp "github.com/MangoSociety/know_api/internal/notes/delivery/http"
	notesRepo "github.com/MangoSociety/know_api/internal/notes/repository"
	notesService "github.com/MangoSociety/know_api/internal/notes/service"

	categoriesHttp "github.com/MangoSociety/know_api/internal/categories/delivery/http"
	categoriesRepo "github.com/MangoSociety/know_api/internal/categories/repository"
	categoriesService "github.com/MangoSociety/know_api/internal/categories/service"

	spheresHttp "github.com/MangoSociety/know_api/internal/spheres/delivery/http"
	spheresRepo "github.com/MangoSociety/know_api/internal/spheres/repository"
	spheresService "github.com/MangoSociety/know_api/internal/spheres/service"

	"github.com/MangoSociety/know_api/pkg/mongodb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация MongoDB
	dbClient, err := mongodb.NewClient("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}

	// Инициализация репозиториев
	noteRepo := notesRepo.NewNoteRepository(dbClient)
	categoryRepo := categoriesRepo.NewCategoryRepository(dbClient)
	sphereRepo := spheresRepo.NewSphereRepository(dbClient)

	// Инициализация сервисов
	noteService := notesService.NewNoteService(noteRepo)
	categoryService := categoriesService.NewCategoryService(categoryRepo)
	sphereService := spheresService.NewSphereService(sphereRepo)

	// Инициализация HTTP обработчиков
	noteHandler := notesHttp.NewNoteHandler(noteService)
	categoryHandler := categoriesHttp.NewCategoryHandler(categoryService)
	sphereHandler := spheresHttp.NewSphereHandler(sphereService)

	// Инициализация маршрутов
	router := gin.Default()

	// Настройка CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Укажите ваш фронтенд URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/notes", noteHandler.GetNotes)
	router.POST("/notes", noteHandler.CreateNote)
	router.PUT("/notes/:id", noteHandler.UpdateNote)
	router.DELETE("/notes/:id", noteHandler.DeleteNote)

	router.GET("/categories", categoryHandler.GetCategories)
	router.POST("/categories", categoryHandler.CreateCategory)
	router.PUT("/categories/:id", categoryHandler.UpdateCategory)
	router.DELETE("/categories/:id", categoryHandler.DeleteCategory)

	router.GET("/spheres", sphereHandler.GetSpheres)
	router.POST("/spheres", sphereHandler.CreateSphere)
	router.PUT("/spheres/:id", sphereHandler.UpdateSphere)
	router.DELETE("/spheres/:id", sphereHandler.DeleteSphere)

	// Запуск сервера
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

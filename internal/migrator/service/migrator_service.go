package service

import (
	categoryDomain "github.com/MangoSociety/know_api/internal/categories/domain"
	categoryRepository "github.com/MangoSociety/know_api/internal/categories/repository"
	"github.com/MangoSociety/know_api/internal/migrator/domain"
	"github.com/MangoSociety/know_api/internal/migrator/repository"
	notesDomain "github.com/MangoSociety/know_api/internal/notes/domain"
	"github.com/MangoSociety/know_api/internal/notes/models"
	notesRepository "github.com/MangoSociety/know_api/internal/notes/repository"
	sphereDomain "github.com/MangoSociety/know_api/internal/spheres/domain"
	sphereRepository "github.com/MangoSociety/know_api/internal/spheres/repository"
	"github.com/MangoSociety/know_api/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type MigratorService interface {
	Migrate() error
}

type migratorService struct {
	gitHubRepo   repository.GitHubRepository
	noteRepo     notesRepository.NoteRepository
	categoryRepo categoryRepository.CategoryRepository
	sphereRepo   sphereRepository.SphereRepository
}

func NewMigratorService(gitHubRepo repository.GitHubRepository, noteRepo notesRepository.NoteRepository, categoryRepo categoryRepository.CategoryRepository, sphereRepo sphereRepository.SphereRepository) MigratorService {
	return &migratorService{
		gitHubRepo:   gitHubRepo,
		noteRepo:     noteRepo,
		categoryRepo: categoryRepo,
		sphereRepo:   sphereRepo,
	}
}

func (s *migratorService) Migrate() error {
	files, err := s.gitHubRepo.FetchFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		// Используем парсер для создания временной заметки из файла
		tempNote := utils.ParseMarkdownFile(file.Content)

		// Преобразуем временную заметку в доменную
		note, err := s.convertToDomainModel(tempNote)
		if err != nil {
			log.Println("Error converting temp note to domain model:", tempNote.Title, err)
			continue
		}

		existingNote, err := s.noteRepo.GetByTitle(note.Title)
		if err != nil {
			log.Println("Error fetching note by title:", note.Title, err)
			continue
		}
		if existingNote != nil {
			note.ID = existingNote.ID
			if err := s.noteRepo.Update(note); err != nil {
				log.Println("Error updating note:", note.Title, err)
			}
		} else {
			note.ID = primitive.NewObjectID()
			if err := s.noteRepo.Create(note); err != nil {
				log.Println("Error creating note:", note.Title, err)
			}
		}
	}

	return nil
}

func (s *migratorService) convertToDomainModel(tempNote models.Note) (notesDomain.Note, error) {
	var note notesDomain.Note

	note.Title = tempNote.Title
	note.Content = tempNote.Content
	note.ExternalLink = tempNote.ExternalLink
	note.InternalLink = tempNote.InternalLink

	// Получаем идентификатор сферы или создаем новую
	sphere, err := s.sphereRepo.GetByName(tempNote.Sphere)
	if err != nil {
		return note, err
	}
	if sphere == nil {
		newSphere := sphereDomain.Sphere{
			ID:   primitive.NewObjectID(),
			Name: tempNote.Sphere,
		}
		if err := s.sphereRepo.Create(newSphere); err != nil {
			return note, err
		}
		note.SphereID = newSphere.ID
	} else {
		note.SphereID = sphere.ID
	}

	// Получаем идентификаторы категорий или создаем новые
	prevTheme := ""
	for _, theme := range tempNote.Theme {
		category, err := s.categoryRepo.GetByName(theme)
		if err != nil {
			return note, err
		}
		if category == nil {
			newCategory := categoryDomain.Category{
				ID:       primitive.NewObjectID(),
				Name:     theme,
				SphereID: note.SphereID,
			}
			if err := s.categoryRepo.Create(&newCategory, prevTheme); err != nil {
				return note, err
			}
			note.CategoryIDs = append(note.CategoryIDs, newCategory.ID)
		} else {
			note.CategoryIDs = append(note.CategoryIDs, category.ID)
		}
		prevTheme = theme
	}

	return note, nil
}

func parseFileToNote(file domain.GitHubFile) notesDomain.Note {
	sample := utils.ParseMarkdownFile(file.Content)
	// Пример парсинга файла и создания заметки
	return notesDomain.Note{
		Title:       sample.Title,
		SphereID:    primitive.NewObjectID(),                       // Нужно установить правильный SphereID
		CategoryIDs: []primitive.ObjectID{primitive.NewObjectID()}, // Нужно установить правильные CategoryIDs
		Content:     file.Content,
	}
}

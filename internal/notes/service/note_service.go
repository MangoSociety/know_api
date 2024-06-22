package service

import (
	"github.com/MangoSociety/know_api/internal/notes/domain"
	"github.com/MangoSociety/know_api/internal/notes/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoteService interface {
	GetNotes() ([]domain.Note, error)
	CreateNote(note domain.Note) error
	UpdateNote(note domain.Note) error
	DeleteNote(id primitive.ObjectID) error
	GetRandomNoteByCategory(categoryIDs []primitive.ObjectID) (*domain.Note, error)
	GetByHex(hexID string) (*domain.Note, error)
}

type noteService struct {
	repo repository.NoteRepository
}

func NewNoteService(repo repository.NoteRepository) NoteService {
	return &noteService{repo: repo}
}

func (s *noteService) GetNotes() ([]domain.Note, error) {
	return s.repo.GetAll()
}

func (s *noteService) CreateNote(note domain.Note) error {
	return s.repo.Create(note)
}

func (s *noteService) UpdateNote(note domain.Note) error {
	return s.repo.Update(note)
}

func (s *noteService) DeleteNote(id primitive.ObjectID) error {
	return s.repo.Delete(id)
}

func (s *noteService) GetRandomNoteByCategory(categoryIDs []primitive.ObjectID) (*domain.Note, error) {
	return s.repo.GetRandomNoteByCategory(categoryIDs)
}

func (s *noteService) GetByHex(hexID string) (*domain.Note, error) {
	return s.repo.GetByHex(hexID)
}

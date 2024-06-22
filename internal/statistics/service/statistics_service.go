package service

import (
	"github.com/MangoSociety/know_api/internal/statistics/domain"
	"github.com/MangoSociety/know_api/internal/statistics/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatisticsService interface {
	GetAll() ([]domain.Statistics, error)
	GetNotesIDsAll() ([]primitive.ObjectID, error)
	Create(statistics domain.Statistics) error
}

type statisticsService struct {
	repo repository.StatisticsRepository
}

func NewStatisticsService(repo repository.StatisticsRepository) StatisticsService {
	return &statisticsService{repo: repo}
}

func (s *statisticsService) GetAll() ([]domain.Statistics, error) {
	return s.repo.GetAll()
}

func (s *statisticsService) GetNotesIDsAll() ([]primitive.ObjectID, error) {
	allItems, _ := s.repo.GetAll()
	var notesIDs []primitive.ObjectID
	for _, selectedParam := range allItems {
		notesIDs = append(notesIDs, selectedParam.NoteID)
	}
	return notesIDs, nil
}

func (s *statisticsService) Create(statistics domain.Statistics) error {
	return s.repo.Create(statistics)
}

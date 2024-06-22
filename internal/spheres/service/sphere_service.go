package service

import (
	"github.com/MangoSociety/know_api/internal/spheres/domain"
	"github.com/MangoSociety/know_api/internal/spheres/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SphereService interface {
	GetSpheres() ([]domain.Sphere, error)
	CreateSphere(sphere domain.Sphere) error
	UpdateSphere(sphere domain.Sphere) error
	DeleteSphere(id primitive.ObjectID) error
	GetById(id string) (*domain.Sphere, error)
}

type sphereService struct {
	repo repository.SphereRepository
}

func NewSphereService(repo repository.SphereRepository) SphereService {
	return &sphereService{repo: repo}
}

func (s *sphereService) GetSpheres() ([]domain.Sphere, error) {
	return s.repo.GetAll()
}

func (s *sphereService) CreateSphere(sphere domain.Sphere) error {
	return s.repo.Create(sphere)
}

func (s *sphereService) UpdateSphere(sphere domain.Sphere) error {
	return s.repo.Update(sphere)
}

func (s *sphereService) DeleteSphere(id primitive.ObjectID) error {
	return s.repo.Delete(id)
}

func (s *sphereService) GetById(id string) (*domain.Sphere, error) {
	return s.repo.GetById(id)
}

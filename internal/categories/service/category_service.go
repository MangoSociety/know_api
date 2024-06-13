package service

import (
	"github.com/MangoSociety/know_api/internal/categories/domain"
	"github.com/MangoSociety/know_api/internal/categories/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryService interface {
	GetCategories() ([]domain.Category, error)
	CreateCategory(category domain.Category) error
	UpdateCategory(category domain.Category) error
	DeleteCategory(id primitive.ObjectID) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetCategories() ([]domain.Category, error) {
	return s.repo.GetAll()
}

func (s *categoryService) CreateCategory(category domain.Category) error {
	return s.repo.Create(category)
}

func (s *categoryService) UpdateCategory(category domain.Category) error {
	return s.repo.Update(category)
}

func (s *categoryService) DeleteCategory(id primitive.ObjectID) error {
	return s.repo.Delete(id)
}

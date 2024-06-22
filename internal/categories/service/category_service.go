package service

import (
	"github.com/MangoSociety/know_api/internal/categories/domain"
	"github.com/MangoSociety/know_api/internal/categories/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryService interface {
	GetCategories() ([]domain.Category, error)
	GetCategoriesTree(parentName string) ([]*domain.Category, error)
	CreateCategory(category *domain.Category, prevCategory string) error
	UpdateCategory(category domain.Category) error
	DeleteCategory(id primitive.ObjectID) error
	GetCategoriesByParentID(parentID string) ([]*domain.Category, error)
	GetCategoriesByName(name string) (*domain.Category, error)
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

func (s *categoryService) GetCategoriesTree(parentName string) ([]*domain.Category, error) {
	return s.repo.GetCategoriesByParentName(parentName)
}

func (s *categoryService) CreateCategory(category *domain.Category, prevCategory string) error {
	return s.repo.Create(category, prevCategory)
}

func (s *categoryService) UpdateCategory(category domain.Category) error {
	return s.repo.Update(category)
}

func (s *categoryService) DeleteCategory(id primitive.ObjectID) error {
	return s.repo.Delete(id)
}

func (s *categoryService) GetCategoriesByParentID(parentID string) ([]*domain.Category, error) {
	return s.repo.GetCategoriesByParentID(parentID)
}

func (s *categoryService) GetCategoriesByName(name string) (*domain.Category, error) {
	return s.repo.GetByName(name)
}

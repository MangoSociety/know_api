package service

import (
	"github.com/MangoSociety/know_api/internal/user/domain"
	"github.com/MangoSociety/know_api/internal/user/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSelectionService struct {
	repo *repository.UserSelectionRepository
}

func NewUserSelectionService(repo *repository.UserSelectionRepository) UserSelectionService {
	return UserSelectionService{repo: repo}
}

func (s *UserSelectionService) GetCategoriesByChatID(chatID int) []primitive.ObjectID {
	userSelection, _ := s.repo.GetByChatID(chatID)
	var categoryIDs []primitive.ObjectID
	for _, selectedParam := range userSelection.SelectedCategories {
		for _, category := range selectedParam.Categories {
			categoryIDs = append(categoryIDs, category.CategoryID)
		}
	}
	return categoryIDs
}

func (s *UserSelectionService) GetUserSelection(chatID int) (*domain.UserSelection, error) {
	return s.repo.GetByChatID(chatID)
}

func (s *UserSelectionService) AddCategoryToUser(
	chatID int,
	sphereID primitive.ObjectID,
	sphereTitle string,
	categoryID primitive.ObjectID,
	categoryTitle string,
) error {
	selection, err := s.repo.GetByChatID(chatID)
	if err != nil {
		return err
	}

	if selection == nil {
		selection = &domain.UserSelection{
			ChatID: chatID,
			SelectedCategories: []domain.SelectedParameter{
				{
					SphereID:    sphereID,
					SphereTitle: sphereTitle,
					Categories: []domain.SelectedCategory{
						{
							CategoryID:    categoryID,
							CategoryTitle: categoryTitle,
						},
					},
				},
			},
		}
	} else {
		found := false
		for i, sel := range selection.SelectedCategories {
			if sel.SphereID == sphereID {
				selection.SelectedCategories[i].Categories = append(selection.SelectedCategories[i].Categories, domain.SelectedCategory{
					CategoryID:    categoryID,
					CategoryTitle: categoryTitle,
				})
				found = true
				break
			}
		}
		if !found {
			selection.SelectedCategories = append(selection.SelectedCategories, domain.SelectedParameter{
				SphereID:    sphereID,
				SphereTitle: sphereTitle,
				Categories: []domain.SelectedCategory{
					{
						CategoryID:    categoryID,
						CategoryTitle: categoryTitle,
					},
				},
			})
		}
	}

	return s.repo.CreateOrUpdate(selection)
}

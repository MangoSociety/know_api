package memory

import (
	"know_api/internal_1/repository/category"
	questions2 "know_api/internal_1/repository/questions"
	"know_api/internal_1/repository/theme"
)

type UseCaseMemory struct {
	themeRepository theme.MgRepository
	category        category.MGRepository
	questGhRepo     questions2.GHRepository
	questMgRepo     questions2.MGRepository
}

func NewUseCaseMemory(themeRepository theme.MgRepository, category category.MGRepository) *UseCaseMemory {
	return &UseCaseMemory{
		themeRepository: themeRepository,
		category:        category,
	}
}

func (u *UseCaseMemory) MigrateData() {

}

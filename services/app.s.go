package services

import (
	"waheim.api/models"
	"waheim.api/repositories"
)

type AppService struct{}

func NewAppService() *AppService {
	return &AppService{}
}

func (s *AppService) CreateApp(app *models.App) error {
	return repositories.CreateApp(app)
}

func (s *AppService) GetAppById(id string) (models.App, error) {
	return repositories.GetAppById(id)
}

func (s *AppService) GetAllApps(limit, offset int) ([]models.App, error) {
	return repositories.GetAllApps(limit, offset)
}

func (s *AppService) UpdateApp(id string, updates map[string]interface{}) error {
	return repositories.UpdateApp(id, updates)
}

func (s *AppService) DeleteApp(id string) error {
	return repositories.DeleteApp(id)
}

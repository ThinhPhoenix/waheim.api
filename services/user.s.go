package services

import (
	"waheim.api/models"
	"waheim.api/repositories"
)

type UserService interface {
	SignUp(request map[string]string) error
	SignIn(request map[string]string) (string, error)
	AuthMe(token string) (models.User, error)
	GetAllUsers(filters map[string]string, limit, offset int) ([]models.User, error)
	GetUserById(id string) (models.User, error)
	UpdateUser(id string, updates map[string]interface{}) error
	DeleteUser(id string) error
}
func (u *userServiceImpl) GetUserById(id string) (models.User, error) {
	return repositories.GetUserById(id)
}

func (u *userServiceImpl) UpdateUser(id string, updates map[string]interface{}) error {
	return repositories.UpdateUser(id, updates)
}

func (u *userServiceImpl) DeleteUser(id string) error {
	return repositories.DeleteUser(id)
}

type userServiceImpl struct{}

func (u *userServiceImpl) SignUp(request map[string]string) error {
	return repositories.SignUp(request)
}

func (u *userServiceImpl) SignIn(request map[string]string) (string, error) {
	return repositories.SignIn(request)
}

func (u *userServiceImpl) AuthMe(token string) (models.User, error) {
	return repositories.AuthMe(token)
}

func (u *userServiceImpl) GetAllUsers(filters map[string]string, limit, offset int) ([]models.User, error) {
	return repositories.GetAllUsers(filters, limit, offset)
}

func NewUserService() UserService {
	return &userServiceImpl{}
}

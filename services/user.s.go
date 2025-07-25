package services

import (
	"waheim.api/models"
	"waheim.api/repositories"
)

type UserService interface {
	SignUp(request map[string]string) error
	SignIn(request map[string]string) (string, error)
	AuthMe(token string) (models.User, error)
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

func NewUserService() UserService {
	return &userServiceImpl{}
}

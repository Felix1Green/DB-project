package users

import "github.com/Felix1Green/DB-project/internal/pkg/models"

type UseCase interface{
	CreateUser(nickname string, user *models.UserRequestBody) (*[]models.User, error)
	GetProfile(nickname string) (*models.User, error)
	UpdateProfile(nickname string, user *models.UserRequestBody) (*models.User, error)
}

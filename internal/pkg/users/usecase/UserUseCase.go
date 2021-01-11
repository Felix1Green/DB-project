package usecase
import (
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"github.com/Felix1Green/DB-project/internal/pkg/users"
)

type UserUseCase struct{
	Repository users.Repository
}

func NewUserUseCase(repository users.Repository) *UserUseCase{
	return &UserUseCase{
		Repository: repository,
	}
}

func (t *UserUseCase) CreateUser(nickname string, user *models.UserRequestBody) (*[]models.User, error){
	if nickname == "" || user.Email == ""{
		return nil, models.IncorrectInputParams
	}
	return t.Repository.CreateUser(nickname, user)
}


func (t *UserUseCase) GetProfile(nickname string) (*models.User, error){
	if nickname == ""{
		return nil, models.IncorrectInputParams
	}

	return t.Repository.GetProfile(nickname)
}

func (t *UserUseCase) UpdateProfile(nickname string, user *models.UserRequestBody) (*models.User, error){
	if nickname == "" || user.Email == "" && user.About == "" && user.FullName == ""{
		return t.Repository.GetProfile(nickname)
	}
	_, err := t.Repository.GetProfile(nickname)
	if err != nil{
		return nil, models.NoSuchUser
	}
	return t.Repository.UpdateProfile(nickname, user)
}
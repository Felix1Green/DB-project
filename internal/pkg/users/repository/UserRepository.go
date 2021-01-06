package repository

import (
	"database/sql"
	"github.com/Felix1Green/DB-project/internal/pkg/models"
)

type UserRepository struct {
	dbConnection *sql.DB
}

func NewUsersRepository(conn *sql.DB) *UserRepository{
	return &UserRepository{
		dbConnection: conn,
	}
}

func (t *UserRepository) CreateUser(nickname string, user *models.UserRequestBody) (*[]models.User, error) {
	CreationQuery := "INSERT INTO users (nickname, fullname, about, email) VALUES($1,$2,$3,$4)"
	UniqueErrorQuery := "SELECT nickname, fullname, about, email FROM users WHERE nickname = $1 or email = $2"
	_, dbErr := t.dbConnection.Exec(CreationQuery, nickname, user.FullName, user.About, user.Email)
	resultArr := make([]models.User, 0)
	if dbErr != nil {
		result, dbErr := t.dbConnection.Query(UniqueErrorQuery, nickname, user.Email)
		if dbErr != nil || result.Err() != nil {
			return nil, models.IncorrectInputParams
		}
		defer func(){_ = result.Close()}()
		for result.Next() {
			instance := new(models.User)
			_ = result.Scan(&instance.Nickname, &instance.FullName, &instance.About, &instance.Email)
			resultArr = append(resultArr, *instance)
		}
		return &resultArr, models.UserAlreadyExists
	}
	resultArr = append(resultArr, models.User{
		Nickname: nickname,
		FullName: user.FullName,
		About:    user.About,
		Email:    user.Email,
	})
	return &resultArr, nil
}

func (t *UserRepository) GetProfile(nickname string) (*models.User, error) {
	query := "SELECT fullname, about, email FROM users where nickname = $1"
	result := t.dbConnection.QueryRow(query, nickname)
	if result.Err() != nil {
		return nil, models.NoSuchUser
	}
	UserInstance := new(models.User)
	scanErr := result.Scan(&UserInstance.FullName, &UserInstance.About, &UserInstance.Email)
	if scanErr != nil {
		return nil, models.NoSuchUser
	}
	UserInstance.Nickname = nickname
	return UserInstance, nil
}

func (t *UserRepository) UpdateProfile(nickname string, user *models.UserRequestBody) (*models.User, error) {
	query := "UPDATE users SET fullname=$1, about=$2, email=$3 where nickname=$4 returning nickname"
	result := t.dbConnection.QueryRow(query, nickname, user.FullName, user.About, user.Email)
	if result.Err() != nil {
		return nil, models.UserConflict
	}
	nickname = ""
	scanErr := result.Scan(&nickname)
	if scanErr != nil || nickname == "" {
		return nil, models.NoSuchUser
	}

	return &models.User{
		Nickname: nickname,
		FullName: user.FullName,
		About:    user.About,
		Email:    user.Email,
	}, nil
}

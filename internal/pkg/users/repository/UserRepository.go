package repository

import (
	"fmt"
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"github.com/jackc/pgx"
)

type UserRepository struct {
	dbConnection *pgx.ConnPool
}

func NewUsersRepository(conn *pgx.ConnPool) *UserRepository{
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
		defer func(){result.Close()}()
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
	query := "SELECT nickname, fullname, about, email FROM users where nickname = $1"
	result := t.dbConnection.QueryRow(query, nickname)
	if result == nil {
		return nil, models.NoSuchUser
	}
	UserInstance := new(models.User)
	scanErr := result.Scan(&UserInstance.Nickname,&UserInstance.FullName, &UserInstance.About, &UserInstance.Email)
	if scanErr != nil {
		return nil, models.NoSuchUser
	}
	return UserInstance, nil
}

func (t *UserRepository) UpdateProfile(nickname string, user *models.UserRequestBody) (*models.User, error) {
	queryArgs :=make([]interface{},0)
	query := "UPDATE users SET "
	counter := 1
	if user.FullName != ""{
		query += fmt.Sprintf("fullname=$%d ", counter)
		queryArgs = append(queryArgs, user.FullName)
		counter++
	}
	if user.About != ""{
		if counter > 1{
			query += fmt.Sprintf(", about=$%d ", counter)
		}else{
			query += fmt.Sprintf("about=$%d ", counter)
		}
		queryArgs = append(queryArgs, user.About)
		counter++
	}
	if user.Email != ""{
		if counter > 1{
			query += fmt.Sprintf(", email=$%d ", counter)
		}else{
			query += fmt.Sprintf("email=$%d ", counter)
		}
		queryArgs = append(queryArgs, user.Email)
		counter++
	}
	query += fmt.Sprintf("where nickname=$%d returning nickname, fullname, about, email", counter)
	queryArgs = append(queryArgs, nickname)
	item := new(models.User)
	scanErr := t.dbConnection.QueryRow(query, queryArgs...).Scan(&item.Nickname, &item.FullName, &item.About, &item.Email)
	if scanErr != nil || nickname == "" {
		return nil, models.UserConflict
	}

	return item, nil
}

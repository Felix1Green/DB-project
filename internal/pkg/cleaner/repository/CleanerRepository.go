package repository

import (
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"github.com/jackc/pgx"
)

type CleanerRepository struct {
	DBConnection *pgx.ConnPool
}

func NewCleanerRepository(conn *pgx.ConnPool) *CleanerRepository{
	return &CleanerRepository{
		DBConnection: conn,
	}
}


func (t *CleanerRepository) Status()(*models.Status, error){
	result := new(models.Status)
	ScanErr := t.DBConnection.QueryRow("SELECT (SELECT COUNT(id) FROM thread), (SELECT COUNT(id) FROM post), (SELECT COUNT(id) FROM forum), (SELECT COUNT(id) FROM users)").Scan(&result.Thread,
		&result.Post, &result.Forum, &result.User)
	if ScanErr != nil{
		return nil, models.InternalDBError
	}
	return result, nil
}

func (t *CleanerRepository) Clear() error{
	_, err := t.DBConnection.Exec("TRUNCATE users, forum, thread, post, vote CASCADE ")
	if err != nil{
		return models.InternalDBError
	}
	return nil
}
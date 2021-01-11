package repository

import (
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"github.com/jackc/pgx"
)

type PostRepository struct {
	dbConnection *pgx.ConnPool
}

func NewPostRepository(connection *pgx.ConnPool) *PostRepository{
	return &PostRepository{dbConnection: connection}
}

func (t *PostRepository) GetPostDetails(id uint64) (*models.PostModel, error) {
	query := "SELECT id, parent, author, message, isEdited, forum, thread, created FROM post WHERE id = $1"
	result := new(models.PostModel)
	err := t.dbConnection.QueryRow(query, id).Scan(&result.ID, &result.Parent, &result.Author,&result.Message, &result.IsEdited,
		&result.Forum, &result.Thread, &result.Created)
	if err != nil{
		return nil, models.ThreadAbsentsError
	}
	return result, nil
}

func (t *PostRepository) UpdatePost(id uint64, input *models.PostUpdateRequestInput) (*models.PostModel, error){
	query := "UPDATE post SET message = $1, isedited = true WHERE id = $2 RETURNING id, parent, author, message, isEdited, forum, thread, created"
	result := new(models.PostModel)
	err := t.dbConnection.QueryRow(query, input.Message, id).Scan(&result.ID, &result.Parent, &result.Author,&result.Message, &result.IsEdited,
		&result.Forum, &result.Thread, &result.Created)
	if err != nil{
		return nil, models.PostDoesntExists
	}
	return result, nil
}

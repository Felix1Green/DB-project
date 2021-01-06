package repository

import (
	"database/sql"
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"log"
)

type ForumRepository struct {
	dbConnection *sql.DB
}

func NewForumRepository(connection *sql.DB) *ForumRepository{
	return &ForumRepository{
		dbConnection: connection,
	}
}


func (t *ForumRepository) CreateForum(input *models.ForumRequestInput) (*models.Forum, error) {
	query := "INSERT INTO forum (title, user_id, slug) VALUES ($1, $2, $3)"
	uniqueErrQuery := "SELECT v1.title, v1.user, v1.slug, COUNT(v2.id), COUNT(v3.id) FROM ( " +
		"SELECT v1.title, v1.user, v1.slug, COUNT(v2.id) FROM forum " +
		"join thread v2 on (v2.forum=v1.slug) " +
		"where v1.slug = $1 group by v1.title, v1.user, v1.slug ) " +
		"join post v3 on (v3.forum = v1.slug) " +
		"group by v1.title, v1.user, v1.slug"
	_, DBErr := t.dbConnection.Exec(query, input.Title, input.User, input.Slug)
	if DBErr != nil {
		forumInstance := new(models.Forum)
		scanErr := t.dbConnection.QueryRow(uniqueErrQuery, input.Slug).Scan(&forumInstance.Title, &forumInstance.User,
			&forumInstance.Slug, &forumInstance.Threads, &forumInstance.Posts)
		if scanErr != nil {
			return nil, models.InternalDBError
		}
		return forumInstance, models.ForumAlreadyExists
	}
	return &models.Forum{
		Title: input.Title,
		User:  input.User,
		Slug:  input.Slug,
	}, nil
}

func (t *ForumRepository) GetForum(slug string) (*models.Forum, error) {
	query := "SELECT v1.title, v1.user, v1.slug, COUNT(v2.id), COUNT(v3.id) FROM (" +
		"SELECT v1.title, v1.user, v1.slug, COUNT(v2.id) FROM forum " +
		"join thread v2 on (v2.forum=v1.slug) " +
		"where v1.slug = $1 group by v1.title, v1.user, v1.slug ) " +
		"join post v3 on (v3.forum = v1.slug) " +
		"group by v1.title, v1.user, v1.slug "
	res := new(models.Forum)
	ScanErr := t.dbConnection.QueryRow(query, slug).Scan(&res.Title, &res.User, &res.Slug, &res.Threads, &res.Posts)
	if ScanErr != nil{
		log.Println(ScanErr)
		return nil, models.ForumDoesntExists
	}

	return res,nil
}


func (t *ForumRepository) CreateForumThread(slug string, thread *models.ThreadRequestInput) (*models.ThreadModel, error){
	query := "INSERT INTO Thread (title,author,message,forum) VALUES ($1,$2,$3,$4) RETURNING ID, CREATED"
	uniqueErrQuery := "SELECT v1.id, v1.title, v1.author, v1.forum, v1.message, COUNT(v2.id), v1.created " +
		"FROM Thread v1 JOIN post v2 on(v2.thread = v1.id) WHERE v1.title = $1"
	result := t.dbConnection.QueryRow(query, thread.Title, thread.Author, thread.Message, slug)
	var resultID uint64
	var created string
	if result.Err() != nil || result.Scan(&resultID, created) != nil{
		resultThread := new(models.ThreadModel)
		scanErr := t.dbConnection.QueryRow(uniqueErrQuery, thread.Title).Scan(&resultThread.ID, &resultThread.Title,
			&resultThread.Author, &resultThread.Forum,&resultThread.Message, &resultThread.Votes, &resultThread.Created)
		if scanErr != nil{
			return nil, models.ForumDoesntExists
		}
		return resultThread, models.ForumAlreadyExists
	}
	return &models.ThreadModel{
		ID: resultID,
		Title: thread.Title,
		Author: thread.Author,
		Message: thread.Message,
		Created: created,
		Votes: 0,
		Forum: slug,
	}, nil
}


func (t *ForumRepository) GetForumUsers(slug string, limit, since int, desc bool) (*[]models.User, error){
	query := "SELECT DISTINCT ON (v1.id) v1.nickname, v1.fullname, v1.about, v1.email FROM users v1 " +
		"JOIN post v2 on (v2.forum = $1) " +
		"JOIN thread v3 on (v3.forum = $1) " +
		"WHERE v1.id >= $2 " +
		"ORDER BY v1.nickname "
	if desc {
		query += "DESC "
	}
	query += "LIMIT $3"
	rows, DBErr := t.dbConnection.Query(query, slug, since, limit)
	if DBErr != nil{
		return nil, models.ForumDoesntExists
	}

	defer func(){_ = rows.Close()}()
	resultList := make([]models.User, 0)
	for rows.Next(){
		userModel := new(models.User)
		scanErr := rows.Scan(userModel.Nickname, userModel.FullName, userModel.About, userModel.Email)
		if scanErr != nil{
			return nil, models.ForumDoesntExists
		}
		resultList = append(resultList, *userModel)
	}
	return &resultList, nil
}

func (t *ForumRepository) GetForumThreads(slug string, limit, since int, desc bool) (*[]models.ThreadModel, error){
	query := "SELECT id, title, author, forum, message, COUNT(v2.id), created FROM thread " +
		"join post v2 on (v2.thread = id) " +
		"where forum = $1 and created >= $2 " +
		"ORDER BY created "
	if desc{
		query += "DESC "
	}
	query += "LIMIT $3"
	rows, DBErr := t.dbConnection.Query(query, slug, since, limit)
	if DBErr != nil || rows.Err() != nil{
		return nil, models.ForumDoesntExists
	}
	resultList := make([]models.ThreadModel, 0)
	for rows.Next(){
		model := new(models.ThreadModel)
		ScanErr := rows.Scan(&model.ID, &model.Title, &model.Author,&model.Forum,&model.Message,&model.Votes, &model.Created)
		if ScanErr != nil{
			return nil, models.ForumDoesntExists
		}
		resultList = append(resultList, *model)
	}
	return &resultList, nil
}
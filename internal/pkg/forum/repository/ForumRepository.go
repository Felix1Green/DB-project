package repository

import (
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"github.com/Felix1Green/DB-project/internal/pkg/utils"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx"
	"strings"
)

type ForumRepository struct {
	dbConnection *pgx.ConnPool
}

func NewForumRepository(connection *pgx.ConnPool) *ForumRepository{
	return &ForumRepository{
		dbConnection: connection,
	}
}


func (t *ForumRepository) CreateForum(input *models.ForumRequestInput) (*models.Forum, error) {
	query := "INSERT INTO forum (title, user_id, slug) VALUES ($1, $2, $3) returning title, user_id, slug"
	uniqueErrQuery := "SELECT v1.title, v1.user_id, v1.slug, v1.threads, v1.posts FROM forum v1 where v1.slug = $1"
	item := new(models.Forum)
	scanErr := t.dbConnection.QueryRow(query, input.Title, input.User, input.Slug).Scan(&item.Title, &item.User, &item.Slug)
	if scanErr != nil {
		forumInstance := new(models.Forum)
		scanErr := t.dbConnection.QueryRow(uniqueErrQuery, input.Slug).Scan(&forumInstance.Title, &forumInstance.User,
			&forumInstance.Slug, &forumInstance.Threads, &forumInstance.Posts)
		if scanErr != nil {
			return nil, models.InternalDBError
		}
		return forumInstance, models.ForumAlreadyExists
	}
	return item, nil
}

func (t *ForumRepository) GetForum(slug string) (*models.Forum, error) {
	query := "SELECT v1.title, v1.user_id, v1.slug, v1.threads, v1.posts FROM forum v1 where v1.slug = $1"
	res := new(models.Forum)
	ScanErr := t.dbConnection.QueryRow(query, slug).Scan(&res.Title, &res.User, &res.Slug, &res.Threads, &res.Posts)
	if ScanErr != nil{
		return nil, models.ForumDoesntExists
	}

	return res,nil
}

func (t *ForumRepository) GetForumSimple(slug string) (*models.Forum, error){
	query := "SELECT v1.slug, v1.title, v1.user_id FROM forum v1 where v1.slug = $1"
	res := new(models.Forum)
	ScanErr := t.dbConnection.QueryRow(query, slug).Scan(&res.Slug, &res.Title, &res.User)
	if ScanErr != nil{
		return nil, models.ForumDoesntExists
	}
	return res, nil
}


func (t *ForumRepository) CreateForumThread(slug string, thread *models.ThreadRequestInput) (*models.ThreadModel, error){
	query := "INSERT INTO Thread (title,author,message,forum, slug) VALUES ($1,$2,$3,$4,$5) RETURNING ID, CREATED"
	queryArgs := []interface{}{
		thread.Title, thread.Author, thread.Message, slug,thread.Slug,
	}
	if thread.Created != ""{
		query = "INSERT INTO Thread (title, author, message, forum,slug, created) VALUES ($1,$2,$3,$4,$5,$6) RETURNING ID, CREATED"
		queryArgs = append(queryArgs, thread.Created)
	}
	uniqueErrQuery := "SELECT v1.id, v1.title, v1.author, v1.forum, v1.message, v1.votes_counter, v1.created, v1.slug " +
		"FROM Thread v1 WHERE v1.slug = $1 GROUP BY v1.id"

	result := t.dbConnection.QueryRow(query, queryArgs...)
	var resultID uint64
	var created strfmt.DateTime
	if err := result.Scan(&resultID, &created); err != nil{
		resultThread := new(models.ThreadModel)
		scanErr := t.dbConnection.QueryRow(uniqueErrQuery, thread.Slug).Scan(&resultThread.ID, &resultThread.Title,
			&resultThread.Author, &resultThread.Forum,&resultThread.Message, &resultThread.Votes, &resultThread.Created, &resultThread.Slug)
		if scanErr != nil{
			return nil, models.ForumDoesntExists
		}
		if strings.HasPrefix(resultThread.Slug, utils.SlugCreatedSign){
			resultThread.Slug = ""
		}
		return resultThread, models.ThreadUniqueErr
	}
	return &models.ThreadModel{
		ID: resultID,
		Title: thread.Title,
		Author: thread.Author,
		Message: thread.Message,
		Slug: thread.Slug,
		Votes: 0,
		Created: created,
		Forum: slug,
	}, nil
}


func (t *ForumRepository) GetForumUsers(slug string, limit int, since string, desc bool) (*[]models.User, error){
	query := "SELECT v1.nickname, v1.fullname, v1.about, v1.email FROM forum_users v2 JOIN users v1 on(v1.nickname = v2.user_nickname) " +
		"where v2.forum = $1 and v1.nickname > $2 COLLATE \"C\" ORDER BY v1.nickname COLLATE \"C\" "
	if desc {
		if since != ""{
			query = "SELECT v1.nickname, v1.fullname, v1.about, v1.email FROM forum_users v2 JOIN users v1 on(v1.nickname = v2.user_nickname) " +
				"where v2.forum = $1 and v1.nickname < $2 COLLATE \"C\" ORDER BY v1.nickname COLLATE \"C\" DESC "
		}else{
			query = "SELECT v1.nickname, v1.fullname, v1.about, v1.email FROM forum_users v2 JOIN users v1 on(v1.nickname = v2.user_nickname) " +
				"where v2.forum = $1 and v1.nickname > $2 COLLATE \"C\" ORDER BY v1.nickname COLLATE \"C\" DESC "
		}
	}
	query += "LIMIT $3"
	rows, DBErr := t.dbConnection.Query(query, slug, since, limit)
	if DBErr != nil{
		return nil, models.ForumDoesntExists
	}

	defer func(){rows.Close()}()
	resultList := make([]models.User, 0)
	for rows.Next(){
		userModel := new(models.User)
		scanErr := rows.Scan(&userModel.Nickname, &userModel.FullName, &userModel.About, &userModel.Email)
		if scanErr != nil{
			return nil, models.ForumDoesntExists
		}
		resultList = append(resultList, *userModel)
	}
	return &resultList, nil
}

func (t *ForumRepository) GetForumThreads(slug string, limit int, since string, desc bool) (*[]models.ThreadModel, error){
	queryArgs := []interface{}{
		slug, limit,
	}
	query := "SELECT v1.id, v1.title, v1.author, v1.forum, v1.message, v1.votes_counter, v1.created, v1.slug FROM thread v1 " +
		"where v1.forum = $1 " +
		"GROUP BY v1.id, v1.created ORDER BY v1.created "
	if since != ""{
		if desc{
			query = "SELECT v1.id, v1.title, v1.author, v1.forum, v1.message, v1.votes_counter, v1.created, v1.slug FROM thread v1 " +
				"where v1.forum = $1 and v1.created <= $3 " +
				"GROUP BY v1.id, v1.created ORDER BY v1.created "
		}else{
			query = "SELECT v1.id, v1.title, v1.author, v1.forum, v1.message, v1.votes_counter, v1.created, v1.slug FROM thread v1 " +
				"where v1.forum = $1 and v1.created >= $3 " +
				"GROUP BY v1.id, v1.created ORDER BY v1.created "
		}
		queryArgs = append(queryArgs, since)
	}
	if desc{
		query += "DESC "
	}
	query += "LIMIT $2 "
	rows, DBErr := t.dbConnection.Query(query, queryArgs...)
	if DBErr != nil || rows.Err() != nil{
		return nil, models.ForumDoesntExists
	}
	resultList := make([]models.ThreadModel, 0)
	defer func() {
		rows.Close()
	}()
	for rows.Next(){
		model := new(models.ThreadModel)
		ScanErr := rows.Scan(&model.ID, &model.Title, &model.Author,&model.Forum,&model.Message,&model.Votes, &model.Created, &model.Slug)
		if ScanErr != nil{
			return nil, models.ForumDoesntExists
		}
		if strings.HasPrefix(model.Slug, utils.SlugCreatedSign){
			model.Slug = ""
		}
		resultList = append(resultList, *model)
	}
	return &resultList, nil
}
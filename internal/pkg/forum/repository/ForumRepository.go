package repository

import (
	"database/sql"
	"github.com/Felix1Green/DB-project/internal/pkg/models"
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
	query := "INSERT INTO forum (title, user_id, slug) VALUES ($1, $2, $3) returning title, user_id, slug"
	uniqueErrQuery := "SELECT v1.title, v1.user_id, v1.slug, v1.cnt, COUNT(v3.id) FROM ( " +
		"SELECT v1.title, v1.user_id, v1.slug, COUNT(v2.id) as cnt FROM forum v1 " +
		"left join thread v2 on (v2.forum=v1.slug) " +
		"where v1.slug = $1 group by v1.title, v1.user_id, v1.slug ) v1 " +
		"left join post v3 on (v3.forum = v1.slug) " +
		"group by v1.title, v1.user_id, v1.slug, v1.cnt"
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
	query := "SELECT v1.title, v1.user_id, v1.slug, v1.cnt, COUNT(v3.id) FROM (" +
		"SELECT v1.title, v1.user_id, v1.slug, COUNT(v2.id) as cnt FROM forum v1 " +
		"left join thread v2 on (v2.forum=v1.slug) " +
		"where v1.slug = $1 group by v1.title, v1.user_id, v1.slug ) v1 " +
		"left join post v3 on (v3.forum = v1.slug) " +
		"group by v1.title, v1.user_id, v1.slug, v1.cnt "
	res := new(models.Forum)
	ScanErr := t.dbConnection.QueryRow(query, slug).Scan(&res.Title, &res.User, &res.Slug, &res.Threads, &res.Posts)
	if ScanErr != nil{
		return nil, models.ForumDoesntExists
	}

	return res,nil
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
	uniqueErrQuery := "SELECT v1.id, v1.title, v1.author, v1.forum, v1.message, COUNT(v2.id), v1.created, v1.slug " +
		"FROM Thread v1 LEFT JOIN post v2 on(v2.thread = v1.id) WHERE v1.slug = $1 GROUP BY v1.id"

	result := t.dbConnection.QueryRow(query, queryArgs...)
	var resultID uint64
	var created string
	if result.Err() != nil || result.Scan(&resultID, &created) != nil{
		resultThread := new(models.ThreadModel)
		scanErr := t.dbConnection.QueryRow(uniqueErrQuery, thread.Slug).Scan(&resultThread.ID, &resultThread.Title,
			&resultThread.Author, &resultThread.Forum,&resultThread.Message, &resultThread.Votes, &resultThread.Created, &resultThread.Slug)
		if scanErr != nil{
			return nil, models.ForumDoesntExists
		}
		return resultThread, models.ThreadUniqueErr
	}
	return &models.ThreadModel{
		ID: resultID,
		Title: thread.Title,
		Author: thread.Author,
		Message: thread.Message,
		Created: created,
		Slug: thread.Slug,
		Votes: 0,
		Forum: slug,
	}, nil
}


func (t *ForumRepository) GetForumUsers(slug string, limit int, since string, desc bool) (*[]models.User, error){
	query := "SELECT v1.nickname, v1.fullname, v1.about, v1.email FROM " +
		"(SELECT DISTINCT v2.nickname, v2. fullname, v2.about, v2.email FROM users v2 " +
		"left JOIN post v3 on (v3.author = v2.nickname) " +
		"left JOIN thread v4 on (v4.author = v2.nickname) " +
		"where (v3.forum = $1 or v4.forum = $1) and v2.nickname > $2 COLLATE \"C\") as v1 " +
		"ORDER BY v1.nickname COLLATE \"C\" "
	if desc {
		if since != ""{
			query = "SELECT v1.nickname, v1.fullname, v1.about, v1.email FROM " +
				"(SELECT DISTINCT v2.nickname, v2. fullname, v2.about, v2.email FROM users v2 " +
				"left JOIN post v3 on (v3.author = v2.nickname) " +
				"left JOIN thread v4 on (v4.author = v2.nickname) " +
				"where (v3.forum = $1 or v4.forum = $1) and v2.nickname < $2 COLLATE \"C\") as v1 " +
				"ORDER BY v1.nickname COLLATE \"C\" DESC "
		}else{
			query = "SELECT v1.nickname, v1.fullname, v1.about, v1.email FROM " +
				"(SELECT DISTINCT v2.nickname, v2. fullname, v2.about, v2.email FROM users v2 " +
				"left JOIN post v3 on (v3.author = v2.nickname) " +
				"left JOIN thread v4 on (v4.author = v2.nickname) " +
				"where (v3.forum = $1 or v4.forum = $1) and v2.nickname > $2 COLLATE \"C\") as v1 " +
				"ORDER BY v1.nickname COLLATE \"C\" DESC "
		}
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
	query := "SELECT v1.id, v1.title, v1.author, v1.forum, v1.message, COUNT(v2.id), v1.created, v1.slug FROM thread v1 " +
		"left join post v2 on (v2.thread = v1.id) " +
		"where v1.forum = $1 " +
		"GROUP BY v1.id, v1.created ORDER BY v1.created "
	if since != ""{
		if desc{
			query = "SELECT v1.id, v1.title, v1.author, v1.forum, v1.message, COUNT(v2.id), v1.created, v1.slug FROM thread v1 " +
				"left join post v2 on (v2.thread = v1.id) " +
				"where v1.forum = $1 and v1.created <= $3 " +
				"GROUP BY v1.id, v1.created ORDER BY v1.created "
		}else{
			query = "SELECT v1.id, v1.title, v1.author, v1.forum, v1.message, COUNT(v2.id), v1.created, v1.slug FROM thread v1 " +
				"left join post v2 on (v2.thread = v1.id) " +
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
		_ = rows.Close()
	}()
	for rows.Next(){
		model := new(models.ThreadModel)
		ScanErr := rows.Scan(&model.ID, &model.Title, &model.Author,&model.Forum,&model.Message,&model.Votes, &model.Created, &model.Slug)
		if ScanErr != nil{
			return nil, models.ForumDoesntExists
		}
		resultList = append(resultList, *model)
	}
	return &resultList, nil
}
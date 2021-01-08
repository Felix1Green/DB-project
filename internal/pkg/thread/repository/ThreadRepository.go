package repository

import (
	"database/sql"
	"github.com/lib/pq"
	"fmt"
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"log"
	"strings"
)


type ThreadRepository struct {
	DBConnection *sql.DB
}

func NewThreadRepository(conn *sql.DB) *ThreadRepository{
	return &ThreadRepository{
		DBConnection: conn,
	}
}

func bulkPostInsert(rows []models.PostCreateRequestInput, query, forumName string, threadID uint64) (string, *[]interface{}) {
	ValueStrings := make([]interface{}, 0)
	QueryStrings := make([]string, 0)
	i := 0
	for _, val := range rows {
		QueryStrings = append(QueryStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)",
			i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))
		ValueStrings = append(ValueStrings, val.Parent, val.Author, val.Message, threadID, forumName)
		i++
	}
	smtp := fmt.Sprintf(query, strings.Join(QueryStrings, ","))
	return smtp, &ValueStrings
}


func (t* ThreadRepository) CreatePosts(slug uint64, forumName string, body *[]models.PostCreateRequestInput) (*[]models.PostModel, error){
	if t.DBConnection == nil{
		return nil, models.InternalDBError
	}

	insertQuery := "INSERT INTO ticket (parent,author,message,thread,forum) VALUES %s"
	resultQuery, values := bulkPostInsert(*body, insertQuery, forumName, slug)
	resultQuery += " RETURNING id, created"
	rows, DBErr := t.DBConnection.Query(resultQuery, values)
	if DBErr != nil || rows == nil || rows.Err() != nil{
		return nil, models.ParentPostDoesntExists
	}
	resultList := make([]models.PostModel, 0)
	for i := 0; rows.Next(); i++{
		item := models.PostModel{}
		ScanErr := rows.Scan(&item.ID, &item.Created)
		if ScanErr != nil{
			return nil, ScanErr
		}
		item.Author = (*body)[i].Author
		item.Message = (*body)[i].Message
		item.Parent = (*body)[i].Parent
		item.Forum = forumName
		item.Thread = slug
		resultList = append(resultList, item)
	}
	return &resultList, nil
}

func (t *ThreadRepository) CheckParentsExisting(parentsID []uint64) (bool,error){
	if t.DBConnection == nil{
		return false, models.InternalDBError
	}

	query := "SELECT COUNT(id) from post where id IN($1)"
	counter := t.DBConnection.QueryRow(query, pq.Array(parentsID))
	if counter == nil || counter.Err() != nil{
		return false, models.InternalDBError
	}
	var parentsCounter int = 0
	ScanErr := counter.Scan(&parentsCounter)
	if ScanErr != nil{
		return false, models.InternalDBError
	}

	return parentsCounter == len(parentsID), nil
}


func (t *ThreadRepository) GetThreadDetails(slug uint64) (*models.ThreadModel, error){
	if t.DBConnection == nil{
		return nil, models.InternalDBError
	}

	query := "SELECT v1.ID, v1.title, v1.author, v1.forum, v1.message, v1.votes_counter, v1.created FROM thread v1 where ID = $1"
	resultRow := t.DBConnection.QueryRow(query, slug)
	if resultRow.Err() != nil{
		return nil, models.ThreadAbsentsError
	}
	resultItem := new(models.ThreadModel)
	ScanErr := resultRow.Scan(&resultItem.ID, &resultItem.Title, &resultItem.Author, &resultItem.Forum, &resultItem.Message,
		&resultItem.Votes, &resultItem.Created)
	if ScanErr != nil{
		return nil, models.ThreadAbsentsError
	}
	return resultItem, nil
}

func (t *ThreadRepository) GetThreadDetailsBySlug(slug string) (*models.ThreadModel, error){
	if t.DBConnection == nil{
		return nil, models.InternalDBError
	}

	query := "SELECT v1.ID, v1.title, v1.author, v1.forum, v1.message, v1.votes_counter, v1.created FROM thread v1 where slug = $1"
	resultRow := t.DBConnection.QueryRow(query, slug)
	if resultRow.Err() != nil{
		return nil, models.ThreadAbsentsError
	}
	resultItem := new(models.ThreadModel)
	ScanErr := resultRow.Scan(&resultItem.ID, &resultItem.Title, &resultItem.Author, &resultItem.Forum, &resultItem.Message,
		&resultItem.Votes, &resultItem.Created)
	if ScanErr != nil{
		return nil, models.ThreadAbsentsError
	}
	return resultItem, nil
}

func (t *ThreadRepository) UpdateThreadDetails(slug uint64, input *models.ThreadUpdateInput) (*models.ThreadModel, error){
	if t.DBConnection == nil{
		return nil, models.InternalDBError
	}

	query := "UPDATE thread SET title = $1, message = $2 WHERE ID = $3 RETURNING author, forum, votes_counter, created"
	resultItem := new(models.ThreadModel)
	ScanErr := t.DBConnection.QueryRow(query, input.Title, input.Message, slug).Scan(&resultItem.Author,
		&resultItem.Forum, &resultItem.Votes, &resultItem.Created)
	if ScanErr != nil{
		return nil, models.ThreadAbsentsError
	}
	resultItem.ID = slug
	resultItem.Title = input.Title
	resultItem.Message = input.Message
	return resultItem, nil
}

func (t *ThreadRepository) GetThreadPosts(threadID uint64, limit int, since int64, sort string, desc bool)(*[]models.PostModel, error){
	if t.DBConnection == nil{
		return nil, models.InternalDBError
	}

	var query = "SELECT v1.id, v1.parent, v1.author, v1.message, v1.isEdited, v1.forum, v1.thread, v1.created from post v1 " +
		"WHERE v1.thread = $1 AND v1.id > $2 "
	if sort == "tree"{
		query += "ORDER BY CASE WHEN v1.parent = 0 THEN v1.id ELSE v1.parent END, " +
			"CASE WHEN v1.parent = 0 THEN 0 ELSE v1.id END "
	}else if sort == "parent_tree"{
		query += "AND v1.parent < $3 ORDER BY CASE WHEN v1.parent = 0 THEN v1.id ELSE v1.parent END, " +
			"CASE WHEN v1.parent = 0 THEN 0 ELSE v1.id END "
	}else{
		query += "ORDER BY v1.created "
	}
	if desc {
		query += "DESC "
	}
	if sort != "parent_tree"{
		query += "LIMIT $3 "
	}

	resultRows, DBErr := t.DBConnection.Query(query, threadID, since, limit)
	if DBErr != nil || resultRows == nil || resultRows.Err() != nil{
		return nil, models.ThreadAbsentsError
	}

	resultList := make([]models.PostModel, 0)
	for resultRows.Next(){
		item := models.PostModel{}
		ScanErr := resultRows.Scan(&item.ID, &item.Parent, &item.Author, &item.Message, &item.IsEdited, &item.Forum, &item.Thread, &item.Created)
		if ScanErr != nil{
			return nil, models.InternalDBError
		}
		resultList = append(resultList, item)
	}

	return &resultList, nil
}

func (t *ThreadRepository) IncrementThreadVotes(threadID uint64) error{
	if t.DBConnection == nil{
		return models.InternalDBError
	}
	query := "UPDATE thread SET votes_counter = votes_counter + 1 WHERE id = $1"
	_, DBErr := t.DBConnection.Exec(query, threadID)
	return DBErr
}

func (t *ThreadRepository) SetThreadVote(threadID uint64, input models.ThreadVoteInput) (*models.ThreadModel, error){
	if t.DBConnection == nil{
		return nil, models.InternalDBError
	}

	query := "INSERT INTO vote (thread_id, user_name, rating) VALUES ($1,$2,$3) on conflict(thread_id, user_name) do update set rating = $3 RETURNING ID"
	var result uint64
	ScanErr := t.DBConnection.QueryRow(query,threadID, input.Nickname, input.Voice).Scan(&result)
	if ScanErr != nil{
		UpdatingErr := t.IncrementThreadVotes(threadID)
		if UpdatingErr != nil{
			return nil, models.InternalDBError
		}
	}
	return t.GetThreadDetails(threadID)
}



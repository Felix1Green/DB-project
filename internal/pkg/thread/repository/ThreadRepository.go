package repository

import (
	"database/sql"
	"fmt"
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"github.com/Felix1Green/DB-project/internal/pkg/utils"
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

	insertQuery := "INSERT INTO post (parent,author,message,thread,forum) VALUES %s"
	resultQuery, values := bulkPostInsert(*body, insertQuery, forumName, slug)
	resultQuery += " RETURNING id, created"
	rows, DBErr := t.DBConnection.Query(resultQuery, *values...)
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

func (t *ThreadRepository) CreateSinglePost(slug uint64, forumName string, body models.PostCreateRequestInput)(*models.PostModel, error){
	if t.DBConnection == nil{
		return nil, models.InternalDBError
	}

	insertQuery := ""
	if body.Parent == 0{
		insertQuery = "INSERT INTO post (parent,author,message,thread,forum, created, path) VALUES ($1,$2,$3,$4,$5,$6, array(select currval('post_id_seq')::integer)) RETURNING id, created";
	}else{
		insertQuery = "INSERT INTO post (parent,author,message,thread,forum, created, path) VALUES ($1,$2,$3,$4,$5,$6, (SELECT path from post where id = $1) || (select currval('post_id_seq')::integer)) RETURNING id, created";
	}
	rows := t.DBConnection.QueryRow(insertQuery, body.Parent, body.Author, body.Message, slug, forumName, body.Created)
	if rows == nil || rows.Err() != nil{
		return nil, models.NoSuchUser
	}
	result := new(models.PostModel)
	scanErr := rows.Scan(&result.ID, &result.Created)
	if scanErr != nil{
		return nil, models.InternalDBError
	}
	result.Author = body.Author
	result.Message = body.Message
	result.Parent = body.Parent
	result.Forum = forumName
	result.Thread = slug

	return result, nil
}

func (t *ThreadRepository) CheckParentsExisting(parentsID uint64, slug uint64) (bool,error){
	if t.DBConnection == nil{
		return false, models.InternalDBError
	}

	query := "SELECT id from post where id = $1 and thread = $2"
	counter := t.DBConnection.QueryRow(query, parentsID, slug)
	if counter == nil || counter.Err() != nil{
		return false, models.InternalDBError
	}
	var parentsCounter int = 0
	ScanErr := counter.Scan(&parentsCounter)
	if ScanErr != nil{
		return false, models.ParentPostDoesntExists
	}

	return true, nil
}


func (t *ThreadRepository) GetThreadDetails(slug uint64) (*models.ThreadModel, error){
	if t.DBConnection == nil{
		return nil, models.InternalDBError
	}

	query := "SELECT v1.ID, v1.title, v1.author, v1.forum, v1.message, v1.votes_counter, v1.created, v1.slug FROM thread v1 where ID = $1"
	resultRow := t.DBConnection.QueryRow(query, slug)
	if resultRow.Err() != nil{
		return nil, models.ThreadAbsentsError
	}
	resultItem := new(models.ThreadModel)
	ScanErr := resultRow.Scan(&resultItem.ID, &resultItem.Title, &resultItem.Author, &resultItem.Forum, &resultItem.Message,
		&resultItem.Votes, &resultItem.Created, &resultItem.Slug)
	if ScanErr != nil{
		return nil, models.ThreadAbsentsError
	}
	if strings.HasPrefix(resultItem.Slug, utils.SlugCreatedSign){
		resultItem.Slug = ""
	}
	return resultItem, nil
}

func (t *ThreadRepository) GetThreadDetailsBySlug(slug string) (*models.ThreadModel, error){
	if t.DBConnection == nil{
		return nil, models.InternalDBError
	}

	query := "SELECT v1.ID, v1.title, v1.author, v1.forum, v1.message, v1.votes_counter, v1.created, v1.slug FROM thread v1 where slug = $1"
	resultRow := t.DBConnection.QueryRow(query, slug)
	if resultRow.Err() != nil{
		return nil, models.ThreadAbsentsError
	}
	resultItem := new(models.ThreadModel)
	ScanErr := resultRow.Scan(&resultItem.ID, &resultItem.Title, &resultItem.Author, &resultItem.Forum, &resultItem.Message,
		&resultItem.Votes, &resultItem.Created, &resultItem.Slug)
	if ScanErr != nil{
		return nil, models.ThreadAbsentsError
	}
	if strings.HasPrefix(resultItem.Slug, utils.SlugCreatedSign){
		resultItem.Slug = ""
	}
	return resultItem, nil
}

func (t *ThreadRepository) UpdateThreadDetails(slug uint64, input *models.ThreadUpdateInput) (*models.ThreadModel, error){
	if t.DBConnection == nil{
		return nil, models.InternalDBError
	}
	query := "UPDATE thread SET title = $1, message = $2 WHERE ID = $3 RETURNING author, forum, votes_counter, created, slug, message, title"
	queryArgs := []interface{}{ input.Title, input.Message, slug}
	if input.Title == ""{
		query = "UPDATE thread SET message = $1 WHERE ID = $2 RETURNING author, forum, votes_counter, created, slug, message, title"
		queryArgs = []interface{}{input.Message, slug}
	}else if input.Message == ""{
		query = "UPDATE thread SET title = $1 WHERE ID = $2 RETURNING author, forum, votes_counter, created, slug, message, title"
		queryArgs = []interface{}{input.Title, slug}
	}
	resultItem := new(models.ThreadModel)
	ScanErr := t.DBConnection.QueryRow(query, queryArgs...).Scan(&resultItem.Author,
		&resultItem.Forum, &resultItem.Votes, &resultItem.Created, &resultItem.Slug,&resultItem.Message ,&resultItem.Title)
	if ScanErr != nil{
		return nil, models.ThreadAbsentsError
	}
	resultItem.ID = slug
	if strings.HasPrefix(resultItem.Slug, utils.SlugCreatedSign){
		resultItem.Slug = ""
	}
	return resultItem, nil
}


func (t *ThreadRepository) GetThreadPosts(threadID uint64, limit int, since int64, sort string, desc bool)(*[]models.PostModel, error){
	if t.DBConnection == nil{
		return nil, models.InternalDBError
	}

	query := BuildGetThreadsQuery(sort, limit, since, desc)
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
	query := "UPDATE thread SET votes_counter = (select sum(rating) from vote where thread_id = $1) WHERE id = $1"
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
		return nil, models.NoSuchUser
	}
	return t.GetThreadDetails(threadID)
}


func BuildGetThreadsQuery(sort string, limit int, since int64, desc bool) string{
	var query = "SELECT v1.id, v1.parent, v1.author, v1.message, v1.isEdited, v1.forum, v1.thread, v1.created from post v1 WHERE v1.thread = $1 "
	switch sort{
	case "tree":
		if since > 0{
			if !desc {
				query += "AND v1.path > (select p.path from post p where p.id = $2) ORDER BY v1.path "
			}else{
				query += "AND v1.path < (select p.path from post p where p.id = $2) ORDER BY v1.path DESC "
			}
		}else{
			query += "AND v1.path[1] > $2 ORDER BY v1.path "
			if desc{
				query += "DESC "
			}
		}
		query += "LIMIT $3 "
		break
	case "parent_tree":
		if since > 0{
			if desc{
				query += "and v1.path[1] IN(SELECT p.path[1] from post p where p.thread = $1 and p.parent = 0 AND p.path[1] < (select p2.path[1] from post p2 where p2.id = $2)  order by p.path[1] DESC limit $3) ORDER BY v1.path[1] DESC, v1.path "
			}else{
				query += "and v1.path[1] IN(SELECT p.path[1] from post p where p.thread = $1 and p.parent = 0 AND p.path[1] > (select p2.path[1] from post p2 where p2.id = $2)  order by p.path[1] limit $3) ORDER BY v1.path[1], v1.path "
			}

		}else{
			if desc{
				query += "AND v1.path[1] > $2 and v1.path[1] IN(SELECT p.path[1] from post p where p.thread = $1 and p.id > $2 and p.parent = 0 order by p.path[1] DESC limit $3) ORDER BY v1.path[1] DESC, v1.path"
			}else{
				query += "AND v1.path[1] > $2 and v1.path[1] IN(SELECT p.path[1] from post p where p.thread = $1 and p.id > $2 and p.parent = 0 order by p.path[1] limit $3) ORDER BY v1.path  "
			}
		}
		break
	default:
		if since > 0{
			if desc{
				query += "AND v1.id < $2 ORDER BY v1.created DESC,v1.id DESC "
			}else{
				query += "AND v1.id > $2 ORDER BY v1.created,v1.id "
			}
		}else{
			if desc{
				query += "AND v1.id > $2 ORDER BY v1.created DESC,v1.id DESC "
			}else{
				query += "AND v1.id > $2 ORDER BY v1.created,v1.id "
			}
		}
		query += "LIMIT $3 "
		break
	}
	return query
}


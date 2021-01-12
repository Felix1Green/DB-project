package usecase

import (
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"github.com/Felix1Green/DB-project/internal/pkg/thread"
	"strconv"
)

type ThreadUseCase struct {
	repository thread.Repository
}

func NewThreadUseCase(repository thread.Repository) *ThreadUseCase{
	return &ThreadUseCase{
		repository: repository,
	}
}

func (t *ThreadUseCase) CreatePosts(slug string, body *[]models.PostCreateRequestInput) (*[]models.PostModel, error) {
	th, err := t.GetThreadID(slug)
	if err != nil {
		return nil, err
	}
	if len(*body) < 1{
		res := make([]models.PostModel, 0)
		return &res, nil
	}
	post, err := t.repository.CreatePosts(th.ID, th.Forum, body)
	if err != nil{
		return nil, err
	}
	return post, nil
}
func (t *ThreadUseCase) GetThreadID(slug string) (*models.ThreadModel, error){
	th, err := strconv.Atoi(slug)
	threadID := uint64(th)
	if err != nil{
		resp, err :=  t.repository.CheckThreadExistingBySlug(slug)
		return resp, err
	}
	resp, err := t.repository.CheckThreadExisting(threadID)
	return resp, err
}

func (t *ThreadUseCase) GetThreadDetails(slug string) (*models.ThreadModel, error){
	th, err := strconv.Atoi(slug)
	threadID := uint64(th)
	if err != nil{
		resp, err :=  t.repository.GetThreadDetailsBySlug(slug)
		return resp, err
	}
	resp, err := t.repository.GetThreadDetails(threadID)
	return resp, err
}

func (t *ThreadUseCase) UpdateThreadDetails(slug string, input *models.ThreadUpdateInput) (*models.ThreadModel, error){
	threadObj, err := t.GetThreadDetails(slug)
	if err != nil{
		return nil, models.ThreadAbsentsError
	}
	if input.Title == "" && input.Message == ""{
		return threadObj, nil
	}

	resp, err :=  t.repository.UpdateThreadDetails(threadObj.ID, input)
	return resp, err
}

func (t *ThreadUseCase) GetThreadPosts(threadSlug string, limit int, since int64, sort string, desc bool) (*[]models.PostModel, error){
	if limit < 1{
		limit = 100
	}
	if sort == ""{
		sort = "flat"
	}

	th, err := strconv.Atoi(threadSlug)
	threadID := uint64(th)
	if err != nil{
		threadObj, err := t.repository.GetThreadDetailsBySlug(threadSlug)
		if err != nil{
			return nil, models.ThreadDoesntExist
		}
		threadID = threadObj.ID
	}else{
		_, err := t.repository.GetThreadDetails(threadID)
		if err != nil{
			return nil, models.ThreadDoesntExist
		}
	}
	return t.repository.GetThreadPosts(threadID, limit, since, sort, desc)
}

func (t *ThreadUseCase) SetThreadVote(threadSlug string, input models.ThreadVoteInput) (*models.ThreadModel, error){
	th, err := t.GetThreadID(threadSlug)
	if err != nil{
		return nil, models.ThreadAbsentsError
	}
	return t.repository.SetThreadVote(th.ID, input)
}
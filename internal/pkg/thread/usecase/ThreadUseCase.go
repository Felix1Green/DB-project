package usecase

import (
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"github.com/Felix1Green/DB-project/internal/pkg/thread"
	"strconv"
	"time"
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
	result := make([]models.PostModel,0)
	if len(*body) < 1{
		return &result, nil
	}

	threadID, castErr := strconv.Atoi(slug)
	th := new(models.ThreadModel)
	if castErr != nil{
		thr, err := t.repository.GetThreadDetailsBySlug(slug)
		if err != nil{
			return nil, models.ThreadAbsentsError
		}
		th = thr
	}else{
		thr, err := t.repository.GetThreadDetails(uint64(threadID))
		if err != nil{
			return nil, models.ThreadAbsentsError
		}
		th = thr
	}
	timeString := time.Now()
	for _, val := range *body{
		if val.Parent != 0{
			avail, err := t.repository.CheckParentsExisting(val.Parent)
			if !avail || err != nil{
				return nil, err
			}
		}
		val.Created = timeString
		post, err := t.repository.CreateSinglePost(th.ID, th.Forum, val)
		if err != nil{
			return nil, err
		}
		result = append(result, *post)
	}

	return &result, nil
}

func (t *ThreadUseCase) GetThreadDetails(slug string) (*models.ThreadModel, error){
	th, err := strconv.Atoi(slug)
	threadID := uint64(th)
	if err != nil{
		threadObj, err := t.repository.GetThreadDetailsBySlug(slug)
		if err != nil{
			return nil, models.ThreadDoesntExist
		}
		threadID = threadObj.ID
	}
	return t.repository.GetThreadDetails(threadID)
}

func (t *ThreadUseCase) UpdateThreadDetails(slug string, input *models.ThreadUpdateInput) (*models.ThreadModel, error){
	thr, castErr := strconv.Atoi(slug)
	threadID := uint64(thr)
	if castErr != nil{
		th, err := t.repository.GetThreadDetailsBySlug(slug)
		if err != nil{
			return nil, models.ThreadDoesntExist
		}
		threadID = th.ID
	}
	return t.repository.UpdateThreadDetails(threadID, input)
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
	}
	return t.repository.GetThreadPosts(threadID, limit, since, sort, desc)
}

func (t *ThreadUseCase) SetThreadVote(threadSlug string, input models.ThreadVoteInput) (*models.ThreadModel, error){
	th, castErr := strconv.Atoi(threadSlug)
	threadID := uint64(th)
	if castErr != nil{
		threadObj, err := t.repository.GetThreadDetailsBySlug(threadSlug)
		if err != nil{
			return nil, models.ThreadDoesntExist
		}
		threadID = threadObj.ID
	}
	return t.repository.SetThreadVote(threadID, input)
}
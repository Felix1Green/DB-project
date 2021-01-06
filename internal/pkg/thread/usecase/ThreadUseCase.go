package usecase

import (
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"github.com/Felix1Green/DB-project/internal/pkg/thread"
)

type ThreadUseCase struct {
	repository thread.Repository
}

func NewThreadUseCase(repository thread.Repository) *ThreadUseCase{
	return &ThreadUseCase{
		repository: repository,
	}
}

func (t *ThreadUseCase) CreatePosts(slug uint64, body *[]models.PostCreateRequestInput) (*[]models.PostModel, error) {
	parentsIDs := make([]uint64, 0)
	parentsMap := make(map[uint64]bool, 0)
	for _, val := range *body{
		if _, ok := parentsMap[val.Parent]; !ok{
			parentsIDs = append(parentsIDs, val.Parent)
			parentsMap[val.Parent] = true
		}
	}
	avail, err := t.repository.CheckParentsExisting(parentsIDs)
	if err != nil || !avail{
		return nil, models.ParentPostDoesntExists
	}

	th, err := t.repository.GetThreadDetails(slug)
	if err != nil{
		return nil, models.ThreadAbsentsError
	}
	return t.repository.CreatePosts(slug, th.Forum, body)
}

func (t *ThreadUseCase) GetThreadDetails(slug uint64) (*models.ThreadModel, error){
	return t.repository.GetThreadDetails(slug)
}

func (t *ThreadUseCase) UpdateThreadDetails(slug uint64, input *models.ThreadUpdateInput) (*models.ThreadModel, error){
	return t.repository.UpdateThreadDetails(slug, input)
}

func (t *ThreadUseCase) GetThreadPosts(threadID uint64, limit int, since int64, sort string, desc bool) (*[]models.PostModel, error){
	if limit < 1{
		limit = 100
	}
	if sort == ""{
		sort = "flat"
	}
	return t.repository.GetThreadPosts(threadID, limit, since, sort, desc)
}

func (t *ThreadUseCase) SetThreadVote(threadID uint64, input models.ThreadVoteInput) (*models.ThreadModel, error){
	return t.repository.SetThreadVote(threadID, input)
}
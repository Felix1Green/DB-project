package thread

import "github.com/Felix1Green/DB-project/internal/pkg/models"

type UseCase interface{
	CreatePosts(slug string, body *[]models.PostCreateRequestInput) (*[]models.PostModel, error)
	GetThreadDetails(slug string) (*models.ThreadModel, error)
	UpdateThreadDetails(slug uint64, input *models.ThreadUpdateInput) (*models.ThreadModel, error)
	GetThreadPosts(threadID string, limit int, since int64, sort string, desc bool) (*[]models.PostModel, error)
	SetThreadVote(threadID string, input models.ThreadVoteInput) (*models.ThreadModel, error)
}

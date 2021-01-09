package thread

import "github.com/Felix1Green/DB-project/internal/pkg/models"

type Repository interface{
	CreatePosts(slug uint64, forumName string, body *[]models.PostCreateRequestInput) (*[]models.PostModel, error)
	GetThreadDetails(slug uint64) (*models.ThreadModel, error)
	UpdateThreadDetails(slug uint64, input *models.ThreadUpdateInput) (*models.ThreadModel, error)
	GetThreadPosts(threadID uint64, limit int, since int64, sort string, desc bool)(*[]models.PostModel, error)
	SetThreadVote(threadID uint64, input models.ThreadVoteInput) (*models.ThreadModel, error)
	GetThreadDetailsBySlug(slug string) (*models.ThreadModel, error)
	CreateSinglePost(slug uint64, forumName string, body models.PostCreateRequestInput)(*models.PostModel, error)
	CheckParentsExisting(parentsID uint64, slug uint64) (bool,error)
}
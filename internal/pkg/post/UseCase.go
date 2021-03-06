package post

import "github.com/Felix1Green/DB-project/internal/pkg/models"

type UseCase interface{
	GetPostDetails(id uint64, author, thread, forum bool) (*models.PostDetails, error)
	UpdatePost(id uint64, input *models.PostUpdateRequestInput) (*models.PostModel, error)
}

package post

import "github.com/Felix1Green/DB-project/internal/pkg/models"

type Repository interface{
	GetPostDetails(id uint64) (*models.PostModel, error)
	UpdatePost(id uint64, input *models.PostUpdateRequestInput) (*models.PostModel, error)
}

package usecase

import (
	"github.com/Felix1Green/DB-project/internal/pkg/forum"
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"github.com/Felix1Green/DB-project/internal/pkg/post"
	"github.com/Felix1Green/DB-project/internal/pkg/users"
)

type PostUseCase struct {
	repository post.Repository
	userRepository users.Repository
	forumRepository forum.Repository
}


func NewPostUseCase(repository post.Repository) *PostUseCase{
	return &PostUseCase{repository: repository}
}

func (t *PostUseCase) GetPostDetails(id uint64) (*models.PostModel, error){
	return t.repository.GetPostDetails(id)
}

func (t *PostUseCase) UpdatePost(id uint64, input *models.PostUpdateRequestInput) (*models.PostModel, error){
	return t.repository.UpdatePost(id, input)
}
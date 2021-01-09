package usecase

import (
	"github.com/Felix1Green/DB-project/internal/pkg/forum"
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"github.com/Felix1Green/DB-project/internal/pkg/post"
	"github.com/Felix1Green/DB-project/internal/pkg/thread"
	"github.com/Felix1Green/DB-project/internal/pkg/users"
)

type PostUseCase struct {
	repository post.Repository
	userRepository users.Repository
	forumRepository forum.Repository
	threadRepository thread.Repository
}


func NewPostUseCase(repository post.Repository, us users.Repository, forum forum.Repository, th thread.Repository) *PostUseCase{
	return &PostUseCase{
		repository: repository,
		userRepository: us,
		forumRepository: forum,
		threadRepository: th,
	}
}

func (t *PostUseCase) GetPostDetails(id uint64, author, thread, forum bool) (*models.PostDetails, error){
	result := new(models.PostDetails)
	result.Author = nil
	result.Forum = nil
	result.Thread = nil
	p, err := t.repository.GetPostDetails(id)
	if err != nil{
		return nil, err
	}
	result.Post = p
	if author{
		result.Author, _ = t.userRepository.GetProfile(result.Post.Author)
	}
	if thread{
		result.Thread, _ =t.threadRepository.GetThreadDetails(result.Post.Thread)
	}
	if forum{
		result.Forum, _ = t.forumRepository.GetForum(result.Post.Forum)
	}
	return result, nil
}

func (t *PostUseCase) UpdatePost(id uint64, input *models.PostUpdateRequestInput) (*models.PostModel, error){
	postObj, err := t.repository.GetPostDetails(id)
	if postObj == nil || input.Message == "" || input.Message == postObj.Message{
		return postObj, err
	}
	return t.repository.UpdatePost(id, input)
}
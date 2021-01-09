package usecase

import (
	"github.com/Felix1Green/DB-project/internal/pkg/forum"
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"github.com/Felix1Green/DB-project/internal/pkg/users"
	"github.com/Felix1Green/DB-project/internal/pkg/utils"
)

type ForumUseCase struct {
	repository     forum.Repository
	userRepository users.Repository
}

func NewForumUseCase(repository forum.Repository, usersRepository users.Repository) *ForumUseCase{
	return &ForumUseCase{
		repository: repository,
		userRepository: usersRepository,
	}
}

func (t *ForumUseCase) CreateForum(input *models.ForumRequestInput) (*models.Forum, error) {
	if input == nil || input.Slug == "" || input.Title == "" || input.User == "" {
		return nil, models.IncorrectInputParams
	}
	user, err := t.userRepository.GetProfile(input.User)
	if err != nil {
		return nil, models.NoSuchUser
	}

	input.User = user.Nickname
	return t.repository.CreateForum(input)
}

func (t *ForumUseCase) GetForum(slug string) (*models.Forum, error) {
	if slug == "" {
		return nil, models.IncorrectInputParams
	}

	return t.repository.GetForum(slug)
}

func (t *ForumUseCase) CreateForumThread(slug string, thread *models.ThreadRequestInput) (*models.ThreadModel, error) {
	if slug == "" || thread == nil || thread.Title == "" || thread.Author == "" {
		return nil, models.IncorrectInputParams
	}

	isEmpty := false
	if thread.Slug == ""{
		isEmpty = true
		thread.Slug = utils.RandStringRunes(13)
	}

	user, err := t.userRepository.GetProfile(thread.Author)
	if err != nil{
		return nil, models.NoSuchUser
	}

	thread.Author = user.Nickname

	forumObj, err := t.repository.GetForum(slug)
	if err != nil{
		return nil, models.ForumDoesntExists
	}
	slug = forumObj.Slug
	result, err :=  t.repository.CreateForumThread(slug, thread)
	if result != nil && isEmpty{
		result.Slug = ""
	}
	return result, err
}

func (t *ForumUseCase) GetForumUsers(slug string, limit int , since string, desc bool) (*[]models.User, error) {
	if slug == "" {
		return nil, models.IncorrectInputParams
	}

	_, err := t.repository.GetForum(slug)
	if err != nil{
		return nil, models.ForumDoesntExists
	}
	return t.repository.GetForumUsers(slug, limit, since, desc)
}

func (t *ForumUseCase) GetForumThreads(slug string, limit int, since string, desc bool) (*[]models.ThreadModel, error) {
	if slug == "" {
		return nil, models.IncorrectInputParams
	}

	_, err := t.repository.GetForum(slug)
	if err != nil{
		return nil, models.ForumDoesntExists
	}
	return t.repository.GetForumThreads(slug, limit, since, desc)
}

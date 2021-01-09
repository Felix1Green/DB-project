package forum

import "github.com/Felix1Green/DB-project/internal/pkg/models"

type UseCase interface{
	CreateForum(input *models.ForumRequestInput) (*models.Forum, error)
	GetForum(slug string) (*models.Forum, error)
	CreateForumThread(slug string, thread *models.ThreadRequestInput) (*models.ThreadModel, error)
	GetForumUsers(slug string, limit int, since string, desc bool) (*[]models.User, error)
	GetForumThreads(slug string, limit int, since string, desc bool) (*[]models.ThreadModel, error)
}
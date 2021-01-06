package forum

import (
	"database/sql"
	"fmt"
	"github.com/Felix1Green/DB-project/internal/pkg/forum"
	"github.com/Felix1Green/DB-project/internal/pkg/forum/delivery"
	"github.com/Felix1Green/DB-project/internal/pkg/forum/repository"
	"github.com/Felix1Green/DB-project/internal/pkg/forum/usecase"
	Users "github.com/Felix1Green/DB-project/internal/pkg/users"
	"github.com/gorilla/mux"
)

type Service struct {
	Handler *delivery.ForumDelivery
	Repository *repository.ForumRepository
	UseCase *usecase.ForumUseCase
	Router *mux.Router
}


func configureRouter(handler *delivery.ForumDelivery) *mux.Router{
	router := mux.NewRouter()

	router.HandleFunc("/api/forum/create/", handler.CreateForum)
	router.HandleFunc(fmt.Sprintf("/api/forum/{%s:.+}/details/", forum.SlugPathName), handler.GetForum)
	router.HandleFunc(fmt.Sprintf("/api/forum/{%s:.+}/create/", forum.SlugPathName), handler.CreateForumThread)
	router.HandleFunc(fmt.Sprintf("/api/forum/{%s:.+}/users/", forum.SlugPathName), handler.GetForumUsers)
	router.HandleFunc(fmt.Sprintf("/api/forum/{%s:.+}/threads/", forum.SlugPathName), handler.GetForumThreads)

	return router
}


func Start(DBConnection *sql.DB, usersRepository Users.Repository) *Service{
	rep := repository.NewForumRepository(DBConnection)
	uc := usecase.NewForumUseCase(rep, usersRepository)
	handler := delivery.NewForumDelivery(uc)
	router := configureRouter(handler)
	return &Service{
		Handler: handler,
		Router: router,
		UseCase: uc,
		Repository: rep,
	}
}
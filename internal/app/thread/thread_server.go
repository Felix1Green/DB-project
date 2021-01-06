package thread

import (
	"database/sql"
	"fmt"
	"github.com/Felix1Green/DB-project/internal/pkg/thread"
	"github.com/Felix1Green/DB-project/internal/pkg/thread/Delivery"
	"github.com/Felix1Green/DB-project/internal/pkg/thread/repository"
	"github.com/Felix1Green/DB-project/internal/pkg/thread/usecase"
	"github.com/gorilla/mux"
	"net/http"
)

type Service struct {
	Handler *Delivery.ThreadDelivery
	Repository thread.Repository
	UseCase thread.UseCase
	Router *mux.Router
}


func configureRouter(handler *Delivery.ThreadDelivery) *mux.Router{
	router := mux.NewRouter()

	router.HandleFunc(fmt.Sprintf("/api/thread/{%s:[0-9]+}/create/", thread.PathThreadName), handler.CreateNewPosts)
	router.HandleFunc(fmt.Sprintf("/api/thread/{%s:[0-9]+}/details/", thread.PathThreadName), handler.GetThreadDetails).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("/api/thread/{%s:[0-9]+}/details/", thread.PathThreadName), handler.UpdateThreadInfo).Methods(http.MethodPost)
	router.HandleFunc(fmt.Sprintf("/api/thread/{%s:[0-9]+}/posts/", thread.PathThreadName), handler.GetThreadPosts)
	router.HandleFunc(fmt.Sprintf("/api/thread/{%s:[0-9]+}/vote/", thread.PathThreadName), handler.SetThreadVote)

	return router
}


func Start(conn *sql.DB) *Service{
	rep := repository.NewThreadRepository(conn)
	uc := usecase.NewThreadUseCase(rep)
	handler := Delivery.NewThreadDelivery(uc)
	router := configureRouter(handler)

	return &Service{
		Handler: handler,
		Repository: rep,
		UseCase: uc,
		Router: router,
	}
}

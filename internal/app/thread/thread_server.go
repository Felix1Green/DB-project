package thread

import (
	"fmt"
	"github.com/Felix1Green/DB-project/internal/pkg/thread"
	"github.com/Felix1Green/DB-project/internal/pkg/thread/Delivery"
	"github.com/Felix1Green/DB-project/internal/pkg/thread/repository"
	"github.com/Felix1Green/DB-project/internal/pkg/thread/usecase"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
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

	router.HandleFunc(fmt.Sprintf("/api/thread/{%s:.+}/create", thread.PathThreadName), handler.CreateNewPosts)
	router.HandleFunc(fmt.Sprintf("/api/thread/{%s:.+}/details", thread.PathThreadName), handler.GetThreadDetails).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("/api/thread/{%s:.+}/details", thread.PathThreadName), handler.UpdateThreadInfo).Methods(http.MethodPost)
	router.HandleFunc(fmt.Sprintf("/api/thread/{%s:.+}/posts", thread.PathThreadName), handler.GetThreadPosts)
	router.HandleFunc(fmt.Sprintf("/api/thread/{%s:.+}/vote", thread.PathThreadName), handler.SetThreadVote)

	return router
}


func Start(conn *pgx.ConnPool) *Service{
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

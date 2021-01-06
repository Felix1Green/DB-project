package post

import (
	"database/sql"
	"fmt"
	"github.com/Felix1Green/DB-project/internal/pkg/post"
	"github.com/Felix1Green/DB-project/internal/pkg/post/delivery"
	"github.com/Felix1Green/DB-project/internal/pkg/post/repository"
	"github.com/Felix1Green/DB-project/internal/pkg/post/usecase"
	"github.com/gorilla/mux"
	"net/http"
)

type Service struct {
	Repository post.Repository
	UseCase post.UseCase
	Handler *delivery.PostDelivery
	Router *mux.Router
}


func configureRouter(handler *delivery.PostDelivery) *mux.Router{
	router := mux.NewRouter()

	router.HandleFunc(fmt.Sprintf("/api/post/{%s:[0-9]+}/details/", post.PathPostName), handler.GetPostDetails).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("/api/post/{%s:[0-9]+}/details/",post.PathPostName), handler.UpdatePostMessage).Methods(http.MethodPost)

	return router
}

func Start(DBConnection *sql.DB) *Service{
	rep := repository.NewPostRepository(DBConnection)
	uc := usecase.NewPostUseCase(rep)
	handler := delivery.NewPostDelivery(uc)
	router := configureRouter(handler)

	return &Service{
		Handler: handler,
		Router: router,
		UseCase: uc,
		Repository: rep,
	}
}
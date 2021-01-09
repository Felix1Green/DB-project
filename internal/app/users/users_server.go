package users

import (
	"database/sql"
	"fmt"
	"github.com/Felix1Green/DB-project/internal/pkg/users"
	"github.com/Felix1Green/DB-project/internal/pkg/users/delivery"
	"github.com/Felix1Green/DB-project/internal/pkg/users/repository"
	"github.com/Felix1Green/DB-project/internal/pkg/users/usecase"
	"github.com/gorilla/mux"
	"net/http"
)

type Service struct {
	Handler *delivery.UserDelivery
	Router *mux.Router
	UseCase users.UseCase
	Repository users.Repository
}

func configureRouter(handler *delivery.UserDelivery) *mux.Router{
	router := mux.NewRouter()

	router.HandleFunc(fmt.Sprintf("/api/user/{%s:.+}/create", users.NickNamePath),handler.CreateUser)
	router.HandleFunc(fmt.Sprintf("/api/user/{%s:.+}/profile", users.NickNamePath),handler.GetUser).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("/api/user/{%s:.+}/profile", users.NickNamePath),handler.UpdateUser).Methods(http.MethodPost)

	return router
}

func Start(sqlConn *sql.DB) *Service{
	rep := repository.NewUsersRepository(sqlConn)
	uc := usecase.NewUserUseCase(rep)
	handler := delivery.NewUserDelivery(uc)

	router := configureRouter(handler)
	return &Service{
		Handler: handler,
		Repository: rep,
		UseCase: uc,
		Router: router,
	}
}
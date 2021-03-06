package cleaner

import (
	"github.com/Felix1Green/DB-project/internal/pkg/cleaner/delivery"
	"github.com/Felix1Green/DB-project/internal/pkg/cleaner/repository"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

type Service struct {
	Handler *delivery.CleanerDelivery
	Repository *repository.CleanerRepository
	Router *mux.Router
}

func configureRouter(handler *delivery.CleanerDelivery)*mux.Router{
	router := mux.NewRouter()

	router.HandleFunc("/api/service/clear", handler.ClearDB)
	router.HandleFunc("/api/service/status", handler.GetStatus)
	return router
}

func Start(DbConnection *pgx.ConnPool) *Service{
	rep := repository.NewCleanerRepository(DbConnection)
	handler := delivery.NewCleanerDelivery(rep)
	router := configureRouter(handler)
	return &Service{
		Handler: handler,
		Repository: rep,
		Router: router,
	}
}

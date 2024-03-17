package httpServer

import (
	authHttp "film_library/internal/auth/delivery/http"
	repository2 "film_library/internal/auth/repository"
	usecase2 "film_library/internal/auth/usecase"
	serviceHttp "film_library/internal/service/delivery/http"
	"film_library/internal/service/repository"
	"film_library/internal/service/usecase"
	"film_library/pkg/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (s *Server) MapHandlers() error {
	db, err := storage.InitPsqlDB(s.cfg)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	if err = storage.CreateTables(db); err != nil {
		log.Printf(err.Error())
		return err
	}

	serviceRepo := repository.NewPostgresRepository(db)
	authRepo := repository2.NewPostgresRepository(db)

	serviceUC := usecase.NewServiceUsecase(serviceRepo)
	authUC := usecase2.NewAuthUsecase(authRepo)

	authR := authHttp.NewAuthHandler(authUC)
	serviceR := serviceHttp.NewServiceHandler(serviceUC, authUC)

	rtr := mux.NewRouter()
	serviceHttp.MapRoutes(rtr, serviceR)
	authHttp.MapRoutes(rtr, authR)
	http.Handle("/", rtr)

	return nil
}

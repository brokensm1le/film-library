package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

func MapRoutes(rtr *mux.Router, s *ServiceHandler) {
	api := rtr.PathPrefix("/api").Subrouter()
	api.Use(s.userIdentity)
	api.HandleFunc("/actor/add", s.CreateActor).Methods(http.MethodPost)
	api.HandleFunc("/actor/get/{actor_name:[A-Za-z+]+}", s.GetActor).Methods(http.MethodGet)
	api.HandleFunc("/actor/get_all", s.GetActors).Methods(http.MethodGet)
	api.HandleFunc("/actor/delete/{actor_name:[A-Za-z+]+}", s.DeleteActor).Methods(http.MethodDelete)
	api.HandleFunc("/actor/update/{actor_name:[A-Za-z+]+}", s.UpdateActor).Methods(http.MethodPatch)
	api.HandleFunc("/actor/search/{actor_name:[A-Za-z+]+}", s.SearchActor).Methods(http.MethodGet)

	api.HandleFunc("/film/add", s.CreateFilm).Methods(http.MethodPost)
	api.HandleFunc("/film/get/{film_name:[0-9A-Za-z.?+]+}", s.GetFilm).Methods(http.MethodGet)
	api.HandleFunc("/film/get_all", s.GetFilms).Methods(http.MethodGet)
	api.HandleFunc("/film/delete/{film_name:[0-9A-Za-z.?+]+}", s.DeleteFilm).Methods(http.MethodDelete)
	api.HandleFunc("/film/update/{film_name:[0-9A-Za-z.?+]+}", s.UpdateFilm).Methods(http.MethodPatch)
	api.HandleFunc("/film/search/{film_name:[0-9A-Za-z.?+]+}", s.SearchFilms).Methods(http.MethodGet)

	api.HandleFunc("/relation/films_by_actor", s.AddFilmsByActor).Methods(http.MethodPost)
	api.HandleFunc("/relation/actors_by_film", s.AddActorsByFilm).Methods(http.MethodPost)
	api.HandleFunc("/relation/delete", s.DeleteActorFilm).Methods(http.MethodDelete)
}

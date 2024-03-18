package http

import (
	"encoding/json"
	"film_library/internal/auth"
	"film_library/internal/cconstant"
	"film_library/internal/service"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"slices"
	"strings"
)

type ServiceHandler struct {
	serviceUC service.Usecase
	authUC    auth.Usecase
}

func NewServiceHandler(serviceUC service.Usecase, authUC auth.Usecase) *ServiceHandler {
	return &ServiceHandler{
		serviceUC: serviceUC,
		authUC:    authUC,
	}
}

// @Summary      CreateActor
// @Description  Add actor
// @Tags         actor
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string true  "Authorization"
// @Param        input	body	service.Actor  true  "actor data"
// @Success      200  {object}	auth.ResponseModel
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /actor/add [post]
func (s *ServiceHandler) CreateActor(rw http.ResponseWriter, r *http.Request) {
	var (
		data service.Actor
		resp *service.ResponseModel = &service.ResponseModel{Status: "OK"}
	)

	tokenData := r.Context().Value(cconstant.ContextValue).(auth.TokenData)
	if tokenData.Role == 0 {
		log.Printf("Request: CreateActor. Error: %s", "Don't have permission")
		http.Error(rw, fmt.Sprintf("You don't have permission for this operation."), http.StatusForbidden)
		return
	}
	log.Printf("Request: CreateActor. User with ID:%d", tokenData.Id)

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Request: CreateActor. Error: %s", "Uncorrect data")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.validateActor(&data); err != nil {
		log.Printf("Request: CreateActor. Error: %s", err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err := s.serviceUC.CreateActor(&data)
	if err != nil {
		log.Printf("Request: CreateActor. Error: %s", resp.Error)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rawResponse, _ := json.Marshal(resp)
	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write(rawResponse)
}

// @Summary      GetActor
// @Description  Get actor
// @Tags         actor
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string true  "Authorization"
// @Param        q	query string  false  "actor name"
// @Success      200  {object}	service.Actor
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /actor/get/{actor_name} [get]
func (s *ServiceHandler) GetActor(rw http.ResponseWriter, r *http.Request) {

	tokenData := r.Context().Value(cconstant.ContextValue).(auth.TokenData)
	log.Printf("Request: GetActor. User with ID:%d", tokenData.Id)

	name := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	name = strings.ReplaceAll(name, "+", " ")

	actor, err := s.serviceUC.GetActor(name)
	if err != nil {
		log.Printf("Request: GetActor. Error: %s", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rawResponse, _ := json.Marshal(actor)
	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write(rawResponse)
}

// @Summary      GetActors
// @Description  Get actors
// @Tags         actor
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string true  "Authorization"
// @Param        Sort 			header  string true	 "Sort"
// @Success      200  {array}	service.Actor
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /actor/get_all [get]
func (s *ServiceHandler) GetActors(rw http.ResponseWriter, r *http.Request) {

	tokenData := r.Context().Value(cconstant.ContextValue).(auth.TokenData)
	log.Printf("Request: GetActors. User with ID:%d", tokenData.Id)

	sort := r.Header.Get("Sort")
	if idx := slices.IndexFunc(cconstant.FieldsActor, func(c string) bool { return c == sort }); idx == -1 {
		sort = "actor_name"
	}

	actor, err := s.serviceUC.GetActors(&service.DetailsParams{Sort: sort})
	if err != nil {
		log.Printf("Request: GetActors. Error: %s", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rawResponse, _ := json.Marshal(actor)
	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write(rawResponse)
}

// @Summary      UpdateActor
// @Description  Update actor
// @Tags         actor
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string true  "Authorization"
// @Param        q		query 	string  	   false  "actor name"
// @Param        input	body	service.Actor  true   "actor data"
// @Success      200  {object}	service.ResponseModel
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /actor/update/{actor_name} [patch]
func (s *ServiceHandler) UpdateActor(rw http.ResponseWriter, r *http.Request) {
	var (
		data service.Actor
		resp *service.ResponseModel = &service.ResponseModel{Status: "OK"}
	)

	tokenData := r.Context().Value(cconstant.ContextValue).(auth.TokenData)
	if tokenData.Role == 0 {
		log.Printf("Request: UpdateActor. Error: %s", "Don't have permission")
		http.Error(rw, fmt.Sprintf("You don't have permission for this operation."), http.StatusForbidden)
		return
	}
	log.Printf("Request: UpdateActor. User with ID:%d", tokenData.Id)

	name := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	name = strings.ReplaceAll(name, "+", " ")

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Request: UpdateActor. Error: %s", err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err := s.serviceUC.UpdateActor(name, &data)
	if err != nil {
		log.Printf("Request: UpdateActor. Error: %s", resp.Error)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rawResponse, _ := json.Marshal(resp)
	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write(rawResponse)
}

// @Summary      DeleteActor
// @Description  Delete actor
// @Tags         actor
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string true  "Authorization"
// @Param        q		query 	string  	   false  "actor name"
// @Success      200  {object}	service.ResponseModel
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /actor/delete/{actor_name} [delete]
func (s *ServiceHandler) DeleteActor(rw http.ResponseWriter, r *http.Request) {
	var (
		resp *service.ResponseModel = &service.ResponseModel{Status: "OK"}
	)

	tokenData := r.Context().Value(cconstant.ContextValue).(auth.TokenData)
	if tokenData.Role == 0 {
		log.Printf("Request: DeleteActor. Error: %s", "Don't have permission")
		http.Error(rw, fmt.Sprintf("You don't have permission for this operation."), http.StatusForbidden)
		return
	}
	log.Printf("Request: DeleteActor. User with ID:%d", tokenData.Id)

	name := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	name = strings.ReplaceAll(name, "+", " ")

	err := s.serviceUC.DeleteActor(name)
	if err != nil {
		log.Printf("Request: DeleteActor. Error: %s", resp.Error)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rawResponse, _ := json.Marshal(resp)
	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write(rawResponse)
}

// @Summary      SearchActor
// @Description  Search actor
// @Tags         actor
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string true  "Authorization"
// @Param        q				query 	string false "pattern"
// @Success      200  {array}	string
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /actor/search/{actor_name} [get]
func (s *ServiceHandler) SearchActor(rw http.ResponseWriter, r *http.Request) {
	tokenData := r.Context().Value(cconstant.ContextValue).(auth.TokenData)
	log.Printf("Request: SearchActor. User with ID:%d", tokenData.Id)

	pattern := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	pattern = strings.ReplaceAll(pattern, "+", " ")

	films, err := s.serviceUC.SearchActor(pattern)
	if err != nil {
		log.Printf("Request: SearchActor. Error: %s", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rawResponse, _ := json.Marshal(films)
	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write(rawResponse)
}

// @Summary      CreateFilm
// @Description  Add film
// @Tags         film
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string true  "Authorization"
// @Param        input	body	service.Film   true  "film data"
// @Success      200  {object}	service.ResponseModel
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /film/add [post]
func (s *ServiceHandler) CreateFilm(rw http.ResponseWriter, r *http.Request) {
	var (
		data service.Film
		resp *service.ResponseModel = &service.ResponseModel{Status: "OK"}
	)

	tokenData := r.Context().Value(cconstant.ContextValue).(auth.TokenData)
	if tokenData.Role == 0 {
		log.Printf("Request: CreateFilm. Error: %s", "Don't have permission")
		http.Error(rw, fmt.Sprintf("You don't have permission for this operation."), http.StatusForbidden)
		return
	}
	log.Printf("Request: CreateFilm. User with ID:%d", tokenData.Id)

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Request: CreateFilm. Error: %s", err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.validateFilm(&data); err != nil {
		log.Printf("Request: CreateFilm. Error: %s", "Uncorrect data")
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	err := s.serviceUC.CreateFilm(&data)
	if err != nil {
		log.Printf("Request: CreateFilm. Error: %s", resp.Error)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rawResponse, _ := json.Marshal(resp)
	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write(rawResponse)
}

// @Summary      GetFilm
// @Description  Get film
// @Tags         film
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string true  "Authorization"
// @Param        q	query string  false  "film name"
// @Success      200  {object}	service.Film
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /film/get/{actor_name} [get]
func (s *ServiceHandler) GetFilm(rw http.ResponseWriter, r *http.Request) {

	tokenData := r.Context().Value(cconstant.ContextValue).(auth.TokenData)
	log.Printf("Request: GetFilm. User with ID:%d", tokenData.Id)

	name := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	name = strings.ReplaceAll(name, "+", " ")

	actor, err := s.serviceUC.GetFilm(name)
	if err != nil {
		log.Printf("GetFilm: CreateFilm. Error: %s", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rawResponse, _ := json.Marshal(actor)
	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write(rawResponse)
}

// @Summary      GetFilms
// @Description  Get films
// @Tags         film
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string true  "Authorization"
// @Param        Sort 			header  string true	 "Sort"
// @Success      200  {array}	service.Film
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /film/get_all [get]
func (s *ServiceHandler) GetFilms(rw http.ResponseWriter, r *http.Request) {

	tokenData := r.Context().Value(cconstant.ContextValue).(auth.TokenData)
	log.Printf("Request: GetFilms. User with ID:%d", tokenData.Id)

	sort := r.Header.Get("Sort")

	if idx := slices.IndexFunc(cconstant.FieldsFilm, func(c string) bool { return c == sort }); idx == -1 {
		sort = "rating"
	}

	films, err := s.serviceUC.GetFilms(&service.DetailsParams{Sort: sort})
	if err != nil {
		log.Printf("Request: GetFilms. Error: %s", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rawResponse, _ := json.Marshal(films)
	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write(rawResponse)
}

// @Summary      UpdateFilm
// @Description  Update film
// @Tags         film
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string true  "Authorization"
// @Param        q		query 	string  	   false  "film name"
// @Param        input	body	service.Film   true   "new film data"
// @Success      200  {object}	service.ResponseModel
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /film/update/{film_name} [patch]
func (s *ServiceHandler) UpdateFilm(rw http.ResponseWriter, r *http.Request) {
	var (
		data service.Film
		resp *service.ResponseModel = &service.ResponseModel{Status: "OK"}
	)

	tokenData := r.Context().Value(cconstant.ContextValue).(auth.TokenData)
	if tokenData.Role == 0 {
		log.Printf("Request: UpdateFilm. Error: %s", "Don't have permission")
		http.Error(rw, fmt.Sprintf("You don't have permission for this operation."), http.StatusForbidden)
		return
	}
	log.Printf("Request: UpdateFilm. User with ID:%d", tokenData.Id)

	name := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	name = strings.ReplaceAll(name, "+", " ")

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Request: UpdateFilm. Error: %s", err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err := s.serviceUC.UpdateFilm(name, &data)
	if err != nil {
		log.Printf("Request: UpdateFilm. Error: %s", resp.Error)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rawResponse, _ := json.Marshal(resp)
	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write(rawResponse)
}

// @Summary      DeleteFilm
// @Description  Delete film
// @Tags         film
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string true   "Authorization"
// @Param        q				query 	string false  "film name"
// @Success      200  {object}	service.ResponseModel
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /film/delete/{film_name} [delete]
func (s *ServiceHandler) DeleteFilm(rw http.ResponseWriter, r *http.Request) {
	var (
		resp *service.ResponseModel = &service.ResponseModel{Status: "OK"}
	)

	tokenData := r.Context().Value(cconstant.ContextValue).(auth.TokenData)
	if tokenData.Role == 0 {
		log.Printf("Request: DeleteFilm. Error: %s", "Don't have permission")
		http.Error(rw, fmt.Sprintf("You don't have permission for this operation."), http.StatusForbidden)
		return
	}
	log.Printf("Request: DeleteFilm. User with ID:%d", tokenData.Id)

	name := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	name = strings.ReplaceAll(name, "+", " ")

	err := s.serviceUC.DeleteFilm(name)
	if err != nil {
		log.Printf("Request: DeleteFilm. Error: %s", resp.Error)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rawResponse, _ := json.Marshal(resp)
	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write(rawResponse)
}

// @Summary      SearchFilms
// @Description  Search films
// @Tags         film
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string true  "Authorization"
// @Param        q				query 	string false "pattern"
// @Success      200  {array}	string
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /film/search/{film_name} [get]
func (s *ServiceHandler) SearchFilms(rw http.ResponseWriter, r *http.Request) {
	tokenData := r.Context().Value(cconstant.ContextValue).(auth.TokenData)
	log.Printf("Request: SearchFilms. User with ID:%d", tokenData.Id)

	pattern := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	pattern = strings.ReplaceAll(pattern, "+", " ")

	films, err := s.serviceUC.SearchFilms(pattern)
	if err != nil {
		log.Printf("Request: DeleteFilm. Error: %s", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rawResponse, _ := json.Marshal(films)
	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write(rawResponse)
}

// @Summary      AddFilmsByActor
// @Description  Add films by actor
// @Tags         relation
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string true  "Authorization"
// @Param        input	body	service.AddFilmsByActorParams  true   "relation data"
// @Success      200  {object}	service.ResponseModel
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /relation/films_by_actor [post]
func (s *ServiceHandler) AddFilmsByActor(rw http.ResponseWriter, r *http.Request) {
	var (
		data service.AddFilmsByActorParams
		resp *service.ResponseModel = &service.ResponseModel{Status: "OK"}
	)

	tokenData := r.Context().Value(cconstant.ContextValue).(auth.TokenData)
	if tokenData.Role == 0 {
		log.Printf("Request: AddFilmsByActor. Error: %s", "Don't have permission")
		http.Error(rw, fmt.Sprintf("You don't have permission for this operation."), http.StatusForbidden)
		return
	}
	log.Printf("Request: AddFilmsByActor. User with ID:%d", tokenData.Id)

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Request: AddFilmsByActor. Error: %s", err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if len(data.Films) == 0 {
		log.Printf("Request: AddFilmsByActor. Error: %s", "Uncorrect data")
		http.Error(rw, fmt.Sprintf("len data should be > 0"), http.StatusBadRequest)
		return
	}

	err := s.serviceUC.AddFilmsByActor(&data)
	if err != nil {
		log.Printf("Request: AddFilmsByActor. Error: %s", resp.Error)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rawResponse, _ := json.Marshal(resp)
	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write(rawResponse)
}

// @Summary      AddActorsByFilm
// @Description  Add actors by film
// @Tags         relation
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string true  "Authorization"
// @Param        input	body	service.AddActorsByFilmParams  true   "relation data"
// @Success      200  {object}	service.ResponseModel
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /relation/actors_by_film [post]
func (s *ServiceHandler) AddActorsByFilm(rw http.ResponseWriter, r *http.Request) {
	var (
		data service.AddActorsByFilmParams
		resp *service.ResponseModel = &service.ResponseModel{Status: "OK"}
	)

	tokenData := r.Context().Value(cconstant.ContextValue).(auth.TokenData)
	if tokenData.Role == 0 {
		log.Printf("Request: AddActorsByFilm. Error: %s", "Don't have permission")
		http.Error(rw, fmt.Sprintf("You don't have permission for this operation."), http.StatusForbidden)
		return
	}
	log.Printf("Request: AddActorsByFilm. User with ID:%d", tokenData.Id)

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Request: AddActorsByFilm. Error: %s", err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if len(data.Actors) == 0 {
		log.Printf("Request: AddActorsByFilm. Error: %s", "Uncorrect data")
		http.Error(rw, fmt.Sprintf("len data should be > 0"), http.StatusBadRequest)
		return
	}

	err := s.serviceUC.AddActorsByFilm(&data)
	if err != nil {
		log.Printf("Request: AddActorsByFilm. Error: %s", resp.Error)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rawResponse, _ := json.Marshal(resp)
	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write(rawResponse)
}

// @Summary      DeleteActorFilm
// @Description  Delete relation actor film
// @Tags         relation
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string true  "Authorization"
// @Param        input	body	service.DeleteActorFilmParams  true   "relation data"
// @Success      200  {object}	service.ResponseModel
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /relation/delete [delete]
func (s *ServiceHandler) DeleteActorFilm(rw http.ResponseWriter, r *http.Request) {
	var (
		data service.DeleteActorFilmParams
		resp *service.ResponseModel = &service.ResponseModel{Status: "OK"}
	)

	tokenData := r.Context().Value(cconstant.ContextValue).(auth.TokenData)
	if tokenData.Role == 0 {
		log.Printf("Request: DeleteActorFilm. Error: %s", "Don't have permission")
		http.Error(rw, fmt.Sprintf("You don't have permission for this operation."), http.StatusForbidden)
		return
	}
	log.Printf("Request: DeleteActorFilm. User with ID:%d", tokenData.Id)

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Request: DeleteActorFilm. Error: %s", err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err := s.serviceUC.DeleteActorFilm(&data)
	if err != nil {
		log.Printf("Request: DeleteActorFilm. Error: %s", resp.Error)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rawResponse, _ := json.Marshal(resp)
	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write(rawResponse)
}

//---------------------------------------------------------------------------------------------------------------------

func (s *ServiceHandler) validateActor(data *service.Actor) error {
	if len(data.Name) == 0 || len(data.Name) > 100 {
		return fmt.Errorf("size Name should be [1;100]")
	}
	if data.Sex != "f" && data.Sex != "m" {
		return fmt.Errorf("sex should be 'm' - male or 'f' - famale")
	}
	patternDate, _ := regexp.Compile("[1-2][0-9][0-9][0-9]-[0-1][0-9]-[0-3][0-9]")
	if !patternDate.MatchString(data.BDate) {
		return fmt.Errorf("bdate should be '2000-01-01' format")
	}

	return nil
}

func (s *ServiceHandler) validateFilm(data *service.Film) error {

	if len(data.Name) == 0 || len(data.Name) > 150 {
		return fmt.Errorf("size Name should be [1;150]")
	}
	if len(data.Desc) > 1000 {
		return fmt.Errorf("size Name should be < 1000 symbols")
	}
	if data.Rating <= 0 || data.Rating > 10 {
		return fmt.Errorf("rating should be (0;10]")
	}
	patternDate, _ := regexp.Compile("[1-2][0-9][0-9][0-9]-[0-1][0-9]-[0-3][0-9]")
	if !patternDate.MatchString(data.RDate) {
		return fmt.Errorf("rdate should be '2000-01-01' format")
	}

	return nil
}

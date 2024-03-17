package usecase

import "film_library/internal/service"

type ServiceUsecase struct {
	repo service.Repository
}

func NewServiceUsecase(repo service.Repository) service.Usecase {
	return &ServiceUsecase{repo: repo}
}

func (s *ServiceUsecase) CreateActor(params *service.Actor) error {
	return s.repo.CreateActor(params)
}

func (s *ServiceUsecase) GetActor(name string) (*service.Actor, error) {
	var (
		resp *service.Actor
		err  error
	)

	resp, err = s.repo.GetActor(name)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (s *ServiceUsecase) GetActors(params *service.DetailsParams) ([]service.Actor, error) {
	var (
		resp []service.Actor
		err  error
	)

	resp, err = s.repo.GetActors(params)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (s *ServiceUsecase) UpdateActor(name string, params *service.Actor) error {
	return s.repo.UpdateActor(name, params)
}

func (s *ServiceUsecase) DeleteActor(name string) error {
	return s.repo.DeleteActor(name)
}

func (s *ServiceUsecase) SearchActor(pattern string) ([]string, error) {
	var (
		resp []string
		err  error
	)

	resp, err = s.repo.SearchActor(pattern)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (s *ServiceUsecase) CreateFilm(params *service.Film) error {
	return s.repo.CreateFilm(params)
}

func (s *ServiceUsecase) GetFilm(name string) (*service.Film, error) {
	var (
		resp *service.Film
		err  error
	)

	resp, err = s.repo.GetFilm(name)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (s *ServiceUsecase) GetFilms(params *service.DetailsParams) ([]service.Film, error) {
	var (
		resp []service.Film
		err  error
	)

	resp, err = s.repo.GetFilms(params)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (s *ServiceUsecase) UpdateFilm(name string, params *service.Film) error {
	return s.repo.UpdateFilm(name, params)
}

func (s *ServiceUsecase) DeleteFilm(name string) error {
	return s.repo.DeleteFilm(name)
}

func (s *ServiceUsecase) SearchFilms(pattern string) ([]string, error) {
	var (
		resp []string
		err  error
	)

	resp, err = s.repo.SearchFilms(pattern)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (s *ServiceUsecase) AddFilmsByActor(params *service.AddFilmsByActorParams) error {
	return s.repo.AddFilmsByActor(params)
}

func (s *ServiceUsecase) AddActorsByFilm(params *service.AddActorsByFilmParams) error {
	return s.repo.AddActorsByFilm(params)
}

func (s *ServiceUsecase) DeleteActorFilm(params *service.DeleteActorFilmParams) error {
	return s.repo.DeleteActorFilm(params)
}

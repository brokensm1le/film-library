package service

type Usecase interface {
	CreateActor(params *Actor) error
	GetActor(name string) (*Actor, error)
	GetActors(params *DetailsParams) ([]Actor, error)
	UpdateActor(name string, params *Actor) error
	DeleteActor(name string) error
	SearchActor(pattern string) ([]string, error)

	CreateFilm(params *Film) error
	GetFilm(name string) (*Film, error)
	GetFilms(params *DetailsParams) ([]Film, error)
	UpdateFilm(name string, params *Film) error
	DeleteFilm(name string) error
	SearchFilms(pattern string) ([]string, error)

	AddFilmsByActor(params *AddFilmsByActorParams) error
	AddActorsByFilm(params *AddActorsByFilmParams) error
	DeleteActorFilm(params *DeleteActorFilmParams) error
}

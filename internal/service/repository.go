package service

type Repository interface {
	CreateActor(params *Actor) error
	GetActor(name string) (*Actor, error)
	GetActors(params *DetailsParams) ([]Actor, error)
	DeleteActor(name string) error
	UpdateActor(name string, params *Actor) error
	SearchActor(pattern string) ([]string, error)

	CreateFilm(params *Film) error
	GetFilm(name string) (*Film, error)
	GetFilms(params *DetailsParams) ([]Film, error)
	DeleteFilm(name string) error
	UpdateFilm(name string, params *Film) error
	SearchFilms(pattern string) ([]string, error)

	AddFilmsByActor(params *AddFilmsByActorParams) error
	AddActorsByFilm(params *AddActorsByFilmParams) error
	DeleteActorFilm(params *DeleteActorFilmParams) error
}

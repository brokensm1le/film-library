package cconstant

import "time"

const (
	ActorDB string = "filmdb.public.actor"
	FilmDB  string = "filmdb.public.film"
	AuthDB  string = "filmdb.public.auth"
)

const (
	Salt      = "xjifcmefdx2oxe3x"
	SignedKey = "efcj34s3dr4cwdxxjuu34"
	TokenTTL  = 6 * time.Hour
)

const (
	AuthHeader   = "Authorization"
	ContextValue = "tokenData"
)

var (
	FieldsActor = []string{"actor_name", "sex", "bdate"}
	FieldsFilm  = []string{"film_name", "release_date", "rating", "description"}
)

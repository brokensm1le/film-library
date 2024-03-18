package service

import (
	"database/sql"
)

type Actor struct {
	Name  string         `json:"name" db:"actor_name"`
	Sex   string         `json:"sex" db:"sex"`
	BDate string         `json:"bdate" db:"bdate"`
	Films sql.NullString `json:"films" db:"list_film"`
}

type Film struct {
	Name   string         `json:"name" db:"film_name"`
	RDate  string         `json:"rdate" db:"release_date"`
	Rating float32        `json:"rating" db:"rating"`
	Desc   string         `json:"desc" db:"description"`
	Actors sql.NullString `json:"actors" db:"list_actor"`
}

type FilmName struct {
	Name string `json:"name" db:"film_name"`
}

type DetailsParams struct {
	Sort string `json:"sort"`
}

type AddFilmsByActorParams struct {
	Actor string   `json:"actor"`
	Films []string `json:"films"`
}

type AddActorsByFilmParams struct {
	Film   string   `json:"film"`
	Actors []string `json:"actors"`
}

type DeleteActorFilmParams struct {
	Film  string `json:"film"`
	Actor string `json:"actor"`
}

type ResponseModel struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

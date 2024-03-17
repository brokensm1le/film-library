package repository

import (
	"context"
	"database/sql"
	"film_library/internal/cconstant"
	"film_library/internal/service"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) service.Repository {
	return &postgresRepository{db: db}
}

// ----------------------------------------------------- Actor ----------------------------------------------------------

func (p *postgresRepository) CreateActor(params *service.Actor) error {
	var (
		query = `
		INSERT INTO %[1]s; (actor_name, sex, bdate)
		VALUES ($1, $2, $3)`

		values = []any{params.Name, params.Sex, params.BDate}
	)

	query = fmt.Sprintf(query, cconstant.ActorDB)

	if _, err := p.db.Exec(query, values...); err != nil {
		return err
	}

	return nil
}

func (p *postgresRepository) GetActor(name string) (*service.Actor, error) {
	var (
		data  []service.Actor
		query = `
		SELECT actor_name, sex, bdate, list_film
		FROM %[1]s 
		WHERE actor_name = $1
		`

		values = []any{name}
	)

	query = fmt.Sprintf(query, cconstant.ActorDB)

	if err := p.db.Select(&data, query, values...); err != nil {
		return &service.Actor{}, err
	}

	if len(data) == 0 {
		return &service.Actor{}, fmt.Errorf("no actor")
	}

	return &data[0], nil
}

func (p *postgresRepository) GetActors(params *service.DetailsParams) ([]service.Actor, error) {
	var (
		data  []service.Actor
		query = `
		SELECT actor_name, sex, bdate, list_film
		FROM %[1]s
		ORDER BY `
	)

	query += params.Sort

	query = fmt.Sprintf(query, cconstant.ActorDB)

	if err := p.db.Select(&data, query); err != nil {
		return data, err
	}

	if len(data) == 0 {
		return data, fmt.Errorf("no data")
	}

	return data, nil
}

func (p *postgresRepository) DeleteActor(name string) error {
	var (
		query = `
		DELETE FROM %[1]s 
		WHERE actor_name = $1
		`

		values = []any{name}
	)

	query = fmt.Sprintf(query, cconstant.ActorDB)

	if _, err := p.db.Exec(query, values...); err != nil {
		return err
	}

	return nil
}

func (p *postgresRepository) UpdateActor(name string, params *service.Actor) error {
	var (
		query  string = `UPDATE %[1]s SET `
		values []any
		cnt    = 1
	)

	var subQuery string
	if params.Name != "" {
		subQuery += " actor_name"
		cnt++
		values = append(values, params.Name)
	}
	if params.Sex != "" {
		if cnt > 1 {
			subQuery += fmt.Sprintf(", sex")
		} else {
			subQuery += fmt.Sprintf(" sex")
		}
		cnt++
		values = append(values, params.Sex)
	}
	if params.BDate != "" {
		if cnt > 1 {
			subQuery += fmt.Sprintf(", bdate")
		} else {
			subQuery += fmt.Sprintf("bdate")
		}
		cnt++
		values = append(values, params.BDate)
	}
	if len(values) == 1 {
		query += fmt.Sprintf(" %s = \n $1 WHERE actor_name = $2", subQuery)
	} else {
		query += fmt.Sprintf("( %s ) = \n(", subQuery)
		for i := 1; i < len(values)+1; i++ {
			if i == len(values) {
				query += fmt.Sprintf("$%d)\n WHERE actor_name = $%d", i, i+1)
				continue
			}
			query += fmt.Sprintf("$%d, ", i)
		}
	}

	values = append(values, name)

	// -----------------------------------------------------------------------------------------------------------------------------

	query = fmt.Sprintf(query, cconstant.ActorDB)

	// -----------------------------------------------------------------------------------------------------------------------------

	if _, err := p.db.Exec(query, values...); err != nil {
		return err
	}

	return nil
}

func (p *postgresRepository) SearchActor(pattern string) ([]string, error) {
	var (
		data  []string
		query = `
		SELECT actor_name
		FROM %[1]s
		WHERE actor_name LIKE `
	)
	query = fmt.Sprintf(query, cconstant.ActorDB)

	query += "'%" + pattern + "%'"

	if err := p.db.Select(&data, query); err != nil {
		return data, err
	}

	if len(data) == 0 {
		return data, fmt.Errorf("no data")
	}

	return data, nil
}

// ----------------------------------------------------- FILM ----------------------------------------------------------

func (p *postgresRepository) CreateFilm(params *service.Film) error {
	var (
		query = `
		INSERT INTO %[1]s (film_name, release_date, rating, description)
		VALUES ($1, $2, $3, $4)`

		values = []any{params.Name, params.RDate, params.Rating, params.Desc}
	)

	query = fmt.Sprintf(query, cconstant.FilmDB)

	if _, err := p.db.Exec(query, values...); err != nil {
		return err
	}

	return nil
}

func (p *postgresRepository) GetFilm(name string) (*service.Film, error) {
	var (
		data  []service.Film
		query = `
		SELECT film_name, release_date, rating, description, list_actor
		FROM %[1]s 
		WHERE film_name = $1
		`

		values = []any{name}
	)

	query = fmt.Sprintf(query, cconstant.FilmDB)

	if err := p.db.Select(&data, query, values...); err != nil {
		return &service.Film{}, err
	}

	if len(data) == 0 {
		return &service.Film{}, fmt.Errorf("no actor")
	}

	return &data[0], nil
}

func (p *postgresRepository) GetFilms(params *service.DetailsParams) ([]service.Film, error) {
	var (
		data  []service.Film
		query = `
		SELECT film_name, release_date, rating, description, list_actor
		FROM %[1]s
		ORDER BY `

		//values = []any{params.Sort}
	)
	query += params.Sort
	if params.Sort == "rating" {
		query += fmt.Sprintf(" DESC")
	}

	query = fmt.Sprintf(query, cconstant.FilmDB)

	//fmt.Println(query, params.Sort)

	if err := p.db.Select(&data, query); err != nil {
		return data, err
	}

	if len(data) == 0 {
		return data, fmt.Errorf("no data")
	}

	return data, nil
}

func (p *postgresRepository) DeleteFilm(name string) error {
	var (
		query = `
		DELETE FROM %[1]s 
		WHERE film_name = $1
		`

		values = []any{name}
	)

	query = fmt.Sprintf(query, cconstant.FilmDB)

	if _, err := p.db.Exec(query, values...); err != nil {
		return err
	}

	return nil
}

func (p *postgresRepository) UpdateFilm(name string, params *service.Film) error {
	var (
		query  string = `UPDATE %[1]s SET `
		values []any
		cnt    = 1
	)

	var subQuery string
	if params.Name != "" {
		subQuery += " film_name"
		cnt++
		values = append(values, params.Name)
	}
	if params.RDate != "" {
		if cnt > 1 {
			subQuery += fmt.Sprintf(", release_date")
		} else {
			subQuery += fmt.Sprintf(" release_date")
		}
		cnt++
		values = append(values, params.RDate)
	}
	if params.Rating != 0 {
		if cnt > 1 {
			subQuery += fmt.Sprintf(", rating")
		} else {
			subQuery += fmt.Sprintf("rating")
		}
		cnt++
		values = append(values, params.Rating)
	}
	if params.Desc != "" {
		if cnt > 1 {
			subQuery += fmt.Sprintf(", description")
		} else {
			subQuery += fmt.Sprintf("description")
		}
		cnt++
		values = append(values, params.Desc)
	}
	if len(values) == 1 {
		query += fmt.Sprintf(" %s = \n $1 WHERE film_name = $2", subQuery)
	} else {
		query += fmt.Sprintf("( %s ) = \n(", subQuery)
		for i := 1; i < len(values)+1; i++ {
			if i == len(values) {
				query += fmt.Sprintf("$%d)\n WHERE film_name = $%d", i, i+1)
				continue
			}
			query += fmt.Sprintf("$%d, ", i)
		}
	}

	values = append(values, name)

	// -----------------------------------------------------------------------------------------------------------------------------

	query = fmt.Sprintf(query, cconstant.FilmDB)

	// -----------------------------------------------------------------------------------------------------------------------------

	if _, err := p.db.Exec(query, values...); err != nil {
		return err
	}

	return nil
}

func (p *postgresRepository) SearchFilms(pattern string) ([]string, error) {
	var (
		data  []string
		query = `
		SELECT film_name
		FROM %[1]s
		WHERE film_name LIKE `
	)
	query = fmt.Sprintf(query, cconstant.FilmDB)

	query += "'%" + pattern + "%'"

	if err := p.db.Select(&data, query); err != nil {
		return data, err
	}

	if len(data) == 0 {
		return data, fmt.Errorf("no data")
	}

	return data, nil
}

// ----------------------------------------------------- Relations ----------------------------------------------------------

func (p *postgresRepository) AddFilmsByActor(params *service.AddFilmsByActorParams) error {
	var (
		query []string = []string{`
			UPDATE %[1]s SET list_film = (
							 SELECT array_agg(distinct e) FROM UNNEST(list_film || '{"%[3]s"}') e) 
			WHERE actor_name = '%[4]s';
		`}
		err error
		res sql.Result
	)

	for _, film := range params.Films {
		subQuery := `
		UPDATE %[2]s SET list_actor = (
		                 SELECT array_agg(distinct e) FROM UNNEST(list_actor || '{"%[4]s"}') e) 
		WHERE film_name = '`
		subQuery += film + "';\n"
		query = append(query, subQuery)
	}

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	for _, q := range query {
		q = fmt.Sprintf(q, cconstant.ActorDB, cconstant.FilmDB, strings.Join(params.Films[:], "\",\""), params.Actor)
		fmt.Println(q)
		if res, err = tx.Exec(q); err != nil {
			tx.Rollback()
			return err
		}
		if affected, _ := res.RowsAffected(); affected == 0 {
			tx.Rollback()
			return fmt.Errorf("couldn't find film or actor")
		}
	}

	return tx.Commit()
}

func (p *postgresRepository) AddActorsByFilm(params *service.AddActorsByFilmParams) error {
	var (
		query []string = []string{`
			UPDATE %[1]s SET list_actor = (
							 SELECT array_agg(distinct e) FROM UNNEST(list_actor || '{"%[3]s"}') e) 
			WHERE film_name = '%[4]s';
		`}
		//values = []any{params.Film}
		err error
		res sql.Result
	)

	for _, actor := range params.Actors {
		subQuery := `
		UPDATE %[2]s SET list_film = (
		                 SELECT array_agg(distinct e) FROM UNNEST(list_film || '{"%[4]s"}') e) 
		WHERE actor_name = '`
		subQuery += actor + "';\n"
		query = append(query, subQuery)
	}

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	for _, q := range query {
		q = fmt.Sprintf(q, cconstant.FilmDB, cconstant.ActorDB, strings.Join(params.Actors[:], "\",\""), params.Film)
		fmt.Println(q)
		if res, err = tx.Exec(q); err != nil {
			tx.Rollback()
			return err
		}

		if affected, _ := res.RowsAffected(); affected == 0 {
			tx.Rollback()
			return fmt.Errorf("couldn't find film or actor")
		}
	}

	return tx.Commit()
}

func (p *postgresRepository) DeleteActorFilm(params *service.DeleteActorFilmParams) error {
	var (
		query []string = []string{`
			UPDATE %[1]s SET list_actor = array_remove(list_actor, '%[3]s') 
			WHERE film_name = '%[4]s';`, `
			UPDATE %[2]s SET list_film = array_remove(list_film, '%[4]s') 
			WHERE actor_name = '%[3]s';
		`}
		err error
		res sql.Result
	)
	ctx := context.Background()
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, q := range query {
		q = fmt.Sprintf(q, cconstant.FilmDB, cconstant.ActorDB, params.Actor, params.Film)
		fmt.Println(q)
		if res, err = tx.ExecContext(ctx, q); err != nil {
			tx.Rollback()
			return err
		}
		if affected, _ := res.RowsAffected(); affected == 0 {
			tx.Rollback()
			return fmt.Errorf("couldn't find film or actor")
		}
	}

	return tx.Commit()
}

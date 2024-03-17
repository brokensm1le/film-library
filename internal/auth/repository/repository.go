package repository

import (
	"film_library/internal/auth"
	"film_library/internal/cconstant"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) auth.Repository {
	return &postgresRepository{db: db}
}

func (p *postgresRepository) CreateUser(user *auth.User) error {
	var (
		query = `
		INSERT INTO %[1]s (login, password)
		VALUES ($1, $2)`

		values = []any{user.Login, user.Password}
	)

	query = fmt.Sprintf(query, cconstant.AuthDB)

	if _, err := p.db.Exec(query, values...); err != nil {
		return err
	}

	return nil
}

func (p *postgresRepository) GetUser(params *auth.SignInParams) (*auth.User, error) {
	var (
		data  []auth.User
		query = `
		SELECT *
		FROM %[1]s 
		WHERE login=$1 AND password=$2
		`

		values = []any{params.Login, params.Password}
	)

	query = fmt.Sprintf(query, cconstant.AuthDB)

	if err := p.db.Select(&data, query, values...); err != nil {
		return &auth.User{}, err
	}

	if len(data) == 0 {
		return &auth.User{}, fmt.Errorf("uncorrect login or password")
	}

	return &data[0], nil

}

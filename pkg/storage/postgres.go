package storage

import (
	"film_library/config"
	"fmt"
	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
)

func InitPsqlDB(c *config.Config) (*sqlx.DB, error) {
	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Postgres.Host, c.Postgres.Port, c.Postgres.User, c.Postgres.Password, c.Postgres.DBName, c.Postgres.SSLMode)

	return sqlx.Connect(c.Postgres.PgDriver, connectionUrl)
}

func CreateTables(db *sqlx.DB) error {
	var (
		query = `
		CREATE TABLE IF NOT EXISTS "actor"
		(
			id         		serial       not null unique,
			actor_name   	varchar(100) not null,
			sex   	   		varchar(1)   not null,
			bdate      		date		 not null,
		    list_film       text[]
		);
		CREATE TABLE IF NOT EXISTS "film"
		(
			id           serial       not null unique,
			film_name    varchar(150) not null,
			release_date date		  not null,
		    rating       real		  not null,
		    description	 varchar(1000),
		    list_actor   text[]
		);
		CREATE TABLE IF NOT EXISTS "auth"
		(
			id         	serial       not null unique,
			login   	varchar(255)  not null unique,
			password   	text  not null,
			role      	smallint	 default 0
		);
		`
	)
	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}

package database

import (
	"net/url"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	dbhost     = os.Getenv("DB_HOST")
	dbport     = os.Getenv("DB_PORT")
	dbuser     = os.Getenv("DB_USER")
	dbpassword = os.Getenv("DB_PASSWORD")
	dbName     = os.Getenv("DB_NAME")
)

func Open() (*sqlx.DB, error) {
	q := make(url.Values)
	q.Set("sslmode", "disable")

	dbConnection := url.URL{
		Scheme:   "postgres",
		Host:     dbhost,
		User:     url.UserPassword(dbuser, dbpassword),
		Path:     dbName,
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", dbConnection.String())
}

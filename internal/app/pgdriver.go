package app

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func GetPgDB(dbDSN string) (*sql.DB, error) {
	return sql.Open("pgx", dbDSN)
}

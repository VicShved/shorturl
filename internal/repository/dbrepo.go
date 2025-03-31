package repository

import (
	"context"
	"database/sql"
	"time"
)

type DBRow struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

type DBRepository struct {
	db *sql.DB
}

func GetDBRepository(db *sql.DB) (*DBRepository, error) {
	repo := &DBRepository{
		db: db,
	}
	err := repo.CreateTable()
	return repo, err
}

func (r *DBRepository) CreateTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := r.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS public.pract_keyvalue ("key" varchar NOT NULL,  value varchar NULL,CONSTRAINT pract_keyvalues_pk PRIMARY KEY ("key"))`)
	return err

}

func (r DBRepository) Save(short, original string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := r.db.ExecContext(ctx, `INSERT INTO public.pract_keyvalue ("key", value) VALUES($1, $2)`, short, original)
	return err
}

func (r DBRepository) Read(short string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	row := r.db.QueryRowContext(ctx, `SELECT value FROM public.pract_keyvalue WHERE "key" = $1`, short)
	var value sql.NullString
	err := row.Scan(&value)

	if err == sql.ErrNoRows {
		return "", false
	}

	if !value.Valid {
		return "", false
	}

	if err != nil {
		return "", false
	}
	return value.String, true
}

func (r DBRepository) Len() int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	row := r.db.QueryRowContext(ctx, `SELECT count(*) FROM public.pract_keyvalue`)
	var value sql.NullInt32
	err := row.Scan(&value)

	if err == sql.ErrNoRows {
		return 0
	}

	if !value.Valid {
		return 0
	}

	if err != nil {
		return 0
	}
	return int(value.Int32)
}

func (r DBRepository) Ping() error {
	return r.db.Ping()
}

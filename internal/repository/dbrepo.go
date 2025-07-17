// repository
package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBRow struct
type DBRow struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

// DBRepository struct
type DBRepository struct {
	db *sql.DB
}

// GetDBRepo(dsn string)
func GetDBRepo(dsn string) *DBRepository {
	// postgres driver
	pgdriver, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	// defer pgdriver.Close()
	dbrepo, err := GetDBRepository(pgdriver)
	if err != nil {
		panic(err)
	}
	return dbrepo
}

// GetDBRepository(db *sql.DB)
func GetDBRepository(db *sql.DB) (*DBRepository, error) {
	repo := &DBRepository{
		db: db,
	}
	err := repo.CreateTable()
	return repo, err
}

// CreateTable()
func (r *DBRepository) CreateTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := r.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS public.pract_keyvalue ("key" varchar NOT NULL,  value varchar NULL,CONSTRAINT pract_keyvalues_pk PRIMARY KEY ("key"))`)
	return err

}

// Save(short, original string)
func (r DBRepository) Save(short, original string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := r.db.ExecContext(ctx, `INSERT INTO public.pract_keyvalue ("key", value) VALUES($1, $2)`, short, original)
	if err != nil {
		// проверяем, что ошибка сигнализирует о потенциальном нарушении целостности данных
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = ErrPKConflict
		}
	}
	return err
}

// Read(short string)
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

// Len()
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

// Ping()
func (r DBRepository) Ping() error {
	return r.db.Ping()
}

// Batch(data *[]KeyLongURLStr)
func (r DBRepository) SaveBatch(data *[]KeyLongURLStr, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Open transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	for _, element := range *data {
		_, err := tx.ExecContext(ctx, `INSERT INTO public.pract_keyvalue ("key", value) VALUES($1, $2)`, element.Key, element.LongURL)
		if err != nil {
			// rollback transaction
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		// rollback transaction
		tx.Rollback()
		return err
	}
	return err
}

// DeleteUserUrls(shortURLs *[]string, userID string)
func (r DBRepository) DeleteUserUrls(shortURLs *[]string, userID string) error {
	return nil // TODO need realizaion
}

// Close db connection
func (r DBRepository) CloseConnect() {
	r.db.Close()
}

// UsersCount()
func (r DBRepository) UsersCount() (int, error) {
	return 0, nil
}

// UrlsCount
func (r DBRepository) UrlsCount() (int, error) {
	return 0, nil
}

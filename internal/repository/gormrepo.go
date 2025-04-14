package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ShortOriginalURL struct {
	Short    string `json:"short_url" gorm:"primaryKey;size:32"`
	Original string `json:"original_url"`
	UserID   string `json:"user_id" gorm:"primaryKey;size:36"`
}

type GormRepository struct {
	DB *gorm.DB
}

func GetGormDB(dns string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	return db, err
}

func GetGormRepo(dns string) (*GormRepository, error) {
	db, err := GetGormDB(dns)
	if err != nil {
		return nil, err
	}
	repo := &GormRepository{
		DB: db,
	}
	err = repo.Migrate()
	if err != nil {
		return nil, err
	}
	return repo, err
}

func (r *GormRepository) Migrate() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := r.DB.WithContext(ctx).AutoMigrate(&ShortOriginalURL{})
	return err
}

func (r GormRepository) Save(short string, original string, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	row := ShortOriginalURL{Short: short, Original: original, UserID: userID}
	result := r.DB.WithContext(ctx).Create(&row)
	if result.Error != nil {
		// проверяем, что ошибка сигнализирует о потенциальном нарушении целостности данных
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return ErrPKConflict
		}
	}
	return nil
}

func (r GormRepository) Read(short string, userID string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	row := ShortOriginalURL{}
	result := r.DB.WithContext(ctx).First(&row, short, userID)

	if result.Error != nil {
		return "", false
	}

	return row.Original, true
}

func (r GormRepository) Len() int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var count int64
	result := r.DB.WithContext(ctx).Model(&ShortOriginalURL{}).Count(&count)

	if result.Error != nil {
		return 0
	}

	return int(count)
}

func (r GormRepository) Ping() error {
	sqlDB, _ := r.DB.DB()
	return sqlDB.Ping()
}

func (r GormRepository) Batch(data *[]KeyLongURLStr, userID string) error {
	var rows []ShortOriginalURL
	for _, element := range *data {
		rows = append(rows, ShortOriginalURL{Short: element.Key, Original: element.LongURL, UserID: userID})
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	result := r.DB.WithContext(ctx).Create(&rows)
	return result.Error
}

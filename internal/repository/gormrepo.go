package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ShortOriginalURL struct {
	Short    string `json:"short_url" gorm:"primaryKey"`
	Original string `json:"original_url"`
}

type GormRepository struct {
	db *gorm.DB
}

func GetGormDB(dns string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	return db, err
}

func GetGormRepo(dns string) (*GormRepository, error) {
	db, _ := GetGormDB(dns)
	repo := &GormRepository{
		db: db,
	}
	err := repo.Migrate()
	return repo, err
}

func (r *GormRepository) Migrate() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := r.db.WithContext(ctx).AutoMigrate(&ShortOriginalURL{})
	return err
}

func (r GormRepository) Save(short, original string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	row := ShortOriginalURL{Short: short, Original: original}
	result := r.db.WithContext(ctx).Create(&row)
	if result.Error != nil {
		// проверяем, что ошибка сигнализирует о потенциальном нарушении целостности данных
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrPKConflict
		}
	}
	return nil
}

func (r GormRepository) Read(short string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	row := ShortOriginalURL{}
	result := r.db.WithContext(ctx).First(&row, short)

	if result.Error != nil {
		return "", false
	}

	return row.Original, true
}

func (r GormRepository) Len() int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var count int64
	result := r.db.WithContext(ctx).Model(&ShortOriginalURL{}).Count(&count)

	if result.Error != nil {
		return 0
	}

	return int(count)
}

func (r GormRepository) Ping() error {
	sqlDB, _ := r.db.DB()
	return sqlDB.Ping()
}

func (r GormRepository) Batch(data *[]KeyLongURLStr) error {
	var rows []ShortOriginalURL
	for _, element := range *data {
		rows = append(rows, ShortOriginalURL{Short: element.Key, Original: element.LongURL})
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	result := r.db.WithContext(ctx).Create(&rows)

	// // Open transaction
	// tx, err := r.db.BeginTx(ctx, nil)
	// if err != nil {
	// 	return err
	// }
	// for _, element := range *data {
	// 	_, err := tx.ExecContext(ctx, `INSERT INTO public.pract_keyvalue ("key", value) VALUES($1, $2)`, element.Key, element.LongURL)
	// 	if err != nil {
	// 		// rollback transaction
	// 		tx.Rollback()
	// 		return err
	// 	}
	// }
	// err = tx.Commit()
	// if err != nil {
	// 	// rollback transaction
	// 	tx.Rollback()
	// 	return err
	// }
	return result.Error
}

package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/VicShved/shorturl/internal/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// KeyOriginalURL struct
type KeyOriginalURL struct {
	Key       string `json:"short_url" gorm:"primaryKey;size:32"`
	Original  string `json:"original_url"`
	UserID    string `json:"user_id" gorm:"primaryKey;size:36"`
	IsDeleted bool   `json:"is_deleted" gorm:"is_deleted"`
}

// GormRepository struct
type GormRepository struct {
	DB *gorm.DB
}

// GetGormDB(dns string)
func GetGormDB(dns string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	return db, err
}

// GetGormRepo(dns string)
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

// Migrate()
func (r *GormRepository) Migrate() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := r.DB.WithContext(ctx).AutoMigrate(&KeyOriginalURL{})
	return err
}

// Save(short string, original string, userID string)
func (r GormRepository) Save(short string, original string, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	row := KeyOriginalURL{Key: short, Original: original, UserID: userID}
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

// Read(short string, userID string)
func (r GormRepository) Read(short string, userID string) (string, bool, bool) {
	logger.Log.Debug("Read", zap.String("UserID", userID))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	row := KeyOriginalURL{}
	result := r.DB.WithContext(ctx).Where(KeyOriginalURL{Key: short}).First(&row) //
	if result.Error != nil {
		return "", false, false
	}
	original := row.Original

	result = r.DB.WithContext(ctx).Where(KeyOriginalURL{Key: short, UserID: userID}).First(&row)
	if result.Error != nil {
		logger.Log.Error("Err", zap.String("err", result.Error.Error()))
		return original, true, false
	}

	logger.Log.Debug("Read row result", zap.Any("row", row))
	return original, true, row.IsDeleted
}

// Len()
func (r GormRepository) Len() int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var count int64
	result := r.DB.WithContext(ctx).Model(&KeyOriginalURL{}).Count(&count)

	if result.Error != nil {
		return 0
	}

	return int(count)
}

// Ping()
func (r GormRepository) Ping() error {
	sqlDB, _ := r.DB.DB()
	return sqlDB.Ping()
}

// Batch(data *[]KeyLongURLStr, userID string)
func (r GormRepository) Batch(data *[]KeyLongURLStr, userID string) error {
	var rows []KeyOriginalURL
	for _, element := range *data {
		rows = append(rows, KeyOriginalURL{Key: element.Key, Original: element.LongURL, UserID: userID})
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	result := r.DB.WithContext(ctx).Create(&rows)
	if result.Error != nil {
		// проверяем, что ошибка сигнализирует о потенциальном нарушении целостности данных
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return ErrPKConflict
		}
	}
	return nil
}

// GetUserUrls(userID string)
func (r GormRepository) GetUserUrls(userID string) (*[]KeyOriginalURL, error) {
	var rows []KeyOriginalURL
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	result := r.DB.WithContext(ctx).Select("Key", "Original").Where("user_id = ?", userID).Find(&rows)
	if result.Error != nil {
		return nil, result.Error
	}
	return &rows, nil
}

// DelUserUrls(shortURLs *[]string, userID string)
func (r GormRepository) DelUserUrls(shortURLs *[]string, userID string) error {
	// Создаю буферизированный канал
	ch := make(chan string, 10)
	defer close(ch)

	// Инициализирую счетчик горутин
	var wg sync.WaitGroup
	for _, short := range *shortURLs {
		wg.Add(1)
		// Для каждого ключа запускаю горутину, которая записывает ключ в канал
		go func(shortURL string) {
			defer wg.Done()
			ch <- shortURL
		}(short)
	}

	// Запускаю горутину, которая считывает из канала ключи и далее помечает строки с данными ключами как уделенные
	go func(ch chan string) {
		var shorts []string
		for sh := range ch {
			shorts = append(shorts, sh)
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		result := r.DB.WithContext(ctx).Model(
			KeyOriginalURL{}).Where("user_id = ?", userID).Where(
			"key IN ?", shorts,
		).Updates(KeyOriginalURL{IsDeleted: true})
		if result.Error != nil {
			logger.Log.Error("Error", zap.String("Err", result.Error.Error()))
		}
		logger.Log.Debug("DELETE DONE", zap.Any("shorts", shorts))
	}(ch)
	// Жду пока все горутины записи в канал выполнятся
	wg.Wait()
	logger.Log.Debug("return from Func Delete ")
	return nil
}

// Close connection
func (r GormRepository) Close() {
	sqlDB, _ := r.DB.DB()
	sqlDB.Close()
}

// UsersCount()
func (r GormRepository) UsersCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var count int64
	result := r.DB.WithContext(ctx).Model(&KeyOriginalURL{}).Distinct("UserID").Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(count), nil
}

// UrlsCount
func (r GormRepository) UrlsCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var count int64
	result := r.DB.WithContext(ctx).Model(&KeyOriginalURL{}).Distinct("Key").Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(count), nil
}

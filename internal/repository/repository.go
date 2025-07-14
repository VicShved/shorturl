package repository

import (
	"errors"
)

// KeyLongURLStr struct
type KeyLongURLStr struct {
	Key     string
	LongURL string
}

// RepoInterface interface
type RepoInterface interface {
	Save(key string, value string, userID string) error
	Read(key string, userID string) (string, bool, bool)
	Ping() error
	Len() int
	Batch(data *[]KeyLongURLStr, userID string) error
	GetUserUrls(userID string) (*[]KeyOriginalURL, error)
	DelUserUrls(shortURLs *[]string, userID string) error
	Close()
}

// ErrPKConflict
var ErrPKConflict = errors.New("PK conflict")

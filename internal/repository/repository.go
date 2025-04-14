package repository

import (
	"errors"
)

type KeyLongURLStr struct {
	Key     string
	LongURL string
}

type RepoInterface interface {
	Save(key string, value string, userID string) error
	Read(key string, userID string) (string, bool)
	Ping() error
	Len() int
	Batch(data *[]KeyLongURLStr, userID string) error
}

var ErrPKConflict = errors.New("PK conflict")

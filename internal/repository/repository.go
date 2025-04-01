package repository

type KeyLongURLStr struct {
	Key     string
	LongURL string
}

type RepoInterface interface {
	Save(key string, value string) error
	Read(key string) (string, bool)
	Ping() error
	Len() int
	Batch(*[]KeyLongURLStr) error
}

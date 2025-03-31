package service

type SaverReader interface {
	Save(key string, value string) error
	Read(key string) (string, bool)
	Ping() error
	Len() int
}

type Service struct {
	SaverReader
}
type ShortenService struct {
	repo SaverReader
}

func (s *ShortenService) Save(key string, value string) error {
	return s.repo.Save(key, value)
}

func (s *ShortenService) Read(key string) (string, bool) {
	return s.repo.Read(key)
}

func (s *ShortenService) Ping() error {
	return s.repo.Ping()
}

func (s *ShortenService) Len() int {
	return s.repo.Len()
}

func GetService(repo SaverReader) *ShortenService {
	return &ShortenService{repo: repo}
}

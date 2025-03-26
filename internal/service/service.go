package service

import "github.com/VicShved/shorturl/internal/repository"

type SaverReader interface {
	Save(key string, value string) error
	Read(key string) (string, bool)
}

type Service struct {
	SaverReader
}
type ShortenService struct {
	repo *repository.FileRepository
}

func (s *ShortenService) Save(key string, value string) error {
	return s.repo.Save(key, value)
}

func (s *ShortenService) Read(key string) (string, bool) {
	return s.repo.Read(key)
}

func GetService(repo *repository.FileRepository) *ShortenService {
	return &ShortenService{repo: repo}
}

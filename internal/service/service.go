package service

import (
	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/repository"
	// "go.uber.org/zap"
)

type BatchReqJSON struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchRespJSON struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

//	type Service struct {
//		RepoInterface
//	}
type ShortenService struct {
	repo    repository.RepoInterface
	baseurl string
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

func (s *ShortenService) Batch(indata *[]BatchReqJSON) ([]BatchRespJSON, error) {
	var results []BatchRespJSON
	var repodata []repository.KeyLongURLStr
	// Prepare results & data for repo
	for _, element := range *indata {
		shorturl, key := s.GetShortURL(&element.OriginalURL)
		res := BatchRespJSON{
			CorrelationID: element.CorrelationID,
			ShortURL:      *shorturl,
		}
		results = append(results, res)
		repodata = append(repodata, repository.KeyLongURLStr{Key: *key, LongURL: element.OriginalURL})

	}

	// batch on repo layer
	err := s.repo.Batch(&repodata)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *ShortenService) GetShortURL(longURL *string) (*string, *string) {
	key := app.Hash(*longURL)
	newurl := (*s).baseurl + "/" + key
	return &newurl, &key

}

func GetService(repo repository.RepoInterface, baseurl string) *ShortenService {
	return &ShortenService{repo: repo, baseurl: baseurl}
}

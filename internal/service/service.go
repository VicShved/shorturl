// service package
package service

import (
	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/repository"
	// "go.uber.org/zap"
)

// BatchReqJSON struct
type BatchReqJSON struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchRespJSON struct
type BatchRespJSON struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// UserURLRespJSON struct
type UserURLRespJSON struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// ShortenService struct
type ShortenService struct {
	repo    repository.RepoInterface
	baseurl string
}

// Save(key string, value string, userID string)
func (s *ShortenService) Save(key string, value string, userID string) error {
	return s.repo.Save(key, value, userID)
}

// Read(key string, userID string)
func (s *ShortenService) Read(key string, userID string) (string, bool, bool) {
	return s.repo.Read(key, userID)
}

// Ping()
func (s *ShortenService) Ping() error {
	return s.repo.Ping()
}

// Len()
func (s *ShortenService) Len() int {
	return s.repo.Len()
}

// Batch(indata *[]BatchReqJSON, userID string)
func (s *ShortenService) Batch(indata *[]BatchReqJSON, userID string) ([]BatchRespJSON, error) {
	var results []BatchRespJSON
	var repodata []repository.KeyLongURLStr
	// Prepare results & data for repo
	for _, element := range *indata {
		shorturl, key := s.GetShortURLFromLong(&element.OriginalURL)
		res := BatchRespJSON{
			CorrelationID: element.CorrelationID,
			ShortURL:      *shorturl,
		}
		results = append(results, res)
		repodata = append(repodata, repository.KeyLongURLStr{Key: *key, LongURL: element.OriginalURL})

	}

	// batch on repo layer
	err := s.repo.Batch(&repodata, userID)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetUserURLs(userID string)
func (s *ShortenService) GetUserURLs(userID string) (*[]UserURLRespJSON, error) {
	var results []UserURLRespJSON
	localResult, err := s.repo.GetUserUrls(userID)
	if err != nil {
		return nil, err
	}
	for _, elem := range *localResult {
		results = append(results, UserURLRespJSON{ShortURL: *s.GetShortURLFromKey(elem.Key), OriginalURL: elem.Original})
	}
	return &results, nil
}

// DelUserURLs(shortURLs *[]string, userID string)
func (s *ShortenService) DelUserURLs(shortURLs *[]string, userID string) error {
	err := s.repo.DelUserUrls(shortURLs, userID)
	return err
}

// GetShortURLFromKey(key string)
func (s *ShortenService) GetShortURLFromKey(key string) *string {
	newurl := (*s).baseurl + "/" + key
	return &newurl
}

// GetShortURLFromLong(longURL *string)
func (s *ShortenService) GetShortURLFromLong(longURL *string) (*string, *string) {
	key := app.Hash(string(*longURL))
	newurl := s.GetShortURLFromKey(key)
	return newurl, &key
}

// GetService(repo repository.RepoInterface, baseurl string)
func GetService(repo repository.RepoInterface, baseurl string) *ShortenService {
	return &ShortenService{repo: repo, baseurl: baseurl}
}

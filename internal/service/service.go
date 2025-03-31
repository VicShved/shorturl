package service

import (
	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/logger"
	"go.uber.org/zap"
)

type SaverReader interface {
	Save(key string, value string) error
	Read(key string) (string, bool)
	Ping() error
	Len() int
}

type BatchReqJSON struct {
	CorrelationID string `json:correlation_id`
	OriginalURL   string `json:original_url`
}

type BatchRespJSON struct {
	CorrelationID string `json:correlation_id`
	ShortURL      string `json:short_url`
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

func (s *ShortenService) Batch(indata *[]BatchReqJSON) ([]BatchRespJSON, error) {
	var results []BatchRespJSON
	for _, element := range *indata {
		logger.Log.Info("elem", zap.String("id", element.CorrelationID), zap.String("Original", element.OriginalURL))
		res := BatchRespJSON{
			CorrelationID: element.CorrelationID,
			ShortURL:      app.Hash(element.OriginalURL),
		}
		results = append(results, res)
		logger.Log.Info("res", zap.String("id", res.CorrelationID), zap.String("short", res.ShortURL))
		err := s.Save(res.ShortURL, element.OriginalURL)
		if err != nil {
			logger.Log.Error("Error", zap.Error(err))
			return results, err
		}
	}
	return results, nil
}

func GetService(repo SaverReader) *ShortenService {
	return &ShortenService{repo: repo}
}

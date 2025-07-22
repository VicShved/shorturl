package gserver

import (
	"context"
	"errors"

	pb "github.com/VicShved/shorturl/internal/gserver/proto"
	"github.com/VicShved/shorturl/internal/logger"
	"github.com/VicShved/shorturl/internal/repository"
	"github.com/VicShved/shorturl/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func getUserID(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	logger.Log.Debug("Get", zap.Any("md", md))

	users := md.Get("userID")
	logger.Log.Debug("getUserID", zap.Any("users", users))
	userID := users[0]
	return userID
}

// Get реализует получение длинного адреса из короткого ключа (хеша)
func (s *GServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	var response pb.GetResponse

	logger.Log.Debug("Get ", zap.Any("In", in))
	userID := getUserID(ctx)
	logger.Log.Debug("(s *GServer) Get(", zap.String("userID", userID))
	url, exists, isDeleted := s.serv.Read(in.Key, userID)
	if exists {
		return nil, status.Errorf(codes.NotFound, "Нет такого ключа у пользователя %s", userID)
	}
	if isDeleted {
		return nil, status.Errorf(codes.Unavailable, "Ключ удален")
	}
	response.Url = url
	return &response, nil
}

func (s *GServer) Post(ctx context.Context, in *pb.PostRequest) (*pb.PostResponse, error) {
	var response pb.PostResponse
	userID := getUserID(ctx)

	newurl, key := s.serv.GetShortURLFromLong(&in.Url)

	err := s.serv.Save(*key, in.Url, userID)

	if err != nil && errors.Is(err, repository.ErrPKConflict) {
		return nil, status.Errorf(codes.AlreadyExists, "Данный адрес уже записан")
	}

	response.Result = *newurl
	return &response, nil
}

func convert2ServBatchReqJSON(indata []*pb.BatchRequestElement) []service.BatchReqJSON {
	outDatas := make([]service.BatchReqJSON, len(indata))
	for i, source := range indata {
		logger.Log.Debug("convert2ServBatchReqJSON", zap.Any("source", source.GetCorrelationId()))
		outDatas[i].CorrelationID = source.GetCorrelationId()
		outDatas[i].OriginalURL = source.GetOriginalUrl()
		logger.Log.Debug("convert2ServBatchReqJSON", zap.Any("outDatas[i]", outDatas[i]))
	}
	return outDatas
}

func convert2BatchResponseElement(indata []service.BatchRespJSON) []*pb.BatchResponseElement {
	outData := make([]*pb.BatchResponseElement, len(indata))
	for i, source := range indata {
		outEl := new(pb.BatchResponseElement)
		outEl.CorrelationId = source.CorrelationID
		outEl.ShortUrl = source.ShortURL
		outData[i] = outEl
	}
	return outData
}

func (s *GServer) Batch(ctx context.Context, in *pb.BatchRequest) (*pb.BatchResponse, error) {
	var response pb.BatchResponse
	userID := getUserID(ctx)

	logger.Log.Debug("Batch", zap.Any("in", in.GetData()))
	servData := convert2ServBatchReqJSON(in.GetData())
	logger.Log.Debug("Batch", zap.Any("servData", servData))

	results, err := s.serv.Batch(&servData, userID)
	if err != nil && errors.Is(err, repository.ErrPKConflict) {
		return nil, status.Errorf(codes.AlreadyExists, "")
	}

	response.Data = convert2BatchResponseElement(results)
	logger.Log.Debug("Batch handled", zap.Any("response", response.Data))
	return &response, nil
}

func (s *GServer) PingDB(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	err := s.serv.Ping()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "БД недоступна")
	}
	return nil, nil
}

func convert2GetUserURLsElements(in []service.UserURLRespJSON) []*pb.GetUserURLsElement {
	results := make([]*pb.GetUserURLsElement, len(in))
	for i, source := range in {
		element := new(pb.GetUserURLsElement)
		element.ShortUrl = source.ShortURL
		element.OriginalUrl = source.OriginalURL
		results[i] = element
	}
	return results
}

func (s *GServer) GetUserURLs(ctx context.Context, in *pb.Empty) (*pb.GetUserURLsResponse, error) {
	var response pb.GetUserURLsResponse
	userID := getUserID(ctx)
	outdata, err := s.serv.GetUserURLs(userID)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if len(*outdata) == 0 {
		return nil, status.Errorf(codes.NotFound, "Записи отсутствуют")
	}
	response.Elements = convert2GetUserURLsElements(*outdata)
	return &response, nil
}

func (s *GServer) DelUserURLs(ctx context.Context, in *pb.DelUserURLsRequest) (*pb.Empty, error) {
	userID := getUserID(ctx)
	shorts := in.GetShorts()
	err := s.serv.DelUserURLs(&shorts, userID)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "Ошибка сервера")
	}
	return nil, nil
}

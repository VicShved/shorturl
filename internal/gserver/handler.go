package gserver

import (
	"context"

	pb "github.com/VicShved/shorturl/internal/gserver/proto"
	"github.com/VicShved/shorturl/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Get реализует интефейс
func (s *GServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	var response pb.GetResponse

	logger.Log.Debug("Get ", zap.Any("In", in))
	userID := "userID"
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

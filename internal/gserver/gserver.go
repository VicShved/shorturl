package gserver

import (
	pb "github.com/VicShved/shorturl/internal/gserver/proto"
	"github.com/VicShved/shorturl/internal/service"
	"google.golang.org/grpc"
)

// GServer
type GServer struct {
	pb.UnimplementedShortenerServer
	serv *service.ShortenService
}

func GetServer(serv *service.ShortenService) (*grpc.Server, error) {
	server := grpc.NewServer()
	gServer := GServer{serv: serv}
	pb.RegisterShortenerServer(server, &gServer)
	return server, nil
}

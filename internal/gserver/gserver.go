package gserver

import (
	pb "github.com/VicShved/shorturl/internal/gserver/proto"
	"github.com/VicShved/shorturl/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// GServer
type GServer struct {
	pb.UnimplementedShortenerServiceServer
	serv *service.ShortenService
}

func GetServer(serv *service.ShortenService) (*grpc.Server, error) {
	keepAlive := grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionAgeGrace: 84000})
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(authUnaryInterceptor),
		keepAlive,
		grpc.MaxRecvMsgSize(1024*1024*1000),
		grpc.MaxSendMsgSize(1024*1024*100),
		grpc.ConnectionTimeout(60000),
	)
	gServer := GServer{serv: serv}
	pb.RegisterShortenerServiceServer(server, &gServer)
	return server, nil
}

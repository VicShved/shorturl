package gserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/VicShved/shorturl/internal/app"
	pb "github.com/VicShved/shorturl/internal/gserver/proto"
	"github.com/VicShved/shorturl/internal/middware"
	"github.com/VicShved/shorturl/internal/repository"
	"github.com/VicShved/shorturl/internal/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var baseURL = "localhost:8080"
var testAuthToken = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIyYjc2ZTQwZi0wZjRkLTNlMzEtN2E4ZC1kNDE0NGQ2ZjFlM2QifQ.WsuYghl_U6-651MZekM3ZlNbiwpcJ08K-TSfCstpMGjb3Ev4RsvVsxhWzfYV3iFlIoKFLm9z_rNM6Y747kLrag"
var badAuthToken = "BADeyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIyYjc2ZTQwZi0wZjRkLTNlMzEtN2E4ZC1kNDE0NGQ2ZjFlM2QifQ.WsuYghl_U6-651MZekM3ZlNbiwpcJ08K-TSfCstpMGjb3Ev4RsvVsxhWzfYV3iFlIoKFLm9z_rNM6Y747kLrag"

func setup() *grpc.Server {
	app.ServerConfig.BaseURL = baseURL
	repo := repository.GetFileRepository(app.ServerConfig.FileStoragePath)
	serv := service.GetService(repo, baseURL)
	gserver, _ := GetServer(serv)
	listener, _ := net.Listen("tcp", baseURL)
	go gserver.Serve(listener)
	return gserver
}

func shutdown(server *grpc.Server) {
	server.GracefulStop()
}

func TestMain(m *testing.M) {
	gserver := setup()
	code := m.Run()
	shutdown(gserver)
	os.Exit(code)
}

func getAuthToken() (string, error) {
	conn, err := grpc.NewClient(baseURL, grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewShortenerServiceClient(conn)
	// md := metadata.Pairs(middware.AuthorizationCookName, "")
	ctx := context.Background()
	var header metadata.MD
	_, err = c.PingDB(ctx, &pb.Empty{}, grpc.Header(&header))
	if err != nil {
		log.Print(err)
	}
	authToken := header.Get(middware.AuthorizationCookName)
	if len(authToken) == 0 {
		return "", errors.New("Сервер не возвратил auth token")
	}
	fmt.Println("authToken ", authToken)
	testAuthToken = authToken[0]
	return authToken[0], nil
}

func TestGetAuthToken(t *testing.T) {
	tokenStr, err := getAuthToken()
	if err != nil {
		assert.True(t, false, err.Error())
	}
	token, _, err := middware.ParseTokenUserID(tokenStr)
	if err != nil {
		assert.True(t, false, err.Error())
	}
	assert.True(t, token.Valid, "Токен невалиден")
}

func post(tokenStr string, url string) (*pb.PostResponse, error) {
	conn, err := grpc.NewClient(baseURL, grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewShortenerServiceClient(conn)
	md := metadata.Pairs(middware.AuthorizationCookName, tokenStr)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	var header metadata.MD
	return c.Post(ctx, &pb.PostRequest{Url: url}, grpc.Header(&header))
}
func TestPost(t *testing.T) {
	tokenStr, err := getAuthToken()
	response, err := post(tokenStr, "https://pract.org")
	if err != nil {
		log.Print(err)
	}
	fmt.Println("response = ", response.GetResult())
}

func get(tokenStr string, shortUrl string) (*pb.GetResponse, error) {
	conn, err := grpc.NewClient(baseURL, grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewShortenerServiceClient(conn)
	md := metadata.Pairs(middware.AuthorizationCookName, tokenStr) // middware.AuthorizationCookName, authToken
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	var header metadata.MD
	return c.Get(ctx, &pb.GetRequest{Key: shortUrl}, grpc.Header(&header))

}

func TestGet(t *testing.T) {
	url := "https://pract.org"
	tokenStr, err := getAuthToken()
	postResponse, err := post(tokenStr, url)
	respUrl := postResponse.GetResult()
	splits := strings.Split(respUrl, "/")
	response, err := get(tokenStr, splits[1])
	if err != nil {
		log.Print(err)
	}
	fmt.Println("postResponse ", postResponse.GetResult())
	assert.Equal(t, url, response.GetUrl())
	fmt.Println("response = ", response.String())
}

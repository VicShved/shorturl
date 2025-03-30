package app

import (
	"flag"
	"os"
)

type ServerConfigStruct struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
	DBDSN           string
}

var ServerConfig ServerConfigStruct

func GetServerConfig() *ServerConfigStruct {
	flag.StringVar(&ServerConfig.ServerAddress, "a", "localhost:8080", "start base url")
	flag.StringVar(&ServerConfig.BaseURL, "b", "http://localhost:8080", "result base url")
	flag.StringVar(&ServerConfig.FileStoragePath, "f", "dbtxt.txt", "file storage path")
	flag.StringVar(&ServerConfig.DBDSN, "d", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable", "database dsn")
	flag.Parse()

	value, exists := os.LookupEnv("SERVER_ADDRESS")
	if exists {
		ServerConfig.ServerAddress = value
	}

	value, exists = os.LookupEnv("BASE_URL")
	if exists {
		ServerConfig.BaseURL = value
	}

	value, exists = os.LookupEnv("FILE_STORAGE_PATH")
	if exists {
		ServerConfig.FileStoragePath = value
	}

	value, exists = os.LookupEnv("DATABASE_DSN")
	if exists {
		ServerConfig.DBDSN = value
	}

	return &ServerConfig
}

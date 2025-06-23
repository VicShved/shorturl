// package app
package app

import (
	"flag"
	"os"
)

// type ServerConfigStruct
type ServerConfigStruct struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
	DBDSN           string
	SecretKey       string
	LogLevel        string
}

// var ServerConfig
var ServerConfig ServerConfigStruct

// func GetServerConfig
func GetServerConfig() *ServerConfigStruct {
	flag.StringVar(&ServerConfig.ServerAddress, "a", "localhost:8080", "start base url")
	flag.StringVar(&ServerConfig.BaseURL, "b", "http://localhost:8080", "result base url")
	flag.StringVar(&ServerConfig.FileStoragePath, "f", "", "file storage path")
	// host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable
	flag.StringVar(&ServerConfig.DBDSN, "d", "", "DataBase DSN")
	flag.StringVar(&ServerConfig.SecretKey, "s", "VeryImpotantSecretKey.YesYes", "Secret key")
	flag.StringVar(&ServerConfig.LogLevel, "l", "INFO", "Log level")
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

	value, exists = os.LookupEnv("SECRET_KEY")
	if exists {
		ServerConfig.SecretKey = value
	}

	value, exists = os.LookupEnv("LOG_LEVEL")
	if exists {
		ServerConfig.LogLevel = value
	}

	return &ServerConfig
}

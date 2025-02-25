package app

import "flag"

type ServerConfigStruct struct {
	ServerAddress string
	BaseURL       string
}

var ServerConfig ServerConfigStruct

func InitServerConfig() *ServerConfigStruct {
	flag.StringVar(&ServerConfig.ServerAddress, "a", "localhost:8080", "start base url")
	flag.StringVar(&ServerConfig.BaseURL, "b", "http://localhost:8080", "result base url")
	flag.Parse()
	return &ServerConfig
}

package app

import "flag"

type ServerConfigStruct struct {
	StartBaseURL  string
	ResultBaseURL string
}

var ServerConfig ServerConfigStruct

func InitServerConfig() *ServerConfigStruct {
	flag.StringVar(&ServerConfig.StartBaseURL, "a", "localhost:8080", "start base url")
	flag.StringVar(&ServerConfig.ResultBaseURL, "b", "http://localhost:8080", "result base url")
	flag.Parse()
	return &ServerConfig
}

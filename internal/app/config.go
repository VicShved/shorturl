package app

import "flag"

type ServerConfigStruct struct {
	StartBaseURL  string
	ResultBaseURL string
}

var ServerConfig ServerConfigStruct

func InitServerConfig() *ServerConfigStruct {
	flag.StringVar(&ServerConfig.StartBaseURL, "a", "localhost:0000", "start base url")
	flag.StringVar(&ServerConfig.ResultBaseURL, "b", "localhost:0000", "result base url")
	flag.Parse()
	return &ServerConfig
}

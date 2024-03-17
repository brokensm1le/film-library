package main

import (
	"film_library/config"
	"film_library/internal/httpServer"
	"log"
)

// @title           Film Library App API
// @version         1.0
// @description     This is a sample server film library server.

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apiKey  ApiKeyAuth
// @in header
// @name Authorization

func main() {
	viperInstance, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Cannot load config. Error: {%s}", err.Error())
	}

	cfg, err := config.ParseConfig(viperInstance)
	if err != nil {
		log.Fatalf("Cannot parse config. Error: {%s}", err.Error())
	}

	s := httpServer.NewServer(cfg)
	if err = s.Run(); err != nil {
		log.Print(err)
	}
}

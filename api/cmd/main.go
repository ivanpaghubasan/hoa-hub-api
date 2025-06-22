package main

import (
	"log"

	"github.com/ivanpaghubasan/hoa-hub/internal/config"
	"github.com/ivanpaghubasan/hoa-hub/internal/database"
	"github.com/ivanpaghubasan/hoa-hub/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	// initialize repo

	// initialize service

	// initialize server

	// start server
	var service interface{}
	s := server.New(service, cfg)
	if err := s.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}

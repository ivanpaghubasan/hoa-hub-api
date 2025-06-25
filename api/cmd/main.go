package main

import (
	"log"

	"github.com/ivanpaghubasan/hoa-hub/internal/config"
	"github.com/ivanpaghubasan/hoa-hub/internal/db"
	"github.com/ivanpaghubasan/hoa-hub/internal/repository"
	"github.com/ivanpaghubasan/hoa-hub/internal/server"
	"github.com/ivanpaghubasan/hoa-hub/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := db.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	// initialize repo
	repos := repository.NewRepository(db)

	// initialize service
	services := service.NewService(repos)

	// start server
	s := server.New(services, cfg)
	if err := s.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}

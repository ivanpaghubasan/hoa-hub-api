package main

import (
	"log"

	"github.com/ivanpaghubasan/hoa-hub-api/internal/auth"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/config"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/db"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/repository"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/server"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/service"
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

	jwt := auth.NewJWTAuth(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTAudience, cfg.JWTCookieDomain)

	// initialize service
	services := service.NewService(repos)

	// start server
	s := server.New(services, cfg, jwt)
	if err := s.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}

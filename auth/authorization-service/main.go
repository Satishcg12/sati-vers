package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/satishcg12/sati-vers/auth/authorization-service/api"
	"github.com/satishcg12/sati-vers/auth/authorization-service/config"
	"github.com/satishcg12/sati-vers/auth/authorization-service/db"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	database := db.NewDatabase(cfg.Database)

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	migrationsPath := "db/migrations"
	err = database.AutoMigrate(db, migrationsPath)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(cfg.HTTPServer, db)

	server.Start(ctx)
}

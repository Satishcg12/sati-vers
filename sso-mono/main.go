package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/Satishcg12/sati-vers/sso-mono/api"
	"github.com/Satishcg12/sati-vers/sso-mono/config"
	"github.com/Satishcg12/sati-vers/sso-mono/db"
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

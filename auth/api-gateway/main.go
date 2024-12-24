package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/satishcg12/sati-vers/auth/api-gateway/api"
	"github.com/satishcg12/sati-vers/auth/api-gateway/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(cfg.HTTPServer)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	server.Start(ctx)
}

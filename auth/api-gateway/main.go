package main

import (
	"context"
	"log"

	"github.com/satishcg12/sati-vers/auth/api-gateway/api"
	"github.com/satishcg12/sati-vers/auth/api-gateway/config"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(cfg.HTTPServer)
	server.Start(ctx)
}

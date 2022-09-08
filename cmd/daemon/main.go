package main

import (
	"log"

	"github.com/hetonei/arcanery-go-backend/config"
	"github.com/hetonei/arcanery-go-backend/internal/app"
)

func main() {
	cfg, err := config.LoadConfig()
	log.Println(cfg.HttpServer.Host)
	if err != nil {
		log.Fatal(err)
	}
	app.Run(&cfg)
}

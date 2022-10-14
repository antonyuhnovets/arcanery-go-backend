package main

import (
	"log"

	"github.com/hetonei/arcanery-go-backend/config"
	"github.com/hetonei/arcanery-go-backend/internal/app"
)

func main() {
	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(&cfg)
}

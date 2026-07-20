package main

import (
	"catalog-service/config"
	"catalog-service/internal/app"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	app.Run(cfg)
}

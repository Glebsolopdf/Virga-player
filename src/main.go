package main

import (
	"log"

	"virga-player/app"
)

func main() {
	if err := app.New().Run(); err != nil {
		log.Fatalf("application error: %v", err)
	}
}

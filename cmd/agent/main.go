package main

import (
	"log"

	"github.com/bbquite/mca-server/internal/app"
)

func main() {
	if err := app.RunAgent(); err != nil {
		log.Fatal(err)
	}
}

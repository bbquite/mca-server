package main

import (
	"log"

	"github.com/bbquite/mca-server/internal/app"
)

func main() {

	err := app.RunAgentAsync()
	if err != nil {
		log.Fatal(err)
	}

}

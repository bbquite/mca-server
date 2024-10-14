package main

import (
	"log"

	"github.com/bbquite/mca-server/internal/app"
)

func main() {
	if err := app.RunAgentAsync(); err != nil {
		log.Fatal(err)
	}
	//if err := app.RunAgent(); err != nil {
	//	log.Fatal(err)
	//}
}

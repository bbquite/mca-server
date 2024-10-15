package main

import (
	"log"

	"github.com/bbquite/mca-server/internal/app"
)

const asyncAgent = true

var err error

func main() {

	if asyncAgent {
		err = app.RunAgentAsync()
	} else {
		err = app.RunAgent()
	}

	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"github.com/bbquite/mca-server/internal/handlers"
	"log"
)

func main() {
	if err := handlers.AgentRun(); err != nil {
		log.Fatal(err)
	}
}

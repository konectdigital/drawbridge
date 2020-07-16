package main

import (
	"os"

	"github.com/konectdigital/drawbridge/config"
	"github.com/konectdigital/drawbridge/log"
	"github.com/konectdigital/drawbridge/server"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: drawbridge /path/to/config.yml")
	}

	c, err := config.Load(os.Args[1])
	if err != nil {
		log.Panicf("Failed to load config: %v", err)
	}

	if err := server.ListenAndServe(c); err != nil {
		log.Printf("Failed to run server: %v", err)
	}
}

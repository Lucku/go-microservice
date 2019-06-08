package main

import (
	"log"

	"github.com/lucku/otto-coding-challenge/server"
)

// Start server, nothing more
func main() {
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

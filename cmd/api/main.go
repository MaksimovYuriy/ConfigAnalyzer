package main

import (
	"log"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/api"
)

func main() {
	if err := api.Run(); err != nil {
		log.Fatal(err)
	}
}

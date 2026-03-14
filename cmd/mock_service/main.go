package main

import (
	"log"

	"github.com/nikzarch/mock-currency-service/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

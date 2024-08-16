package main

import (
	"log"

	"github.com/opoccomaxao/tgx-api/pkg/app"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

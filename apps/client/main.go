package main

import (
	"log"

	"github.com/Clayal10/enders_game/apps/client/client"
)

func main() {
	if err := client.Start(); err != nil {
		log.Fatal(err)
	}
}

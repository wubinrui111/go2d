package main

import (
	"log"
	game "github.com/wubinrui111/2d-game/internal/engine"
)

func main() {
	if err := game.Run(); err != nil {
		log.Fatal(err)
	}
}
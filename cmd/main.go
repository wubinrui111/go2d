package main

import (
	"log"
	game "github.com/yourusername/2d-game/internal/engine"
)

func main() {
	if err := game.Run(); err != nil {
		log.Fatal(err)
	}
}
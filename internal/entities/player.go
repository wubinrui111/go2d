package entities

import (
	"github.com/yourusername/2d-game/internal/components"
	"image/color"
)

type Player struct {
	*BaseBlock           // Embed BaseBlock for common properties
	components.Velocity  // Player's velocity
	components.Gravity   // Player's gravity
	OnGround bool        // Whether the player is on the ground
}

func NewPlayer(x, y float64) *Player {
	player := &Player{
		BaseBlock: NewBaseBlock(x, y, 32, 32), // Default player size
		Velocity: components.Velocity{
			X: 0,
			Y: 0,
		},
		Gravity: *components.NewGravity(),
		OnGround: false,
	}
	
	// Set player-specific properties
	player.BaseBlock.SetColor(color.RGBA{0, 0, 255, 255}) // Blue color
	player.BaseBlock.SetName("Player")
	
	return player
}
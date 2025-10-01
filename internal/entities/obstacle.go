package entities

import (
	"image/color"
)

// Obstacle represents a solid object in the game that blocks player movement
type Obstacle struct {
	*BaseBlock // Embed BaseBlock for common properties
}

// NewObstacle creates a new obstacle with the given position and dimensions
func NewObstacle(x, y, width, height float64) *Obstacle {
	obstacle := &Obstacle{
		BaseBlock: NewBaseBlock(x, y, width, height),
	}
	
	// Set obstacle-specific properties
	obstacle.BaseBlock.SetColor(color.RGBA{100, 100, 100, 255}) // Gray color
	obstacle.BaseBlock.SetName("Obstacle")
	
	return obstacle
}
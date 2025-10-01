package entities

import (
	"image/color"
)

// Ground represents the ground in the game, derived from BaseBlock
type Ground struct {
	*BaseBlock // Embed BaseBlock for common properties
}

// NewGround creates a new ground with the given position and dimensions
func NewGround(x, y, width, height float64) *Ground {
	ground := &Ground{
		BaseBlock: NewBaseBlock(x, y, width, height),
	}
	
	// Set ground-specific properties
	ground.BaseBlock.SetColor(color.RGBA{0, 180, 0, 255}) // Darker green color for ground surface
	ground.BaseBlock.SetName("Ground")
	
	return ground
}

// IsMouseOver checks if the mouse cursor is over the ground
func (g *Ground) IsMouseOver(mouseX, mouseY float64) bool {
	return g.Box.X <= mouseX && mouseX <= g.Box.X+g.Box.Width &&
		g.Box.Y <= mouseY && mouseY <= g.Box.Y+g.Box.Height
}
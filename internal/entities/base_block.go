package entities

import (
	"github.com/yourusername/2d-game/internal/components"
	"image/color"
)

const (
	// Default block size, matching player size
	DefaultBlockSize = 32.0
)

// BaseBlock represents a basic block with common properties
// All other block types should embed this struct
type BaseBlock struct {
	components.Position
	components.Box
	Color color.RGBA  // Color of the block
	Name  string      // Name of the block type
}

// NewBaseBlock creates a new base block with the given position and default size
func NewBaseBlock(x, y, width, height float64) *BaseBlock {
	return &BaseBlock{
		Position: components.Position{
			X: x,
			Y: y,
		},
		Box: components.Box{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		},
		Color: color.RGBA{128, 128, 128, 255}, // Default gray color
		Name:  "BaseBlock",
	}
}

// UpdateBoxPosition updates the box position to match the block's position
func (bb *BaseBlock) UpdateBoxPosition() {
	bb.Box.X = bb.Position.X
	bb.Box.Y = bb.Position.Y
}

// IsMouseOver checks if the mouse cursor is over the block
func (bb *BaseBlock) IsMouseOver(mouseX, mouseY float64) bool {
	return bb.Box.X <= mouseX && mouseX <= bb.Box.X+bb.Box.Width &&
		bb.Box.Y <= mouseY && mouseY <= bb.Box.Y+bb.Box.Height
}

// GetColor returns the color of the block
func (bb *BaseBlock) GetColor() color.RGBA {
	return bb.Color
}

// SetColor sets the color of the block
func (bb *BaseBlock) SetColor(c color.RGBA) {
	bb.Color = c
}

// GetName returns the name of the block type
func (bb *BaseBlock) GetName() string {
	return bb.Name
}

// SetName sets the name of the block type
func (bb *BaseBlock) SetName(name string) {
	bb.Name = name
}
package entities

import (
	"image/color"
)

// SmallBlock represents a small block with size similar to the player
type SmallBlock struct {
	*BaseBlock // Embed BaseBlock for common properties
}

// NewSmallBlock creates a new small block with the given position
// The size is set to be similar to the player (32x32)
func NewSmallBlock(x, y float64) *SmallBlock {
	block := &SmallBlock{
		BaseBlock: NewBaseBlock(x, y, DefaultBlockSize, DefaultBlockSize),
	}
	
	// Set small block-specific properties
	block.BaseBlock.SetColor(color.RGBA{200, 100, 100, 255}) // Reddish color
	block.BaseBlock.SetName("SmallBlock")
	
	return block
}

// NewSmallBlockWithColor creates a new small block with the given position and color
func NewSmallBlockWithColor(x, y float64, c color.RGBA) *SmallBlock {
	block := NewSmallBlock(x, y)
	block.BaseBlock.SetColor(c)
	return block
}

// NewRedBlock creates a new red block
func NewRedBlock(x, y float64) *SmallBlock {
	block := &SmallBlock{
		BaseBlock: NewBaseBlock(x, y, DefaultBlockSize, DefaultBlockSize),
	}
	
	block.BaseBlock.SetColor(color.RGBA{255, 100, 100, 255}) // Red color
	block.BaseBlock.SetName("RedBlock")
	
	return block
}

// NewBlueBlock creates a new blue block
func NewBlueBlock(x, y float64) *SmallBlock {
	block := &SmallBlock{
		BaseBlock: NewBaseBlock(x, y, DefaultBlockSize, DefaultBlockSize),
	}
	
	block.BaseBlock.SetColor(color.RGBA{100, 100, 255, 255}) // Blue color
	block.BaseBlock.SetName("BlueBlock")
	
	return block
}

// NewGreenBlock creates a new green block
func NewGreenBlock(x, y float64) *SmallBlock {
	block := &SmallBlock{
		BaseBlock: NewBaseBlock(x, y, DefaultBlockSize, DefaultBlockSize),
	}
	
	block.BaseBlock.SetColor(color.RGBA{100, 255, 100, 255}) // Green color
	block.BaseBlock.SetName("GreenBlock")
	
	return block
}
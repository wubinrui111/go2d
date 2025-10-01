package entities

import (
	"image/color"
	"testing"
)

func TestSmallBlockCreation(t *testing.T) {
	// Test creating a new small block
	block := NewSmallBlock(100, 200)
	
	// Check that the position is set correctly
	if block.Position.X != 100 {
		t.Errorf("Expected X position to be 100, got %f", block.Position.X)
	}
	
	if block.Position.Y != 200 {
		t.Errorf("Expected Y position to be 200, got %f", block.Position.Y)
	}
	
	// Check that the block has the correct size (same as player)
	if block.Box.Width != DefaultBlockSize {
		t.Errorf("Expected block width to be %f, got %f", DefaultBlockSize, block.Box.Width)
	}
	
	if block.Box.Height != DefaultBlockSize {
		t.Errorf("Expected block height to be %f, got %f", DefaultBlockSize, block.Box.Height)
	}
	
	// Check that the block has the correct color
	expectedColor := color.RGBA{200, 100, 100, 255} // Reddish color
	if block.GetColor() != expectedColor {
		t.Errorf("Expected block color to be %v, got %v", expectedColor, block.GetColor())
	}
	
	// Check that the block has the correct name
	if block.GetName() != "SmallBlock" {
		t.Errorf("Expected block name to be 'SmallBlock', got '%s'", block.GetName())
	}
}

func TestSmallBlockCreationWithColor(t *testing.T) {
	// Test creating a new small block with custom color
	customColor := color.RGBA{50, 150, 200, 255}
	block := NewSmallBlockWithColor(50, 75, customColor)
	
	// Check that the position is set correctly
	if block.Position.X != 50 {
		t.Errorf("Expected X position to be 50, got %f", block.Position.X)
	}
	
	if block.Position.Y != 75 {
		t.Errorf("Expected Y position to be 75, got %f", block.Position.Y)
	}
	
	// Check that the block has the correct size (same as player)
	if block.Box.Width != DefaultBlockSize {
		t.Errorf("Expected block width to be %f, got %f", DefaultBlockSize, block.Box.Width)
	}
	
	if block.Box.Height != DefaultBlockSize {
		t.Errorf("Expected block height to be %f, got %f", DefaultBlockSize, block.Box.Height)
	}
	
	// Check that the block has the correct custom color
	if block.GetColor() != customColor {
		t.Errorf("Expected block color to be %v, got %v", customColor, block.GetColor())
	}
	
	// Check that the block has the correct name
	if block.GetName() != "SmallBlock" {
		t.Errorf("Expected block name to be 'SmallBlock', got '%s'", block.GetName())
	}
}
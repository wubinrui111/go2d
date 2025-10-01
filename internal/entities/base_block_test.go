package entities

import (
	"image/color"
	"testing"
)

func TestBaseBlockCreation(t *testing.T) {
	// Test creating a new base block
	block := NewBaseBlock(10, 20, 30, 40)
	
	// Check that the position is set correctly
	if block.Position.X != 10 {
		t.Errorf("Expected X position to be 10, got %f", block.Position.X)
	}
	
	if block.Position.Y != 20 {
		t.Errorf("Expected Y position to be 20, got %f", block.Position.Y)
	}
	
	// Check that the box dimensions are set correctly
	if block.Box.X != 10 {
		t.Errorf("Expected box X to be 10, got %f", block.Box.X)
	}
	
	if block.Box.Y != 20 {
		t.Errorf("Expected box Y to be 20, got %f", block.Box.Y)
	}
	
	if block.Box.Width != 30 {
		t.Errorf("Expected box width to be 30, got %f", block.Box.Width)
	}
	
	if block.Box.Height != 40 {
		t.Errorf("Expected box height to be 40, got %f", block.Box.Height)
	}
	
	// Check default color
	expectedColor := color.RGBA{128, 128, 128, 255}
	if block.GetColor() != expectedColor {
		t.Errorf("Expected default color to be %v, got %v", expectedColor, block.GetColor())
	}
	
	// Check default name
	if block.GetName() != "BaseBlock" {
		t.Errorf("Expected default name to be 'BaseBlock', got '%s'", block.GetName())
	}
}

func TestBaseBlockColorMethods(t *testing.T) {
	block := NewBaseBlock(0, 0, 10, 10)
	
	// Test setting and getting color
	newColor := color.RGBA{255, 0, 0, 255}
	block.SetColor(newColor)
	
	if block.GetColor() != newColor {
		t.Errorf("Expected color to be %v, got %v", newColor, block.GetColor())
	}
}

func TestBaseBlockNameMethods(t *testing.T) {
	block := NewBaseBlock(0, 0, 10, 10)
	
	// Test setting and getting name
	newName := "TestBlock"
	block.SetName(newName)
	
	if block.GetName() != newName {
		t.Errorf("Expected name to be '%s', got '%s'", newName, block.GetName())
	}
}

func TestUpdateBoxPosition(t *testing.T) {
	block := NewBaseBlock(0, 0, 10, 10)
	
	// Change position
	block.Position.X = 50
	block.Position.Y = 60
	
	// Box should not be updated yet
	if block.Box.X == 50 || block.Box.Y == 60 {
		t.Error("Box position was updated before calling UpdateBoxPosition")
	}
	
	// Update box position
	block.UpdateBoxPosition()
	
	// Now box should be updated
	if block.Box.X != 50 {
		t.Errorf("Expected box X to be 50, got %f", block.Box.X)
	}
	
	if block.Box.Y != 60 {
		t.Errorf("Expected box Y to be 60, got %f", block.Box.Y)
	}
}

func TestIsMouseOver(t *testing.T) {
	block := NewBaseBlock(10, 20, 30, 40)
	
	// Test mouse position inside the block
	if !block.IsMouseOver(15, 25) {
		t.Error("Expected mouse position (15, 25) to be over the block")
	}
	
	// Test mouse position outside the block
	if block.IsMouseOver(5, 10) {
		t.Error("Expected mouse position (5, 10) to not be over the block")
	}
	
	// Test mouse position on the edge of the block (should be considered over)
	if !block.IsMouseOver(10, 20) {
		t.Error("Expected mouse position (10, 20) to be over the block (on edge)")
	}
	
	if !block.IsMouseOver(40, 60) {
		t.Error("Expected mouse position (40, 60) to be over the block (on edge)")
	}
	
	// Test mouse position just outside the block
	if block.IsMouseOver(9, 19) {
		t.Error("Expected mouse position (9, 19) to not be over the block")
	}
	
	if block.IsMouseOver(41, 61) {
		t.Error("Expected mouse position (41, 61) to not be over the block")
	}
}
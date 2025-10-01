package entities

import (
	"image/color"
	"testing"
)

func TestGroundCreation(t *testing.T) {
	// Test creating a new ground (surface only)
	ground := NewGround(0, 550, 800, 32)
	
	// Check that the position is set correctly
	if ground.Position.X != 0 {
		t.Errorf("Expected X position to be 0, got %f", ground.Position.X)
	}
	
	if ground.Position.Y != 550 {
		t.Errorf("Expected Y position to be 550, got %f", ground.Position.Y)
	}
	
	// Check that the ground has the correct size (surface only)
	if ground.Box.Width != 800 {
		t.Errorf("Expected ground width to be 800, got %f", ground.Box.Width)
	}
	
	if ground.Box.Height != 32 {
		t.Errorf("Expected ground height to be 32, got %f", ground.Box.Height)
	}
	
	// Check that the ground has the correct color
	expectedColor := color.RGBA{0, 180, 0, 255} // Darker green color
	if ground.GetColor() != expectedColor {
		t.Errorf("Expected ground color to be %v, got %v", expectedColor, ground.GetColor())
	}
	
	// Check that the ground has the correct name
	if ground.GetName() != "Ground" {
		t.Errorf("Expected ground name to be 'Ground', got '%s'", ground.GetName())
	}
}

func TestGroundMouseOver(t *testing.T) {
	ground := NewGround(0, 550, 800, 32)
	
	// Test mouse position inside the ground
	if !ground.IsMouseOver(100, 560) {
		t.Error("Expected mouse position (100, 560) to be over the ground")
	}
	
	// Test mouse position outside the ground
	if ground.IsMouseOver(100, 500) {
		t.Error("Expected mouse position (100, 500) to not be over the ground")
	}
	
	// Test mouse position on the edge of the ground (should be considered over)
	if !ground.IsMouseOver(0, 550) {
		t.Error("Expected mouse position (0, 550) to be over the ground (on edge)")
	}
	
	if !ground.IsMouseOver(800, 582) {
		t.Error("Expected mouse position (800, 582) to be over the ground (on edge)")
	}
	
	// Test mouse position just outside the ground
	if ground.IsMouseOver(-1, 549) {
		t.Error("Expected mouse position (-1, 549) to not be over the ground")
	}
	
	if ground.IsMouseOver(801, 583) {
		t.Error("Expected mouse position (801, 583) to not be over the ground")
	}
}
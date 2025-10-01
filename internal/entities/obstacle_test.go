package entities

import (
	"image/color"
	"testing"
)

func TestObstacleCreation(t *testing.T) {
	// Test creating a new obstacle
	obstacle := NewObstacle(50, 75, 25, 30)
	
	// Check that the position is set correctly
	if obstacle.Position.X != 50 {
		t.Errorf("Expected X position to be 50, got %f", obstacle.Position.X)
	}
	
	if obstacle.Position.Y != 75 {
		t.Errorf("Expected Y position to be 75, got %f", obstacle.Position.Y)
	}
	
	// Check that the obstacle has the correct size
	if obstacle.Box.Width != 25 {
		t.Errorf("Expected obstacle width to be 25, got %f", obstacle.Box.Width)
	}
	
	if obstacle.Box.Height != 30 {
		t.Errorf("Expected obstacle height to be 30, got %f", obstacle.Box.Height)
	}
	
	// Check that the obstacle has the correct color
	expectedColor := color.RGBA{100, 100, 100, 255} // Gray
	if obstacle.GetColor() != expectedColor {
		t.Errorf("Expected obstacle color to be %v, got %v", expectedColor, obstacle.GetColor())
	}
	
	// Check that the obstacle has the correct name
	if obstacle.GetName() != "Obstacle" {
		t.Errorf("Expected obstacle name to be 'Obstacle', got '%s'", obstacle.GetName())
	}
}
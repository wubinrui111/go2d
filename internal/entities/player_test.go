package entities

import (
	"image/color"
	"testing"
)

func TestPlayerCreation(t *testing.T) {
	// Test creating a new player
	player := NewPlayer(100, 200)
	
	// Check that the position is set correctly
	if player.Position.X != 100 {
		t.Errorf("Expected X position to be 100, got %f", player.Position.X)
	}
	
	if player.Position.Y != 200 {
		t.Errorf("Expected Y position to be 200, got %f", player.Position.Y)
	}
	
	// Check that the player has the correct size
	if player.Box.Width != 32 {
		t.Errorf("Expected player width to be 32, got %f", player.Box.Width)
	}
	
	if player.Box.Height != 32 {
		t.Errorf("Expected player height to be 32, got %f", player.Box.Height)
	}
	
	// Check that the player has the correct color
	expectedColor := color.RGBA{0, 0, 255, 255} // Blue
	if player.GetColor() != expectedColor {
		t.Errorf("Expected player color to be %v, got %v", expectedColor, player.GetColor())
	}
	
	// Check that the player has the correct name
	if player.GetName() != "Player" {
		t.Errorf("Expected player name to be 'Player', got '%s'", player.GetName())
	}
	
	// Check that velocity is initialized
	if player.Velocity.X != 0 {
		t.Errorf("Expected initial X velocity to be 0, got %f", player.Velocity.X)
	}
	
	if player.Velocity.Y != 0 {
		t.Errorf("Expected initial Y velocity to be 0, got %f", player.Velocity.Y)
	}
	
	// Check that gravity is enabled
	if !player.Gravity.Enabled {
		t.Error("Expected gravity to be enabled by default")
	}
	
	// Check that player is not on ground initially
	if player.OnGround {
		t.Error("Expected player to not be on ground initially")
	}
}
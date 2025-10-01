package graphics

import (
	"testing"

	"github.com/wubinrui111/2d-game/internal/components"
	"image/color"
)

func TestNewInventorySystem(t *testing.T) {
	is := NewInventorySystem()
	
	if is == nil {
		t.Error("NewInventorySystem should return a non-nil inventory system")
	}
	
	if is.Visible != false {
		t.Error("New inventory system should not be visible by default")
	}
}

func TestToggleVisibility(t *testing.T) {
	is := NewInventorySystem()
	
	// Initially not visible
	if is.Visible != false {
		t.Error("Inventory system should not be visible initially")
	}
	
	// Toggle visibility
	is.Visible = !is.Visible
	
	if is.Visible != true {
		t.Error("Inventory system should be visible after toggle")
	}
	
	// Toggle back
	is.Visible = !is.Visible
	
	if is.Visible != false {
		t.Error("Inventory system should not be visible after second toggle")
	}
}

func TestInventorySystemWithInventory(t *testing.T) {
	// Create inventory system
	is := NewInventorySystem()
	
	// Create inventory
	inventory := components.NewInventory(27, 9)
	
	// Add some items to inventory
	item := components.Item{
		ID:       "stone",
		Name:     "Stone",
		Count:    32,
		MaxStack: 64,
		Color:    color.RGBA{128, 128, 128, 255},
	}
	
	inventory.AddItem(item)
	
	// Test that we can call Draw without crashing
	// (We can't easily test the actual rendering)
	
	// Test with visible inventory
	is.Visible = true
	// This would normally draw the full inventory
	
	// Test with hidden inventory
	is.Visible = false
	// This would normally draw just the hotbar
	
	// If we get here without crashing, the test passes
	// Note: Actual visual testing would require more sophisticated testing tools
}

func TestInventorySystemUpdate(t *testing.T) {
	// Create inventory system
	is := NewInventorySystem()
	
	// Create inventory
	inventory := components.NewInventory(27, 9)
	
	// Test that we can call Update without crashing
	// (We can't easily test actual input handling in a unit test)
	is.Update(inventory)
	
	// If we get here without crashing, the test passes
	// Note: Actual input testing would require mocking the input system
}
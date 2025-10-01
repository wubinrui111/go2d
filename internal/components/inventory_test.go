package components

import (
	"image/color"
	"testing"
)

func TestNewInventory(t *testing.T) {
	inv := NewInventory(27, 9)
	
	if len(inv.Slots) != 27 {
		t.Errorf("Expected 27 slots, got %d", len(inv.Slots))
	}
	
	if inv.SelectedSlot != 0 {
		t.Errorf("Expected selected slot to be 0, got %d", inv.SelectedSlot)
	}
	
	if inv.HotbarSize != 9 {
		t.Errorf("Expected hotbar size to be 9, got %d", inv.HotbarSize)
	}
	
	// Check that all slots are initialized as empty
	for i, slot := range inv.Slots {
		if slot.Item != nil {
			t.Errorf("Expected slot %d to be empty, but found item %v", i, slot.Item)
		}
		if slot.Count != 0 {
			t.Errorf("Expected slot %d count to be 0, but got %d", i, slot.Count)
		}
	}
}

func TestAddItem(t *testing.T) {
	inv := NewInventory(9, 9)
	
	item := Item{
		ID:       "stone",
		Name:     "Stone",
		Count:    5,
		MaxStack: 64,
		Color:    color.RGBA{128, 128, 128, 255},
	}
	
	// Add items to inventory
	success := inv.AddItem(item)
	if !success {
		t.Error("Failed to add item to inventory")
	}
	
	// Check that the item was added correctly
	slot := inv.Slots[0]
	if slot.Item == nil {
		t.Error("Expected item to be added to slot 0")
	} else {
		if slot.Item.ID != "stone" {
			t.Errorf("Expected item ID to be 'stone', got '%s'", slot.Item.ID)
		}
		if slot.Count != 5 {
			t.Errorf("Expected count to be 5, got %d", slot.Count)
		}
	}
}

func TestAddItemStacking(t *testing.T) {
	inv := NewInventory(9, 9)
	
	item1 := Item{
		ID:       "stone",
		Name:     "Stone",
		Count:    32,
		MaxStack: 64,
		Color:    color.RGBA{128, 128, 128, 255},
	}
	
	item2 := Item{
		ID:       "stone",
		Name:     "Stone",
		Count:    16,
		MaxStack: 64,
		Color:    color.RGBA{128, 128, 128, 255},
	}
	
	// Add first stack
	inv.AddItem(item1)
	
	// Add second stack - should stack with the first
	inv.AddItem(item2)
	
	// Check that items were stacked
	slot := inv.Slots[0]
	if slot.Item == nil {
		t.Error("Expected item in slot 0")
	} else {
		if slot.Count != 48 {
			t.Errorf("Expected count to be 48 (32+16), got %d", slot.Count)
		}
	}
	
	// Check that no other slots were used
	for i := 1; i < len(inv.Slots); i++ {
		if inv.Slots[i].Count != 0 {
			t.Errorf("Expected slot %d to be empty, but found %d items", i, inv.Slots[i].Count)
		}
	}
}

func TestAddItemMultipleSlots(t *testing.T) {
	inv := NewInventory(9, 9)
	
	item := Item{
		ID:       "stone",
		Name:     "Stone",
		Count:    100, // More than one stack
		MaxStack: 64,
		Color:    color.RGBA{128, 128, 128, 255},
	}
	
	inv.AddItem(item)
	
	// Check first slot (should be full)
	if inv.Slots[0].Count != 64 {
		t.Errorf("Expected first slot to have 64 items, got %d", inv.Slots[0].Count)
	}
	
	// Check second slot (should have remaining items)
	if inv.Slots[1].Count != 36 {
		t.Errorf("Expected second slot to have 36 items, got %d", inv.Slots[1].Count)
	}
	
	// Check that other slots are empty
	for i := 2; i < len(inv.Slots); i++ {
		if inv.Slots[i].Count != 0 {
			t.Errorf("Expected slot %d to be empty, but found %d items", i, inv.Slots[i].Count)
		}
	}
}

func TestRemoveItem(t *testing.T) {
	inv := NewInventory(9, 9)
	
	// Add items to inventory
	item := Item{
		ID:       "stone",
		Name:     "Stone",
		Count:    32,
		MaxStack: 64,
		Color:    color.RGBA{128, 128, 128, 255},
	}
	inv.AddItem(item)
	
	// Remove some items
	success := inv.RemoveItem("stone", 10)
	if !success {
		t.Error("Failed to remove items")
	}
	
	// Check remaining count
	slot := inv.Slots[0]
	if slot.Count != 22 {
		t.Errorf("Expected 22 items remaining, got %d", slot.Count)
	}
}

func TestRemoveItemMultipleSlots(t *testing.T) {
	inv := NewInventory(9, 9)
	
	// Add items that span multiple slots
	item := Item{
		ID:       "stone",
		Name:     "Stone",
		Count:    100,
		MaxStack: 64,
		Color:    color.RGBA{128, 128, 128, 255},
	}
	inv.AddItem(item)
	
	// Remove more items than in the first slot
	success := inv.RemoveItem("stone", 80)
	if !success {
		t.Error("Failed to remove items from multiple slots")
	}
	
	// First slot should be empty
	if inv.Slots[0].Count != 0 {
		t.Errorf("Expected first slot to be empty, got %d", inv.Slots[0].Count)
	}
	
	// Second slot should have 20 items left (100 - 80)
	if inv.Slots[1].Count != 20 {
		t.Errorf("Expected second slot to have 20 items, got %d", inv.Slots[1].Count)
	}
}

func TestGetItemCount(t *testing.T) {
	inv := NewInventory(9, 9)
	
	// Add items to multiple slots
	item := Item{
		ID:       "stone",
		Name:     "Stone",
		Count:    100,
		MaxStack: 64,
		Color:    color.RGBA{128, 128, 128, 255},
	}
	inv.AddItem(item)
	
	// Check total count
	count := inv.GetItemCount("stone")
	if count != 100 {
		t.Errorf("Expected total count to be 100, got %d", count)
	}
}

func TestGetSelectedItem(t *testing.T) {
	inv := NewInventory(9, 9)
	
	// Add an item
	item := Item{
		ID:       "stone",
		Name:     "Stone",
		Count:    32,
		MaxStack: 64,
		Color:    color.RGBA{128, 128, 128, 255},
	}
	inv.AddItem(item)
	
	// Select the first slot
	inv.SelectSlot(0)
	
	// Get selected item
	selected := inv.GetSelectedItem()
	if selected == nil {
		t.Error("Expected selected item, got nil")
	} else {
		if selected.Item.ID != "stone" {
			t.Errorf("Expected selected item ID to be 'stone', got '%s'", selected.Item.ID)
		}
		if selected.Count != 32 {
			t.Errorf("Expected selected item count to be 32, got %d", selected.Count)
		}
	}
}

func TestSlotSelection(t *testing.T) {
	inv := NewInventory(27, 9) // 27 total slots, 9 in hotbar
	
	// Test selecting next slot
	inv.SelectNextSlot()
	if inv.SelectedSlot != 1 {
		t.Errorf("Expected selected slot to be 1, got %d", inv.SelectedSlot)
	}
	
	// Test wrapping around
	inv.SelectedSlot = 8
	inv.SelectNextSlot()
	if inv.SelectedSlot != 0 {
		t.Errorf("Expected selected slot to wrap to 0, got %d", inv.SelectedSlot)
	}
	
	// Test selecting previous slot
	inv.SelectedSlot = 1
	inv.SelectPreviousSlot()
	if inv.SelectedSlot != 0 {
		t.Errorf("Expected selected slot to be 0, got %d", inv.SelectedSlot)
	}
	
	// Test wrapping around backwards
	inv.SelectedSlot = 0
	inv.SelectPreviousSlot()
	if inv.SelectedSlot != 8 {
		t.Errorf("Expected selected slot to wrap to 8, got %d", inv.SelectedSlot)
	}
	
	// Test direct slot selection
	inv.SelectSlot(5)
	if inv.SelectedSlot != 5 {
		t.Errorf("Expected selected slot to be 5, got %d", inv.SelectedSlot)
	}
	
	// Test invalid slot selection (should not change)
	inv.SelectSlot(10) // Invalid slot (beyond hotbar)
	if inv.SelectedSlot != 5 {
		t.Errorf("Expected selected slot to remain 5, got %d", inv.SelectedSlot)
	}
}

func TestInventoryState(t *testing.T) {
	inv := NewInventory(9, 9)
	
	// Test empty inventory
	if !inv.IsEmpty() {
		t.Error("Expected new inventory to be empty")
	}
	
	if inv.IsFull() {
		t.Error("Expected new inventory to not be full")
	}
	
	// Add items
	item := Item{
		ID:       "stone",
		Name:     "Stone",
		Count:    32,
		MaxStack: 64,
		Color:    color.RGBA{128, 128, 128, 255},
	}
	inv.AddItem(item)
	
	// Test non-empty inventory
	if inv.IsEmpty() {
		t.Error("Expected inventory with items to not be empty")
	}
	
	if inv.IsFull() {
		t.Error("Expected inventory with one item to not be full")
	}
	
	// Fill inventory completely
	fullStack := Item{
		ID:       "dirt",
		Name:     "Dirt",
		Count:    64,
		MaxStack: 64,
		Color:    color.RGBA{100, 50, 0, 255},
	}
	
	for i := 0; i < 9; i++ {
		stack := fullStack
		inv.AddItem(stack)
	}
	
	// Test full inventory
	if inv.IsEmpty() {
		t.Error("Expected full inventory to not be empty")
	}
	
	// Note: This might not be completely full if the first item is still there
	// depending on inventory implementation details
}
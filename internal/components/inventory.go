package components

// Inventory represents a player's inventory for storing items
type Inventory struct {
	// Slots is a slice of inventory slots
	Slots []ItemStack

	// SelectedSlot is the index of the currently selected slot
	SelectedSlot int

	// HotbarSize is the number of slots in the hotbar
	HotbarSize int
}

// NewInventory creates a new inventory with the specified number of slots
func NewInventory(slotCount int, hotbarSize int) *Inventory {
	slots := make([]ItemStack, slotCount)
	
	// Initialize all slots as empty
	for i := range slots {
		slots[i] = ItemStack{
			Item:  nil,
			Count: 0,
		}
	}

	return &Inventory{
		Slots:        slots,
		SelectedSlot: 0,
		HotbarSize:   hotbarSize,
	}
}

// AddItem adds items to the inventory, stacking as appropriate
func (inv *Inventory) AddItem(item Item) bool {
	// Try to stack with existing items of the same type
	for i := range inv.Slots {
		slot := &inv.Slots[i]
		if slot.Item != nil && slot.Item.ID == item.ID && slot.Count < slot.Item.MaxStack {
			// Calculate how many items we can add to this slot
			canAdd := slot.Item.MaxStack - slot.Count
			willAdd := item.Count
			if willAdd > canAdd {
				willAdd = canAdd
			}

			// Add items to this slot
			slot.Count += willAdd
			item.Count -= willAdd

			// If we've added all items, we're done
			if item.Count <= 0 {
				return true
			}
		}
	}

	// If there are still items left, try to put them in an empty slot
	if item.Count > 0 {
		for i := range inv.Slots {
			slot := &inv.Slots[i]
			if slot.Item == nil {
				// Put as many items as possible in this slot
				maxStack := item.MaxStack
				if item.Count > maxStack {
					slot.Item = &item
					slot.Count = maxStack
					item.Count -= maxStack
				} else {
					slot.Item = &item
					slot.Count = item.Count
					item.Count = 0
				}
				if item.Count <= 0 {
					return true
				}
			}
		}
	}

	// If we get here, we couldn't fit all items
	return item.Count <= 0
}

// RemoveItem removes a specific number of items from the inventory
func (inv *Inventory) RemoveItem(itemID string, count int) bool {
	removed := 0
	// Find slots with matching items and remove them
	for i := range inv.Slots {
		slot := &inv.Slots[i]
		if slot.Item != nil && slot.Item.ID == itemID {
			// Calculate how many we can remove from this slot
			toRemove := count - removed
			if toRemove > slot.Count {
				toRemove = slot.Count
			}

			// Remove items from this slot
			slot.Count -= toRemove
			removed += toRemove

			// If the slot is empty, clear it
			if slot.Count == 0 {
				slot.Item = nil
			}

			// If we've removed enough items, we're done
			if removed >= count {
				return true
			}
		}
	}
	return removed >= count
}

// GetItemCount returns the total count of a specific item in the inventory
func (inv *Inventory) GetItemCount(itemID string) int {
	count := 0
	for _, slot := range inv.Slots {
		if slot.Item != nil && slot.Item.ID == itemID {
			count += slot.Count
		}
	}
	return count
}

// GetSelectedItem returns the currently selected item stack
func (inv *Inventory) GetSelectedItem() *ItemStack {
	if inv.SelectedSlot >= 0 && inv.SelectedSlot < len(inv.Slots) {
		slot := &inv.Slots[inv.SelectedSlot]
		if slot.Count > 0 {
			return slot
		}
	}
	return nil
}

// SelectNextSlot moves the selection to the next slot in the hotbar
func (inv *Inventory) SelectNextSlot() {
	inv.SelectedSlot = (inv.SelectedSlot + 1) % inv.HotbarSize
}

// SelectPreviousSlot moves the selection to the previous slot in the hotbar
func (inv *Inventory) SelectPreviousSlot() {
	inv.SelectedSlot--
	if inv.SelectedSlot < 0 {
		inv.SelectedSlot = inv.HotbarSize - 1
	}
}

// SelectSlot selects a specific slot in the hotbar
func (inv *Inventory) SelectSlot(slotIndex int) {
	if slotIndex >= 0 && slotIndex < inv.HotbarSize {
		inv.SelectedSlot = slotIndex
	}
}

// IsFull checks if the inventory is completely full
func (inv *Inventory) IsFull() bool {
	for _, slot := range inv.Slots {
		if slot.Item == nil || slot.Count < slot.Item.MaxStack {
			return false
		}
	}
	return true
}

// IsEmpty checks if the inventory is completely empty
func (inv *Inventory) IsEmpty() bool {
	for _, slot := range inv.Slots {
		if slot.Count > 0 {
			return false
		}
	}
	return true
}
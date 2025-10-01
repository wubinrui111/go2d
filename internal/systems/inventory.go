package graphics

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/wubinrui111/2d-game/internal/components"
)

const (
	// Slot size and spacing
	SlotSize   = 32
	SlotMargin = 4

	// Inventory position
	InventoryX = 10
	InventoryY = 10

	// Hotbar position (bottom of screen)
	HotbarX      = 10
	HotbarY      = 550
	HotbarWidth  = 9 * (SlotSize + SlotMargin)
	HotbarHeight = SlotSize
)

// InventorySystem handles the rendering and interaction with the player's inventory
type InventorySystem struct {
	// Visible indicates whether the full inventory is visible (not just the hotbar)
	Visible bool
	
	// MouseAttachedSlot indicates which slot item is attached to mouse (-1 if none)
	MouseAttachedSlot int
	
	// MouseAttachedItem stores the item attached to the mouse
	MouseAttachedItem *components.ItemStack
	
	// GameMode indicates the current game mode (0 = survival, 1 = creative)
	GameMode int
	
	// BlockSprites stores references to block sprites for rendering
	BlockSprites map[string]*ebiten.Image
}

// NewInventorySystem creates a new inventory system
func NewInventorySystem() *InventorySystem {
	return &InventorySystem{
		Visible:           false,
		MouseAttachedSlot: -1, // -1 indicates no item is attached to mouse
		MouseAttachedItem: nil,
		GameMode:          0, // 0 = survival mode by default
		BlockSprites:      nil, // Initially no sprites
	}
}

// SetBlockSprites sets the block sprites for the inventory system
func (is *InventorySystem) SetBlockSprites(sprites map[string]*ebiten.Image) {
	is.BlockSprites = sprites
}

// Update handles input for the inventory system
func (is *InventorySystem) Update(inventory *components.Inventory) {
	// Toggle full inventory visibility with 'E' key
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		is.Visible = !is.Visible
	}
	
	// Toggle game mode with 'G' key (creative/survival)
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		is.GameMode = 1 - is.GameMode // Toggle between 0 and 1
	}

	// Handle hotbar slot selection with number keys
	for i := 1; i <= inventory.HotbarSize; i++ {
		key := ebiten.Key(int(ebiten.Key1) + i - 1)
		if inpututil.IsKeyJustPressed(key) {
			inventory.SelectSlot(i - 1)
			break
		}
	}

	// Handle mouse wheel for slot selection (reversed direction)
	_, wheelY := ebiten.Wheel()
	if wheelY > 0 {
		inventory.SelectPreviousSlot()
	} else if wheelY < 0 {
		inventory.SelectNextSlot()
	}
	
	// Handle mouse attachment for inventory slots
	is.handleMouseAttachment(inventory)
}

// Draw renders the inventory based on its visibility state
func (is *InventorySystem) Draw(screen *ebiten.Image, inventory *components.Inventory) {
	if is.Visible {
		// Draw full inventory
		is.drawFullInventory(screen, inventory)
	} else {
		// Draw hotbar only
		is.drawHotbar(screen, inventory)
	}
}

// drawHotbar renders only the hotbar (first N slots) at the bottom of the screen
func (is *InventorySystem) drawHotbar(screen *ebiten.Image, inventory *components.Inventory) {
	// Draw slots first
	for i := 0; i < inventory.HotbarSize && i < len(inventory.Slots); i++ {
		// Skip drawing the slot that has an attached item
		if is.MouseAttachedSlot == i {
			continue
		}
		
		slot := inventory.Slots[i]
		x := float64(HotbarX + i*(SlotSize+SlotMargin))
		y := float64(HotbarY)

		// Draw slot background
		slotColor := color.RGBA{100, 100, 100, 200}
		if i == inventory.SelectedSlot {
			slotColor = color.RGBA{150, 150, 150, 255} // Highlight selected slot
		}
		ebitenutil.DrawRect(screen, x, y, SlotSize, SlotSize, slotColor)

		// Draw item if present
		if slot.Item != nil && slot.Count > 0 {
			// Draw item with sprite if available, otherwise use color
			if is.BlockSprites != nil {
				// Try to get sprite for this item
				var itemSprite *ebiten.Image
				
				// Map item ID to sprite
				switch slot.Item.ID {
				case "stone":
					itemSprite = is.BlockSprites["stone"]
				case "dirt":
					itemSprite = is.BlockSprites["dirt"]
				case "wood":
					itemSprite = is.BlockSprites["wood"]
				case "small_block":
					itemSprite = is.BlockSprites["small_block"]
				case "red_block":
					itemSprite = is.BlockSprites["RedBlock"]
				case "blue_block":
					itemSprite = is.BlockSprites["BlueBlock"]
				case "green_block":
					itemSprite = is.BlockSprites["GreenBlock"]
				default:
					// Try to find sprite by item ID directly
					if sprite, exists := is.BlockSprites[slot.Item.ID]; exists {
						itemSprite = sprite
					}
				}
				
				if itemSprite != nil {
					// Draw sprite
					opts := &ebiten.DrawImageOptions{}
					opts.GeoM.Translate(x+2, y+2)
					screen.DrawImage(itemSprite, opts)
				} else {
					// Fallback to color
					ebitenutil.DrawRect(screen, x+2, y+2, SlotSize-4, SlotSize-4, slot.Item.Color)
				}
			} else {
				// Draw item color (as a placeholder for actual sprite)
				ebitenutil.DrawRect(screen, x+2, y+2, SlotSize-4, SlotSize-4, slot.Item.Color)
			}

			// Draw item count
			countText := fmt.Sprintf("%d", slot.Count)
			textWidth := len(countText) * 6 // Approximate character width
			ebitenutil.DebugPrintAt(screen, countText, int(x)+SlotSize-textWidth-2, int(y)+SlotSize-12)
		}

		// Draw slot number
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", (i+1)%10), int(x+2), int(y+2))
	}
	
	// Draw attached item if there is one
	if is.MouseAttachedItem != nil {
		// Get mouse position
		mouseX, mouseY := ebiten.CursorPosition()
		
		// Draw item at mouse position
		x := float64(mouseX)
		y := float64(mouseY)
		
		// Draw semi-transparent background
		ebitenutil.DrawRect(screen, x, y, SlotSize, SlotSize, color.RGBA{100, 100, 100, 150})
		
		// Draw item with sprite if available, otherwise use color
		if is.BlockSprites != nil && is.MouseAttachedItem.Item != nil {
			// Try to get sprite for this item
			var itemSprite *ebiten.Image
			
			// Map item ID to sprite
			switch is.MouseAttachedItem.Item.ID {
			case "stone":
				itemSprite = is.BlockSprites["stone"]
			case "dirt":
				itemSprite = is.BlockSprites["dirt"]
			case "wood":
				itemSprite = is.BlockSprites["wood"]
			case "small_block":
				itemSprite = is.BlockSprites["small_block"]
			case "red_block":
				itemSprite = is.BlockSprites["RedBlock"]
			case "blue_block":
				itemSprite = is.BlockSprites["BlueBlock"]
			case "green_block":
				itemSprite = is.BlockSprites["GreenBlock"]
			default:
				// Try to find sprite by item ID directly
				if sprite, exists := is.BlockSprites[is.MouseAttachedItem.Item.ID]; exists {
					itemSprite = sprite
				}
			}
			
			if itemSprite != nil {
				// Draw sprite
				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Translate(x+2, y+2)
				screen.DrawImage(itemSprite, opts)
			} else {
				// Fallback to color
				ebitenutil.DrawRect(screen, x+2, y+2, SlotSize-4, SlotSize-4, is.MouseAttachedItem.Item.Color)
			}
		} else if is.MouseAttachedItem.Item != nil {
			// Draw item color
			ebitenutil.DrawRect(screen, x+2, y+2, SlotSize-4, SlotSize-4, is.MouseAttachedItem.Item.Color)
		}
		
		// Draw item count
		countText := fmt.Sprintf("%d", is.MouseAttachedItem.Count)
		textWidth := len(countText) * 6 // Approximate character width
		ebitenutil.DebugPrintAt(screen, countText, int(x)+SlotSize-textWidth-2, int(y)+SlotSize-12)
	}
}

// drawFullInventory renders the full inventory grid
func (is *InventorySystem) drawFullInventory(screen *ebiten.Image, inventory *components.Inventory) {
	// Draw semi-transparent background
	ebitenutil.DrawRect(screen, 0, 0, 800, 600, color.RGBA{0, 0, 0, 150})

	// Draw title
	if is.GameMode == 0 {
		ebitenutil.DebugPrintAt(screen, "Inventory (Survival Mode)", 300, 20)
	} else {
		ebitenutil.DebugPrintAt(screen, "Inventory (Creative Mode)", 300, 20)
	}

	// Define creative items
	creativeItems := []components.Item{
		{
			ID:       "stone",
			Name:     "Stone",
			Count:    1,
			MaxStack: 64,
			Color:    color.RGBA{128, 128, 128, 255},
		},
		{
			ID:       "dirt",
			Name:     "Dirt",
			Count:    1,
			MaxStack: 64,
			Color:    color.RGBA{100, 50, 0, 255},
		},
		{
			ID:       "wood",
			Name:     "Wood",
			Count:    1,
			MaxStack: 64,
			Color:    color.RGBA{100, 70, 30, 255},
		},
		{
			ID:       "small_block",
			Name:     "Small Block",
			Count:    1,
			MaxStack: 64,
			Color:    color.RGBA{200, 200, 50, 255},
		},
		{
			ID:       "red_block",
			Name:     "Red Block",
			Count:    1,
			MaxStack: 64,
			Color:    color.RGBA{200, 50, 50, 255},
		},
		{
			ID:       "blue_block",
			Name:     "Blue Block",
			Count:    1,
			MaxStack: 64,
			Color:    color.RGBA{50, 50, 200, 255},
		},
		{
			ID:       "green_block",
			Name:     "Green Block",
			Count:    1,
			MaxStack: 64,
			Color:    color.RGBA{50, 200, 50, 255},
		},
	}

	// Calculate total slots needed (player slots + creative items in creative mode)
	totalSlots := len(inventory.Slots)
	if is.GameMode == 1 {
		// Only include creative items in creative mode
		totalSlots += len(creativeItems)
	}
	
	cols := 9
	
	// Draw all slots including creative items (if in creative mode)
	for i := 0; i < totalSlots; i++ {
		var slot *components.ItemStack
		var isCreativeItem bool
		var creativeItemIndex int
		
		// Determine if this is a player inventory slot or a creative item slot
		if i < len(inventory.Slots) {
			// Player inventory slot
			slot = &inventory.Slots[i]
		} else {
			// Creative item slot (only in creative mode)
			isCreativeItem = true
			creativeItemIndex = i - len(inventory.Slots)
			slot = &components.ItemStack{
				Item:  &creativeItems[creativeItemIndex],
				Count: creativeItems[creativeItemIndex].MaxStack,
			}
		}
		
		// Skip drawing the slot that has an attached item (only for player slots)
		if !isCreativeItem && is.MouseAttachedSlot == i {
			continue
		}
		
		row := i / cols
		col := i % cols

		x := float64(InventoryX + col*(SlotSize+SlotMargin))
		y := float64(60 + row*(SlotSize+SlotMargin))

		// Draw slot background
		slotColor := color.RGBA{100, 100, 100, 200}
		if !isCreativeItem && i == inventory.SelectedSlot {
			slotColor = color.RGBA{150, 150, 150, 255} // Highlight selected slot
		} else if isCreativeItem {
			slotColor = color.RGBA{80, 80, 120, 200} // Different color for creative items
		}
		ebitenutil.DrawRect(screen, x, y, SlotSize, SlotSize, slotColor)

		// Draw item if present
		if slot.Item != nil && slot.Count > 0 {
			// Draw item with sprite if available, otherwise use color
			if is.BlockSprites != nil {
				// Try to get sprite for this item
				var itemSprite *ebiten.Image
				
				// Map item ID to sprite
				switch slot.Item.ID {
				case "stone":
					itemSprite = is.BlockSprites["stone"]
				case "dirt":
					itemSprite = is.BlockSprites["dirt"]
				case "wood":
					itemSprite = is.BlockSprites["wood"]
				case "small_block":
					itemSprite = is.BlockSprites["small_block"]
				case "red_block":
					itemSprite = is.BlockSprites["RedBlock"]
				case "blue_block":
					itemSprite = is.BlockSprites["BlueBlock"]
				case "green_block":
					itemSprite = is.BlockSprites["GreenBlock"]
				default:
					// Try to find sprite by item ID directly
					if sprite, exists := is.BlockSprites[slot.Item.ID]; exists {
						itemSprite = sprite
					}
				}
				
				if itemSprite != nil {
					// Draw sprite
					opts := &ebiten.DrawImageOptions{}
					opts.GeoM.Translate(x+2, y+2)
					screen.DrawImage(itemSprite, opts)
				} else {
					// Fallback to color
					ebitenutil.DrawRect(screen, x+2, y+2, SlotSize-4, SlotSize-4, slot.Item.Color)
				}
			} else {
				// Draw item color (as a placeholder for actual sprite)
				ebitenutil.DrawRect(screen, x+2, y+2, SlotSize-4, SlotSize-4, slot.Item.Color)
			}

			// Draw item count
			countText := fmt.Sprintf("%d", slot.Count)
			textWidth := len(countText) * 6 // Approximate character width
			ebitenutil.DebugPrintAt(screen, countText, int(x)+SlotSize-textWidth-2, int(y)+SlotSize-12)
			
			// For creative items, also draw a label
			if isCreativeItem {
				name := slot.Item.Name
				if len(name) > 10 {
					name = name[:10] + "..."
				}
				ebitenutil.DebugPrintAt(screen, name, int(x), int(y+SlotSize+2))
			}
		}
	}

	// Draw instructions
	ebitenutil.DebugPrintAt(screen, "Press 'E' to close inventory", 10, 570)
	if is.GameMode == 0 {
		ebitenutil.DebugPrintAt(screen, "Press 'G' to switch to Creative Mode", 10, 550)
	} else {
		ebitenutil.DebugPrintAt(screen, "Press 'G' to switch to Survival Mode", 10, 550)
	}
	
	// Draw attached item if there is one
	if is.MouseAttachedItem != nil {
		// Get mouse position
		mouseX, mouseY := ebiten.CursorPosition()
		
		// Draw item at mouse position
		x := float64(mouseX)
		y := float64(mouseY)
		
		// Draw semi-transparent background
		ebitenutil.DrawRect(screen, x, y, SlotSize, SlotSize, color.RGBA{100, 100, 100, 150})
		
		// Draw item with sprite if available, otherwise use color
		if is.BlockSprites != nil && is.MouseAttachedItem.Item != nil {
			// Try to get sprite for this item
			var itemSprite *ebiten.Image
			
			// Map item ID to sprite
			switch is.MouseAttachedItem.Item.ID {
			case "stone":
				itemSprite = is.BlockSprites["stone"]
			case "dirt":
				itemSprite = is.BlockSprites["dirt"]
			case "wood":
				itemSprite = is.BlockSprites["wood"]
			case "small_block":
				itemSprite = is.BlockSprites["small_block"]
			case "red_block":
				itemSprite = is.BlockSprites["RedBlock"]
			case "blue_block":
				itemSprite = is.BlockSprites["BlueBlock"]
			case "green_block":
				itemSprite = is.BlockSprites["GreenBlock"]
			default:
				// Try to find sprite by item ID directly
				if sprite, exists := is.BlockSprites[is.MouseAttachedItem.Item.ID]; exists {
					itemSprite = sprite
				}
			}
			
			if itemSprite != nil {
				// Draw sprite
				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Translate(x+2, y+2)
				screen.DrawImage(itemSprite, opts)
			} else {
				// Fallback to color
				ebitenutil.DrawRect(screen, x+2, y+2, SlotSize-4, SlotSize-4, is.MouseAttachedItem.Item.Color)
			}
		} else if is.MouseAttachedItem.Item != nil {
			// Draw item color
			ebitenutil.DrawRect(screen, x+2, y+2, SlotSize-4, SlotSize-4, is.MouseAttachedItem.Item.Color)
		}
		
		// Draw item count
		countText := fmt.Sprintf("%d", is.MouseAttachedItem.Count)
		textWidth := len(countText) * 6 // Approximate character width
		ebitenutil.DebugPrintAt(screen, countText, int(x)+SlotSize-textWidth-2, int(y)+SlotSize-12)
	}
}

// handleMouseAttachment handles the mouse attachment logic for inventory slots
func (is *InventorySystem) handleMouseAttachment(inventory *components.Inventory) {
	// Get mouse position
	mouseX, mouseY := ebiten.CursorPosition()
	
	// Check if we are attaching an item to the mouse
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if is.MouseAttachedSlot == -1 { // Only attach if nothing is already attached
			if is.Visible {
				// Define creative items
				creativeItems := []components.Item{
					{
						ID:       "stone",
						Name:     "Stone",
						Count:    1,
						MaxStack: 64,
						Color:    color.RGBA{128, 128, 128, 255},
					},
					{
						ID:       "dirt",
						Name:     "Dirt",
						Count:    1,
						MaxStack: 64,
						Color:    color.RGBA{100, 50, 0, 255},
					},
					{
						ID:       "wood",
						Name:     "Wood",
						Count:    1,
						MaxStack: 64,
						Color:    color.RGBA{100, 70, 30, 255},
					},
					{
						ID:       "small_block",
						Name:     "Small Block",
						Count:    1,
						MaxStack: 64,
						Color:    color.RGBA{200, 200, 50, 255},
					},
					{
						ID:       "red_block",
						Name:     "Red Block",
						Count:    1,
						MaxStack: 64,
						Color:    color.RGBA{200, 50, 50, 255},
					},
					{
						ID:       "blue_block",
						Name:     "Blue Block",
						Count:    1,
						MaxStack: 64,
						Color:    color.RGBA{50, 50, 200, 255},
					},
					{
						ID:       "green_block",
						Name:     "Green Block",
						Count:    1,
						MaxStack: 64,
						Color:    color.RGBA{50, 200, 50, 255},
					},
				}
				
				// Calculate total slots (only include creative items in creative mode)
				totalSlots := len(inventory.Slots)
				if is.GameMode == 1 {
					totalSlots += len(creativeItems)
				}
				
				cols := 9
				
				// Check all slots including creative items (if in creative mode)
				for i := 0; i < totalSlots; i++ {
					var slot *components.ItemStack
					var isCreativeItem bool
					
					// Determine if this is a player inventory slot or a creative item slot
					if i < len(inventory.Slots) {
						// Player inventory slot
						slot = &inventory.Slots[i]
					} else {
						// Creative item slot (only in creative mode)
						isCreativeItem = true
						creativeItemIndex := i - len(inventory.Slots)
						slot = &components.ItemStack{
							Item:  &creativeItems[creativeItemIndex],
							Count: creativeItems[creativeItemIndex].MaxStack,
						}
					}
					
					row := i / cols
					col := i % cols

					x := float64(InventoryX + col*(SlotSize+SlotMargin))
					y := float64(60 + row*(SlotSize+SlotMargin))
					
					// Check if mouse is within slot bounds
					if float64(mouseX) >= x && float64(mouseX) <= x+SlotSize && float64(mouseY) >= y && float64(mouseY) <= y+SlotSize {
						// Only attach if slot has an item
						if slot.Item != nil && slot.Count > 0 {
							if isCreativeItem {
								// For creative items, attach a copy with max stack count
								is.MouseAttachedItem = &components.ItemStack{
									Item:  slot.Item,
									Count: slot.Count,
								}
								// Use special slot value to indicate this is a creative item
								is.MouseAttachedSlot = -2
							} else {
								// For player inventory slots
								is.MouseAttachedSlot = i
								// Create a copy of the item stack to attach to mouse
								is.MouseAttachedItem = &components.ItemStack{
									Item:  slot.Item,
									Count: slot.Count,
								}
							}
							return
						}
					}
				}
			} else {
				// Check hotbar slots
				is.checkHotbarSlotClick(inventory, float64(mouseX), float64(mouseY))
			}
		} else {
			// Place item in clicked slot
			if is.Visible {
				// Handle placement in full inventory
				is.handleFullInventoryPlacement(inventory, float64(mouseX), float64(mouseY))
			} else {
				// Handle placement in hotbar
				is.handleHotbarPlacement(inventory, float64(mouseX), float64(mouseY))
			}
		}
	}
	
	// Handle right-click to detach creative items
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) && is.MouseAttachedSlot == -2 {
		is.MouseAttachedSlot = -1
		is.MouseAttachedItem = nil
	}
}

// checkHotbarSlotClick checks if a hotbar slot was clicked for attachment
func (is *InventorySystem) checkHotbarSlotClick(inventory *components.Inventory, mouseX, mouseY float64) {
	for i := 0; i < inventory.HotbarSize && i < len(inventory.Slots); i++ {
		slot := inventory.Slots[i]
		x := float64(HotbarX + i*(SlotSize+SlotMargin))
		y := float64(HotbarY)
		
		// Check if mouse is within slot bounds
		if mouseX >= x && mouseX <= x+SlotSize && mouseY >= y && mouseY <= y+SlotSize {
			// Only attach if slot has an item
			if slot.Item != nil && slot.Count > 0 {
				is.MouseAttachedSlot = i
				// Create a copy of the item stack to attach to mouse
				is.MouseAttachedItem = &components.ItemStack{
					Item:  slot.Item,
					Count: slot.Count,
				}
				return
			}
		}
	}
}

// checkFullInventorySlotClick checks if a full inventory slot was clicked for attachment
func (is *InventorySystem) checkFullInventorySlotClick(inventory *components.Inventory, mouseX, mouseY float64) {
	cols := 9
	for i, slot := range inventory.Slots {
		row := i / cols
		col := i % cols
		
		x := float64(InventoryX + col*(SlotSize+SlotMargin))
		y := float64(60 + row*(SlotSize+SlotMargin))
		
		// Check if mouse is within slot bounds
		if mouseX >= x && mouseX <= x+SlotSize && mouseY >= y && mouseY <= y+SlotSize {
			// Only attach if slot has an item
			if slot.Item != nil && slot.Count > 0 {
				is.MouseAttachedSlot = i
				// Create a copy of the item stack to attach to mouse
				is.MouseAttachedItem = &components.ItemStack{
					Item:  slot.Item,
					Count: slot.Count,
				}
				return
			}
		}
	}
}

// handleHotbarPlacement handles placing an attached item in the hotbar
func (is *InventorySystem) handleHotbarPlacement(inventory *components.Inventory, mouseX, mouseY float64) {
	// Find which slot we're placing onto
	for i := 0; i < inventory.HotbarSize && i < len(inventory.Slots); i++ {
		x := float64(HotbarX + i*(SlotSize+SlotMargin))
		y := float64(HotbarY)
		
		// Check if mouse is within slot bounds
		if mouseX >= x && mouseX <= x+SlotSize && mouseY >= y && mouseY <= y+SlotSize {
			// Place item in slot
			if is.MouseAttachedItem != nil {
				// If it's the same slot, just detach the item
				if i == is.MouseAttachedSlot {
					// But don't detach creative items
					if is.MouseAttachedSlot != -2 {
						is.MouseAttachedSlot = -1
						is.MouseAttachedItem = nil
					}
					return
				}
				
				// For creative mode items, just place them in the slot (don't swap)
				if is.MouseAttachedSlot == -2 { // -2 indicates creative item
					targetSlot := &inventory.Slots[i]
					// Make sure the target slot is properly initialized
					if targetSlot == nil {
						inventory.Slots[i] = components.ItemStack{}
						targetSlot = &inventory.Slots[i]
					}
					targetSlot.Item = is.MouseAttachedItem.Item
					targetSlot.Count = is.MouseAttachedItem.Count
					
					// Keep the item attached to mouse for multiple placements
					// is.MouseAttachedItem = nil
					return
				}
				
				// Handle regular item placement/swapping
				if is.MouseAttachedSlot >= 0 && is.MouseAttachedSlot < len(inventory.Slots) {
					// Swap items between slots
					attachedSlot := &inventory.Slots[is.MouseAttachedSlot]
					targetSlot := &inventory.Slots[i]
					
					// Store the target slot's item
					var targetItem *components.ItemStack
					if targetSlot.Item != nil && targetSlot.Count > 0 {
						targetItem = &components.ItemStack{
							Item:  targetSlot.Item,
							Count: targetSlot.Count,
						}
					}
					
					// Place attached item in target slot
					targetSlot.Item = is.MouseAttachedItem.Item
					targetSlot.Count = is.MouseAttachedItem.Count
					
					// Place target item in attached slot (or clear if target was empty)
					if targetItem != nil {
						attachedSlot.Item = targetItem.Item
						attachedSlot.Count = targetItem.Count
					} else {
						attachedSlot.Item = nil
						attachedSlot.Count = 0
					}
				} else {
					// Just place the item in the slot (for items attached from outside the inventory)
					targetSlot := &inventory.Slots[i]
					targetSlot.Item = is.MouseAttachedItem.Item
					targetSlot.Count = is.MouseAttachedItem.Count
				}
				
				// Detach item from mouse
				is.MouseAttachedSlot = -1
				is.MouseAttachedItem = nil
				return
			}
		}
	}
	
	// If clicked outside a slot, drop the item (detach it) - but not for creative items
	if is.MouseAttachedSlot != -2 {
		is.MouseAttachedSlot = -1
		is.MouseAttachedItem = nil
	}
}

// handleFullInventoryPlacement handles placing an attached item in the full inventory
func (is *InventorySystem) handleFullInventoryPlacement(inventory *components.Inventory, mouseX, mouseY float64) {
	// Define creative items
	creativeItems := []components.Item{
		{
			ID:       "stone",
			Name:     "Stone",
			Count:    1,
			MaxStack: 64,
			Color:    color.RGBA{128, 128, 128, 255},
		},
		{
			ID:       "dirt",
			Name:     "Dirt",
			Count:    1,
			MaxStack: 64,
			Color:    color.RGBA{100, 50, 0, 255},
		},
		{
			ID:       "wood",
			Name:     "Wood",
			Count:    1,
			MaxStack: 64,
			Color:    color.RGBA{100, 70, 30, 255},
		},
		{
			ID:       "small_block",
			Name:     "Small Block",
			Count:    1,
			MaxStack: 64,
			Color:    color.RGBA{200, 200, 50, 255},
		},
		{
			ID:       "red_block",
			Name:     "Red Block",
			Count:    1,
			MaxStack: 64,
			Color:    color.RGBA{200, 50, 50, 255},
		},
		{
			ID:       "blue_block",
			Name:     "Blue Block",
			Count:    1,
			MaxStack: 64,
			Color:    color.RGBA{50, 50, 200, 255},
		},
		{
			ID:       "green_block",
			Name:     "Green Block",
			Count:    1,
			MaxStack: 64,
			Color:    color.RGBA{50, 200, 50, 255},
		},
	}
	
	// Calculate total slots (only include creative items in creative mode)
	totalSlots := len(inventory.Slots)
	if is.GameMode == 1 {
		totalSlots += len(creativeItems)
	}
	
	cols := 9
	
	// Find which slot we're placing onto
	for i := 0; i < totalSlots; i++ {
		var isCreativeItem bool
		
		// Determine if this is a player inventory slot or a creative item slot
		if i >= len(inventory.Slots) {
			// Creative item slot (only in creative mode)
			isCreativeItem = true
		}
		
		row := i / cols
		col := i % cols

		x := float64(InventoryX + col*(SlotSize+SlotMargin))
		y := float64(60 + row*(SlotSize+SlotMargin))
		
		// Check if mouse is within slot bounds
		if float64(mouseX) >= x && float64(mouseX) <= x+SlotSize && float64(mouseY) >= y && float64(mouseY) <= y+SlotSize {
			// Place item in slot
			if is.MouseAttachedItem != nil {
				// If it's the same slot, just detach the item
				if i == is.MouseAttachedSlot {
					// But don't detach creative items
					if is.MouseAttachedSlot != -2 {
						is.MouseAttachedSlot = -1
						is.MouseAttachedItem = nil
					}
					return
				}
				
				// For creative mode items, just place them in the slot (don't swap)
				if is.MouseAttachedSlot == -2 { // -2 indicates creative item
					if !isCreativeItem {
						// Can only place creative items in player inventory slots
						targetSlot := &inventory.Slots[i]
						// Make sure the target slot is properly initialized
						if targetSlot == nil {
							inventory.Slots[i] = components.ItemStack{}
							targetSlot = &inventory.Slots[i]
						}
						
						// Create a copy of the item to place in the slot
						itemCopy := *is.MouseAttachedItem.Item
						targetSlot.Item = &itemCopy
						targetSlot.Count = is.MouseAttachedItem.Count
						
						// Keep the item attached to mouse for multiple placements
						// is.MouseAttachedItem = nil
						return
					}
				}
				
				// Handle regular item placement/swapping
				if is.MouseAttachedSlot >= 0 && is.MouseAttachedSlot < len(inventory.Slots) {
					// Can only swap items in player inventory slots
					if !isCreativeItem {
						// Swap items between slots
						attachedSlot := &inventory.Slots[is.MouseAttachedSlot]
						targetSlot := &inventory.Slots[i]
						
						// Store the target slot's item
						var targetItem *components.ItemStack
						if targetSlot.Item != nil && targetSlot.Count > 0 {
							targetItem = &components.ItemStack{
								Item:  targetSlot.Item,
								Count: targetSlot.Count,
							}
						}
						
						// Place attached item in target slot
						targetSlot.Item = is.MouseAttachedItem.Item
						targetSlot.Count = is.MouseAttachedItem.Count
						
						// Place target item in attached slot (or clear if target was empty)
						if targetItem != nil {
							attachedSlot.Item = targetItem.Item
							attachedSlot.Count = targetItem.Count
						} else {
							attachedSlot.Item = nil
							attachedSlot.Count = 0
						}
					}
				} else if is.MouseAttachedSlot < 0 {
					// Just place the item in the slot (for items attached from outside the inventory)
					if !isCreativeItem {
						targetSlot := &inventory.Slots[i]
						targetSlot.Item = is.MouseAttachedItem.Item
						targetSlot.Count = is.MouseAttachedItem.Count
					}
				}
				
				// Detach item from mouse (but not for creative items)
				if is.MouseAttachedSlot != -2 {
					is.MouseAttachedSlot = -1
					is.MouseAttachedItem = nil
				}
				return
			}
		}
	}
	
	// If clicked outside a slot, drop the item (detach it) - but not for creative items
	if is.MouseAttachedSlot != -2 {
		is.MouseAttachedSlot = -1
		is.MouseAttachedItem = nil
	}
}

// handleCreativeItemClick checks if a creative item was clicked and attaches it to the mouse
func (is *InventorySystem) handleCreativeItemClick(mouseX, mouseY float64, creativeItems []components.Item) {
	// Only handle clicks in creative mode
	if is.GameMode != 1 {
		return
	}
	
	// Check if a creative item was clicked
	cols := 9
	startX := 10
	startY := 420
	
	for i, item := range creativeItems {
		col := i % cols
		row := i / cols
		
		x := float64(startX + col*(SlotSize+SlotMargin))
		y := float64(startY + row*(SlotSize+SlotMargin))
		
		// Check if mouse is within item bounds
		if mouseX >= x && mouseX <= x+SlotSize && mouseY >= y && mouseY <= y+SlotSize {
			// Create a copy of the item
			itemCopy := item
			
			// Attach a copy of this item to the mouse with max stack count
			is.MouseAttachedItem = &components.ItemStack{
				Item:  &itemCopy,
				Count: item.MaxStack,
			}
			// Use special slot value to indicate this is a creative item
			is.MouseAttachedSlot = -2
			return
		}
	}
}

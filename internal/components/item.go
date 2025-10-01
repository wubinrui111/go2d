package components

import "image/color"

// Item represents an item that can be collected and stored in an inventory
type Item struct {
	// ID is the unique identifier for the item type
	ID string

	// Name is the display name of the item
	Name string

	// Count is the quantity of this item
	Count int

	// MaxStack is the maximum number of this item that can be stacked
	MaxStack int

	// Color represents the item's color for rendering purposes
	Color color.RGBA
}

// ItemStack represents a stack of items in an inventory slot
type ItemStack struct {
	// Item is the type of item in this stack
	Item *Item

	// Count is the number of items in this stack
	Count int
}
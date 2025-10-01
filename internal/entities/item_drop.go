package entities

import (
	"image/color"
	"math"
	"github.com/wubinrui111/2d-game/internal/components"
)

const (
	// ItemDropSize is the size of an item drop
	ItemDropSize = 16.0
	
	// PickupDistance is the distance at which items are automatically picked up
	PickupDistance = 32.0
	
	// AttractionDistance is the distance at which items start being attracted to the player
	AttractionDistance = 100.0
)

// ItemDrop represents a dropped item on the ground
type ItemDrop struct {
	components.Position
	components.Box
	Item     *components.Item // The item this drop represents
	Color    color.RGBA       // Color of the item for rendering
	Velocity components.Position // Velocity for movement
}

// NewItemDrop creates a new item drop at the given position
func NewItemDrop(x, y float64, item *components.Item) *ItemDrop {
	drop := &ItemDrop{
		Position: components.Position{
			X: x,
			Y: y,
		},
		Box: components.Box{
			X:      x,
			Y:      y,
			Width:  ItemDropSize,
			Height: ItemDropSize,
		},
		Item: item,
		Color: item.Color,
		Velocity: components.Position{
			X: 0,
			Y: 0,
		},
	}
	
	return drop
}

// Update updates the item drop's position and behavior
func (id *ItemDrop) Update(playerPos components.Position) {
	// Calculate distance to player
	dx := playerPos.X - id.Position.X
	dy := playerPos.Y - id.Position.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	
	// If player is close enough, attract the item
	if distance <= AttractionDistance && distance > PickupDistance/2 {
		// Normalize direction vector
		if distance > 0 {
			dx /= distance
			dy /= distance
		}
		
		// Apply attraction force (weaker when farther away)
		attractionStrength := 100.0 * (AttractionDistance - distance) / AttractionDistance
		id.Velocity.X += dx * attractionStrength
		id.Velocity.Y += dy * attractionStrength
	}
	
	// Apply some drag to slow down movement
	id.Velocity.X *= 0.9
	id.Velocity.Y *= 0.9
	
	// Update position based on velocity
	id.Position.X += id.Velocity.X
	id.Position.Y += id.Velocity.Y
	
	// Update box position to match
	id.Box.X = id.Position.X
	id.Box.Y = id.Position.Y
}

// ShouldPickup checks if the item should be picked up by the player
func (id *ItemDrop) ShouldPickup(playerPos components.Position) bool {
	// Calculate distance to player
	dx := playerPos.X - id.Position.X
	dy := playerPos.Y - id.Position.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	
	return distance <= PickupDistance
}

// GetItem returns the item this drop represents
func (id *ItemDrop) GetItem() *components.Item {
	return id.Item
}
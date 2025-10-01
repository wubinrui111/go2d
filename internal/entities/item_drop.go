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
	
	// GravityForce is the force of gravity applied to item drops
	GravityForce = 300.0
	
	// ShrinkDuration is the time in seconds it takes for an item to shrink to nothing
	ShrinkDuration = 30.0
	
	// Lifetime is the total time in seconds an item drop exists before disappearing
	Lifetime = 60.0
	
	// GroundLevel is the Y position where the ground is located
	GroundLevel = 550.0
	
	// BounceFactor is how much velocity is retained after bouncing
	BounceFactor = 0.3
	
	// FrictionFactor is how much horizontal velocity is reduced when on ground
	FrictionFactor = 0.8
	
	// MinBounceVelocity is the minimum velocity required to bounce
	MinBounceVelocity = 50.0
)

// ItemDrop represents a dropped item on the ground
type ItemDrop struct {
	components.Position
	components.Box
	Item     *components.Item // The item this drop represents
	Color    color.RGBA       // Color of the item for rendering
	Velocity components.Position // Velocity for movement
	Gravity  *components.Gravity // Gravity component
	Life     float64          // Current life in seconds
	OnGround bool             // Whether the item drop is on the ground
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
		Gravity: &components.Gravity{
			Enabled: true,
			Force:   GravityForce,
		},
		Life: 0,
	}
	
	return drop
}

// Update updates the item drop's position and behavior
func (id *ItemDrop) Update(playerPos components.Position, deltaTime float64, blocks []components.BoxHolder) {
	// Update life
	id.Life += deltaTime
	
	// Apply gravity continuously (don't check OnGround to allow infinite falling)
	if id.Gravity.Enabled {
		id.Velocity.Y += id.Gravity.Force * deltaTime
	}
	
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
		id.Velocity.X += dx * attractionStrength * deltaTime
		id.Velocity.Y += dy * attractionStrength * deltaTime
	}
	
	// Apply some drag to slow down movement
	id.Velocity.X *= math.Pow(0.9, deltaTime*60)
	id.Velocity.Y *= math.Pow(0.9, deltaTime*60)
	
	// Store previous position for collision detection
	prevX := id.Position.X
	prevY := id.Position.Y
	
	// Update position based on velocity
	id.Position.X += id.Velocity.X * deltaTime * 60
	id.Position.Y += id.Velocity.Y * deltaTime * 60
	
	// Check block collisions
	id.checkBlockCollisions(blocks, prevX, prevY)
	
	// 移除地面碰撞检测调用，确保持续下落
	// id.checkGroundCollision()
	
	// Update box position to match
	id.Box.X = id.Position.X
	id.Box.Y = id.Position.Y
	
	// Update size based on life (shrink over time)
	lifeRatio := id.Life / Lifetime
	if lifeRatio > 0.5 { // Start shrinking after half the lifetime
		shrinkProgress := (lifeRatio - 0.5) / 0.5 // From 0 to 1 in the second half of life
		id.Box.Width = ItemDropSize * (1 - shrinkProgress)
		id.Box.Height = ItemDropSize * (1 - shrinkProgress)
	}
}

// checkBlockCollisions checks if the item drop collides with any blocks
func (id *ItemDrop) checkBlockCollisions(blocks []components.BoxHolder, prevX, prevY float64) {
	// Create a temporary box for the new position
	tempBox := components.Box{
		X:      id.Position.X,
		Y:      id.Position.Y,
		Width:  ItemDropSize,
		Height: ItemDropSize,
	}
	
	// Check collision with all blocks
	for _, block := range blocks {
		if tempBox.Intersects(block.GetBox()) {
			// Calculate intersection depth
			xDepth, yDepth := tempBox.GetIntersectionDepth(block.GetBox())
			
			// Determine the minimum translation vector
			if math.Abs(xDepth) < math.Abs(yDepth) {
				// Horizontal collision - stop horizontal movement
				id.Position.X = prevX // Revert to previous X position
				id.Velocity.X = 0
			} else {
				// Vertical collision - stop vertical movement
				id.Position.Y = prevY // Revert to previous Y position
				id.Velocity.Y = 0
			}
		}
	}
}

// checkGroundCollision disables ground collision, allowing items to fall infinitely
func (id *ItemDrop) checkGroundCollision() {
	// Ground collision is disabled to allow infinite falling
	id.OnGround = false
}

// ShouldPickup checks if the item should be picked up by the player
func (id *ItemDrop) ShouldPickup(playerPos components.Position) bool {
	// If item has exceeded its lifetime, it should disappear
	if id.Life >= Lifetime {
		return false // Don't pick up, just disappear
	}
	
	// Calculate distance to player
	dx := playerPos.X - id.Position.X
	dy := playerPos.Y - id.Position.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	
	return distance <= PickupDistance
}

// ShouldDisappear checks if the item should disappear (lifetime exceeded)
func (id *ItemDrop) ShouldDisappear() bool {
	return id.Life >= Lifetime
}

// GetItem returns the item this drop represents
func (id *ItemDrop) GetItem() *components.Item {
	return id.Item
}

// GetCurrentSize returns the current size of the item drop (considering shrink effect)
func (id *ItemDrop) GetCurrentSize() (width, height float64) {
	return id.Box.Width, id.Box.Height
}
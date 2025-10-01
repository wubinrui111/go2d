// internal/components/box.go
package components

import (
	"math"
)

// Box represents a rectangular collision box
type Box struct {
	X, Y, Width, Height float64
}

// BoxHolder interface for anything that has a Box
type BoxHolder interface {
	GetBox() *Box
}

// GetBox returns the box itself
func (b *Box) GetBox() *Box {
	return b
}

// Intersects checks if this box intersects with another box
func (b *Box) Intersects(other *Box) bool {
	return b.X < other.X+other.Width &&
		b.X+b.Width > other.X &&
		b.Y < other.Y+other.Height &&
		b.Y+b.Height > other.Y
}

// GetIntersectionDepth calculates the depth of intersection between two boxes
// Returns the x and y depths needed to separate the boxes
func (b *Box) GetIntersectionDepth(other *Box) (float64, float64) {
	// Calculate centers
	centerX1 := b.X + b.Width/2
	centerY1 := b.Y + b.Height/2
	centerX2 := other.X + other.Width/2
	centerY2 := other.Y + other.Height/2
	
	// Calculate minimum translation distances
	minTransX := (b.Width + other.Width) / 2 - math.Abs(centerX1-centerX2)
	minTransY := (b.Height + other.Height) / 2 - math.Abs(centerY1-centerY2)
	
	// Determine direction to push
	if centerX1 < centerX2 {
		minTransX = -minTransX
	}
	
	if centerY1 < centerY2 {
		minTransY = -minTransY
	}
	
	return minTransX, minTransY
}
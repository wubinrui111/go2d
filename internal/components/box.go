// internal/components/box.go
package components

// Box represents a rectangular box for collision detection
type Box struct {
	X      float64 // X position
	Y      float64 // Y position
	Width  float64 // Width of the box
	Height float64 // Height of the box
}

// Intersects checks if this box intersects with another box
func (b *Box) Intersects(other *Box) bool {
	return b.X < other.X+other.Width &&
		b.X+b.Width > other.X &&
		b.Y < other.Y+other.Height &&
		b.Y+b.Height > other.Y
}

// GetIntersectionDepth calculates how much this box penetrates into another box
// Returns the penetration depth on X and Y axes
// Positive values indicate penetration in the positive direction, negative in the negative direction
func (b *Box) GetIntersectionDepth(other *Box) (float64, float64) {
	// Calculate overlap on each axis
	xOverlap := 0.0
	yOverlap := 0.0
	
	if b.X < other.X+other.Width && b.X+b.Width > other.X {
		// Calculate the minimum translation distance on X axis
		leftOverlap := (b.X + b.Width) - other.X
		rightOverlap := (other.X + other.Width) - b.X
		
		if leftOverlap < rightOverlap {
			xOverlap = -leftOverlap // Negative means move left
		} else {
			xOverlap = rightOverlap // Positive means move right
		}
	}
	
	if b.Y < other.Y+other.Height && b.Y+b.Height > other.Y {
		// Calculate the minimum translation distance on Y axis
		topOverlap := (b.Y + b.Height) - other.Y
		bottomOverlap := (other.Y + other.Height) - b.Y
		
		if topOverlap < bottomOverlap {
			yOverlap = -topOverlap // Negative means move up
		} else {
			yOverlap = bottomOverlap // Positive means move down
		}
	}
	
	return xOverlap, yOverlap
}
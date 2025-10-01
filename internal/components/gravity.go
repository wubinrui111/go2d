// internal/components/gravity.go
package components

// Gravity represents the gravity component
type Gravity struct {
	Enabled bool    // Whether gravity is enabled for this entity
	Force   float64 // The force of gravity to apply
}

// NewGravity creates a new gravity component with default values
func NewGravity() *Gravity {
	return &Gravity{
		Enabled: true,
		Force:   300.0, // Pixels per second squared
	}
}
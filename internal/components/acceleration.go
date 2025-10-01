package components

// Acceleration represents the acceleration component for movement
type Acceleration struct {
	GroundSpeed    float64 // Speed on the ground (pixels per second)
	AirSpeed       float64 // Speed in the air (pixels per second)
	GroundFriction float64 // Friction when on the ground (0.0 to 1.0)
	AirResistance  float64 // Air resistance when in the air (0.0 to 1.0)
	JumpForce      float64 // Upward force when jumping (pixels per second)
}

// NewAcceleration creates a new acceleration component with default values
func NewAcceleration() *Acceleration {
	return &Acceleration{
		GroundSpeed:    10000.0, // Faster ground movement
		AirSpeed:       1500.0,  // Slower air movement
		GroundFriction: 0.8,     // High friction on ground for quick stops
		AirResistance:  0.95,    // Low air resistance for floaty movement
		JumpForce:      300.0,   // Jump force
	}
}

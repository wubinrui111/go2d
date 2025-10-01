package components

// Health represents the health component for entities
type Health struct {
	// Current health points
	Current int
	
	// Maximum health points
	Max int
	
	// Indicates if the entity is alive
	Alive bool
}

// NewHealth creates a new health component with the specified max health
func NewHealth(maxHealth int) *Health {
	return &Health{
		Current: maxHealth,
		Max:     maxHealth,
		Alive:   true,
	}
}

// TakeDamage reduces the current health by the specified amount
func (h *Health) TakeDamage(damage int) {
	h.Current -= damage
	if h.Current <= 0 {
		h.Current = 0
		h.Alive = false
	}
}

// Heal increases the current health by the specified amount
// Cannot exceed maximum health
func (h *Health) Heal(amount int) {
	h.Current += amount
	if h.Current > h.Max {
		h.Current = h.Max
	}
}

// IsAlive returns whether the entity is alive
func (h *Health) IsAlive() bool {
	return h.Alive && h.Current > 0
}

// SetMaxHealth sets the maximum health and adjusts current health if needed
func (h *Health) SetMaxHealth(max int) {
	h.Max = max
	if h.Current > h.Max {
		h.Current = h.Max
	}
}

// GetHealthPercentage returns the current health as a percentage of maximum health
func (h *Health) GetHealthPercentage() float64 {
	if h.Max <= 0 {
		return 0
	}
	return float64(h.Current) / float64(h.Max) * 100
}
package quality

// Quality is a hint for controlling the amount of interpolation a graphics
// context does when scaling an image.
type Quality int

const (
	// Default lets the context decide.
	Default Quality = iota
	// None turns off interpolation.
	None
	// Low quality, fast interpolation.
	Low
	// High quality, slower than Medium.
	High
	// Medium quality, slower than Low.
	Medium
)

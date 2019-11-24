package side

// Side constants.
const (
	Top Side = iota
	Left
	Bottom
	Right
)

// Side specifies which side an object should be on.
type Side uint8

// Horizontal returns true if the side is to the left or right.
func (s Side) Horizontal() bool {
	return s == Left || s == Right
}

// Vertical returns true if the side is to the top or bottom.
func (s Side) Vertical() bool {
	return s == Top || s == Bottom
}

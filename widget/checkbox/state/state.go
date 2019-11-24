package state

// Possible values for State.
const (
	Unchecked State = iota
	Mixed
	Checked
)

// State represents the current state of a checkbox.
type State uint8

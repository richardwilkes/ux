package undo

// Edit defines the required methods an undoable edit must implement.
type Edit interface {
	// Name returns the localized name of the edit, suitable for displaying in
	// a user interface menu. Note that no leading "Undo " or "Redo " should
	// be part of this name, as the Manager will add this.
	Name() string
	// Cost returns a cost factor for this edit. When the cost values of the
	// edits within a given Manager exceed the Manager's defined cost limit,
	// the oldest edits will be discarded until the cost values are less than
	// or equal to the Manager's defined limit. Note that if this method
	// returns a value less than 1, it will be set to 1 for purposes of this
	// calculation.
	Cost() int
	// Undo the state.
	Undo()
	// Redo the state.
	Redo()
	// Absorb gives this edit a chance to absorb a new edit that is about to
	// be added to the manager. If this method returns true, it is assumed
	// this edit has incorporated any necessary state into itself to perform
	// an undo/redo and the other edit will be discarded.
	Absorb(other Edit) bool
	// Release is called when this edit is no longer needed by the Manager.
	Release()
}

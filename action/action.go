package action

import (
	"github.com/richardwilkes/ux/keys"
)

// Action describes an action that can be performed.
type Action interface {
	// ID returns a unique ID for the action. This value should be suitable
	// for use as a menu item ID.
	ID() int
	// Title returns the text to display for this action. Typically used in a
	// menu item title or tooltip for a button.
	Title() string
	// HotKey is the key that will trigger the action. Returns nil if no hot
	// key is set.
	HotKey() *keys.Key
	// HotKeyModifiers returns the modifier keys that must be pressed for the
	// hot key to be recognized.
	HotKeyModifiers() keys.Modifiers
	// Enabled returns true if the action can be used. Care should be made to
	// keep this method fast to avoid slowing down the user interface.
	Enabled(source interface{}) bool
	// Execute the action. Will only be called if the action has been
	// triggered and Enabled() returns true.
	Execute(source interface{})
}

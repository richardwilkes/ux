// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

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

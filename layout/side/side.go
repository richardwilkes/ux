// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

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

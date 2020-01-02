// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

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

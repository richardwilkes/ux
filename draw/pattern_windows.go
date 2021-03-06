// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package draw

type osPattern = int

func osNewPattern(img *Image) osPattern {
	return 0 // RAW: Implement
}

func (p *Pattern) osPrepareForFill(gc Context) {
	// RAW: Implement
}

func (p *Pattern) osFill(gc Context) {
	// RAW: Implement
}

func (p *Pattern) osFillEvenOdd(gc Context) {
	// RAW: Implement
}

func (p *Pattern) osStroke(gc Context) {
	// RAW: Implement
}

func (r *patternRef) osIsValid() bool {
	return false // RAW: Implement
}

func (r *patternRef) osDispose() {
	// RAW: Implement
}

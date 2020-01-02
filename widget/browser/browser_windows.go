// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package browser

import (
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
)

type osBrowser = int

func osNewBrowser(wnd *ux.Window) (osBrowser, error) {
	// RAW: Implement
	return 0, errs.New("browser panel not supported")
}

func (b *Browser) osSetFrame(rect geom.Rect) {
	// RAW: Implement
}

func (b *Browser) osLoadURL(url string) {
	// RAW: Implement
}

func (b *Browser) osDispose() {
	// RAW: Implement
}

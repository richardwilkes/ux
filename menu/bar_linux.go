// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package menu

import "github.com/richardwilkes/ux"

func osMenuBarForWindow(wnd *ux.Window, updater Updater) (bar *Bar, isGlobal, isFirst bool) {
	if !wnd.IsValid() {
		return nil, false, false
	}
	// RAW: Implement
	return nil, false, false
}

func osMenuBarHeightInWindow() float64 {
	// RAW: Implement
	return 0
}

func (bar *Bar) osInsertMenu(atIndex, id int, menu *Menu) {
	// RAW: Implement
}

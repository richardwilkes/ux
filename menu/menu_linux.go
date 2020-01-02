// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package menu

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/keys"
)

type osMenu = int

func osNewMenu(title string, updater Updater) osMenu {
	// RAW: Implement
	return 0
}

func (menu *Menu) osIsSame(other *Menu) bool {
	// RAW: Implement
	return menu.native == other.native
}

func (menu *Menu) osItemAtIndex(index int) *Item {
	// RAW: Implement
	return nil
}

func (menu *Menu) osInsertSeparator(atIndex int) {
	// RAW: Implement
}

func (menu *Menu) osInsertItem(atIndex, id int, title string, key *keys.Key, keyModifiers keys.Modifiers, validator ItemValidator, handler ItemHandler) *Item {
	// RAW: Implement
	return nil
}

func (menu *Menu) osInsertNewMenu(atIndex, id int, title string, updater Updater) *Menu {
	// RAW: Implement
	return nil
}

func (menu *Menu) osInsertMenu(atIndex, id int, subMenu *Menu) {
	// RAW: Implement
}

func (menu *Menu) osRemoveItem(index int) {
	// RAW: Implement
}

func (menu *Menu) osItemCount() int {
	// RAW: Implement
	return 0
}

func (menu *Menu) osPopup(wnd *ux.Window, where geom.Rect, currentIndex int) {
	// RAW: Implement
}

func (menu *Menu) osDispose() {
	// RAW: Implement
}

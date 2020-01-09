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
	"github.com/richardwilkes/ux/keys"
	"github.com/richardwilkes/ux/widget/checkbox/state"
)

type osItemData struct {
	validator ItemValidator
	handler   ItemHandler
	key       *keys.Key
	modifiers keys.Modifiers
}

type osItem = *osItemData

func (item *Item) osIsSame(other *Item) bool {
	// RAW: Implement
	return item.native == other.native
}

func (item *Item) osMenu() *Menu {
	// RAW: Implement
	return nil
}

func (item *Item) osIsSeparator() bool {
	// RAW: Implement
	return false
}

func (item *Item) osID() int {
	// RAW: Implement
	return 0
}

func (item *Item) osTitle() string {
	// RAW: Implement
	return ""
}

func (item *Item) osSetTitle(title string) {
	// RAW: Implement
	// win32.SetMenuItemInfo(menu.menu, uint32(index), true, &win32.MENUITEMINFO{
	// 	Size:     uint32(unsafe.Sizeof(win32.MENUITEMINFO{})), //nolint:gosec
	// 	Mask:     win32.MIIM_STRING,
	// 	TypeData: win32.ToSysWin32Str(title, false),
	// })
}

func (item *Item) osSubMenu() *Menu {
	// RAW: Implement
	return nil
}

func (item *Item) osCheckState() state.State {
	// RAW: Implement
	return state.Off
}

func (item *Item) osSetCheckState(s state.State) {
	// RAW: Implement
}

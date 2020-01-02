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
	"runtime"

	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/ux/action"
	"github.com/richardwilkes/ux/ids"
	"github.com/richardwilkes/ux/keys"
)

// NewEditMenu creates a standard 'Edit' menu.
func NewEditMenu(prefsHandler ItemHandler, updater Updater) *Menu {
	menu := New(i18n.Text("Edit"), updater)
	menu.InsertActionItem(-1, action.Cut)
	menu.InsertActionItem(-1, action.Copy)
	menu.InsertActionItem(-1, action.Paste)
	menu.InsertActionItem(-1, action.Delete)
	menu.InsertActionItem(-1, action.SelectAll)
	if runtime.GOOS != toolbox.MacOS && prefsHandler != nil {
		menu.InsertSeparator(-1)
		menu.InsertItem(-1, ids.PreferencesItemID, i18n.Text("Preferences…"), keys.Comma, keys.OSMenuCmdModifier(), nil, prefsHandler)
	}
	return menu
}

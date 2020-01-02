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
	"fmt"
	"runtime"

	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/ids"
	"github.com/richardwilkes/ux/keys"
)

// NewAppMenu creates a standard 'App' menu. Really only intended for macOS,
// although other platforms can use it if desired.
func NewAppMenu(aboutHandler, prefsHandler ItemHandler, updater Updater) *Menu {
	menu := New(cmdline.AppName, updater)
	menu.InsertItem(-1, ids.AboutItemID, fmt.Sprintf(i18n.Text("About %s"), cmdline.AppName), nil, 0, func(*Item) bool { return aboutHandler != nil }, aboutHandler)
	if prefsHandler != nil {
		menu.InsertSeparator(-1)
		menu.InsertItem(-1, ids.PreferencesItemID, i18n.Text("Preferences…"), keys.Comma, keys.OSMenuCmdModifier(), nil, prefsHandler)
	}
	if runtime.GOOS == toolbox.MacOS {
		menu.InsertSeparator(-1)
		menu.InsertMenu(-1, ids.ServicesMenuID, i18n.Text("Services"), nil)
		menu.InsertSeparator(-1)
		menu.InsertItem(-1, ids.HideItemID, fmt.Sprintf(i18n.Text("Hide %s"), cmdline.AppName), keys.H, keys.OSMenuCmdModifier(), nil, func(*Item) { ux.HideApp() })
		menu.InsertItem(-1, ids.HideOthersItemID, i18n.Text("Hide Others"), keys.H, keys.OptionModifier|keys.OSMenuCmdModifier(), nil, func(*Item) { ux.HideOtherApps() })
		menu.InsertItem(-1, ids.ShowAllItemID, i18n.Text("Show All"), nil, 0, nil, func(*Item) { ux.ShowAllApps() })
	}
	menu.InsertSeparator(-1)
	InsertQuitItem(menu, -1)
	return menu
}

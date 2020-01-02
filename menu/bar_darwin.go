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
	"github.com/richardwilkes/macos/ns"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/ids"
)

var menuBar *Bar

func osMenuBarForWindow(_ *ux.Window, updater Updater) (bar *Bar, isGlobal, isFirst bool) {
	first := false
	if menuBar == nil {
		menuBar = &Bar{bar: New("", updater)}
		ns.SharedApplication().SetMainMenu(menuBar.bar.native)
		first = true
	}
	return menuBar, true, first
}

func osMenuBarHeightInWindow() float64 {
	return 0
}

func (bar *Bar) osInsertMenu(atIndex, id int, menu *Menu) {
	bar.bar.osInsertMenu(atIndex, id, menu)
	switch id {
	case ids.AppMenuID:
		if servicesMenu := bar.Menu(ids.ServicesMenuID); servicesMenu != nil {
			ns.SharedApplication().SetServicesMenu(servicesMenu.native)
		}
	case ids.WindowMenuID:
		ns.SharedApplication().SetWindowsMenu(menu.native)
	case ids.HelpMenuID:
		ns.SharedApplication().SetHelpMenu(menu.native)
	}
}

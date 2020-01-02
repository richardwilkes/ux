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
	"github.com/richardwilkes/ux/ids"
)

// NewHelpMenu creates a standard 'Help' menu.
func NewHelpMenu(aboutHandler ItemHandler, updater Updater) *Menu {
	menu := New(i18n.Text("Help"), updater)
	if runtime.GOOS != toolbox.MacOS {
		menu.InsertItem(-1, ids.AboutItemID, fmt.Sprintf(i18n.Text("About %s"), cmdline.AppName), nil, 0, nil, aboutHandler)
	}
	return menu
}

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
func NewHelpMenu(aboutHandler func(), updater func(*Menu)) *Menu {
	menu := New(ids.HelpMenuID, i18n.Text("Help"), updater)
	if runtime.GOOS != toolbox.MacOS {
		menu.InsertItem(-1, ids.AboutItemID, fmt.Sprintf(i18n.Text("About %s"), cmdline.AppName), nil, 0, nil, aboutHandler)
	}
	return menu
}

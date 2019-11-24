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
func NewAppMenu(aboutHandler, prefsHandler func(), updater func(*Menu)) *Menu {
	menu := New(ids.AppMenuID, cmdline.AppName, updater)
	menu.InsertItem(-1, ids.AboutItemID, fmt.Sprintf(i18n.Text("About %s"), cmdline.AppName), nil, 0, func() bool { return aboutHandler != nil }, aboutHandler)
	if prefsHandler != nil {
		menu.InsertSeparator(-1)
		menu.InsertItem(-1, ids.PreferencesItemID, i18n.Text("Preferences…"), keys.Comma, keys.OSMenuCmdModifier(), nil, prefsHandler)
	}
	if runtime.GOOS == toolbox.MacOS {
		menu.InsertSeparator(-1)
		menu.InsertMenu(-1, ids.ServicesMenuID, i18n.Text("Services"), nil)
		menu.InsertSeparator(-1)
		menu.InsertItem(-1, ids.HideItemID, fmt.Sprintf(i18n.Text("Hide %s"), cmdline.AppName), keys.H, keys.OSMenuCmdModifier(), nil, ux.HideApp)
		menu.InsertItem(-1, ids.HideOthersItemID, i18n.Text("Hide Others"), keys.H, keys.OptionModifier|keys.OSMenuCmdModifier(), nil, ux.HideOtherApps)
		menu.InsertItem(-1, ids.ShowAllItemID, i18n.Text("Show All"), nil, 0, nil, ux.ShowAllApps)
	}
	menu.InsertSeparator(-1)
	InsertQuitItem(menu, -1)
	return menu
}

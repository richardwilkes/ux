package menu

import (
	"runtime"

	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/ids"
	"github.com/richardwilkes/ux/keys"
)

// NewFileMenu creates a standard 'File' menu.
func NewFileMenu(updater Updater) *Menu {
	menu := New(i18n.Text("File"), updater)
	InsertCloseKeyWindowItem(menu, -1)
	if runtime.GOOS != toolbox.MacOS {
		menu.InsertSeparator(-1)
		InsertQuitItem(menu, -1)
	}
	return menu
}

// InsertCloseKeyWindowItem creates the standard "Close" menu item that will
// close the current key window when chosen.
func InsertCloseKeyWindowItem(menu *Menu, atIndex int) {
	menu.InsertItem(atIndex, ids.CloseItemID, i18n.Text("Close"), keys.W, keys.OSMenuCmdModifier(), CloseKeyWindowValidator, CloseKeyWindowHandler)
}

// CloseKeyWindowValidator provides the standard validation function for the
// "Close" menu.
func CloseKeyWindowValidator(item *Item) bool {
	wnd := ux.WindowWithFocus()
	return wnd != nil && wnd.Closable()
}

// CloseKeyWindowHandler provides the standard handler function for the
// "Close" menu.
func CloseKeyWindowHandler(item *Item) {
	if wnd := ux.WindowWithFocus(); wnd != nil && wnd.Closable() {
		wnd.AttemptClose()
	}
}

// InsertQuitItem creates the standard "Quit"/"Exit" menu item that will
// issue the Quit command when chosen.
func InsertQuitItem(menu *Menu, atIndex int) {
	var title string
	if runtime.GOOS == toolbox.MacOS {
		title = i18n.Text("Quit")
	} else {
		title = i18n.Text("Exit")
	}
	menu.InsertItem(atIndex, ids.QuitItemID, title, keys.Q, keys.OSMenuCmdModifier(), nil, func(*Item) { ux.AttemptQuit() })
}

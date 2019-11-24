package menu

import (
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/ids"
	"github.com/richardwilkes/ux/keys"
)

// NewWindowMenu creates a standard 'Window' menu.
func NewWindowMenu(updater func(*Menu)) *Menu {
	menu := New(ids.WindowMenuID, i18n.Text("Window"), updater)
	InsertMinimizeItem(menu, -1)
	InsertZoomItem(menu, -1)
	menu.InsertSeparator(-1)
	menu.InsertItem(-1, ids.BringAllWindowsToFrontItemID, i18n.Text("Bring All to Front"), nil, 0, nil, ux.AppWindowsToFront)
	return menu
}

// InsertMinimizeItem creates the standard "Minimize" menu item that will
// issue the Minimize command to the current key window when chosen.
func InsertMinimizeItem(menu *Menu, atIndex int) {
	menu.InsertItem(atIndex, ids.MinimizeItemID, i18n.Text("Minimize"), keys.M, keys.OSMenuCmdModifier(), MinimizeValidator, MinimizeHandler)
}

// MinimizeValidator provides the standard validation function for the
// "Minimize" menu item.
func MinimizeValidator() bool {
	w := ux.WindowWithFocus()
	return w != nil && w.Minimizable()
}

// MinimizeHandler provides the standard handler function for the "Minimize"
// menu item.
func MinimizeHandler() {
	if wnd := ux.WindowWithFocus(); wnd != nil {
		wnd.Minimize()
	}
}

// InsertZoomItem creates the standard "Zoom" menu item that will issue the
// Zoom command to the current key window when chosen.
func InsertZoomItem(menu *Menu, atIndex int) {
	menu.InsertItem(atIndex, ids.ZoomItemID, i18n.Text("Zoom"), keys.Z, keys.ShiftModifier|keys.OSMenuCmdModifier(), ZoomValidator, ZoomHandler)
}

// ZoomValidator provides the standard validation function for the "Zoom" menu
// item.
func ZoomValidator() bool {
	w := ux.WindowWithFocus()
	return w != nil && w.Resizable()
}

// ZoomHandler provides the standard handler function for the "Zoom" menu
// item.
func ZoomHandler() {
	if wnd := ux.WindowWithFocus(); wnd != nil {
		wnd.Zoom()
	}
}

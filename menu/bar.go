package menu

import (
	"runtime"

	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/ids"
)

// Bar represents a set of menus.
type Bar struct {
	bar *Menu
}

// BarForWindow returns the menu bar for the given window. On macOS, the menu
// bar is a global entity and the same value will be returned for all windows.
// On macOS, you may pass nil for the window parameter. If isGlobal returns as
// true, the first time this function is called, isFirst will be returned as
// true as well, allowing you to only initialize the global menu bar once.
// 'updater' is optional. If present, it will be called prior to showing a
// top-level menu, giving a chance to modify that menu.
func BarForWindow(window *ux.Window, updater Updater) (bar *Bar, isGlobal, isFirst bool) {
	return osMenuBarForWindow(window, updater)
}

// BarHeight returns the height of the Bar when displayed in a window.
func BarHeight() float64 {
	return osMenuBarHeightInWindow()
}

// InsertStdMenus adds the standard menus to the menu bar.
func (bar *Bar) InsertStdMenus(aboutHandler, prefsHandler ItemHandler, updater Updater) {
	if runtime.GOOS == toolbox.MacOS {
		bar.InsertMenu(-1, ids.AppMenuID, NewAppMenu(aboutHandler, prefsHandler, updater))
	}
	bar.InsertMenu(-1, ids.FileMenuID, NewFileMenu(updater))
	bar.InsertMenu(-1, ids.EditMenuID, NewEditMenu(prefsHandler, updater))
	bar.InsertMenu(-1, ids.WindowMenuID, NewWindowMenu(updater))
	bar.InsertMenu(-1, ids.HelpMenuID, NewHelpMenu(aboutHandler, updater))
}

// Menu returns the menu with the specified id anywhere in the menu bar and
// its sub-menus.
func (bar *Bar) Menu(id int) *Menu {
	if item := bar.bar.Item(id); item != nil {
		return item.SubMenu()
	}
	return nil
}

// MenuAtIndex returns the menu at the specified index within the menu bar.
func (bar *Bar) MenuAtIndex(index int) *Menu {
	if item := bar.bar.ItemAtIndex(index); item != nil {
		return item.SubMenu()
	}
	return nil
}

// MenuItem returns the menu item with the specified id anywhere in the menu
// bar and its sub-menus.
func (bar *Bar) MenuItem(id int) *Item {
	return bar.bar.Item(id)
}

// InsertMenu inserts a menu at the specified item index within this menu bar.
// Pass in a negative index to append to the end.
func (bar *Bar) InsertMenu(atIndex, id int, menu *Menu) {
	bar.osInsertMenu(atIndex, id, menu)
}

// IndexOf returns the index of the menu within this menu bar, or -1.
func (bar *Bar) IndexOf(menu *Menu) int {
	for i := bar.Count() - 1; i >= 0; i-- {
		if m := bar.MenuAtIndex(i); m != nil && m.IsSame(menu) {
			return i
		}
	}
	return -1
}

// Remove the menu at the specified index from this menu bar.
func (bar *Bar) Remove(index int) {
	bar.bar.RemoveItem(index)
}

// Count of menus in this menu bar.
func (bar *Bar) Count() int {
	return bar.bar.Count()
}

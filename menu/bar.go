package menu

import (
	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/ux"
	"runtime"
)

// Bar represents a set of menus.
type Bar struct {
	bar osMenuBar
}

// BarForWindow returns the menu bar for the given window. On macOS, the menu
// bar is a global entity and the same value will be returned for all windows.
// On macOS, you may pass nil for the window parameter. If isGlobal is true,
// the first time this function is called, isFirst will be true as well,
// allowing you to only initialize the global menu bar once. 'updater' is
// optional. If present, it will be called prior to showing a top-level menu,
// giving a chance to modify that menu.
func BarForWindow(window *ux.Window, updater func(*Menu)) (bar *Bar, isGlobal, isFirst bool) {
	if window.IsValid() {
		return osMenuBarForWindow(window, updater)
	}
	return nil, false, false
}

// BarHeight returns the height of the Bar when displayed in a window.
func BarHeight() float64 {
	return osMenuBarHeightInWindow()
}

// InsertStdMenus adds the standard menus to the menu bar.
func (bar *Bar) InsertStdMenus(aboutHandler, prefsHandler func(), updater func(*Menu)) {
	if runtime.GOOS == toolbox.MacOS {
		bar.InsertMenu(-1, NewAppMenu(aboutHandler, prefsHandler, updater))
	}
	bar.InsertMenu(-1, NewFileMenu(updater))
	bar.InsertMenu(-1, NewEditMenu(prefsHandler, updater))
	bar.InsertMenu(-1, NewWindowMenu(updater))
	bar.InsertMenu(-1, NewHelpMenu(aboutHandler, updater))
}

// Menu returns the menu with the specified id anywhere in the menu bar and
// its sub-menus.
func (bar *Bar) Menu(id int) *Menu {
	for i := bar.Count() - 1; i >= 0; i-- {
		menu := bar.MenuAtIndex(i)
		if menu.ID == id {
			return menu
		}
		if item := menu.Item(id); item != nil {
			return item.SubMenu
		}
	}
	return nil
}

// MenuAtIndex returns the menu at the specified index within the menu bar.
func (bar *Bar) MenuAtIndex(index int) *Menu {
	return bar.osMenuAtIndex(index)
}

// MenuItem returns the menu item with the specified id anywhere in the menu
// bar and its sub-menus.
func (bar *Bar) MenuItem(id int) *Item {
	for i := bar.Count() - 1; i >= 0; i-- {
		menu := bar.MenuAtIndex(i)
		if menu.ID == id {
			return &Item{
				Index:   i,
				ID:      id,
				Title:   menu.Title,
				SubMenu: menu,
			}
		}
		if item := menu.Item(id); item != nil {
			return item
		}
	}
	return nil
}

// InsertMenu inserts a menu at the specified item index within this menu bar.
// Pass in a negative index to append to the end.
func (bar *Bar) InsertMenu(atIndex int, menu *Menu) {
	bar.osInsertMenu(atIndex, menu)
}

// IndexOf returns the index of the menu within this menu bar, or -1.
func (bar *Bar) IndexOf(menu *Menu) int {
	for i := bar.Count() - 1; i >= 0; i-- {
		if bar.MenuAtIndex(i) == menu {
			return i
		}
	}
	return -1
}

// Remove the menu at the specified index from this menu bar.
func (bar *Bar) Remove(index int) {
	if index >= 0 && index < bar.Count() {
		bar.osRemoveMenu(index)
	}
}

// Count of menus in this menu bar.
func (bar *Bar) Count() int {
	return bar.osItemCount()
}

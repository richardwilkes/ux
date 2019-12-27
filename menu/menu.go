package menu

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/action"
	"github.com/richardwilkes/ux/ids"
	"github.com/richardwilkes/ux/keys"
)

// Updater is a function called to update a menu before it is shown.
type Updater func(menu *Menu)

// Menu represents a set of menu items.
type Menu struct {
	native osMenu
}

// New creates a new menu. 'updater' is optional. If present, it will be
// called prior to showing the menu, giving a chance to modify the menu.
func New(title string, updater Updater) *Menu {
	return &Menu{native: osNewMenu(title, updater)}
}

// IsSame returns true if the two menus represent the same object. Do not use
// == to test for equality.
func (menu *Menu) IsSame(other *Menu) bool {
	return menu.osIsSame(other)
}

// ItemAtIndex returns the menu item at the specified index within the menu.
func (menu *Menu) ItemAtIndex(index int) *Item {
	return menu.osItemAtIndex(index)
}

// Item returns the menu item with the specified id anywhere in the menu and
// and its sub-menus.
func (menu *Menu) Item(id int) *Item {
	for i := menu.Count() - 1; i >= 0; i-- {
		item := menu.ItemAtIndex(i)
		if item.ID() == id {
			return item
		}
		if sub := item.SubMenu(); sub != nil {
			if item = sub.Item(id); item != nil {
				return item
			}
		}
	}
	return nil
}

// InsertSeparator inserts a menu separator at the specified item index within
// this menu. Pass in a negative index to append to the end.
func (menu *Menu) InsertSeparator(atIndex int) {
	menu.osInsertSeparator(atIndex)
}

// InsertSeparatorIfNeeded inserts a menu separator at the specified item
// index within this menu if the item that would precede it is not a
// separator. Pass in a negative index to append to the end.
func (menu *Menu) InsertSeparatorIfNeeded(atIndex int) {
	if count := menu.Count(); count != 0 {
		if atIndex < 0 {
			atIndex = count
		}
		if atIndex != 0 {
			if !menu.ItemAtIndex(atIndex - 1).IsSeparator() {
				menu.InsertSeparator(atIndex)
			}
		}
	}
}

// InsertActionItem inserts a menu item using the action at the specified item
// index within this menu. Pass in a negative index to append to the end.
func (menu *Menu) InsertActionItem(atIndex int, cmd action.Action) *Item {
	return menu.InsertItem(atIndex, cmd.ID(), cmd.Title(), cmd.HotKey(), cmd.HotKeyModifiers(), func(item *Item) bool { return cmd.Enabled(item) }, func(item *Item) { cmd.Execute(item) })
}

// InsertActionItemForContextMenu inserts a menu item for a context menu using
// the action at the specified item index within this menu. Pass in a negative
// index to append to the end. If the item would be disabled, it is not added
// and nil is returned.
func (menu *Menu) InsertActionItemForContextMenu(atIndex int, cmd action.Action) *Item {
	id := cmd.ID() | ids.ContextMenuIDFlag
	if cmd.Enabled(nil) {
		return menu.InsertItem(atIndex, id, cmd.Title(), nil, 0, func(item *Item) bool { return cmd.Enabled(item) }, func(item *Item) { cmd.Execute(item) })
	}
	return nil
}

// InsertItem inserts a menu item at the specified item index within this
// menu. Pass in a negative index to append to the end. Both 'validator' and
// 'handler' may be nil for default behavior.
func (menu *Menu) InsertItem(atIndex, id int, title string, key *keys.Key, keyModifiers keys.Modifiers, validator ItemValidator, handler ItemHandler) *Item {
	if validator == nil {
		validator = func(*Item) bool { return true }
	}
	if handler == nil {
		handler = func(*Item) {}
	}
	return menu.osInsertItem(atIndex, id, title, key, keyModifiers, validator, handler)
}

// InsertMenu inserts a new sub-menu at the specified item index within this
// menu. Pass in a negative index to append to the end. 'updater' is optional.
// If present, it will be called prior to showing the menu, giving a chance to
// modify the menu.
func (menu *Menu) InsertMenu(atIndex, id int, title string, updater Updater) *Menu {
	return menu.osInsertNewMenu(atIndex, id, title, updater)
}

// RemoveItem removes the menu item at the specified index from this menu.
func (menu *Menu) RemoveItem(index int) {
	if index >= 0 && index < menu.Count() {
		menu.osRemoveItem(index)
	}
}

// Count of menu items in this menu.
func (menu *Menu) Count() int {
	return menu.osItemCount()
}

// Popup the menu at the specified position within the window.
func (menu *Menu) Popup(wnd *ux.Window, where geom.Rect, currentIndex int) {
	if wnd.IsValid() {
		menu.osPopup(wnd, where, currentIndex)
	}
}

// Dispose releases any OS resources associated with this menu.
func (menu *Menu) Dispose() {
	menu.osDispose()
}

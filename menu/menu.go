package menu

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/action"
	"github.com/richardwilkes/ux/ids"
	"github.com/richardwilkes/ux/keys"
)

// Menu represents a set of menu items.
type Menu struct {
	ID    int
	Title string
	menu  osMenu
	valid bool
}

// Item holds information about menu items.
type Item struct {
	Owner   *Menu
	Index   int
	ID      int
	Title   string
	SubMenu *Menu
}

// New creates a new menu. 'updater' is optional. If present, it will be
// called prior to showing the menu, giving a chance to modify the menu.
func New(id int, title string, updater func(*Menu)) *Menu {
	menu := &Menu{
		ID:    id,
		Title: title,
		menu:  osNewMenu(title, updater),
		valid: true,
	}
	menu.osAddNativeMenu()
	return menu
}

// IsValid returns true if the menu is still valid (i.e. hasn't been
// disposed).
func (menu *Menu) IsValid() bool {
	return menu.valid
}

// ItemAtIndex returns the menu item at the specified index within the menu.
func (menu *Menu) ItemAtIndex(index int) *Item {
	if menu.IsValid() {
		return menu.osItemAtIndex(index)
	}
	return nil
}

// Item returns the menu item with the specified id anywhere in the menu and
// and its sub-menus.
func (menu *Menu) Item(id int) *Item {
	for i := menu.Count() - 1; i >= 0; i-- {
		item := menu.ItemAtIndex(i)
		if item.ID == id {
			return item
		}
		if item.SubMenu != nil {
			if item = item.SubMenu.Item(id); item != nil {
				return item
			}
		}
	}
	return nil
}

// SetItemAtIndexTitle sets the title of the menu item at the specified index
// within the menu.
func (menu *Menu) SetItemAtIndexTitle(index int, title string) {
	if menu.IsValid() {
		menu.osSetItemTitle(index, title)
	}
}

// SetItemTitle sets the title of the menu item with the specified id anywhere
// in the menu and its sub-menus.
func (menu *Menu) SetItemTitle(id int, title string) {
	if item := menu.Item(id); item != nil {
		item.Owner.SetItemAtIndexTitle(item.Index, title)
	}
}

// InsertSeparator inserts a menu separator at the specified item index within
// this menu. Pass in a negative index to append to the end.
func (menu *Menu) InsertSeparator(atIndex int) {
	if menu.IsValid() {
		menu.osInsertSeparator(atIndex)
	}
}

// InsertSeparatorIfNeeded inserts a menu separator at the specified item
// index within this menu if the item that would precede it is not a
// separator. Pass in a negative index to append to the end.
func (menu *Menu) InsertSeparatorIfNeeded(atIndex int) {
	if menu.IsValid() {
		if count := menu.osItemCount(); count != 0 {
			if atIndex < 0 {
				atIndex = count
			}
			if atIndex != 0 {
				if menu.osItemAtIndex(atIndex-1).ID != 0 {
					menu.osInsertSeparator(atIndex)
				}
			}
		}
	}
}

// InsertActionItem inserts a menu item using the action at the specified item
// index within this menu. Pass in a negative index to append to the end.
func (menu *Menu) InsertActionItem(atIndex int, cmd action.Action) {
	menu.InsertItem(atIndex, cmd.ID(), cmd.Title(), cmd.HotKey(), cmd.HotKeyModifiers(), cmd.Enabled, cmd.Execute)
}

// InsertActionItemForContextMenu inserts a menu item for a context menu using
// the action at the specified item index within this menu. Pass in a negative
// index to append to the end. If the item would be disabled, it is not added.
func (menu *Menu) InsertActionItemForContextMenu(atIndex int, cmd action.Action) {
	if cmd.Enabled() {
		menu.InsertItem(atIndex, cmd.ID()|ids.ContextMenuIDFlag, cmd.Title(), nil, 0, nil, cmd.Execute)
	}
}

// InsertItem inserts a menu item at the specified item index within this
// menu. Pass in a negative index to append to the end. Both 'validator' and
// 'handler' may be nil for default behavior.
func (menu *Menu) InsertItem(atIndex, id int, title string, key *keys.Key, keyModifiers keys.Modifiers, validator func() bool, handler func()) {
	if menu.IsValid() {
		if validator == nil {
			validator = func() bool { return true }
		}
		if handler == nil {
			handler = func() {}
		}
		menu.osInsertItem(atIndex, id, title, key, keyModifiers, validator, handler)
	}
}

// InsertMenu inserts a new sub-menu at the specified item index within this
// menu. Pass in a negative index to append to the end. 'updater' is optional.
// If present, it will be called prior to showing the menu, giving a chance to
// modify the menu.
func (menu *Menu) InsertMenu(atIndex, id int, title string, updater func(*Menu)) *Menu {
	if menu.IsValid() {
		return menu.osInsertMenu(atIndex, id, title, updater)
	}
	return nil
}

// RemoveItem removes the menu item at the specified index from this menu.
func (menu *Menu) RemoveItem(index int) {
	if menu.IsValid() && index >= 0 && index < menu.Count() {
		menu.osRemoveItem(index)
	}
}

// Count of menu items in this menu.
func (menu *Menu) Count() int {
	if menu.IsValid() {
		return menu.osItemCount()
	}
	return 0
}

// Popup the menu at the specified position within the window.
func (menu *Menu) Popup(wnd *ux.Window, where geom.Rect, currentIndex int) {
	if menu.IsValid() && wnd.IsValid() {
		menu.osPopup(wnd, where, currentIndex)
	}
}

// Dispose releases any OS resources associated with this menu.
func (menu *Menu) Dispose() {
	menu.osRemoveNativeMenu()
	if menu.IsValid() {
		menu.valid = false
		menu.osDispose()
	}
}

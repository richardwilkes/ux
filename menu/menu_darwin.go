package menu

import (
	"github.com/richardwilkes/macos/ns"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/keys"
)

type osMenu = *ns.Menu

func osNewMenu(title string, updater Updater) osMenu {
	var u func(*ns.Menu)
	if updater != nil {
		u = func(nsmenu *ns.Menu) {
			updater(&Menu{native: nsmenu})
		}
	}
	return ns.MenuInitWithTitle(title, u)
}

func (menu *Menu) osIsSame(other *Menu) bool {
	return menu.native.Native() == other.native.Native()
}

func (menu *Menu) osItemAtIndex(index int) *Item {
	if index < 0 || index >= menu.native.NumberOfItems() {
		return nil
	}
	return &Item{native: menu.native.ItemAtIndex(index)}
}

func (menu *Menu) osiInsertItemAtIndex(item *ns.MenuItem, index int) {
	if index < 0 {
		index = menu.native.NumberOfItems()
	}
	menu.native.InsertItemAtIndex(item, index)
}

func (menu *Menu) osInsertSeparator(atIndex int) {
	menu.osiInsertItemAtIndex(ns.MenuSeparatorItem(), atIndex)
}

func (menu *Menu) osInsertItem(atIndex, id int, title string, key *keys.Key, keyModifiers keys.Modifiers, validator ItemValidator, handler ItemHandler) *Item {
	var keyCodeStr string
	if key != nil {
		keyCodeStr = key.RuneStr()
	}
	item := ns.MenuItemInitWithTitleActionKeyEquivalent(id, title, keyCodeStr, int(keyModifiers)<<16, func(item *ns.MenuItem) bool { return validator(&Item{native: item}) }, func(item *ns.MenuItem) { handler(&Item{native: item}) })
	menu.osiInsertItemAtIndex(item, atIndex)
	return &Item{native: item}
}

func (menu *Menu) osInsertNewMenu(atIndex, id int, title string, updater Updater) *Menu {
	item := ns.MenuItemInitWithTitleActionKeyEquivalent(id, title, "", 0, nil, nil)
	subMenu := New(title, updater)
	item.SetSubMenu(subMenu.native)
	menu.osiInsertItemAtIndex(item, atIndex)
	return subMenu
}

func (menu *Menu) osInsertMenu(atIndex, id int, subMenu *Menu) {
	item := ns.MenuItemInitWithTitleActionKeyEquivalent(id, subMenu.native.Title(), "", 0, nil, nil)
	item.SetSubMenu(subMenu.native)
	menu.osiInsertItemAtIndex(item, atIndex)
}

func (menu *Menu) osRemoveItem(index int) {
	menu.native.RemoveItem(index)
}

func (menu *Menu) osItemCount() int {
	return menu.native.NumberOfItems()
}

func (menu *Menu) osPopup(wnd *ux.Window, where geom.Rect, currentIndex int) {
	if item := menu.native.ItemAtIndex(currentIndex); item != nil {
		menu.native.PopupMenuPositioningItemAtLocationInView(item, where.X, where.Y, where.Width, where.Height, wnd.OSWindow().ContentView())
	}
}

func (menu *Menu) osDispose() {
	menu.native.Release()
	menu.native = nil
}

package menu

import (
	"github.com/richardwilkes/macos/ns"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/keys"
)

var nativeMenuMap = make(map[ns.MenuNative]*Menu)

type osMenu = *ns.Menu

func osNewMenu(title string, updater func(*Menu)) osMenu {
	var u func(*ns.Menu)
	if updater != nil {
		u = func(nsmenu *ns.Menu) {
			if m, ok := nativeMenuMap[nsmenu.Native()]; ok {
				updater(m)
			}
		}
	}
	return ns.MenuInitWithTitle(title, u)
}

func (menu *Menu) osAddNativeMenu() {
	nativeMenuMap[menu.menu.Native()] = menu
}

func (menu *Menu) osRemoveNativeMenu() {
	delete(nativeMenuMap, menu.menu.Native())
}

func (menu *Menu) osItemAtIndex(index int) *Item {
	if index < 0 || index >= menu.menu.NumberOfItems() {
		return nil
	}
	item := menu.menu.ItemAtIndex(index)
	var subMenu *Menu
	if sub := item.Submenu(); sub != nil {
		if subMenu = nativeMenuMap[sub.Native()]; subMenu == nil {
			subMenu = &Menu{
				ID:    item.Tag(),
				Title: sub.Title(),
				menu:  sub,
				valid: true,
			}
			subMenu.osAddNativeMenu()
		}
	}
	return &Item{
		Owner:   menu,
		Index:   index,
		ID:      item.Tag(),
		Title:   item.Title(),
		SubMenu: subMenu,
	}
}

func (menu *Menu) osSetItemTitle(index int, title string) {
	if index >= 0 && index < menu.menu.NumberOfItems() {
		menu.menu.ItemAtIndex(index).SetTitle(title)
	}
}

func (menu *Menu) osInsertSeparator(atIndex int) {
	if atIndex < 0 {
		atIndex = menu.menu.NumberOfItems()
	}
	menu.menu.InsertItemAtIndex(ns.MenuSeparatorItem(), atIndex)
}

func (menu *Menu) osInsertItem(atIndex, id int, title string, key *keys.Key, keyModifiers keys.Modifiers, validator func() bool, handler func()) {
	var keyCodeStr string
	if key != nil {
		keyCodeStr = key.RuneStr()
	}
	item := ns.MenuItemInitWithTitleActionKeyEquivalent(id, title, keyCodeStr, int(keyModifiers)<<16, validator, handler)
	if atIndex < 0 {
		atIndex = menu.menu.NumberOfItems()
	}
	menu.menu.InsertItemAtIndex(item, atIndex)
}

func (menu *Menu) osInsertMenu(atIndex, id int, title string, updater func(*Menu)) *Menu {
	item := ns.MenuItemInitWithTitleActionKeyEquivalent(id, title, "", 0, nil, nil)
	subMenu := New(id, title, updater)
	item.SetSubmenu(subMenu.menu)
	if atIndex < 0 {
		atIndex = menu.menu.NumberOfItems()
	}
	menu.menu.InsertItemAtIndex(item, atIndex)
	return subMenu
}

func (menu *Menu) osRemoveItem(index int) {
	menu.menu.RemoveItem(index)
}

func (menu *Menu) osItemCount() int {
	return menu.menu.NumberOfItems()
}

func (menu *Menu) osPopup(wnd *ux.Window, where geom.Rect, currentIndex int) {
	if item := menu.menu.ItemAtIndex(currentIndex); item != nil {
		menu.menu.PopupMenuPositioningItemAtLocationInView(item, where.X, where.Y, where.Width, where.Height, wnd.OSWindow().ContentView())
	}
}

func (menu *Menu) osDispose() {
	menu.menu.Release()
}

func insertMenuNoCreate(menu, subMenu *Menu, atIndex int) {
	item := ns.MenuItemInitWithTitleActionKeyEquivalent(subMenu.ID, subMenu.Title, "", 0, nil, nil)
	item.SetSubmenu(subMenu.menu)
	if atIndex < 0 {
		atIndex = menu.menu.NumberOfItems()
	}
	menu.menu.InsertItemAtIndex(item, atIndex)
}

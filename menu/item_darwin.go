package menu

import "github.com/richardwilkes/macos/ns"

type osItem = *ns.MenuItem

func (item *Item) osIsSame(other *Item) bool {
	return item.native.Native() == other.native.Native()
}

func (item *Item) osMenu() *Menu {
	if menu := item.native.Menu(); menu != nil {
		return &Menu{native: menu}
	}
	return nil
}

func (item *Item) osIsSeparator() bool {
	return item.native.IsSeparator()
}

func (item *Item) osID() int {
	return item.native.Tag()
}

func (item *Item) osTitle() string {
	return item.native.Title()
}

func (item *Item) osSetTitle(title string) {
	item.native.SetTitle(title)
}

func (item *Item) osSubMenu() *Menu {
	if !item.native.HasSubMenu() {
		return nil
	}
	menu := item.native.SubMenu()
	if menu == nil {
		return nil
	}
	return &Menu{native: menu}
}

func (item *Item) osCheckState() CheckState {
	switch item.native.State() {
	case ns.MenuItemStateOn:
		return On
	case ns.MenuItemStateMixed:
		return Mixed
	default:
		return Off
	}
}

func (item *Item) osSetCheckState(state CheckState) {
	switch state {
	case On:
		item.native.SetState(ns.MenuItemStateOn)
	case Mixed:
		item.native.SetState(ns.MenuItemStateMixed)
	default:
		item.native.SetState(ns.MenuItemStateOff)
	}
}

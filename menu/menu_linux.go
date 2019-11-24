package menu

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/keys"
)

type osMenu = int

func osNewMenu(title string, updater func(*Menu)) osMenu {
	// RAW: Implement
	return 0
}

func (menu *Menu) osAddNativeMenu() {
	// RAW: Implement
}

func (menu *Menu) osRemoveNativeMenu() {
	// RAW: Implement
}

func (menu *Menu) osItemAtIndex(index int) *Item {
	// RAW: Implement
	return nil
}

func (menu *Menu) osSetItemTitle(index int, title string) {
	// RAW: Implement
}

func (menu *Menu) osInsertSeparator(atIndex int) {
	// RAW: Implement
}

func (menu *Menu) osInsertItem(atIndex, id int, title string, key *keys.Key, keyModifiers keys.Modifiers, validator func() bool, handler func()) {
	// RAW: Implement
}

func (menu *Menu) osInsertMenu(atIndex, id int, title string, updater func(*Menu)) *Menu {
	// RAW: Implement
	return nil
}

func (menu *Menu) osRemoveItem(index int) {
	// RAW: Implement
}

func (menu *Menu) osItemCount() int {
	// RAW: Implement
	return 0
}

func (menu *Menu) osPopup(wnd *ux.Window, where geom.Point, currentIndex int) {
	// RAW: Implement
}

func (menu *Menu) osDispose() {
	// RAW: Implement
}

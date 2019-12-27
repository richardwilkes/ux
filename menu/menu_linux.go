package menu

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/keys"
)

type osMenu = int

func osNewMenu(title string, updater Updater) osMenu {
	// RAW: Implement
	return 0
}

func (menu *Menu) osIsSame(other *Menu) bool {
	// RAW: Implement
	return menu.native == other.native
}

func (menu *Menu) osItemAtIndex(index int) *Item {
	// RAW: Implement
	return nil
}

func (menu *Menu) osInsertSeparator(atIndex int) {
	// RAW: Implement
}

func (menu *Menu) osInsertItem(atIndex, id int, title string, key *keys.Key, keyModifiers keys.Modifiers, validator ItemValidator, handler ItemHandler) *Item {
	// RAW: Implement
	return nil
}

func (menu *Menu) osInsertNewMenu(atIndex, id int, title string, updater Updater) *Menu {
	// RAW: Implement
	return nil
}

func (menu *Menu) osInsertMenu(atIndex, id int, subMenu *Menu) {
	// RAW: Implement
}

func (menu *Menu) osRemoveItem(index int) {
	// RAW: Implement
}

func (menu *Menu) osItemCount() int {
	// RAW: Implement
	return 0
}

func (menu *Menu) osPopup(wnd *ux.Window, where geom.Rect, currentIndex int) {
	// RAW: Implement
}

func (menu *Menu) osDispose() {
	// RAW: Implement
}

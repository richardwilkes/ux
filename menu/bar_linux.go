package menu

import "github.com/richardwilkes/ux"

func osMenuBarForWindow(wnd *ux.Window, updater Updater) (bar *Bar, isGlobal, isFirst bool) {
	if !wnd.IsValid() {
		return nil, false, false
	}
	// RAW: Implement
	return nil, false, false
}

func osMenuBarHeightInWindow() float64 {
	// RAW: Implement
	return 0
}

func (bar *Bar) osInsertMenu(atIndex, id int, menu *Menu) {
	// RAW: Implement
}

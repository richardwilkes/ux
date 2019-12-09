package menu

import "github.com/richardwilkes/ux"

type osMenuBar = int

func osMenuBarForWindow(wnd *ux.Window, updater func(*Menu)) (bar *Bar, isGlobal, isFirst bool) {
	if !window.IsValid() {
		return nil, false, false
	}
	// RAW: Implement
	return nil, false, false
}

func osMenuBarHeightInWindow() float64 {
	// RAW: Implement
	return 0
}

func (bar *Bar) osMenuAtIndex(index int) *Menu {
	// RAW: Implement
	return nil
}

func (bar *Bar) osInsertMenu(atIndex int, menu *Menu) {
	// RAW: Implement
}

func (bar *Bar) osItemCount() int {
	// RAW: Implement
	return 0
}

func (bar *Bar) osRemoveMenu(index int) {
	// RAW: Implement
}

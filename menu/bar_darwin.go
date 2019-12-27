package menu

import (
	"github.com/richardwilkes/macos/ns"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/ids"
)

var menuBar *Bar

func osMenuBarForWindow(_ *ux.Window, updater Updater) (bar *Bar, isGlobal, isFirst bool) {
	first := false
	if menuBar == nil {
		menuBar = &Bar{bar: New("", updater)}
		ns.SharedApplication().SetMainMenu(menuBar.bar.native)
		first = true
	}
	return menuBar, true, first
}

func osMenuBarHeightInWindow() float64 {
	return 0
}

func (bar *Bar) osInsertMenu(atIndex, id int, menu *Menu) {
	bar.bar.osInsertMenu(atIndex, id, menu)
	switch id {
	case ids.AppMenuID:
		if servicesMenu := bar.Menu(ids.ServicesMenuID); servicesMenu != nil {
			ns.SharedApplication().SetServicesMenu(servicesMenu.native)
		}
	case ids.WindowMenuID:
		ns.SharedApplication().SetWindowsMenu(menu.native)
	case ids.HelpMenuID:
		ns.SharedApplication().SetHelpMenu(menu.native)
	}
}

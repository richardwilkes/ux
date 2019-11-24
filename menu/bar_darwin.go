package menu

import (
	"github.com/richardwilkes/macos/ns"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/ids"
)

var menuBar *Bar

type osMenuBar = *Menu

func osMenuBarForWindow(_ *ux.Window, updater func(*Menu)) (bar *Bar, isGlobal, isFirst bool) {
	first := false
	if menuBar == nil {
		menuBar = &Bar{bar: New(ids.BarID, "", updater)}
		ns.SharedApplication().SetMainMenu(menuBar.bar.menu)
		first = true
	}
	return menuBar, true, first
}

func osMenuBarHeightInWindow() float64 {
	return 0
}

func (bar *Bar) osMenuAtIndex(index int) *Menu {
	if item := menuBar.bar.ItemAtIndex(index); item != nil {
		return item.SubMenu
	}
	return nil
}

func (bar *Bar) osInsertMenu(atIndex int, menu *Menu) {
	if bar.bar.IsValid() && menu.IsValid() {
		insertMenuNoCreate(bar.bar, menu, atIndex)
		switch menu.ID {
		case ids.AppMenuID:
			if servicesMenu := bar.Menu(ids.ServicesMenuID); servicesMenu != nil {
				ns.SharedApplication().SetServicesMenu(servicesMenu.menu)
			}
		case ids.WindowMenuID:
			ns.SharedApplication().SetWindowsMenu(menu.menu)
		case ids.HelpMenuID:
			ns.SharedApplication().SetHelpMenu(menu.menu)
		}
	}
}

func (bar *Bar) osItemCount() int {
	return bar.bar.Count()
}

func (bar *Bar) osRemoveMenu(index int) {
	bar.bar.RemoveItem(index)
}

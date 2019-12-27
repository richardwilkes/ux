package menu

import (
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/win32"
)

var menuBarMap = make(map[win32.HMENU]*Bar)

type menuBarData struct {
	menu               win32.HMENU
	wnd                win32.HWND
	menuKeys           map[string]*osItem
	needMenuKeyRefresh bool
	needRedraw         bool
}

func osMenuBarForWindow(wnd *ux.Window, updater Updater) (bar *Bar, isGlobal, isFirst bool) {
	if !wnd.IsValid() {
		return nil, false, false
	}
	if ux.MenuItemSelectionCallback == nil {
		ux.MenuItemSelectionCallback = osiHandleMenuItemSelection
		ux.MenuValidationCallback = osiValidateMenu
	}
	// RAW: Implement
	// hwnd := wnd.OSWindow()
	// hmenu := win32.GetMenu(hwnd)
	// if hmenu == win32.NULL {
	// 	if hmenu = win32.CreateMenu(); hmenu != win32.NULL {
	// 		win32.SetMenu(hwnd, hmenu)
	// 		bar = &Bar{
	// 			bar: &menuBarData{
	// 				menu:     hmenu,
	// 				wnd:      hwnd,
	// 				menuKeys: make(map[string]*osItem),
	// 			},
	// 		}
	// 		menuBarMap[hmenu] = bar
	// 		bar.osiMarkForUpdate()
	// 		return bar, false, true
	// 	}
	// }
	// return menuBarMap[hmenu], false, false
	return nil, false, false
}

func osMenuBarHeightInWindow() float64 {
	return float64(win32.GetSystemMetrics(win32.SM_CYMENU))
}

func (bar *Bar) osInsertMenu(atIndex, id int, menu *Menu) {
	// RAW: Implement
	// win32.InsertMenuItem(bar.bar.menu, uint32(atIndex), true, &win32.MENUITEMINFO{
	// 	Size:     uint32(unsafe.Sizeof(win32.MENUITEMINFO{})), //nolint:gosec
	// 	Mask:     win32.MIIM_ID | win32.MIIM_FTYPE | win32.MIIM_STRING | win32.MIIM_SUBMENU,
	// 	Type:     win32.MFT_STRING,
	// 	ID:       uint32(id),
	// 	TypeData: win32.ToSysWin32Str(menu.Title, false),
	// 	SubMenu:  menu.menu,
	// })
	// bar.osiMarkForUpdate()
}

// -- From here down are specific to Windows

func (bar *Bar) osMenuAtIndex(index int) *Menu {
	// RAW: Implement
	return nil
	// var data [512]uint16
	// info := &win32.MENUITEMINFO{
	// 	Size:     uint32(unsafe.Sizeof(win32.MENUITEMINFO{})), //nolint:gosec
	// 	Mask:     win32.MIIM_ID | win32.MIIM_FTYPE | win32.MIIM_STRING | win32.MIIM_SUBMENU,
	// 	TypeData: uintptr(unsafe.Pointer(&data[0])), //nolint:gosec
	// 	CCH:      uint32(len(data) - 1),
	// }
	// if !win32.GetMenuItemInfo(bar.bar.menu, uint32(index), true, info) {
	// 	return nil
	// }
	// return nativeMenuMap[info.SubMenu]
}

func (bar *Bar) osItemCount() int {
	// RAW: Implement
	return 0
	// return win32.GetMenuItemCount(bar.bar.menu)
}

func (bar *Bar) osRemoveMenu(index int) {
	// RAW: Implement
	// win32.DeleteMenu(bar.bar.menu, uint32(index), win32.MF_BYPOSITION)
	// bar.osiMarkForUpdate()
}

func (bar *Bar) osiMarkForUpdate() {
	// RAW: Implement
	// bar.bar.needMenuKeyRefresh = true
	// if !bar.bar.needRedraw {
	// 	bar.bar.needRedraw = true
	// 	ux.Invoke(func() {
	// 		bar.bar.needRedraw = false
	// 		win32.DrawMenuBar(bar.bar.wnd)
	// 	})
	// }
}

func (bar *Bar) osiRefreshMenuKeysForMenu(menu *Menu) {
	// RAW: Implement
	// if menu != nil {
	// 	for i := menu.Count() - 1; i >= 0; i-- {
	// 		mi := menu.ItemAtIndex(i)
	// 		if mi.SubMenu != nil {
	// 			bar.osiRefreshMenuKeysForMenu(mi.SubMenu)
	// 		} else if data, exists := menuItemDataMap[mi.ID]; exists && data.key != nil {
	// 			bar.bar.menuKeys[data.modifiers.String()+data.key.Name] = data
	// 		}
	// 	}
	// }
}

func osiMarkAllForMenuKeyRefresh() {
	// RAW: Implement
	// for _, bar := range menuBarMap {
	// 	bar.bar.needMenuKeyRefresh = true
	// }
}

func osiRefreshMenuKeyForWindow(wnd *ux.Window) map[string]*osItem {
	// RAW: Implement
	// if bar := osiMenuBarForWindowNoCreate(wnd); bar != nil {
	// 	if bar.bar.needMenuKeyRefresh {
	// 		bar.bar.needMenuKeyRefresh = false
	// 		bar.bar.menuKeys = make(map[string]*osItem)
	// 		for i := bar.Count() - 1; i >= 0; i-- {
	// 			bar.osiRefreshMenuKeysForMenu(bar.MenuAtIndex(i))
	// 		}
	// 	}
	// 	return bar.bar.menuKeys
	// }
	return nil
}

func osiMenuBarForWindowNoCreate(wnd *ux.Window) *Bar {
	// RAW: Implement
	// if wnd != nil && wnd.IsValid() {
	// 	if hmenu := win32.GetMenu(wnd.OSWindow()); hmenu != win32.NULL {
	// 		return menuBarMap[hmenu]
	// 	}
	// }
	return nil
}

package menu

import (
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/win32"
	"unsafe"
)

var menuBarMap = make(map[win32.HMENU]*Bar)

type menuBarData struct {
	menu               win32.HMENU
	wnd                win32.HWND
	menuKeys           map[string]*menuItemData
	needMenuKeyRefresh bool
	needRedraw         bool
}

type osMenuBar = *menuBarData

func osMenuBarForWindow(wnd *ux.Window, updater func(*Menu)) (bar *Bar, isGlobal, isFirst bool) {
	if !window.IsValid() {
		return nil, false, false
	}
	if ux.MenuItemSelectionCallback == nil {
		ux.MenuItemSelectionCallback = handleMenuItemSelection
		ux.MenuValidationCallback = validateMenu
	}
	hwnd := wnd.OSWindow()
	hmenu := win32.GetMenu(hwnd)
	if hmenu == win32.NULL {
		if hmenu = win32.CreateMenu(); hmenu != win32.NULL {
			win32.SetMenu(hwnd, hmenu)
			bar = &Bar{
				bar: &menuBarData{
					menu:     hmenu,
					wnd:      hwnd,
					menuKeys: make(map[string]*menuItemData),
				},
			}
			menuBarMap[hmenu] = bar
			bar.markForUpdate()
			return bar, false, true
		}
	}
	return menuBarMap[hmenu], false, false
}

func osMenuBarHeightInWindow() float64 {
	return float64(win32.GetSystemMetrics(win32.SM_CYMENU))
}

func (bar *Bar) osMenuAtIndex(index int) *Menu {
	var data [512]uint16
	info := &win32.MENUITEMINFO{
		Size:     uint32(unsafe.Sizeof(win32.MENUITEMINFO{})), //nolint:gosec
		Mask:     win32.MIIM_ID | win32.MIIM_FTYPE | win32.MIIM_STRING | win32.MIIM_SUBMENU,
		TypeData: uintptr(unsafe.Pointer(&data[0])), //nolint:gosec
		CCH:      uint32(len(data) - 1),
	}
	if !win32.GetMenuItemInfo(bar.bar.menu, uint32(index), true, info) {
		return nil
	}
	return nativeMenuMap[info.SubMenu]
}

func (bar *Bar) osInsertMenu(atIndex int, menu *Menu) {
	win32.InsertMenuItem(bar.bar.menu, uint32(atIndex), true, &win32.MENUITEMINFO{
		Size:     uint32(unsafe.Sizeof(win32.MENUITEMINFO{})), //nolint:gosec
		Mask:     win32.MIIM_ID | win32.MIIM_FTYPE | win32.MIIM_STRING | win32.MIIM_SUBMENU,
		Type:     win32.MFT_STRING,
		ID:       uint32(menu.ID),
		TypeData: win32.ToSysWin32Str(menu.Title, false),
		SubMenu:  menu.menu,
	})
	bar.markForUpdate()
}

func (bar *Bar) osItemCount() int {
	return win32.GetMenuItemCount(bar.bar.menu)
}

func (bar *Bar) osRemoveMenu(index int) {
	win32.DeleteMenu(bar.bar.menu, uint32(index), win32.MF_BYPOSITION)
	bar.markForUpdate()
}

func (bar *Bar) markForUpdate() {
	bar.bar.needMenuKeyRefresh = true
	if !bar.bar.needRedraw {
		bar.bar.needRedraw = true
		ux.Invoke(func() {
			bar.bar.needRedraw = false
			win32.DrawMenuBar(bar.bar.wnd)
		})
	}
}

func (bar *Bar) refreshMenuKeysForMenu(menu *Menu) {
	if menu != nil {
		for i := menu.Count() - 1; i >= 0; i-- {
			mi := menu.ItemAtIndex(i)
			if mi.SubMenu != nil {
				bar.refreshMenuKeysForMenu(mi.SubMenu)
			} else if data, exists := menuItemDataMap[mi.ID]; exists && data.key != nil {
				bar.bar.menuKeys[data.modifiers.String()+data.key.Name] = data
			}
		}
	}
}

func markAllForMenuKeyRefresh() {
	for _, bar := range menuBarMap {
		bar.bar.needMenuKeyRefresh = true
	}
}

func refreshMenuKeyForWindow(wnd *ux.Window) map[string]*menuItemData {
	if bar := menuBarForWindowNoCreate(wnd); bar != nil {
		if bar.bar.needMenuKeyRefresh {
			bar.bar.needMenuKeyRefresh = false
			bar.bar.menuKeys = make(map[string]*menuItemData)
			for i := bar.Count() - 1; i >= 0; i-- {
				bar.refreshMenuKeysForMenu(bar.MenuAtIndex(i))
			}
		}
		return bar.bar.menuKeys
	}
	return nil
}

func menuBarForWindowNoCreate(wnd *ux.Window) *Bar {
	if wnd != nil && wnd.IsValid() {
		if hmenu := win32.GetMenu(wnd.OSWindow()); hmenu != win32.NULL {
			return menuBarMap[hmenu]
		}
	}
	return nil
}

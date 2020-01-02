// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package menu

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/keys"
	"github.com/richardwilkes/win32"
)

var (
	nativeMenuMap   = make(map[win32.HMENU]*Menu)
	menuItemDataMap = make(map[int]*osItem)
)

type osMenu = win32.HMENU

func osNewMenu(title string, updater Updater) osMenu {
	// RAW: Implement
	// return win32.CreatePopupMenu()
	return 0
}

func (menu *Menu) osIsSame(other *Menu) bool {
	// RAW: Implement
	return menu.native == other.native
}

func (menu *Menu) osItemAtIndex(index int) *Item {
	// RAW: Implement
	// var data [512]uint16
	// info := &win32.MENUITEMINFO{
	// 	Size:     uint32(unsafe.Sizeof(win32.MENUITEMINFO{})), //nolint:gosec
	// 	Mask:     win32.MIIM_ID | win32.MIIM_FTYPE | win32.MIIM_STRING | win32.MIIM_SUBMENU,
	// 	TypeData: uintptr(unsafe.Pointer(&data[0])), //nolint:gosec
	// 	CCH:      uint32(len(data) - 1),
	// }
	// if !win32.GetMenuItemInfo(menu.menu, uint32(index), true, info) {
	// 	return nil
	// }
	// mi := &Item{
	// 	Owner: menu,
	// 	Index: index,
	// 	ID:    int(info.ID),
	// 	// TODO: need check state
	// }
	// if info.Type == win32.MFT_STRING {
	// 	mi.Title = strings.SplitN(syscall.UTF16ToString(data[:info.CCH]), "\t", 2)[0] // Remove any key accelerator info
	// 	mi.SubMenu = nativeMenuMap[info.SubMenu]
	// }
	// return mi
	return nil
}

func (menu *Menu) osInsertSeparator(atIndex int) {
	// RAW: Implement
	// win32.InsertMenuItem(menu.menu, uint32(atIndex), true, &win32.MENUITEMINFO{
	// 	Size: uint32(unsafe.Sizeof(win32.MENUITEMINFO{})), //nolint:gosec
	// 	Mask: win32.MIIM_FTYPE,
	// 	Type: win32.MFT_SEPARATOR,
	// })
}

func (menu *Menu) osInsertItem(atIndex, id int, title string, key *keys.Key, keyModifiers keys.Modifiers, validator ItemValidator, handler ItemHandler) *Item {
	// RAW: Implement
	// title = strings.SplitN(title, "\t", 2)[0] // Remove any pre-existing key accelerator info
	// if key != nil {
	// 	title += "\t" + keyModifiers.String() + key.Name
	// }
	// win32.InsertMenuItem(menu.menu, uint32(atIndex), true, &win32.MENUITEMINFO{
	// 	Size:     uint32(unsafe.Sizeof(win32.MENUITEMINFO{})), //nolint:gosec
	// 	Mask:     win32.MIIM_ID | win32.MIIM_FTYPE | win32.MIIM_STRING,
	// 	Type:     win32.MFT_STRING,
	// 	ID:       uint32(id),
	// 	TypeData: win32.ToSysWin32Str(title, false),
	// })
	// menuItemDataMap[id] = &osItem{
	// 	validator: validator,
	// 	handler:   handler,
	// 	key:       key,
	// 	modifiers: keyModifiers,
	// }
	// osiMarkAllForMenuKeyRefresh()
	return nil
}

func (menu *Menu) osInsertNewMenu(atIndex, id int, title string, updater Updater) *Menu {
	// RAW: Implement
	// subMenu := New(id, title, updater)
	// win32.InsertMenuItem(menu.menu, uint32(atIndex), true, &win32.MENUITEMINFO{
	// 	Size:     uint32(unsafe.Sizeof(win32.MENUITEMINFO{})), //nolint:gosec
	// 	Mask:     win32.MIIM_ID | win32.MIIM_FTYPE | win32.MIIM_STRING | win32.MIIM_SUBMENU,
	// 	Type:     win32.MFT_STRING,
	// 	ID:       uint32(id),
	// 	TypeData: win32.ToSysWin32Str(title, false),
	// 	SubMenu:  subMenu.native,
	// })
	// osiMarkAllForMenuKeyRefresh()
	// return subMenu
	return nil
}

func (menu *Menu) osInsertMenu(atIndex, id int, subMenu *Menu) {
	// RAW: Implement
	// win32.InsertMenuItem(bar.bar.menu, uint32(atIndex), true, &win32.MENUITEMINFO{
	// 	Size:     uint32(unsafe.Sizeof(win32.MENUITEMINFO{})), //nolint:gosec
	// 	Mask:     win32.MIIM_ID | win32.MIIM_FTYPE | win32.MIIM_STRING | win32.MIIM_SUBMENU,
	// 	Type:     win32.MFT_STRING,
	// 	ID:       uint32(id),
	// 	TypeData: win32.ToSysWin32Str(subMenu.Title(), false),
	// 	SubMenu:  subMenu.native,
	// })
	// bar.osiMarkForUpdate()
}

func (menu *Menu) osRemoveItem(index int) {
	// RAW: Implement
	// win32.DeleteMenu(menu.menu, uint32(index), win32.MF_BYPOSITION)
	// osiMarkAllForMenuKeyRefresh()
}

func (menu *Menu) osItemCount() int {
	// RAW: Implement
	// return win32.GetMenuItemCount(menu.menu)
	return 0
}

func (menu *Menu) osPopup(wnd *ux.Window, where geom.Rect, currentIndex int) {
	// RAW: Implement
}

func (menu *Menu) osDispose() {
	// RAW: Implement
	// win32.DestroyMenu(menu.menu)
	// osiMarkAllForMenuKeyRefresh()
}

func osiHandleMenuItemSelection(id int) {
	// RAW: Implement
	// if data, ok := menuItemDataMap[id]; ok {
	// 	if data.handler != nil {
	// 		data.handler(id)
	// 	}
	// }
}

func osiValidateMenu(hmenu win32.HMENU) {
	// RAW: Implement
	// if menu, ok := nativeMenuMap[hmenu]; ok && menu.IsValid() {
	// 	for i := menu.Count() - 1; i >= 0; i-- {
	// 		state := win32.MF_ENABLED
	// 		if item := menu.ItemAtIndex(i); item.ID != 0 {
	// 			if info, exists := menuItemDataMap[item.ID]; exists && info.validator != nil {
	// 				if !info.validator(item.ID) {
	// 					state = win32.MF_DISABLED
	// 				}
	// 			}
	// 		}
	// 		win32.EnableMenuItem(hmenu, i, state|win32.MF_BYPOSITION)
	// 	}
	// }
}

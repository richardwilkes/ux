package menu

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/keys"
	"github.com/richardwilkes/win32"
	"strings"
	"syscall"
	"unsafe"
)

var (
	nativeMenuMap   = make(map[win32.HMENU]*Menu)
	menuItemDataMap = make(map[int]*menuItemData)
)

type osMenu = win32.HMENU

type menuItemData struct {
	validator func() bool
	handler   func()
	key       *keys.Key
	modifiers keys.Modifiers
}

func osNewMenu(title string, updater func(*Menu)) osMenu {
	return win32.CreatePopupMenu()
}

func (menu *Menu) osAddNativeMenu() {
	nativeMenuMap[menu.menu] = menu
}

func (menu *Menu) osRemoveNativeMenu() {
	delete(nativeMenuMap, menu.menu)
}

func (menu *Menu) osItemAtIndex(index int) *Item {
	var data [512]uint16
	info := &win32.MENUITEMINFO{
		Size:     uint32(unsafe.Sizeof(win32.MENUITEMINFO{})), //nolint:gosec
		Mask:     win32.MIIM_ID | win32.MIIM_FTYPE | win32.MIIM_STRING | win32.MIIM_SUBMENU,
		TypeData: uintptr(unsafe.Pointer(&data[0])), //nolint:gosec
		CCH:      uint32(len(data) - 1),
	}
	if !win32.GetMenuItemInfo(menu.menu, uint32(index), true, info) {
		return nil
	}
	mi := &Item{
		Owner: menu,
		Index: index,
		ID:    int(info.ID),
	}
	if info.Type == win32.MFT_STRING {
		mi.Title = strings.SplitN(syscall.UTF16ToString(data[:info.CCH]), "\t", 2)[0] // Remove any key accelerator info
		mi.SubMenu = nativeMenuMap[info.SubMenu]
	}
	return mi
}

func (menu *Menu) osSetItemTitle(index int, title string) {
	win32.SetMenuItemInfo(menu.menu, uint32(index), true, &win32.MENUITEMINFO{
		Size:     uint32(unsafe.Sizeof(win32.MENUITEMINFO{})), //nolint:gosec
		Mask:     win32.MIIM_STRING,
		TypeData: win32.ToSysWin32Str(title, false),
	})
}

func (menu *Menu) osInsertSeparator(atIndex int) {
	win32.InsertMenuItem(menu.menu, uint32(atIndex), true, &win32.MENUITEMINFO{
		Size: uint32(unsafe.Sizeof(win32.MENUITEMINFO{})), //nolint:gosec
		Mask: win32.MIIM_FTYPE,
		Type: win32.MFT_SEPARATOR,
	})
}

func (menu *Menu) osInsertItem(atIndex, id int, title string, key *keys.Key, keyModifiers keys.Modifiers, validator func() bool, handler func()) {
	title = strings.SplitN(title, "\t", 2)[0] // Remove any pre-existing key accelerator info
	if key != nil {
		title += "\t" + keyModifiers.String() + key.Name
	}
	win32.InsertMenuItem(menu.menu, uint32(atIndex), true, &win32.MENUITEMINFO{
		Size:     uint32(unsafe.Sizeof(win32.MENUITEMINFO{})), //nolint:gosec
		Mask:     win32.MIIM_ID | win32.MIIM_FTYPE | win32.MIIM_STRING,
		Type:     win32.MFT_STRING,
		ID:       uint32(id),
		TypeData: win32.ToSysWin32Str(title, false),
	})
	menuItemDataMap[id] = &menuItemData{
		validator: validator,
		handler:   handler,
		key:       key,
		modifiers: keyModifiers,
	}
	markAllForMenuKeyRefresh()
}

func (menu *Menu) osInsertMenu(atIndex, id int, title string, updater func(*Menu)) *Menu {
	subMenu := New(id, title, updater)
	minfo := &win32.MENUITEMINFO{
		Size:     uint32(unsafe.Sizeof(win32.MENUITEMINFO{})), //nolint:gosec
		Mask:     win32.MIIM_ID | win32.MIIM_FTYPE | win32.MIIM_STRING | win32.MIIM_SUBMENU,
		Type:     win32.MFT_STRING,
		ID:       uint32(subMenu.ID),
		TypeData: win32.ToSysWin32Str(subMenu.Title, false),
		SubMenu:  subMenu.menu,
	}
	win32.InsertMenuItem(menu.menu, uint32(atIndex), true, minfo)
	markAllForMenuKeyRefresh()
	return subMenu
}

func (menu *Menu) osRemoveItem(index int) {
	win32.DeleteMenu(menu.menu, uint32(index), win32.MF_BYPOSITION)
	markAllForMenuKeyRefresh()
}

func (menu *Menu) osItemCount() int {
	return win32.GetMenuItemCount(menu.menu)
}

func (menu *Menu) osPopup(wnd *ux.Window, where geom.Point, currentIndex int) {
	// RAW: Implement
}

func (menu *Menu) osDispose() {
	win32.DestroyMenu(menu.menu)
	markAllForMenuKeyRefresh()
}

func handleMenuItemSelection(id int) {
	if data, ok := menuItemDataMap[id]; ok {
		if data.handler != nil {
			data.handler()
		}
	}
}

func validateMenu(hmenu win32.HMENU) {
	if menu, ok := nativeMenuMap[hmenu]; ok && menu.IsValid() {
		for i := menu.Count() - 1; i >= 0; i-- {
			state := win32.MF_ENABLED
			if item := menu.ItemAtIndex(i); item.ID != 0 {
				if info, exists := menuItemDataMap[item.ID]; exists && info.validator != nil {
					if !info.validator() {
						state = win32.MF_DISABLED
					}
				}
			}
			win32.EnableMenuItem(hmenu, i, state|win32.MF_BYPOSITION)
		}
	}
}

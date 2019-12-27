package menu

import "github.com/richardwilkes/ux/keys"

type osItemData struct {
	validator ItemValidator
	handler   ItemHandler
	key       *keys.Key
	modifiers keys.Modifiers
}

type osItem = *osItemData

func (item *Item) osIsSame(other *Item) bool {
	// RAW: Implement
	return item.native == other.native
}

func (item *Item) osMenu() *Menu {
	// RAW: Implement
	return nil
}

func (item *Item) osIsSeparator() bool {
	// RAW: Implement
	return false
}

func (item *Item) osID() int {
	// RAW: Implement
	return 0
}

func (item *Item) osTitle() string {
	// RAW: Implement
	return ""
}

func (item *Item) osSetTitle(title string) {
	// RAW: Implement
	// win32.SetMenuItemInfo(menu.menu, uint32(index), true, &win32.MENUITEMINFO{
	// 	Size:     uint32(unsafe.Sizeof(win32.MENUITEMINFO{})), //nolint:gosec
	// 	Mask:     win32.MIIM_STRING,
	// 	TypeData: win32.ToSysWin32Str(title, false),
	// })
}

func (item *Item) osSubMenu() *Menu {
	// RAW: Implement
	return nil
}

func (item *Item) osCheckState() CheckState {
	// RAW: Implement
	return Off
}

func (item *Item) osSetCheckState(state CheckState) {
	// RAW: Implement
}

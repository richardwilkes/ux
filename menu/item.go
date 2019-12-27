package menu

// CheckState holds a menu item's check state.
type CheckState uint8

// Possible menu item check states.
const (
	Off CheckState = iota
	On
	Mixed
)

// ItemValidator is a function called to validate a menu item.
type ItemValidator func(item *Item) bool

// ItemHandler is a function called to handle a menu item that was selected.
type ItemHandler func(item *Item)

// Item holds information about menu items.
type Item struct {
	native osItem
}

// IsSame returns true if the two items represent the same object. Do not use
// == to test for equality.
func (item *Item) IsSame(other *Item) bool {
	return item.osIsSame(other)
}

// Menu returns the owning menu.
func (item *Item) Menu() *Menu {
	return item.osMenu()
}

// ID returns the id of this item.
func (item *Item) ID() int {
	return item.osID()
}

// Index returns the index of the item within its menu. Returns -1 if it is
// not yet attached to a menu.
func (item *Item) Index() int {
	if menu := item.Menu(); menu != nil {
		count := menu.Count()
		for i := 0; i < count; i++ {
			other := menu.ItemAtIndex(i)
			if item.IsSame(other) {
				return i
			}
		}
	}
	return -1
}

// IsSeparator returns true if this item is a separator.
func (item *Item) IsSeparator() bool {
	return item.osIsSeparator()
}

// Title returns the item's title.
func (item *Item) Title() string {
	return item.osTitle()
}

// SetTitle sets the item's title.
func (item *Item) SetTitle(title string) {
	item.osSetTitle(title)
}

// SubMenu returns the item's sub-menu, if any.
func (item *Item) SubMenu() *Menu {
	return item.osSubMenu()
}

// CheckState returns the item's current check state.
func (item *Item) CheckState() CheckState {
	return item.osCheckState()
}

// SetCheckState sets the item's check state.
func (item *Item) SetCheckState(state CheckState) {
	item.osSetCheckState(state)
}

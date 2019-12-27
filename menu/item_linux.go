package menu

type osItem = int

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

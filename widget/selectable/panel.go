package selectable

import "github.com/richardwilkes/ux"

// Panel wraps a ux.Panel for use in a Group.
type Panel struct {
	ux.Panel
	group    *Group
	selected bool
}

// AsSelectable returns the object as a selectable panel.
func (p *Panel) AsSelectable() *Panel {
	return p
}

// Selected returns true if the panel is currently selected.
func (p *Panel) Selected() bool {
	return p.selected
}

// SetSelected sets the panel's selected state.
func (p *Panel) SetSelected(selected bool) {
	if p.group != nil {
		p.group.Select(p)
	} else {
		p.setSelected(selected)
	}
}

func (p *Panel) setSelected(selected bool) {
	if p.selected != selected {
		p.selected = selected
		p.MarkForRedraw()
	}
}

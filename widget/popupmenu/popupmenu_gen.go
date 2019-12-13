// Code created from "widget.go.tmpl" - don't edit by hand

package popupmenu

import (
	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/draw"
)

type managed struct {
	font                 *draw.Font
	backgroundInk        draw.Ink
	focusedBackgroundInk draw.Ink
	pressedBackgroundInk draw.Ink
	edgeInk              draw.Ink
	textInk              draw.Ink
	pressedTextInk       draw.Ink
	cornerRadius         float64
	hMargin              float64
	vMargin              float64
}

func (m *managed) initialize() {
	m.font = draw.SystemFont
	m.backgroundInk = draw.ControlBackgroundInk
	m.focusedBackgroundInk = draw.ControlFocusedBackgroundInk
	m.pressedBackgroundInk = draw.ControlPressedBackgroundInk
	m.edgeInk = draw.ControlEdgeAdjColor
	m.textInk = draw.ControlTextColor
	m.pressedTextInk = draw.AlternateSelectedControlTextColor
	m.cornerRadius = 4
	m.hMargin = 8
	m.vMargin = 1
}

// Font returns the font that will be used when drawing text content.
func (p *PopupMenu) Font() *draw.Font {
	return p.font
}

// SetFont sets the font that will be used when drawing text content. Pass in
// nil to use the default.
func (p *PopupMenu) SetFont(value *draw.Font) *PopupMenu {
	if value == nil {
		value = draw.SystemFont
	}
	if p.font != value {
		p.font = value
		p.MarkForLayoutAndRedraw()
	}
	return p
}

// BackgroundInk returns the ink that will be used for the background when
// enabled but not pressed or focused.
func (p *PopupMenu) BackgroundInk() draw.Ink {
	return p.backgroundInk
}

// SetBackgroundInk sets the ink that will be used for the background when
// enabled but not pressed or focused. Pass in nil to use the default.
func (p *PopupMenu) SetBackgroundInk(value draw.Ink) *PopupMenu {
	if value == nil {
		value = draw.ControlBackgroundInk
	}
	if p.backgroundInk != value {
		p.backgroundInk = value
		p.MarkForRedraw()
	}
	return p
}

// FocusedBackgroundInk returns the ink that will be used for the background
// when enabled and focused.
func (p *PopupMenu) FocusedBackgroundInk() draw.Ink {
	return p.focusedBackgroundInk
}

// SetFocusedBackgroundInk sets the ink that will be used for the background
// when enabled and focused. Pass in nil to use the default.
func (p *PopupMenu) SetFocusedBackgroundInk(value draw.Ink) *PopupMenu {
	if value == nil {
		value = draw.ControlFocusedBackgroundInk
	}
	if p.focusedBackgroundInk != value {
		p.focusedBackgroundInk = value
		p.MarkForRedraw()
	}
	return p
}

// PressedBackgroundInk returns the ink that will be used for the background
// when enabled and pressed.
func (p *PopupMenu) PressedBackgroundInk() draw.Ink {
	return p.pressedBackgroundInk
}

// SetPressedBackgroundInk sets the ink that will be used for the background
// when enabled and pressed. Pass in nil to use the default.
func (p *PopupMenu) SetPressedBackgroundInk(value draw.Ink) *PopupMenu {
	if value == nil {
		value = draw.ControlPressedBackgroundInk
	}
	if p.pressedBackgroundInk != value {
		p.pressedBackgroundInk = value
		p.MarkForRedraw()
	}
	return p
}

// EdgeInk returns the ink that will be used for the edges.
func (p *PopupMenu) EdgeInk() draw.Ink {
	return p.edgeInk
}

// SetEdgeInk sets the ink that will be used for the edges. Pass in nil to
// use the default.
func (p *PopupMenu) SetEdgeInk(value draw.Ink) *PopupMenu {
	if value == nil {
		value = draw.ControlEdgeAdjColor
	}
	if p.edgeInk != value {
		p.edgeInk = value
		p.MarkForRedraw()
	}
	return p
}

// TextInk returns the ink that will be used for the text when disabled or
// not pressed.
func (p *PopupMenu) TextInk() draw.Ink {
	return p.textInk
}

// SetTextInk sets the ink that will be used for the text when disabled or
// not pressed. Pass in nil to use the default.
func (p *PopupMenu) SetTextInk(value draw.Ink) *PopupMenu {
	if value == nil {
		value = draw.ControlTextColor
	}
	if p.textInk != value {
		p.textInk = value
		p.MarkForRedraw()
	}
	return p
}

// PressedTextInk returns the ink that will be used for the text when enabled
// and pressed.
func (p *PopupMenu) PressedTextInk() draw.Ink {
	return p.pressedTextInk
}

// SetPressedTextInk sets the ink that will be used for the text when enabled
// and pressed. Pass in nil to use the default.
func (p *PopupMenu) SetPressedTextInk(value draw.Ink) *PopupMenu {
	if value == nil {
		value = draw.AlternateSelectedControlTextColor
	}
	if p.pressedTextInk != value {
		p.pressedTextInk = value
		p.MarkForRedraw()
	}
	return p
}

// CornerRadius returns the amount of rounding to use on the corners.
func (p *PopupMenu) CornerRadius() float64 {
	return p.cornerRadius
}

// SetCornerRadius sets the amount of rounding to use on the corners.
func (p *PopupMenu) SetCornerRadius(value float64) *PopupMenu {
	if value < 0 {
		value = 0
	}
	if p.cornerRadius != value {
		p.cornerRadius = value
		p.MarkForRedraw()
	}
	return p
}

// HMargin returns the margin on the left and right side of the content.
func (p *PopupMenu) HMargin() float64 {
	return p.hMargin
}

// SetHMargin sets the margin on the left and right side of the content.
func (p *PopupMenu) SetHMargin(value float64) *PopupMenu {
	if value < 0 {
		value = 0
	}
	if p.hMargin != value {
		p.hMargin = value
		p.MarkForRedraw()
	}
	return p
}

// VMargin returns the margin on the top and bottom side of the content.
func (p *PopupMenu) VMargin() float64 {
	return p.vMargin
}

// SetVMargin sets the margin on the top and bottom side of the content.
func (p *PopupMenu) SetVMargin(value float64) *PopupMenu {
	if value < 0 {
		value = 0
	}
	if p.vMargin != value {
		p.vMargin = value
		p.MarkForRedraw()
	}
	return p
}

// SetBorder sets the border. May be nil.
func (p *PopupMenu) SetBorder(value border.Border) *PopupMenu {
	p.Panel.SetBorder(value)
	return p
}

// SetEnabled sets enabled state.
func (p *PopupMenu) SetEnabled(enabled bool) *PopupMenu {
	p.Panel.SetEnabled(enabled)
	return p
}

// SetFocusable whether it can have the keyboard focus.
func (p *PopupMenu) SetFocusable(focusable bool) *PopupMenu {
	p.Panel.SetFocusable(focusable)
	return p
}

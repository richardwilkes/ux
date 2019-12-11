// Code created from "widget.go.tmpl" - don't edit by hand

package checkbox

import (
	"time"

	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/layout/align"
	"github.com/richardwilkes/ux/layout/side"
)

type managed struct {
	image                *draw.Image //nolint:structcheck
	text                 string      //nolint:structcheck
	font                 *draw.Font
	backgroundInk        draw.Ink
	focusedBackgroundInk draw.Ink
	pressedBackgroundInk draw.Ink
	edgeInk              draw.Ink
	textInk              draw.Ink
	pressedTextInk       draw.Ink
	gap                  float64
	cornerRadius         float64
	clickAnimationTime   time.Duration
	hAlign               align.Alignment
	vAlign               align.Alignment
	side                 side.Side
}

func (m *managed) initialize() {
	m.font = draw.SystemFont
	m.backgroundInk = draw.ControlBackgroundInk
	m.focusedBackgroundInk = draw.ControlFocusedBackgroundInk
	m.pressedBackgroundInk = draw.ControlPressedBackgroundInk
	m.edgeInk = draw.ControlEdgeAdjColor
	m.textInk = draw.ControlTextColor
	m.pressedTextInk = draw.AlternateSelectedControlTextColor
	m.gap = 3
	m.cornerRadius = 4
	m.clickAnimationTime = time.Millisecond * 100
	m.hAlign = align.Start
	m.vAlign = align.Middle
	m.side = side.Left
}

// Image returns the image. May be nil.
func (c *CheckBox) Image() *draw.Image {
	return c.image
}

// SetImage sets the image. May be nil.
func (c *CheckBox) SetImage(value *draw.Image) *CheckBox {
	if c.image != value {
		c.image = value
		c.MarkForLayoutAndRedraw()
	}
	return c
}

// Text returns the text content.
func (c *CheckBox) Text() string {
	return c.text
}

// SetText sets the text content.
func (c *CheckBox) SetText(value string) *CheckBox {
	if c.text != value {
		c.text = value
		c.MarkForLayoutAndRedraw()
	}
	return c
}

// Font returns the font that will be used when drawing text content.
func (c *CheckBox) Font() *draw.Font {
	return c.font
}

// SetFont sets the font that will be used when drawing text content. Pass in
// nil to use the default.
func (c *CheckBox) SetFont(value *draw.Font) *CheckBox {
	if value == nil {
		value = draw.SystemFont
	}
	if c.font != value {
		c.font = value
		c.MarkForLayoutAndRedraw()
	}
	return c
}

// BackgroundInk returns the ink that will be used for the background when
// enabled but not pressed or focused.
func (c *CheckBox) BackgroundInk() draw.Ink {
	return c.backgroundInk
}

// SetBackgroundInk sets the ink that will be used for the background when
// enabled but not pressed or focused. Pass in nil to use the default.
func (c *CheckBox) SetBackgroundInk(value draw.Ink) *CheckBox {
	if value == nil {
		value = draw.ControlBackgroundInk
	}
	if c.backgroundInk != value {
		c.backgroundInk = value
		c.MarkForRedraw()
	}
	return c
}

// FocusedBackgroundInk returns the ink that will be used for the background
// when enabled and focused.
func (c *CheckBox) FocusedBackgroundInk() draw.Ink {
	return c.focusedBackgroundInk
}

// SetFocusedBackgroundInk sets the ink that will be used for the background
// when enabled and focused. Pass in nil to use the default.
func (c *CheckBox) SetFocusedBackgroundInk(value draw.Ink) *CheckBox {
	if value == nil {
		value = draw.ControlFocusedBackgroundInk
	}
	if c.focusedBackgroundInk != value {
		c.focusedBackgroundInk = value
		c.MarkForRedraw()
	}
	return c
}

// PressedBackgroundInk returns the ink that will be used for the background
// when enabled and pressed.
func (c *CheckBox) PressedBackgroundInk() draw.Ink {
	return c.pressedBackgroundInk
}

// SetPressedBackgroundInk sets the ink that will be used for the background
// when enabled and pressed. Pass in nil to use the default.
func (c *CheckBox) SetPressedBackgroundInk(value draw.Ink) *CheckBox {
	if value == nil {
		value = draw.ControlPressedBackgroundInk
	}
	if c.pressedBackgroundInk != value {
		c.pressedBackgroundInk = value
		c.MarkForRedraw()
	}
	return c
}

// EdgeInk returns the ink that will be used for the edges.
func (c *CheckBox) EdgeInk() draw.Ink {
	return c.edgeInk
}

// SetEdgeInk sets the ink that will be used for the edges. Pass in nil to
// use the default.
func (c *CheckBox) SetEdgeInk(value draw.Ink) *CheckBox {
	if value == nil {
		value = draw.ControlEdgeAdjColor
	}
	if c.edgeInk != value {
		c.edgeInk = value
		c.MarkForRedraw()
	}
	return c
}

// TextInk returns the ink that will be used for the text when disabled or
// not pressed.
func (c *CheckBox) TextInk() draw.Ink {
	return c.textInk
}

// SetTextInk sets the ink that will be used for the text when disabled or
// not pressed. Pass in nil to use the default.
func (c *CheckBox) SetTextInk(value draw.Ink) *CheckBox {
	if value == nil {
		value = draw.ControlTextColor
	}
	if c.textInk != value {
		c.textInk = value
		c.MarkForRedraw()
	}
	return c
}

// PressedTextInk returns the ink that will be used for the text when enabled
// and pressed.
func (c *CheckBox) PressedTextInk() draw.Ink {
	return c.pressedTextInk
}

// SetPressedTextInk sets the ink that will be used for the text when enabled
// and pressed. Pass in nil to use the default.
func (c *CheckBox) SetPressedTextInk(value draw.Ink) *CheckBox {
	if value == nil {
		value = draw.AlternateSelectedControlTextColor
	}
	if c.pressedTextInk != value {
		c.pressedTextInk = value
		c.MarkForRedraw()
	}
	return c
}

// Gap returns the gap to put between the checkbox, image and text.
func (c *CheckBox) Gap() float64 {
	return c.gap
}

// SetGap sets the gap to put between the checkbox, image and text.
func (c *CheckBox) SetGap(value float64) *CheckBox {
	if value < 0 {
		value = 0
	}
	if c.gap != value {
		c.gap = value
		c.MarkForLayoutAndRedraw()
	}
	return c
}

// CornerRadius returns the amount of rounding to use on the corners.
func (c *CheckBox) CornerRadius() float64 {
	return c.cornerRadius
}

// SetCornerRadius sets the amount of rounding to use on the corners.
func (c *CheckBox) SetCornerRadius(value float64) *CheckBox {
	if value < 0 {
		value = 0
	}
	if c.cornerRadius != value {
		c.cornerRadius = value
		c.MarkForRedraw()
	}
	return c
}

// ClickAnimationTime returns the amount of time to spend animating the click
// action.
func (c *CheckBox) ClickAnimationTime() time.Duration {
	return c.clickAnimationTime
}

// SetClickAnimationTime sets the amount of time to spend animating the click
// action.
func (c *CheckBox) SetClickAnimationTime(value time.Duration) *CheckBox {
	if value < 0 {
		value = 0
	}
	if c.clickAnimationTime != value {
		c.clickAnimationTime = value
	}
	return c
}

// HAlign returns the horizontal alignment.
func (c *CheckBox) HAlign() align.Alignment {
	return c.hAlign
}

// SetHAlign sets the horizontal alignment.
func (c *CheckBox) SetHAlign(value align.Alignment) *CheckBox {
	if c.hAlign != value {
		c.hAlign = value
		c.MarkForRedraw()
	}
	return c
}

// VAlign returns the vertical alignment.
func (c *CheckBox) VAlign() align.Alignment {
	return c.vAlign
}

// SetVAlign sets the vertical alignment.
func (c *CheckBox) SetVAlign(value align.Alignment) *CheckBox {
	if c.vAlign != value {
		c.vAlign = value
		c.MarkForRedraw()
	}
	return c
}

// Side returns the side of the text the image should be on.
func (c *CheckBox) Side() side.Side {
	return c.side
}

// SetSide sets the side of the text the image should be on.
func (c *CheckBox) SetSide(value side.Side) *CheckBox {
	if c.side != value {
		c.side = value
		c.MarkForRedraw()
	}
	return c
}

// SetBorder sets the border. May be nil.
func (c *CheckBox) SetBorder(value border.Border) *CheckBox {
	c.Panel.SetBorder(value)
	return c
}

// SetEnabled sets enabled state.
func (c *CheckBox) SetEnabled(enabled bool) *CheckBox {
	c.Panel.SetEnabled(enabled)
	return c
}

// SetFocusable whether it can have the keyboard focus.
func (c *CheckBox) SetFocusable(focusable bool) *CheckBox {
	c.Panel.SetFocusable(focusable)
	return c
}

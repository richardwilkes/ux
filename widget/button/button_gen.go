// Code created from "widget.go.tmpl" - don't edit by hand

package button

import (
	"time"

	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/layout/align"
	"github.com/richardwilkes/ux/layout/side"
)

type managed struct {
	image                 *draw.Image //nolint:structcheck
	text                  string      //nolint:structcheck
	font                  *draw.Font
	backgroundInk         draw.Ink
	selectedBackgroundInk draw.Ink
	focusedBackgroundInk  draw.Ink
	pressedBackgroundInk  draw.Ink
	edgeInk               draw.Ink
	textInk               draw.Ink
	pressedTextInk        draw.Ink
	gap                   float64
	cornerRadius          float64
	hMargin               float64
	vMargin               float64
	imageOnlyHMargin      float64
	imageOnlyVMargin      float64
	clickAnimationTime    time.Duration
	hAlign                align.Alignment
	vAlign                align.Alignment
	side                  side.Side
	sticky                bool //nolint:structcheck
}

func (m *managed) initialize() {
	m.font = draw.SystemFont
	m.backgroundInk = draw.ControlBackgroundInk
	m.selectedBackgroundInk = draw.ControlSelectedBackgroundInk
	m.focusedBackgroundInk = draw.ControlFocusedBackgroundInk
	m.pressedBackgroundInk = draw.ControlPressedBackgroundInk
	m.edgeInk = draw.ControlEdgeAdjColor
	m.textInk = draw.ControlTextColor
	m.pressedTextInk = draw.AlternateSelectedControlTextColor
	m.gap = 3
	m.cornerRadius = 4
	m.hMargin = 8
	m.vMargin = 1
	m.imageOnlyHMargin = 3
	m.imageOnlyVMargin = 3
	m.clickAnimationTime = time.Millisecond * 100
	m.hAlign = align.Middle
	m.vAlign = align.Middle
	m.side = side.Left
}

// Image returns the image. May be nil.
func (b *Button) Image() *draw.Image {
	return b.image
}

// SetImage sets the image. May be nil.
func (b *Button) SetImage(value *draw.Image) *Button {
	if b.image != value {
		b.image = value
		b.MarkForLayoutAndRedraw()
	}
	return b
}

// Text returns the text content.
func (b *Button) Text() string {
	return b.text
}

// SetText sets the text content.
func (b *Button) SetText(value string) *Button {
	if b.text != value {
		b.text = value
		b.MarkForLayoutAndRedraw()
	}
	return b
}

// Font returns the font that will be used when drawing text content.
func (b *Button) Font() *draw.Font {
	return b.font
}

// SetFont sets the font that will be used when drawing text content. Pass in
// nil to use the default.
func (b *Button) SetFont(value *draw.Font) *Button {
	if value == nil {
		value = draw.SystemFont
	}
	if b.font != value {
		b.font = value
		b.MarkForLayoutAndRedraw()
	}
	return b
}

// BackgroundInk returns the ink that will be used for the background when
// enabled but not selected, pressed or focused.
func (b *Button) BackgroundInk() draw.Ink {
	return b.backgroundInk
}

// SetBackgroundInk sets the ink that will be used for the background when
// enabled but not selected, pressed or focused. Pass in nil to use the
// default.
func (b *Button) SetBackgroundInk(value draw.Ink) *Button {
	if value == nil {
		value = draw.ControlBackgroundInk
	}
	if b.backgroundInk != value {
		b.backgroundInk = value
		b.MarkForRedraw()
	}
	return b
}

// SelectedBackgroundInk returns the ink that will be used for the background
// when enabled and selected, but not pressed or focused.
func (b *Button) SelectedBackgroundInk() draw.Ink {
	return b.selectedBackgroundInk
}

// SetSelectedBackgroundInk sets the ink that will be used for the background
// when enabled and selected, but not pressed or focused. Pass in nil to use
// the default.
func (b *Button) SetSelectedBackgroundInk(value draw.Ink) *Button {
	if value == nil {
		value = draw.ControlSelectedBackgroundInk
	}
	if b.selectedBackgroundInk != value {
		b.selectedBackgroundInk = value
		b.MarkForRedraw()
	}
	return b
}

// FocusedBackgroundInk returns the ink that will be used for the background
// when enabled and focused.
func (b *Button) FocusedBackgroundInk() draw.Ink {
	return b.focusedBackgroundInk
}

// SetFocusedBackgroundInk sets the ink that will be used for the background
// when enabled and focused. Pass in nil to use the default.
func (b *Button) SetFocusedBackgroundInk(value draw.Ink) *Button {
	if value == nil {
		value = draw.ControlFocusedBackgroundInk
	}
	if b.focusedBackgroundInk != value {
		b.focusedBackgroundInk = value
		b.MarkForRedraw()
	}
	return b
}

// PressedBackgroundInk returns the ink that will be used for the background
// when enabled and pressed.
func (b *Button) PressedBackgroundInk() draw.Ink {
	return b.pressedBackgroundInk
}

// SetPressedBackgroundInk sets the ink that will be used for the background
// when enabled and pressed. Pass in nil to use the default.
func (b *Button) SetPressedBackgroundInk(value draw.Ink) *Button {
	if value == nil {
		value = draw.ControlPressedBackgroundInk
	}
	if b.pressedBackgroundInk != value {
		b.pressedBackgroundInk = value
		b.MarkForRedraw()
	}
	return b
}

// EdgeInk returns the ink that will be used for the edges.
func (b *Button) EdgeInk() draw.Ink {
	return b.edgeInk
}

// SetEdgeInk sets the ink that will be used for the edges. Pass in nil to
// use the default.
func (b *Button) SetEdgeInk(value draw.Ink) *Button {
	if value == nil {
		value = draw.ControlEdgeAdjColor
	}
	if b.edgeInk != value {
		b.edgeInk = value
		b.MarkForRedraw()
	}
	return b
}

// TextInk returns the ink that will be used for the text when disabled or
// not pressed.
func (b *Button) TextInk() draw.Ink {
	return b.textInk
}

// SetTextInk sets the ink that will be used for the text when disabled or
// not pressed. Pass in nil to use the default.
func (b *Button) SetTextInk(value draw.Ink) *Button {
	if value == nil {
		value = draw.ControlTextColor
	}
	if b.textInk != value {
		b.textInk = value
		b.MarkForRedraw()
	}
	return b
}

// PressedTextInk returns the ink that will be used for the text when enabled
// and pressed.
func (b *Button) PressedTextInk() draw.Ink {
	return b.pressedTextInk
}

// SetPressedTextInk sets the ink that will be used for the text when enabled
// and pressed. Pass in nil to use the default.
func (b *Button) SetPressedTextInk(value draw.Ink) *Button {
	if value == nil {
		value = draw.AlternateSelectedControlTextColor
	}
	if b.pressedTextInk != value {
		b.pressedTextInk = value
		b.MarkForRedraw()
	}
	return b
}

// Gap returns the gap to put between the image and text.
func (b *Button) Gap() float64 {
	return b.gap
}

// SetGap sets the gap to put between the image and text.
func (b *Button) SetGap(value float64) *Button {
	if value < 0 {
		value = 0
	}
	if b.gap != value {
		b.gap = value
		b.MarkForLayoutAndRedraw()
	}
	return b
}

// CornerRadius returns the amount of rounding to use on the corners.
func (b *Button) CornerRadius() float64 {
	return b.cornerRadius
}

// SetCornerRadius sets the amount of rounding to use on the corners.
func (b *Button) SetCornerRadius(value float64) *Button {
	if value < 0 {
		value = 0
	}
	if b.cornerRadius != value {
		b.cornerRadius = value
		b.MarkForRedraw()
	}
	return b
}

// HMargin returns the margin on the left and right side of the content.
func (b *Button) HMargin() float64 {
	return b.hMargin
}

// SetHMargin sets the margin on the left and right side of the content.
func (b *Button) SetHMargin(value float64) *Button {
	if value < 0 {
		value = 0
	}
	if b.hMargin != value {
		b.hMargin = value
		b.MarkForRedraw()
	}
	return b
}

// VMargin returns the margin on the top and bottom side of the content.
func (b *Button) VMargin() float64 {
	return b.vMargin
}

// SetVMargin sets the margin on the top and bottom side of the content.
func (b *Button) SetVMargin(value float64) *Button {
	if value < 0 {
		value = 0
	}
	if b.vMargin != value {
		b.vMargin = value
		b.MarkForRedraw()
	}
	return b
}

// ImageOnlyHMargin returns the margin on the left and right side of the
// content when only an image is present.
func (b *Button) ImageOnlyHMargin() float64 {
	return b.imageOnlyHMargin
}

// SetImageOnlyHMargin sets the margin on the left and right side of the
// content when only an image is present.
func (b *Button) SetImageOnlyHMargin(value float64) *Button {
	if value < 0 {
		value = 0
	}
	if b.imageOnlyHMargin != value {
		b.imageOnlyHMargin = value
		b.MarkForRedraw()
	}
	return b
}

// ImageOnlyVMargin returns the margin on the top and bottom side of the
// content when only an image is present.
func (b *Button) ImageOnlyVMargin() float64 {
	return b.imageOnlyVMargin
}

// SetImageOnlyVMargin sets the margin on the top and bottom side of the
// content when only an image is present.
func (b *Button) SetImageOnlyVMargin(value float64) *Button {
	if value < 0 {
		value = 0
	}
	if b.imageOnlyVMargin != value {
		b.imageOnlyVMargin = value
		b.MarkForRedraw()
	}
	return b
}

// ClickAnimationTime returns the amount of time to spend animating the click
// action.
func (b *Button) ClickAnimationTime() time.Duration {
	return b.clickAnimationTime
}

// SetClickAnimationTime sets the amount of time to spend animating the click
// action.
func (b *Button) SetClickAnimationTime(value time.Duration) *Button {
	if value < time.Millisecond*50 {
		value = time.Millisecond * 50
	}
	if b.clickAnimationTime != value {
		b.clickAnimationTime = value
	}
	return b
}

// HAlign returns the horizontal alignment.
func (b *Button) HAlign() align.Alignment {
	return b.hAlign
}

// SetHAlign sets the horizontal alignment.
func (b *Button) SetHAlign(value align.Alignment) *Button {
	if b.hAlign != value {
		b.hAlign = value
		b.MarkForRedraw()
	}
	return b
}

// VAlign returns the vertical alignment.
func (b *Button) VAlign() align.Alignment {
	return b.vAlign
}

// SetVAlign sets the vertical alignment.
func (b *Button) SetVAlign(value align.Alignment) *Button {
	if b.vAlign != value {
		b.vAlign = value
		b.MarkForRedraw()
	}
	return b
}

// Side returns the side of the text the image should be on.
func (b *Button) Side() side.Side {
	return b.side
}

// SetSide sets the side of the text the image should be on.
func (b *Button) SetSide(value side.Side) *Button {
	if b.side != value {
		b.side = value
		b.MarkForRedraw()
	}
	return b
}

// Sticky returns whether the button will visually retain its selected state.
func (b *Button) Sticky() bool {
	return b.sticky
}

// SetSticky sets whether the button will visually retain its selected state.
func (b *Button) SetSticky(value bool) *Button {
	if b.sticky != value {
		b.sticky = value
		b.MarkForRedraw()
	}
	return b
}

// SetBorder sets the border. May be nil.
func (b *Button) SetBorder(value border.Border) *Button {
	b.Panel.SetBorder(value)
	return b
}

// SetEnabled sets enabled state.
func (b *Button) SetEnabled(enabled bool) *Button {
	b.Panel.SetEnabled(enabled)
	return b
}

// SetFocusable whether it can have the keyboard focus.
func (b *Button) SetFocusable(focusable bool) *Button {
	b.Panel.SetFocusable(focusable)
	return b
}

// SetSelected sets the panel's selected state.
func (b *Button) SetSelected(selected bool) *Button {
	b.Panel.SetSelected(selected)
	return b
}

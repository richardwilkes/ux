// Code created from "widget.go.tmpl" - don't edit by hand

package textfield

import (
	"time"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/draw"
)

type managed struct {
	font                      *draw.Font
	backgroundInk             draw.Ink
	disabledBackgroundInk     draw.Ink
	invalidBackgroundInk      draw.Ink
	selectedTextBackgroundInk draw.Ink
	textInk                   draw.Ink
	selectedTextInk           draw.Ink
	watermarkInk              draw.Ink
	minimumTextWidth          float64
	blinkRate                 time.Duration
	watermark                 string //nolint:structcheck
	focusedBorder             border.Border
	unfocusedBorder           border.Border
}

func (m *managed) initialize() {
	m.font = draw.UserFont
	m.backgroundInk = draw.TextBackgroundColor
	m.disabledBackgroundInk = draw.WindowBackgroundColor
	m.invalidBackgroundInk = draw.SystemRedColor
	m.selectedTextBackgroundInk = draw.SelectedTextBackgroundColor
	m.textInk = draw.TextColor
	m.selectedTextInk = draw.SelectedTextColor
	m.watermarkInk = draw.PlaceholderTextColor
	m.minimumTextWidth = 10
	m.blinkRate = time.Millisecond * 560
	m.focusedBorder = border.NewCompound(border.NewLine(draw.ControlAccentColor, geom.NewUniformInsets(2), false), border.NewEmpty(geom.Insets{Top: 1, Left: 2, Bottom: 0, Right: 2}))
	m.unfocusedBorder = border.NewCompound(border.NewLine(draw.ControlEdgeAdjColor, geom.NewUniformInsets(1), false), border.NewEmpty(geom.Insets{Top: 2, Left: 3, Bottom: 1, Right: 3}))
}

// Font returns the font that will be used when drawing text content.
func (t *TextField) Font() *draw.Font {
	return t.font
}

// SetFont sets the font that will be used when drawing text content. Pass in
// nil to use the default.
func (t *TextField) SetFont(value *draw.Font) *TextField {
	if value == nil {
		value = draw.UserFont
	}
	if t.font != value {
		t.font = value
		t.MarkForLayoutAndRedraw()
	}
	return t
}

// BackgroundInk returns the ink that will be used for the background when
// enabled.
func (t *TextField) BackgroundInk() draw.Ink {
	return t.backgroundInk
}

// SetBackgroundInk sets the ink that will be used for the background when
// enabled. Pass in nil to use the default.
func (t *TextField) SetBackgroundInk(value draw.Ink) *TextField {
	if value == nil {
		value = draw.TextBackgroundColor
	}
	if t.backgroundInk != value {
		t.backgroundInk = value
		t.MarkForRedraw()
	}
	return t
}

// DisabledBackgroundInk returns the ink that will be used for the background
// when disabled.
func (t *TextField) DisabledBackgroundInk() draw.Ink {
	return t.disabledBackgroundInk
}

// SetDisabledBackgroundInk sets the ink that will be used for the background
// when disabled. Pass in nil to use the default.
func (t *TextField) SetDisabledBackgroundInk(value draw.Ink) *TextField {
	if value == nil {
		value = draw.WindowBackgroundColor
	}
	if t.disabledBackgroundInk != value {
		t.disabledBackgroundInk = value
		t.MarkForRedraw()
	}
	return t
}

// InvalidBackgroundInk returns the ink that will be used for the background
// when marked invalid.
func (t *TextField) InvalidBackgroundInk() draw.Ink {
	return t.invalidBackgroundInk
}

// SetInvalidBackgroundInk sets the ink that will be used for the background
// when marked invalid. Pass in nil to use the default.
func (t *TextField) SetInvalidBackgroundInk(value draw.Ink) *TextField {
	if value == nil {
		value = draw.SystemRedColor
	}
	if t.invalidBackgroundInk != value {
		t.invalidBackgroundInk = value
		t.MarkForRedraw()
	}
	return t
}

// SelectedTextBackgroundInk returns the ink that will be used for the
// background of selected text.
func (t *TextField) SelectedTextBackgroundInk() draw.Ink {
	return t.selectedTextBackgroundInk
}

// SetSelectedTextBackgroundInk sets the ink that will be used for the
// background of selected text. Pass in nil to use the default.
func (t *TextField) SetSelectedTextBackgroundInk(value draw.Ink) *TextField {
	if value == nil {
		value = draw.SelectedTextBackgroundColor
	}
	if t.selectedTextBackgroundInk != value {
		t.selectedTextBackgroundInk = value
		t.MarkForRedraw()
	}
	return t
}

// TextInk returns the ink that will be used for the text content when not
// selected.
func (t *TextField) TextInk() draw.Ink {
	return t.textInk
}

// SetTextInk sets the ink that will be used for the text content when not
// selected. Pass in nil to use the default.
func (t *TextField) SetTextInk(value draw.Ink) *TextField {
	if value == nil {
		value = draw.TextColor
	}
	if t.textInk != value {
		t.textInk = value
		t.MarkForRedraw()
	}
	return t
}

// SelectedTextInk returns the ink that will be used for the text content
// when selected.
func (t *TextField) SelectedTextInk() draw.Ink {
	return t.selectedTextInk
}

// SetSelectedTextInk sets the ink that will be used for the text content
// when selected. Pass in nil to use the default.
func (t *TextField) SetSelectedTextInk(value draw.Ink) *TextField {
	if value == nil {
		value = draw.SelectedTextColor
	}
	if t.selectedTextInk != value {
		t.selectedTextInk = value
		t.MarkForRedraw()
	}
	return t
}

// WatermarkInk returns the ink that will be used for the watermark text
// content.
func (t *TextField) WatermarkInk() draw.Ink {
	return t.watermarkInk
}

// SetWatermarkInk sets the ink that will be used for the watermark text
// content. Pass in nil to use the default.
func (t *TextField) SetWatermarkInk(value draw.Ink) *TextField {
	if value == nil {
		value = draw.PlaceholderTextColor
	}
	if t.watermarkInk != value {
		t.watermarkInk = value
		t.MarkForRedraw()
	}
	return t
}

// MinimumTextWidth returns the minimum horizontal space to permit for text.
func (t *TextField) MinimumTextWidth() float64 {
	return t.minimumTextWidth
}

// SetMinimumTextWidth sets the minimum horizontal space to permit for text.
func (t *TextField) SetMinimumTextWidth(value float64) *TextField {
	if value < 10 {
		value = 10
	}
	if t.minimumTextWidth != value {
		t.minimumTextWidth = value
		t.MarkForLayoutAndRedraw()
	}
	return t
}

// BlinkRate returns the rate at which the cursor blinks.
func (t *TextField) BlinkRate() time.Duration {
	return t.blinkRate
}

// SetBlinkRate sets the rate at which the cursor blinks.
func (t *TextField) SetBlinkRate(value time.Duration) *TextField {
	if value < time.Millisecond*50 {
		value = time.Millisecond * 50
	}
	if t.blinkRate != value {
		t.blinkRate = value
	}
	return t
}

// Watermark returns the help text that will show up in an empty field.
func (t *TextField) Watermark() string {
	return t.watermark
}

// SetWatermark sets the help text that will show up in an empty field.
func (t *TextField) SetWatermark(value string) *TextField {
	if t.watermark != value {
		t.watermark = value
		t.MarkForRedraw()
	}
	return t
}

// FocusedBorder returns the border to use when focused. Note that the border
// should present the same insets as the unfocused border or the display will
// not appear correct.
func (t *TextField) FocusedBorder() border.Border {
	return t.focusedBorder
}

// SetFocusedBorder sets the border to use when focused. Note that the border
// should present the same insets as the unfocused border or the display will
// not appear correct. Pass in nil to use the default.
func (t *TextField) SetFocusedBorder(value border.Border) *TextField {
	if value == nil {
		value = border.NewCompound(border.NewLine(draw.ControlAccentColor, geom.NewUniformInsets(2), false), border.NewEmpty(geom.Insets{Top: 1, Left: 2, Bottom: 0, Right: 2}))
	}
	if t.focusedBorder != value {
		t.focusedBorder = value
		t.MarkForLayoutAndRedraw()
	}
	return t
}

// UnfocusedBorder returns the border to use when not focused. Note that the
// border should present the same insets as the focused border or the display
// will not appear correct.
func (t *TextField) UnfocusedBorder() border.Border {
	return t.unfocusedBorder
}

// SetUnfocusedBorder sets the border to use when not focused. Note that the
// border should present the same insets as the focused border or the display
// will not appear correct. Pass in nil to use the default.
func (t *TextField) SetUnfocusedBorder(value border.Border) *TextField {
	if value == nil {
		value = border.NewCompound(border.NewLine(draw.ControlEdgeAdjColor, geom.NewUniformInsets(1), false), border.NewEmpty(geom.Insets{Top: 2, Left: 3, Bottom: 1, Right: 3}))
	}
	if t.unfocusedBorder != value {
		t.unfocusedBorder = value
		t.MarkForLayoutAndRedraw()
	}
	return t
}

// SetBorder sets the border. May be nil.
func (t *TextField) SetBorder(value border.Border) *TextField {
	t.Panel.SetBorder(value)
	return t
}

// SetEnabled sets enabled state.
func (t *TextField) SetEnabled(enabled bool) *TextField {
	t.Panel.SetEnabled(enabled)
	return t
}

// SetFocusable whether it can have the keyboard focus.
func (t *TextField) SetFocusable(focusable bool) *TextField {
	t.Panel.SetFocusable(focusable)
	return t
}

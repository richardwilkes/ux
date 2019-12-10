package button

import (
	"time"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/keys"
	"github.com/richardwilkes/ux/layout"
	"github.com/richardwilkes/ux/layout/align"
	"github.com/richardwilkes/ux/layout/side"
	"github.com/richardwilkes/ux/widget"
	"github.com/richardwilkes/ux/widget/selectable"
)

// Button represents a clickable button.
type Button struct {
	selectable.Panel
	ClickCallback         func()
	image                 *draw.Image
	text                  string
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
	side                  side.Side
	hAlign                align.Alignment
	vAlign                align.Alignment
	Pressed               bool
	Sticky                bool
}

// NewWithText creates a new button with the specified text.
func NewWithText(text string) *Button {
	return New(text, nil)
}

// NewWithImage creates a new button with the specified image.
func NewWithImage(image *draw.Image) *Button {
	return New("", image)
}

// New creates a new button. Both image and text are optional.
func New(text string, image *draw.Image) *Button {
	b := &Button{
		image:                 image,
		text:                  text,
		font:                  draw.SystemFont,
		backgroundInk:         draw.ControlBackgroundInk,
		selectedBackgroundInk: draw.ControlSelectedBackgroundInk,
		focusedBackgroundInk:  draw.ControlFocusedBackgroundInk,
		pressedBackgroundInk:  draw.ControlPressedBackgroundInk,
		edgeInk:               draw.ControlEdgeAdjColor,
		textInk:               draw.ControlTextColor,
		pressedTextInk:        draw.AlternateSelectedControlTextColor,
		gap:                   3,
		cornerRadius:          4,
		hMargin:               8,
		vMargin:               1,
		imageOnlyHMargin:      3,
		imageOnlyVMargin:      3,
		clickAnimationTime:    time.Millisecond * 100,
		side:                  side.Left,
		hAlign:                align.Middle,
		vAlign:                align.Middle,
	}
	b.InitTypeAndID(b)
	b.SetFocusable(true)
	b.SetSizer(b.DefaultSizes)
	b.DrawCallback = b.DefaultDraw
	b.GainedFocusCallback = b.MarkForRedraw
	b.LostFocusCallback = b.MarkForRedraw
	b.MouseDownCallback = b.DefaultMouseDown
	b.MouseDragCallback = b.DefaultMouseDrag
	b.MouseUpCallback = b.DefaultMouseUp
	b.KeyDownCallback = b.DefaultKeyDown
	return b
}

// Image returns the image. May be nil.
func (b *Button) Image() *draw.Image {
	return b.image
}

// SetImage sets the image. May be nil.
func (b *Button) SetImage(image *draw.Image) *Button {
	if b.image != image {
		b.image = image
		b.MarkForLayoutAndRedraw()
	}
	return b
}

// Text returns the text content.
func (b *Button) Text() string {
	return b.text
}

// SetText sets the text content.
func (b *Button) SetText(text string) *Button {
	if b.text != text {
		b.text = text
		b.MarkForLayoutAndRedraw()
	}
	return b
}

// Font returns the font that will be used when drawing text content.
func (b *Button) Font() *draw.Font {
	return b.font
}

// SetFont sets the font to use when drawing text content. Passing in nil will
// use the default font.
func (b *Button) SetFont(font *draw.Font) *Button {
	if font == nil {
		font = draw.SystemFont
	}
	if b.font != font {
		b.font = font
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
// enabled but not selected, pressed or focused. Passing in nil will use the
// default ink.
func (b *Button) SetBackgroundInk(ink draw.Ink) *Button {
	if ink == nil {
		ink = draw.ControlBackgroundInk
	}
	if b.backgroundInk != ink {
		b.backgroundInk = ink
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
// when enabled and selected, but not pressed or focused. Passing in nil will
// use the default ink.
func (b *Button) SetSelectedBackgroundInk(ink draw.Ink) *Button {
	if ink == nil {
		ink = draw.ControlSelectedBackgroundInk
	}
	if b.selectedBackgroundInk != ink {
		b.selectedBackgroundInk = ink
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
// when enabled and focused. Passing in nil will use the default ink.
func (b *Button) SetFocusedBackgroundInk(ink draw.Ink) *Button {
	if ink == nil {
		ink = draw.ControlFocusedBackgroundInk
	}
	if b.focusedBackgroundInk != ink {
		b.focusedBackgroundInk = ink
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
// when enabled and pressed. Passing in nil will use the default ink.
func (b *Button) SetPressedBackgroundInk(ink draw.Ink) *Button {
	if ink == nil {
		ink = draw.ControlPressedBackgroundInk
	}
	if b.pressedBackgroundInk != ink {
		b.pressedBackgroundInk = ink
		b.MarkForRedraw()
	}
	return b
}

// EdgeInk returns the ink that will be used for the edges.
func (b *Button) EdgeInk() draw.Ink {
	return b.edgeInk
}

// SetEdgeInk sets the ink that will be used for the edges. Passing in nil
// will use the default ink.
func (b *Button) SetEdgeInk(ink draw.Ink) *Button {
	if ink == nil {
		ink = draw.ControlEdgeAdjColor
	}
	if b.edgeInk != ink {
		b.edgeInk = ink
		b.MarkForRedraw()
	}
	return b
}

// TextInk returns the ink that will be used for the text when disabled or not
// pressed.
func (b *Button) TextInk() draw.Ink {
	return b.textInk
}

// SetTextInk sets the ink that will be used for the text when disabled or not
// pressed. Passing in nil will use the default ink.
func (b *Button) SetTextInk(ink draw.Ink) *Button {
	if ink == nil {
		ink = draw.ControlTextColor
	}
	if b.textInk != ink {
		b.textInk = ink
		b.MarkForRedraw()
	}
	return b
}

// PressedTextInk returns the ink that will be used for the text when enabled
// and pressed.
func (b *Button) PressedTextInk() draw.Ink {
	return b.textInk
}

// SetPressedTextInk sets the ink that will be used for the text when enabled
// and pressed. Passing in nil will use the default ink.
func (b *Button) SetPressedTextInk(ink draw.Ink) *Button {
	if ink == nil {
		ink = draw.AlternateSelectedControlTextColor
	}
	if b.pressedTextInk != ink {
		b.pressedTextInk = ink
		b.MarkForRedraw()
	}
	return b
}

// Gap returns the gap to put between the image and text.
func (b *Button) Gap() float64 {
	return b.gap
}

// SetGap sets the gap to put between the image and text.
func (b *Button) SetGap(gap float64) *Button {
	if gap < 0 {
		gap = 0
	}
	if b.gap != gap {
		b.gap = gap
		if b.image != nil && b.text != "" {
			b.MarkForLayoutAndRedraw()
		}
	}
	return b
}

// CornerRadius returns the amount of rounding to use on the corners.
func (b *Button) CornerRadius() float64 {
	return b.cornerRadius
}

// SetCornerRadius sets the amount of rounding to use on the corners.
func (b *Button) SetCornerRadius(cornerRadius float64) *Button {
	if cornerRadius < 0 {
		cornerRadius = 0
	}
	if b.cornerRadius != cornerRadius {
		b.cornerRadius = cornerRadius
		b.MarkForRedraw()
	}
	return b
}

// HMargin returns the margin on the left and right side of the content.
func (b *Button) HMargin() float64 {
	return b.hMargin
}

// SetHMargin sets the margin on the left and right side of the content.
func (b *Button) SetHMargin(hMargin float64) *Button {
	if hMargin < 0 {
		hMargin = 0
	}
	if b.hMargin != hMargin {
		b.hMargin = hMargin
		b.MarkForRedraw()
	}
	return b
}

// VMargin returns the margin on the top and bottom side of the content.
func (b *Button) VMargin() float64 {
	return b.vMargin
}

// SetVMargin sets the margin on the top and bottom side of the content.
func (b *Button) SetVMargin(vMargin float64) *Button {
	if vMargin < 0 {
		vMargin = 0
	}
	if b.vMargin != vMargin {
		b.vMargin = vMargin
		b.MarkForRedraw()
	}
	return b
}

// ImageOnlyHMargin returns the margin on the left and right side of the
// content when only an image is present.
func (b *Button) ImageOnlyHMargin() float64 {
	return b.hMargin
}

// SetImageOnlyHMargin sets the margin on the left and right side of the
// content when only an image is present.
func (b *Button) SetImageOnlyHMargin(hMargin float64) *Button {
	if hMargin < 0 {
		hMargin = 0
	}
	if b.imageOnlyHMargin != hMargin {
		b.imageOnlyHMargin = hMargin
		b.MarkForRedraw()
	}
	return b
}

// ImageOnlyVMargin returns the margin on the top and bottom side of the
// content when only an image is present.
func (b *Button) ImageOnlyVMargin() float64 {
	return b.vMargin
}

// SetImageOnlyVMargin sets the margin on the top and bottom side of the
// content when only an image is present.
func (b *Button) SetImageOnlyVMargin(vMargin float64) *Button {
	if vMargin < 0 {
		vMargin = 0
	}
	if b.imageOnlyVMargin != vMargin {
		b.imageOnlyVMargin = vMargin
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
func (b *Button) SetClickAnimationTime(duration time.Duration) *Button {
	if duration < 0 {
		duration = 0
	}
	b.clickAnimationTime = duration
	return b
}

// Side returns the side of the text the image should be on.
func (b *Button) Side() side.Side {
	return b.side
}

// SetSide sets the side of the text the image should be on.
func (b *Button) SetSide(s side.Side) *Button {
	if b.side != s {
		b.side = s
		b.MarkForRedraw()
	}
	return b
}

// HAlign returns the horizontal alignment.
func (b *Button) HAlign() align.Alignment {
	return b.hAlign
}

// SetHAlign sets the horizontal alignment.
func (b *Button) SetHAlign(hAlign align.Alignment) *Button {
	if b.hAlign != hAlign {
		b.hAlign = hAlign
		b.MarkForRedraw()
	}
	return b
}

// VAlign returns the vertical alignment.
func (b *Button) VAlign() align.Alignment {
	return b.hAlign
}

// SetVAlign sets the vertical alignment.
func (b *Button) SetVAlign(vAlign align.Alignment) *Button {
	if b.vAlign != vAlign {
		b.vAlign = vAlign
		b.MarkForRedraw()
	}
	return b
}

// SetBorder sets the border. May be nil.
func (b *Button) SetBorder(aBorder border.Border) *Button {
	b.Panel.SetBorder(aBorder)
	return b
}

// SetEnabled sets enabled state.
func (b *Button) SetEnabled(enabled bool) *Button {
	b.Panel.SetEnabled(enabled)
	return b
}

// SetFocusable whether the label can have the keyboard focus.
func (b *Button) SetFocusable(focusable bool) *Button {
	b.Panel.SetFocusable(focusable)
	return b
}

// DefaultSizes provides the default sizing.
func (b *Button) DefaultSizes(hint geom.Size) (min, pref, max geom.Size) {
	text := b.text
	if b.image == nil && text == "" {
		text = "M"
	}
	pref = widget.LabelSize(text, b.font, b.image, b.side, b.gap)
	if theBorder := b.Border(); theBorder != nil {
		pref.AddInsets(theBorder.Insets())
	}
	pref.Width += b.horizontalMargin()*2 + 2
	pref.Height += b.verticalMargin()*2 + 2
	pref.GrowToInteger()
	pref.ConstrainForHint(hint)
	return pref, pref, layout.MaxSize(pref)
}

func (b *Button) horizontalMargin() float64 {
	if b.text == "" && b.image != nil {
		return b.imageOnlyHMargin
	}
	return b.hMargin
}

func (b *Button) verticalMargin() float64 {
	if b.text == "" && b.image != nil {
		return b.imageOnlyVMargin
	}
	return b.vMargin
}

// DefaultDraw provides the default drawing.
func (b *Button) DefaultDraw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	if !b.Enabled() {
		gc.SetOpacity(0.33)
	}
	rect := b.ContentRect(false)
	widget.DrawRoundedRectBase(gc, rect, b.cornerRadius, b.currentBackgroundInk(), b.edgeInk)
	rect.InsetUniform(1.5)
	rect.X += b.horizontalMargin()
	rect.Y += b.verticalMargin()
	rect.Width -= b.horizontalMargin() * 2
	rect.Height -= b.verticalMargin() * 2
	widget.DrawLabel(gc, rect, b.hAlign, b.vAlign, b.text, b.font, b.currentTextInk(), b.image, b.side, b.gap, b.Enabled())
}

func (b *Button) currentBackgroundInk() draw.Ink {
	switch {
	case b.Pressed:
		return b.pressedBackgroundInk
	case b.Focused():
		return b.focusedBackgroundInk
	case b.Sticky && b.Selected():
		return b.selectedBackgroundInk
	default:
		return b.backgroundInk
	}
}

func (b *Button) currentTextInk() draw.Ink {
	if b.Pressed || b.Focused() {
		return b.pressedTextInk
	}
	return b.textInk
}

// Click makes the button behave as if a user clicked on it.
func (b *Button) Click() {
	b.SetSelected(true)
	pressed := b.Pressed
	b.Pressed = true
	b.MarkForRedraw()
	b.FlushDrawing()
	b.Pressed = pressed
	time.Sleep(b.clickAnimationTime)
	b.MarkForRedraw()
	if b.ClickCallback != nil {
		b.ClickCallback()
	}
}

// DefaultMouseDown provides the default mouse down handling.
func (b *Button) DefaultMouseDown(where geom.Point, button, clickCount int, mod keys.Modifiers) bool {
	b.Pressed = true
	b.MarkForRedraw()
	return true
}

// DefaultMouseDrag provides the default mouse drag handling.
func (b *Button) DefaultMouseDrag(where geom.Point, button int, mod keys.Modifiers) {
	rect := b.ContentRect(false)
	pressed := rect.ContainsPoint(where)
	if b.Pressed != pressed {
		b.Pressed = pressed
		b.MarkForRedraw()
	}
}

// DefaultMouseUp provides the default mouse up handling.
func (b *Button) DefaultMouseUp(where geom.Point, button int, mod keys.Modifiers) {
	b.Pressed = false
	b.MarkForRedraw()
	rect := b.ContentRect(false)
	if rect.ContainsPoint(where) {
		b.SetSelected(true)
		if b.ClickCallback != nil {
			b.ClickCallback()
		}
	}
}

// DefaultKeyDown provides the default key down handling.
func (b *Button) DefaultKeyDown(keyCode int, ch rune, mod keys.Modifiers, repeat bool) bool {
	if keys.IsControlAction(keyCode) {
		b.Click()
		return true
	}
	return false
}

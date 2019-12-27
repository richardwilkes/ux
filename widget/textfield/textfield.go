package textfield

import (
	"math"
	"strings"
	"time"
	"unicode"

	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/clipboard"
	"github.com/richardwilkes/ux/clipboard/datatypes"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/ids"
	"github.com/richardwilkes/ux/keys"
	"github.com/richardwilkes/ux/layout"
)

// TextField provides a single-line text input control.
type TextField struct {
	ux.Panel
	managed
	ModifiedCallback func()
	ValidateCallback func() bool
	runes            []rune
	selectionStart   int
	selectionEnd     int
	selectionAnchor  int
	forceShowUntil   time.Time
	scrollOffset     float64
	showCursor       bool
	pending          bool
	extendByWord     bool
	invalid          bool
}

// New creates a new, empty, text field.
func New() *TextField {
	t := &TextField{}
	t.managed.initialize()
	t.InitTypeAndID(t)
	t.SetBorder(t.unfocusedBorder)
	t.SetFocusable(true)
	t.SetSizer(t.DefaultSizes)
	t.DrawCallback = t.DefaultDraw
	t.GainedFocusCallback = t.DefaultFocusGained
	t.LostFocusCallback = t.DefaultFocusLost
	t.MouseDownCallback = t.DefaultMouseDown
	t.MouseDragCallback = t.DefaultMouseDrag
	t.UpdateCursorCallback = t.DefaultUpdateCursor
	t.KeyDownCallback = t.DefaultKeyDown
	t.CanPerformCmdCallback = t.DefaultCanPerformCmd
	t.PerformCmdCallback = t.DefaultPerformCmd
	return t
}

// DefaultSizes provides the default sizing.
func (t *TextField) DefaultSizes(hint geom.Size) (min, pref, max geom.Size) {
	var text string
	if len(t.runes) != 0 {
		text = string(t.runes)
	} else {
		text = "M"
	}
	minWidth := t.minimumTextWidth
	pref = t.font.Extents(text)
	if pref.Width < minWidth {
		pref.Width = minWidth
	}
	if b := t.Border(); b != nil {
		insets := b.Insets()
		pref.AddInsets(insets)
		minWidth += insets.Left + insets.Right
	}
	pref.GrowToInteger()
	if hint.Width >= 1 && hint.Width < minWidth {
		hint.Width = minWidth
	}
	pref.ConstrainForHint(hint)
	min = pref
	min.Width = minWidth
	return min, pref, layout.MaxSize(pref)
}

// DefaultDraw provides the default drawing.
func (t *TextField) DefaultDraw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	gc.Rect(dirty)
	gc.Fill(t.currentBackgroundInk())
	rect := t.ContentRect(false)
	gc.Rect(rect)
	gc.Clip()
	textTop := rect.Y + (rect.Height-t.font.Height())/2
	switch {
	case t.HasSelectionRange():
		left := rect.X + t.scrollOffset
		if t.selectionStart > 0 {
			pre := string(t.runes[:t.selectionStart])
			gc.DrawString(left, textTop, t.font, t.textInk, pre)
			left += t.font.Width(pre)
		}
		mid := string(t.runes[t.selectionStart:t.selectionEnd])
		right := rect.X + t.font.Width(string(t.runes[:t.selectionEnd])) + t.scrollOffset
		selRect := geom.Rect{Point: geom.Point{X: left, Y: textTop}, Size: geom.Size{Width: right - left, Height: t.font.Height()}}
		if t.Focused() {
			gc.Rect(selRect)
			gc.Fill(t.selectedTextBackgroundInk)
		} else {
			gc.SetStrokeWidth(2)
			selRect.InsetUniform(0.5)
			gc.Rect(selRect)
			gc.Stroke(t.selectedTextBackgroundInk)
		}
		gc.DrawString(left, textTop, t.font, t.selectedTextInk, mid)
		if t.selectionStart < len(t.runes) {
			gc.DrawString(right, textTop, t.font, t.textInk, string(t.runes[t.selectionEnd:]))
		}
	case len(t.runes) == 0:
		if t.watermark != "" {
			gc.DrawString(rect.X, textTop, t.font, t.watermarkInk, t.watermark)
		}
	default:
		gc.DrawString(rect.X+t.scrollOffset, textTop, t.font, t.textInk, string(t.runes))
	}
	if !t.HasSelectionRange() && t.Focused() {
		if t.showCursor {
			x := rect.X + t.font.Width(string(t.runes[:t.selectionEnd])) + t.scrollOffset
			gc.MoveTo(x, textTop)
			gc.LineTo(x, textTop+t.font.Height()-1)
			gc.Stroke(t.textInk)
		}
		t.scheduleBlink()
	}
}

// Invalid returns true if the field is currently marked as invalid.
func (t *TextField) Invalid() bool {
	return t.invalid
}

func (t *TextField) currentBackgroundInk() draw.Ink {
	switch {
	case t.invalid:
		return t.invalidBackgroundInk
	case !t.Enabled():
		return t.disabledBackgroundInk
	default:
		return t.backgroundInk
	}
}

func (t *TextField) scheduleBlink() {
	window := t.Window()
	if window != nil && window.IsValid() && !t.pending && t.Focused() {
		t.pending = true
		ux.InvokeAfter(t.blink, t.blinkRate)
	}
}

func (t *TextField) blink() {
	window := t.Window()
	if window != nil && window.IsValid() {
		t.pending = false
		if time.Now().After(t.forceShowUntil) {
			t.showCursor = !t.showCursor
			t.MarkForRedraw()
		}
		t.scheduleBlink()
	}
}

// DefaultFocusGained provides the default focus gained handling.
func (t *TextField) DefaultFocusGained() {
	t.SetBorder(t.focusedBorder)
	if !t.HasSelectionRange() {
		t.SelectAll()
	}
	t.showCursor = true
	t.MarkForRedraw()
}

// DefaultFocusLost provides the default focus lost handling.
func (t *TextField) DefaultFocusLost() {
	t.SetBorder(t.unfocusedBorder)
	if !t.CanSelectAll() {
		t.SetSelectionToStart()
	}
	t.MarkForRedraw()
}

// DefaultMouseDown provides the default mouse down handling.
func (t *TextField) DefaultMouseDown(where geom.Point, button, clickCount int, mod keys.Modifiers) bool {
	t.RequestFocus()
	if button == ux.ButtonLeft {
		t.extendByWord = false
		switch clickCount {
		case 2:
			start, end := t.findWordAt(t.ToSelectionIndex(where.X))
			t.SetSelection(start, end)
			t.extendByWord = true
		case 3:
			t.SelectAll()
		default:
			oldAnchor := t.selectionAnchor
			t.selectionAnchor = t.ToSelectionIndex(where.X)
			var start, end int
			if mod.ShiftDown() {
				if oldAnchor > t.selectionAnchor {
					start = t.selectionAnchor
					end = oldAnchor
				} else {
					start = oldAnchor
					end = t.selectionAnchor
				}
			} else {
				start = t.selectionAnchor
				end = t.selectionAnchor
			}
			t.setSelection(start, end, t.selectionAnchor)
		}
		return true
	}
	return false
}

// DefaultMouseDrag provides the default mouse drag handling.
func (t *TextField) DefaultMouseDrag(where geom.Point, button int, mod keys.Modifiers) {
	oldAnchor := t.selectionAnchor
	pos := t.ToSelectionIndex(where.X)
	var start, end int
	if t.extendByWord {
		s1, e1 := t.findWordAt(oldAnchor)
		var dir int
		if pos > s1 {
			dir = -1
		} else {
			dir = 1
		}
		for {
			start, end = t.findWordAt(pos)
			if start != end {
				if start > s1 {
					start = s1
				}
				if end < e1 {
					end = e1
				}
				break
			}
			pos += dir
			if dir > 0 && pos >= s1 || dir < 0 && pos <= e1 {
				start = s1
				end = e1
				break
			}
		}
	} else {
		if pos > oldAnchor {
			start = oldAnchor
			end = pos
		} else {
			start = pos
			end = oldAnchor
		}
	}
	t.setSelection(start, end, oldAnchor)
}

// DefaultUpdateCursor provides the default cursor update handling.
func (t *TextField) DefaultUpdateCursor(where geom.Point) *draw.Cursor {
	if t.Enabled() {
		return draw.TextCursor
	}
	return draw.ArrowCursor
}

// DefaultKeyDown provides the default key down handling.
func (t *TextField) DefaultKeyDown(keyCode int, ch rune, mod keys.Modifiers, repeat bool) bool {
	draw.HideCursorUntilMouseMoves()
	switch keyCode {
	case keys.Backspace.Code:
		t.Delete()
	case keys.Delete.Code, keys.NumpadDelete.Code:
		if t.HasSelectionRange() {
			t.Delete()
		} else if t.selectionStart < len(t.runes) {
			t.runes = append(t.runes[:t.selectionStart], t.runes[t.selectionStart+1:]...)
			t.notifyOfModification()
		}
		t.MarkForRedraw()
	case keys.Left.Code, keys.NumpadLeft.Code:
		extend := mod.ShiftDown()
		if mod.CommandDown() {
			t.handleHome(extend)
		} else {
			t.handleArrowLeft(extend, mod.OptionDown())
		}
	case keys.Right.Code, keys.NumpadRight.Code:
		extend := mod.ShiftDown()
		if mod.CommandDown() {
			t.handleEnd(extend)
		} else {
			t.handleArrowRight(extend, mod.OptionDown())
		}
	case keys.End.Code, keys.NumpadEnd.Code, keys.PageDown.Code, keys.NumpadPageDown.Code, keys.Down.Code, keys.NumpadDown.Code:
		t.handleEnd(mod.ShiftDown())
	case keys.Home.Code, keys.NumpadHome.Code, keys.PageUp.Code, keys.NumpadPageUp.Code, keys.Up.Code, keys.NumpadUp.Code:
		t.handleHome(mod.ShiftDown())
	default:
		if unicode.IsControl(ch) {
			return false
		}
		if t.HasSelectionRange() {
			t.runes = append(t.runes[:t.selectionStart], t.runes[t.selectionEnd:]...)
		}
		t.runes = append(t.runes[:t.selectionStart], append([]rune{ch}, t.runes[t.selectionStart:]...)...)
		t.SetSelectionTo(t.selectionStart + 1)
		t.notifyOfModification()
	}
	return true
}

func (t *TextField) handleHome(extend bool) {
	if extend {
		t.setSelection(0, t.selectionEnd, t.selectionEnd)
	} else {
		t.SetSelectionToStart()
	}
}

func (t *TextField) handleEnd(extend bool) {
	if extend {
		t.SetSelection(t.selectionStart, len(t.runes))
	} else {
		t.SetSelectionToEnd()
	}
}

func (t *TextField) handleArrowLeft(extend, byWord bool) {
	if t.HasSelectionRange() {
		if extend {
			anchor := t.selectionAnchor
			if t.selectionStart == anchor {
				pos := t.selectionEnd - 1
				if byWord {
					start, _ := t.findWordAt(pos)
					pos = xmath.MinInt(xmath.MaxInt(start, anchor), pos)
				}
				t.setSelection(anchor, pos, anchor)
			} else {
				pos := t.selectionStart - 1
				if byWord {
					start, _ := t.findWordAt(pos)
					pos = xmath.MinInt(start, pos)
				}
				t.setSelection(pos, anchor, anchor)
			}
		} else {
			t.SetSelectionTo(t.selectionStart)
		}
	} else {
		pos := t.selectionStart - 1
		if byWord {
			start, _ := t.findWordAt(pos)
			pos = xmath.MinInt(start, pos)
		}
		if extend {
			t.setSelection(pos, t.selectionStart, t.selectionEnd)
		} else {
			t.SetSelectionTo(pos)
		}
	}
}

func (t *TextField) handleArrowRight(extend, byWord bool) {
	if t.HasSelectionRange() {
		if extend {
			anchor := t.selectionAnchor
			if t.selectionEnd == anchor {
				pos := t.selectionStart + 1
				if byWord {
					_, end := t.findWordAt(pos)
					pos = xmath.MaxInt(xmath.MinInt(end, anchor), pos)
				}
				t.setSelection(pos, anchor, anchor)
			} else {
				pos := t.selectionEnd + 1
				if byWord {
					_, end := t.findWordAt(pos)
					pos = xmath.MaxInt(end, pos)
				}
				t.setSelection(anchor, pos, anchor)
			}
		} else {
			t.SetSelectionTo(t.selectionEnd)
		}
	} else {
		pos := t.selectionEnd + 1
		if byWord {
			_, end := t.findWordAt(pos)
			pos = xmath.MaxInt(end, pos)
		}
		if extend {
			t.SetSelection(t.selectionStart, pos)
		} else {
			t.SetSelectionTo(pos)
		}
	}
}

// DefaultCanPerformCmd provides the default can perform command handling.
func (t *TextField) DefaultCanPerformCmd(source interface{}, id int) bool {
	switch id {
	case ids.CutItemID:
		return t.CanCut()
	case ids.CopyItemID:
		return t.CanCopy()
	case ids.PasteItemID:
		return t.CanPaste()
	case ids.DeleteItemID:
		return t.CanDelete()
	case ids.SelectAllItemID:
		return t.CanSelectAll()
	default:
		return false
	}
}

// DefaultPerformCmd provides the default perform command handling.
func (t *TextField) DefaultPerformCmd(source interface{}, id int) {
	switch id {
	case ids.CutItemID:
		t.Cut()
	case ids.CopyItemID:
		t.Copy()
	case ids.PasteItemID:
		t.Paste()
	case ids.DeleteItemID:
		t.Delete()
	case ids.SelectAllItemID:
		t.SelectAll()
	default:
	}
}

// CanCut returns true if the field has a selection that can be cut.
func (t *TextField) CanCut() bool {
	return t.HasSelectionRange()
}

// Cut the selected text to the clipboard.
func (t *TextField) Cut() {
	if t.HasSelectionRange() {
		clipboard.SetDataWithType([]byte(t.SelectedText()), datatypes.PlainText)
		t.Delete()
	}
}

// CanCopy returns true if the field has a selection that can be copied.
func (t *TextField) CanCopy() bool {
	return t.HasSelectionRange()
}

// Copy the selected text to the clipboard.
func (t *TextField) Copy() {
	if t.HasSelectionRange() {
		clipboard.SetDataWithType([]byte(t.SelectedText()), datatypes.PlainText)
	}
}

// CanPaste returns true if the clipboard has content that can be pasted into
// the field.
func (t *TextField) CanPaste() bool {
	return clipboard.HasType(datatypes.PlainText)
}

// Paste any text on the clipboard into the field.
func (t *TextField) Paste() {
	if clipboard.HasType(datatypes.PlainText) {
		runes := []rune(sanitize(string(clipboard.GetFirstData(datatypes.PlainText))))
		if t.HasSelectionRange() {
			t.runes = append(t.runes[:t.selectionStart], t.runes[t.selectionEnd:]...)
		}
		t.runes = append(t.runes[:t.selectionStart], append(runes, t.runes[t.selectionStart:]...)...)
		t.SetSelectionTo(t.selectionStart + len(runes))
		t.notifyOfModification()
	} else if t.HasSelectionRange() {
		t.Delete()
	}
}

// CanDelete returns true if the field has a selection that can be deleted.
func (t *TextField) CanDelete() bool {
	return t.HasSelectionRange() || t.selectionStart > 0
}

// Delete removes the currently selected text, if any.
func (t *TextField) Delete() {
	if t.CanDelete() {
		if t.HasSelectionRange() {
			t.runes = append(t.runes[:t.selectionStart], t.runes[t.selectionEnd:]...)
			t.SetSelectionTo(t.selectionStart)
		} else {
			t.runes = append(t.runes[:t.selectionStart-1], t.runes[t.selectionStart:]...)
			t.SetSelectionTo(t.selectionStart - 1)
		}
		t.notifyOfModification()
		t.MarkForRedraw()
	}
}

// CanSelectAll returns true if the field's selection can be expanded.
func (t *TextField) CanSelectAll() bool {
	return t.selectionStart != 0 || t.selectionEnd != len(t.runes)
}

// SelectAll selects all of the text in the field.
func (t *TextField) SelectAll() {
	t.SetSelection(0, len(t.runes))
}

// Text returns the content of the field.
func (t *TextField) Text() string {
	return string(t.runes)
}

// SetText sets the content of the field.
func (t *TextField) SetText(text string) *TextField {
	text = sanitize(text)
	if string(t.runes) != text {
		t.runes = []rune(text)
		t.SetSelectionToEnd()
		t.notifyOfModification()
	}
	return t
}

func (t *TextField) notifyOfModification() {
	t.MarkForRedraw()
	if t.ModifiedCallback != nil {
		t.ModifiedCallback()
	}
	t.Validate()
}

// Validate forces field content validation to be run.
func (t *TextField) Validate() {
	invalid := false
	if t.ValidateCallback != nil {
		invalid = !t.ValidateCallback()
	}
	if invalid != t.invalid {
		t.invalid = invalid
		t.MarkForRedraw()
	}
}

func sanitize(text string) string {
	return strings.NewReplacer("\n", "", "\r", "").Replace(text)
}

// SelectedText returns the currently selected text.
func (t *TextField) SelectedText() string {
	return string(t.runes[t.selectionStart:t.selectionEnd])
}

// HasSelectionRange returns true is a selection range is currently present.
func (t *TextField) HasSelectionRange() bool {
	return t.selectionStart < t.selectionEnd
}

// SelectionCount returns the number of characters currently selected.
func (t *TextField) SelectionCount() int {
	return t.selectionEnd - t.selectionStart
}

// Selection returns the current start and end selection indexes.
func (t *TextField) Selection() (start, end int) {
	return t.selectionStart, t.selectionEnd
}

// SetSelectionToStart moves the cursor to the beginning of the text and
// removes any range that may have been present.
func (t *TextField) SetSelectionToStart() {
	t.SetSelection(0, 0)
}

// SetSelectionToEnd moves the cursor to the end of the text and removes any
// range that may have been present.
func (t *TextField) SetSelectionToEnd() {
	t.SetSelection(math.MaxInt64, math.MaxInt64)
}

// SetSelectionTo moves the cursor to the specified index and removes any
// range that may have been present.
func (t *TextField) SetSelectionTo(pos int) {
	t.SetSelection(pos, pos)
}

// SetSelection sets the start and end range of the selection. Values beyond
// either end will be constrained to the appropriate end. Likewise, an end
// value less than the start value will be treated as if the start and end
// values were the same.
func (t *TextField) SetSelection(start, end int) {
	t.setSelection(start, end, start)
}

func (t *TextField) setSelection(start, end, anchor int) {
	length := len(t.runes)
	if start < 0 {
		start = 0
	} else if start > length {
		start = length
	}
	if end < start {
		end = start
	} else if end > length {
		end = length
	}
	if anchor < start {
		anchor = start
	} else if anchor > end {
		anchor = end
	}
	if t.selectionStart != start || t.selectionEnd != end || t.selectionAnchor != anchor {
		t.selectionStart = start
		t.selectionEnd = end
		t.selectionAnchor = anchor
		t.forceShowUntil = time.Now().Add(t.blinkRate)
		t.showCursor = true
		t.MarkForRedraw()
		t.ScrollIntoView()
		t.autoScroll()
	}
}

func (t *TextField) autoScroll() {
	rect := t.ContentRect(false)
	if rect.Width > 0 {
		original := t.scrollOffset
		if t.selectionStart == t.selectionAnchor {
			right := t.FromSelectionIndex(t.selectionEnd).X
			if right < rect.X {
				t.scrollOffset = 0
				t.scrollOffset = rect.X - t.FromSelectionIndex(t.selectionEnd).X
			} else if right >= rect.X+rect.Width {
				t.scrollOffset = 0
				t.scrollOffset = rect.X + rect.Width - 1 - t.FromSelectionIndex(t.selectionEnd).X
			}
		} else {
			left := t.FromSelectionIndex(t.selectionStart).X
			if left < rect.X {
				t.scrollOffset = 0
				t.scrollOffset = rect.X - t.FromSelectionIndex(t.selectionStart).X
			} else if left >= rect.X+rect.Width {
				t.scrollOffset = 0
				t.scrollOffset = rect.X + rect.Width - 1 - t.FromSelectionIndex(t.selectionStart).X
			}
		}
		save := t.scrollOffset
		t.scrollOffset = 0
		min := rect.X + rect.Width - 1 - t.FromSelectionIndex(len(t.runes)).X
		if min > 0 {
			min = 0
		}
		max := rect.X - t.FromSelectionIndex(0).X
		if max < 0 {
			max = 0
		}
		if save < min {
			save = min
		} else if save > max {
			save = max
		}
		t.scrollOffset = save
		if original != t.scrollOffset {
			t.MarkForRedraw()
		}
	}
}

// ToSelectionIndex returns the rune index for the specified x-coordinate.
func (t *TextField) ToSelectionIndex(x float64) int {
	rect := t.ContentRect(false)
	return t.font.IndexForPosition(x-(rect.X+t.scrollOffset), string(t.runes))
}

// FromSelectionIndex returns a location in local coordinates for the
// specified rune index.
func (t *TextField) FromSelectionIndex(index int) geom.Point {
	rect := t.ContentRect(false)
	x := rect.X + t.scrollOffset
	top := rect.Y + rect.Height/2
	if index > 0 {
		length := len(t.runes)
		if index > length {
			index = length
		}
		x += t.font.PositionForIndex(index, string(t.runes))
	}
	return geom.Point{X: x, Y: top}
}

func (t *TextField) findWordAt(pos int) (start, end int) {
	length := len(t.runes)
	if pos < 0 {
		pos = 0
	} else if pos >= length {
		pos = length - 1
	}
	start = pos
	end = pos
	if length > 0 && !unicode.IsSpace(t.runes[start]) {
		for start > 0 && !unicode.IsSpace(t.runes[start-1]) {
			start--
		}
		for end < length && !unicode.IsSpace(t.runes[end]) {
			end++
		}
	}
	return start, end
}

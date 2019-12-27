package inkwell

import (
	"time"

	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/clipboard"
	"github.com/richardwilkes/ux/clipboard/datatypes"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/keys"
	"github.com/richardwilkes/ux/widget"
)

// Mask is used to limit the types of ink permitted in the ink well.
type Mask uint8

// Possible ink well masks.
const (
	ColorInkWellMask Mask = 1 << iota
	PatternInkWellMask
	GradientInkWellMask
)

// InkWell represents a control that holds and lets a user choose an ink.
type InkWell struct {
	ux.Panel
	managed
	ink                   draw.Ink
	InkChangedCallback    func()
	ClickCallback         func()
	ValidateImageCallback func(*draw.Image) *draw.Image
	mask                  Mask
	Pressed               bool
	dragInProgress        bool
}

// New creates a new InkWell.
func New() *InkWell {
	well := &InkWell{
		ink:  draw.ControlBackgroundInk,
		mask: ColorInkWellMask | PatternInkWellMask | GradientInkWellMask,
	}
	well.managed.initialize()
	well.InitTypeAndID(well)
	well.SetFocusable(true)
	well.SetSizer(well.DefaultSizes)
	well.ClickCallback = well.DefaultClick
	well.DrawCallback = well.DefaultDraw
	well.GainedFocusCallback = well.MarkForRedraw
	well.LostFocusCallback = well.MarkForRedraw
	well.MouseDownCallback = well.DefaultMouseDown
	well.MouseDragCallback = well.DefaultMouseDrag
	well.MouseUpCallback = well.DefaultMouseUp
	well.KeyDownCallback = well.DefaultKeyDown
	well.DragEnteredCallback = well.DefaultDragEntered
	well.DragUpdatedCallback = well.DefaultDragUpdated
	well.DragExitedCallback = well.DefaultDragExited
	well.DragEndedCallback = well.DefaultDragExited
	well.DropIsAcceptableCallback = well.DefaultIsDropAcceptable
	well.DropCallback = well.DefaultDrop
	well.DropFinishedCallback = well.DefaultDropFinished
	return well
}

// AllowedTypes returns the types of ink allowed to be set via SetInk().
func (well *InkWell) AllowedTypes() Mask {
	return well.mask
}

// SetAllowedTypes sets the types of ink allowed to be set via SetInk().
func (well *InkWell) SetAllowedTypes(mask Mask) *InkWell {
	well.mask = mask
	return well
}

// Ink returns the well's ink.
func (well *InkWell) Ink() draw.Ink {
	return well.ink
}

// SetInk sets the ink well's ink.
func (well *InkWell) SetInk(ink draw.Ink) *InkWell {
	if ink == nil {
		ink = draw.ControlBackgroundInk
	}
	switch ink.(type) {
	case draw.Color, *draw.Color:
		if well.mask&ColorInkWellMask == 0 {
			return well
		}
	case *draw.Pattern:
		if well.mask&PatternInkWellMask == 0 {
			return well
		}
	case *draw.Gradient:
		if well.mask&GradientInkWellMask == 0 {
			return well
		}
	}
	if ink != well.ink {
		well.ink = ink
		well.MarkForRedraw()
		if well.InkChangedCallback != nil {
			well.InkChangedCallback()
		}
	}
	return well
}

// DefaultSizes provides the default sizing.
func (well *InkWell) DefaultSizes(hint geom.Size) (min, pref, max geom.Size) {
	pref.Width = 4 + well.contentSize
	pref.Height = 4 + well.contentSize
	if border := well.Border(); border != nil {
		pref.AddInsets(border.Insets())
	}
	pref.GrowToInteger()
	pref.ConstrainForHint(hint)
	return pref, pref, pref
}

// DefaultDraw provides the default drawing.
func (well *InkWell) DefaultDraw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	if !well.Enabled() {
		gc.SetOpacity(0.33)
	}
	r := well.ContentRect(false)
	widget.DrawRoundedRectBase(gc, r, well.cornerRadius, well.currentBackgroundInk(), well.edgeInk)
	const wellInset = 1.5
	r.InsetUniform(wellInset)
	gc.RoundedRect(r, well.cornerRadius-wellInset)
	if p, ok := well.ink.(*draw.Pattern); ok {
		gc.Save()
		gc.Clip()
		p.Image().DrawInRect(gc, r)
		gc.Restore()
	} else {
		gc.Fill(well.ink)
	}
	gc.RoundedRect(r, well.cornerRadius-wellInset)
	gc.Stroke(well.edgeHighlightInk)
}

func (well *InkWell) currentBackgroundInk() draw.Ink {
	switch {
	case well.Pressed || well.dragInProgress:
		return well.pressedBackgroundInk
	case well.Focused():
		return well.focusedBackgroundInk
	default:
		return well.backgroundInk
	}
}

// DefaultMouseDown provides the default mouse down handling.
func (well *InkWell) DefaultMouseDown(where geom.Point, button, clickCount int, mod keys.Modifiers) bool {
	well.Pressed = true
	well.MarkForRedraw()
	return true
}

// DefaultMouseDrag provides the default mouse drag handling.
func (well *InkWell) DefaultMouseDrag(where geom.Point, button int, mod keys.Modifiers) {
	rect := well.ContentRect(false)
	pressed := rect.ContainsPoint(where)
	if well.Pressed != pressed {
		well.Pressed = pressed
		well.MarkForRedraw()
	}
}

// DefaultMouseUp provides the default mouse up handling.
func (well *InkWell) DefaultMouseUp(where geom.Point, button int, mod keys.Modifiers) {
	well.Pressed = false
	well.MarkForRedraw()
	rect := well.ContentRect(false)
	if rect.ContainsPoint(where) {
		if well.ClickCallback != nil {
			well.ClickCallback()
		}
	}
}

// DefaultKeyDown provides the default key down handling.
func (well *InkWell) DefaultKeyDown(keyCode int, ch rune, mod keys.Modifiers, repeat bool) bool {
	if keys.IsControlAction(keyCode) {
		well.Click()
		return true
	}
	return false
}

// DefaultClick provides the default click handling, which shows a dialog for
// selecting an ink.
func (well *InkWell) DefaultClick() {
	showDialog(well)
}

// Click makes the ink well behave as if a user clicked on it.
func (well *InkWell) Click() {
	pressed := well.Pressed
	well.Pressed = true
	well.MarkForRedraw()
	well.FlushDrawing()
	well.Pressed = pressed
	time.Sleep(well.clickAnimationTime)
	well.MarkForRedraw()
	if well.ClickCallback != nil {
		well.ClickCallback()
	}
}

// DefaultDragEntered provides the default drag entered behavior.
func (well *InkWell) DefaultDragEntered(dragInfo *ux.DragInfo) ux.DragOperation {
	op := well.determineDragOperation(dragInfo)
	if op&dragInfo.SourceOperationMask != 0 {
		well.dragInProgress = true
		well.MarkForRedraw()
	}
	return op
}

// DefaultDragUpdated provides the default drag updated behavior.
func (well *InkWell) DefaultDragUpdated(dragInfo *ux.DragInfo) ux.DragOperation {
	op := well.determineDragOperation(dragInfo)
	dragStillInProgress := op&dragInfo.SourceOperationMask != 0
	if dragStillInProgress != well.dragInProgress {
		well.dragInProgress = dragStillInProgress
		well.MarkForRedraw()
	}
	return op
}

// DefaultDragExited provides the default drag exited behavior.
func (well *InkWell) DefaultDragExited() {
	if well.dragInProgress {
		well.dragInProgress = false
		well.MarkForRedraw()
	}
}

// DefaultIsDropAcceptable provides the default is drop acceptable behavior.
func (well *InkWell) DefaultIsDropAcceptable(dragInfo *ux.DragInfo) bool {
	if ux.DragOperationCopy&dragInfo.SourceOperationMask != ux.DragOperationCopy {
		return false
	}
	dt := well.findBestDataType(dragInfo)
	if dt == datatypes.None {
		return false
	}
	count := 0
	for _, one := range dragInfo.DataForType(dt) {
		if draw.DistillImageURL(clipboard.BytesToURL(one)) != "" {
			count++
		}
	}
	dragInfo.ValidItemsForDrop = count
	return dragInfo.ValidItemsForDrop > 0
}

// DefaultDrop provides the default drop behavior.
func (well *InkWell) DefaultDrop(dragInfo *ux.DragInfo) bool {
	if ux.DragOperationCopy&dragInfo.SourceOperationMask != ux.DragOperationCopy {
		return false
	}
	dt := well.findBestDataType(dragInfo)
	if dt == datatypes.None {
		return false
	}
	for _, one := range dragInfo.DataForType(dt) {
		if urlStr := draw.DistillImageURL(clipboard.BytesToURL(one)); urlStr != "" {
			img, err := draw.NewImageFromURL(urlStr, well.imageScale)
			if err != nil {
				jot.Warn(err)
				continue
			}
			if well.ValidateImageCallback != nil {
				img = well.ValidateImageCallback(img)
			}
			if img != nil {
				p := draw.NewPattern(img)
				well.SetInk(p)
				return true
			}
		}
	}
	return false
}

// DefaultDropFinished provides the default drop finished behavior.
func (well *InkWell) DefaultDropFinished(dragInfo *ux.DragInfo) {
	well.dragInProgress = false
	well.MarkForRedraw()
}

func (well *InkWell) determineDragOperation(dragInfo *ux.DragInfo) ux.DragOperation {
	if well.DefaultIsDropAcceptable(dragInfo) {
		return ux.DragOperationCopy
	}
	return ux.DragOperationNone
}

func (well *InkWell) findBestDataType(dragInfo *ux.DragInfo) datatypes.DataType {
	return dragInfo.FirstTypePresent(datatypes.FileURL, datatypes.URL)
}

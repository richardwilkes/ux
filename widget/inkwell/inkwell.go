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

// InkWell represents a control that holds and lets a user choose an ink.
type InkWell struct {
	ux.Panel
	InkChangedCallback    func()
	ClickCallback         func()
	ValidateImageCallback func(*draw.Image) *draw.Image
	WellInk               draw.Ink
	BackgroundInk         draw.Ink      // The background ink when enabled but not selected, pressed or focused
	FocusedBackgroundInk  draw.Ink      // The background ink when enabled and focused
	PressedBackgroundInk  draw.Ink      // The background ink when enabled and pressed
	EdgeInk               draw.Ink      // The ink to use on the edges
	EdgeHighlightInk      draw.Ink      // The ink to use just inside the edges
	ClickAnimationTime    time.Duration // The amount of time to spend animating the click action
	CornerRadius          float64       // The amount of rounding to use on the corners
	ContentSize           float64       // The content width and height
	ImageScale            float64       // The image scale to use for images dropped onto the well. Defaults to 0.5 to support retina displays.
	Pressed               bool
	dragInProgress        bool
}

// New creates a new InkWell.
func New(ink draw.Ink) *InkWell {
	well := &InkWell{
		WellInk:              ink,
		BackgroundInk:        draw.ControlBackgroundInk,
		FocusedBackgroundInk: draw.ControlFocusedBackgroundInk,
		PressedBackgroundInk: draw.ControlPressedBackgroundInk,
		EdgeInk:              draw.ControlEdgeAdjColor,
		EdgeHighlightInk:     draw.ControlEdgeHighlightAdjColor,
		ClickAnimationTime:   time.Millisecond * 100,
		CornerRadius:         4,
		ContentSize:          20,
		ImageScale:           0.5,
	}
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

// SetInk sets the ink well's ink. Use this rather than directly setting
// WellInk to trigger the InkChangedCallback.
func (well *InkWell) SetInk(ink draw.Ink) {
	if ink != well.WellInk {
		well.WellInk = ink
		well.MarkForRedraw()
		if well.InkChangedCallback != nil {
			well.InkChangedCallback()
		}
	}
}

// DefaultSizes provides the default sizing.
func (well *InkWell) DefaultSizes(hint geom.Size) (min, pref, max geom.Size) {
	pref.Width = 4 + well.ContentSize
	pref.Height = 4 + well.ContentSize
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
	widget.DrawRoundedRectBase(gc, r, well.CornerRadius, well.currentBackgroundInk(), well.EdgeInk)
	r.InsetUniform(1.5)
	gc.RoundedRect(r, well.CornerRadius)
	if p, ok := well.WellInk.(*draw.Pattern); ok {
		gc.Save()
		gc.Clip()
		p.Image().DrawInRect(gc, r)
		gc.Restore()
	} else {
		gc.Fill(well.WellInk)
	}
	gc.RoundedRect(r, well.CornerRadius)
	gc.Stroke(well.EdgeHighlightInk)
}

func (well *InkWell) currentBackgroundInk() draw.Ink {
	switch {
	case well.Pressed || well.dragInProgress:
		return well.PressedBackgroundInk
	case well.Focused():
		return well.FocusedBackgroundInk
	default:
		return well.BackgroundInk
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

// DefaultClick provides the default click handling.
func (well *InkWell) DefaultClick() {
	jot.Debug("clicked on InkWell")
}

// Click makes the ink well behave as if a user clicked on it.
func (well *InkWell) Click() {
	pressed := well.Pressed
	well.Pressed = true
	well.MarkForRedraw()
	well.FlushDrawing()
	well.Pressed = pressed
	time.Sleep(well.ClickAnimationTime)
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
			img, err := draw.NewImageFromURL(urlStr, well.ImageScale)
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

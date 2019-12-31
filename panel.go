package ux

import (
	"fmt"
	"reflect"
	"strings"
	"sync/atomic"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/keys"
	"github.com/richardwilkes/ux/layout"
)

var (
	_            layout.Layoutable = &Panel{}
	nextGlobalID uint64
)

// Panel is the basic user interface element that interacts with the user.
type Panel struct {
	id                                  uint64
	self                                interface{}
	parent                              *Panel
	frame                               geom.Rect
	border                              border.Border
	sizer                               layout.Sizer
	layout                              layout.Layout
	layoutData                          interface{}
	children                            []*Panel
	Tooltip                             *Panel
	data                                map[string]interface{}
	DrawCallback                        func(gc draw.Context, dirty geom.Rect, inLiveResize bool)
	DrawOverCallback                    func(gc draw.Context, dirty geom.Rect, inLiveResize bool)
	GainedFocusCallback                 func()
	LostFocusCallback                   func()
	MouseDownCallback                   func(where geom.Point, button, clickCount int, mod keys.Modifiers) bool
	MouseDragCallback                   func(where geom.Point, button int, mod keys.Modifiers)
	MouseUpCallback                     func(where geom.Point, button int, mod keys.Modifiers)
	MouseEnterCallback                  func(where geom.Point, mod keys.Modifiers)
	MouseMoveCallback                   func(where geom.Point, mod keys.Modifiers)
	MouseExitCallback                   func()
	MouseWheelCallback                  func(where, delta geom.Point, mod keys.Modifiers) bool
	KeyDownCallback                     func(keyCode int, ch rune, mod keys.Modifiers, repeat bool) bool
	KeyUpCallback                       func(keyCode int, mod keys.Modifiers) bool
	UpdateCursorCallback                func(where geom.Point) *draw.Cursor
	UpdateTooltipCallback               func(where geom.Point, suggestedAvoid geom.Rect) geom.Rect
	DragEnteredCallback                 func(dragInfo *DragInfo) DragOperation
	DragUpdatedCallback                 func(dragInfo *DragInfo) DragOperation
	DragExitedCallback                  func()
	DragEndedCallback                   func()
	DropIsAcceptableCallback            func(dragInfo *DragInfo) bool
	DropCallback                        func(dragInfo *DragInfo) bool
	DropFinishedCallback                func(dragInfo *DragInfo)
	CanPerformCmdCallback               func(source interface{}, id int) bool
	PerformCmdCallback                  func(source interface{}, id int)
	FrameChangeCallback                 func()
	FrameChangeInChildHierarchyCallback func(panel *Panel)
	ScrollRectIntoViewCallback          func(rect geom.Rect) bool
	ParentChangedCallback               func()
	NeedsLayout                         bool
	focusable                           bool
	disabled                            bool
}

// NewPanel creates a new panel.
func NewPanel() *Panel {
	p := &Panel{}
	p.InitTypeAndID(p)
	return p
}

// InitTypeAndID initializes the panel with the appropriate
// self-identification information.
func (p *Panel) InitTypeAndID(self interface{}) {
	p.id = atomic.AddUint64(&nextGlobalID, 1)
	p.self = self
}

// ID returns the unique ID for the panel.
func (p *Panel) ID() uint64 {
	if p.id == 0 {
		jot.Fatal(1, errs.New("InitTypeAndID() must be called before use"))
	}
	return p.id
}

// Self returns the actual panel.
func (p *Panel) Self() interface{} {
	return p.self
}

// AsPanel returns this object as a panel.
func (p *Panel) AsPanel() *Panel {
	return p
}

// Is returns true if this panel is the other panel.
func (p *Panel) Is(other *Panel) bool {
	return p != nil && other != nil && p.ID() == other.ID()
}

func (p *Panel) String() string {
	name := reflect.Indirect(reflect.ValueOf(p.self)).Type().String()
	if i := strings.LastIndex(name, "."); i != -1 {
		name = name[i+1:]
	}
	return fmt.Sprintf("%s[%d]", name, p.id)
}

// Children returns the direct descendents of this panel.
func (p *Panel) Children() []*Panel {
	return p.children
}

// ChildrenForLayout is the same as calling Children(), but returns them as
// layout.Layoutable objects instead.
func (p *Panel) ChildrenForLayout() []layout.Layoutable {
	children := make([]layout.Layoutable, len(p.children))
	for i := range p.children {
		children[i] = p.children[i]
	}
	return children
}

// IndexOfChild returns the index of the specified child, or -1 if the
// passed in panel is not a child of this panel.
func (p *Panel) IndexOfChild(child *Panel) int {
	for i, one := range p.children {
		if one.Is(child) {
			return i
		}
	}
	return -1
}

// AddChild adds child to this panel, removing it from any previous parent it
// may have had.
func (p *Panel) AddChild(child *Panel) {
	child.RemoveFromParent()
	p.children = append(p.children, child)
	child.parent = p
	p.NeedsLayout = true
	if child.ParentChangedCallback != nil {
		child.ParentChangedCallback()
	}
}

// AddChildAtIndex adds child to this panel at the index, removing it from any
// previous parent it may have had. Passing in a negative value for the index
// will add it to the end.
func (p *Panel) AddChildAtIndex(child *Panel, index int) {
	child.RemoveFromParent()
	if index < 0 || index >= len(p.children) {
		p.children = append(p.children, child)
	} else {
		p.children = append(p.children, nil)
		copy(p.children[index+1:], p.children[index:])
		p.children[index] = child
	}
	child.parent = p
	p.NeedsLayout = true
	if child.ParentChangedCallback != nil {
		child.ParentChangedCallback()
	}
}

// RemoveAllChildren removes all child panels from this panel.
func (p *Panel) RemoveAllChildren() {
	children := p.children
	for _, child := range children {
		child.parent = nil
	}
	p.children = nil
	p.NeedsLayout = true
	for _, child := range children {
		if child.ParentChangedCallback != nil {
			child.ParentChangedCallback()
		}
	}
}

// RemoveChild removes 'child' from this panel. If 'child' is not a direct
// descendent of this panel, nothing happens.
func (p *Panel) RemoveChild(child *Panel) {
	p.RemoveChildAtIndex(p.IndexOfChild(child))
}

// RemoveChildAtIndex removes the child panel at 'index' from this panel.
// If 'index' is out of range, nothing happens.
func (p *Panel) RemoveChildAtIndex(index int) {
	if index >= 0 && index < len(p.children) {
		child := p.children[index]
		child.parent = nil
		copy(p.children[index:], p.children[index+1:])
		p.children[len(p.children)-1] = nil
		p.children = p.children[:len(p.children)-1]
		p.NeedsLayout = true
		if child.ParentChangedCallback != nil {
			child.ParentChangedCallback()
		}
	}
}

// RemoveFromParent removes this panel from its parent, if any.
func (p *Panel) RemoveFromParent() {
	if p.parent != nil {
		p.parent.RemoveChild(p)
	}
}

// Parent returns the parent panel, if any.
func (p *Panel) Parent() *Panel {
	return p.parent
}

// Window returns the containing window, if any.
func (p *Panel) Window() *Window {
	var prev *Panel
	panel := p
	for {
		if panel == nil {
			if prev != nil {
				if root, ok := prev.Self().(*rootPanel); ok {
					return root.window
				}
			}
			return nil
		}
		prev = panel
		panel = panel.parent
	}
}

// FrameRect returns the location and size of the panel in its parent's
// coordinate system.
func (p *Panel) FrameRect() geom.Rect {
	return p.frame
}

// SetFrameRect sets the location and size of the panel in its parent's
// coordinate system.
func (p *Panel) SetFrameRect(rect geom.Rect) {
	moved := p.frame.X != rect.X || p.frame.Y != rect.Y
	resized := p.frame.Width != rect.Width || p.frame.Height != rect.Height
	if moved || resized {
		p.MarkForRedraw()
		if moved {
			p.frame.Point = rect.Point
		}
		if resized {
			p.frame.Size = rect.Size
			p.NeedsLayout = true
		}
		if p.FrameChangeCallback != nil {
			p.FrameChangeCallback()
		}
		parent := p.parent
		for parent != nil {
			if parent.FrameChangeInChildHierarchyCallback != nil {
				parent.FrameChangeInChildHierarchyCallback(p)
			}
			parent = parent.parent
		}
		p.MarkForRedraw()
	}
}

// ContentRect returns the location and size of the panel in local
// coordinates.
func (p *Panel) ContentRect(includeBorder bool) geom.Rect {
	rect := p.frame.CopyAndZeroLocation()
	if !includeBorder && p.border != nil {
		rect.Inset(p.border.Insets())
	}
	return rect
}

// Border returns the border for this panel, if any.
func (p *Panel) Border() border.Border {
	return p.border
}

// SetBorder sets the border for this panel. May be nil.
func (p *Panel) SetBorder(b border.Border) *Panel {
	if p.border != b {
		p.border = b
		p.MarkForLayoutAndRedraw()
	}
	return p
}

// Sizer returns the sizer for this panel, if any.
func (p *Panel) Sizer() layout.Sizer {
	return p.sizer
}

// SetSizer sets the sizer for this panel. May be nil.
func (p *Panel) SetSizer(sizer layout.Sizer) {
	p.sizer = sizer
	p.NeedsLayout = true
}

// Sizes returns the minimum, preferred, and maximum sizes the panel wishes to
// be. It does this by first asking the panel's layout. If no layout is
// present, then the panel's sizer is asked. If no sizer is present, then it
// finally uses a default set of sizes that are used for all panels.
func (p *Panel) Sizes(hint geom.Size) (min, pref, max geom.Size) {
	if p.layout != nil {
		return p.layout.Sizes(hint)
	}
	if p.sizer != nil {
		return p.sizer(hint)
	}
	return geom.Size{}, geom.Size{}, geom.Size{Width: layout.DefaultMaxSize, Height: layout.DefaultMaxSize}
}

// Layout returns the Layout for this panel, if any.
func (p *Panel) Layout() layout.Layout {
	return p.layout
}

// SetLayout sets the Layout for this panel. May be nil.
func (p *Panel) SetLayout(lay layout.Layout) {
	p.layout = lay
	p.NeedsLayout = true
}

// ValidateLayout performs any layout that needs to be run by this panel or
// its children.
func (p *Panel) ValidateLayout() {
	if p.NeedsLayout {
		if p.layout != nil {
			p.layout.Layout()
			p.MarkForRedraw()
		}
		p.NeedsLayout = false
	}
	for _, child := range p.children {
		child.ValidateLayout()
	}
}

// LayoutData returns the layout data, if any, associated with this panel.
func (p *Panel) LayoutData() interface{} {
	return p.layoutData
}

// SetLayoutData sets layout data on this panel. May be nil.
func (p *Panel) SetLayoutData(data interface{}) {
	p.layoutData = data
	p.NeedsLayout = true
}

// MarkForLayoutAndRedraw marks this panel as needing to be laid out as well
// as redrawn at the next update.
func (p *Panel) MarkForLayoutAndRedraw() {
	p.NeedsLayout = true
	p.MarkForRedraw()
}

// MarkForRedraw marks this panel for drawing at the next update.
func (p *Panel) MarkForRedraw() {
	p.MarkRectForRedraw(p.ContentRect(true))
}

// MarkRectForRedraw marks the rect in local coordinates within the panel for
// drawing at the next update.
func (p *Panel) MarkRectForRedraw(rect geom.Rect) {
	rect.Intersect(p.ContentRect(true))
	if !rect.IsEmpty() {
		if p.parent != nil {
			rect.X += p.frame.X
			rect.Y += p.frame.Y
			p.parent.MarkRectForRedraw(rect)
		} else if w := p.Window(); w != nil {
			w.MarkRectForRedraw(rect)
		}
	}
}

// FlushDrawing is a convenience for calling the parent window's (if any)
// FlushDrawing() method.
func (p *Panel) FlushDrawing() {
	if w := p.Window(); w != nil {
		w.FlushDrawing()
	}
}

// Draw is called by its owning window when a panel needs to be drawn. The
// gc has already had its clip set to the dirty rectangle.
func (p *Panel) Draw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	dirty.Intersect(p.ContentRect(true))
	if !dirty.IsEmpty() {
		gc.Save()
		gc.Rect(dirty)
		gc.Clip()
		if p.DrawCallback != nil {
			gc.Save()
			p.DrawCallback(gc, dirty, inLiveResize)
			gc.Restore()
		}
		for _, child := range p.children {
			adjusted := dirty
			adjusted.Intersect(child.frame)
			if !adjusted.IsEmpty() {
				gc.Save()
				gc.Translate(child.frame.X, child.frame.Y)
				adjusted.X -= child.frame.X
				adjusted.Y -= child.frame.Y
				child.Draw(gc, adjusted, inLiveResize)
				gc.Restore()
			}
		}
		if p.border != nil {
			gc.Save()
			p.border.Draw(gc, p.ContentRect(true), inLiveResize)
			gc.Restore()
		}
		if p.DrawOverCallback != nil {
			p.DrawOverCallback(gc, dirty, inLiveResize)
		}
		gc.Restore()
	}
}

// Enabled returns true if this panel is currently enabled and can receive
// events.
func (p *Panel) Enabled() bool {
	return !p.disabled
}

// SetEnabled sets this panel's enabled state.
func (p *Panel) SetEnabled(enabled bool) *Panel {
	if p.disabled == enabled {
		p.disabled = !enabled
		p.MarkForRedraw()
	}
	return p
}

// Focusable returns true if this panel can have the keyboard focus.
func (p *Panel) Focusable() bool {
	return p.focusable && !p.disabled
}

// SetFocusable sets whether this panel can have the keyboard focus.
func (p *Panel) SetFocusable(focusable bool) *Panel {
	if p.focusable != focusable {
		p.focusable = focusable
	}
	return p
}

// Focused returns true if this panel has the keyboard focus.
func (p *Panel) Focused() bool {
	if wnd := p.Window(); wnd != nil {
		return wnd.Focused() && p.Is(wnd.Focus())
	}
	return false
}

// RequestFocus attempts to make this panel the keyboard focus.
func (p *Panel) RequestFocus() {
	if wnd := p.Window(); wnd != nil {
		wnd.SetFocus(p)
	}
}

// PanelAt returns the leaf-most child panel containing the point, or this
// panel if no child is found.
func (p *Panel) PanelAt(pt geom.Point) *Panel {
	for _, child := range p.children {
		if child.frame.ContainsPoint(pt) {
			pt.Subtract(child.frame.Point)
			return child.PanelAt(pt)
		}
	}
	return p
}

// PointToRoot converts panel-local coordinates into root coordinates, which
// when rooted within a window, will be window-local coordinates.
func (p *Panel) PointToRoot(pt geom.Point) geom.Point {
	pt.Add(p.frame.Point)
	parent := p.parent
	for parent != nil {
		pt.Add(parent.frame.Point)
		parent = parent.parent
	}
	return pt
}

// PointFromRoot converts root coordinates (i.e. window-local, when rooted
// within a window) into panel-local coordinates.
func (p *Panel) PointFromRoot(pt geom.Point) geom.Point {
	pt.Subtract(p.frame.Point)
	parent := p.parent
	for parent != nil {
		pt.Subtract(parent.frame.Point)
		parent = parent.parent
	}
	return pt
}

// RectToRoot converts panel-local coordinates into root coordinates, which
// when rooted within a window, will be window-local coordinates.
func (p *Panel) RectToRoot(rect geom.Rect) geom.Rect {
	rect.Point = p.PointToRoot(rect.Point)
	return rect
}

// RectFromRoot converts root coordinates (i.e. window-local, when rooted
// within a window) into panel-local coordinates.
func (p *Panel) RectFromRoot(rect geom.Rect) geom.Rect {
	rect.Point = p.PointFromRoot(rect.Point)
	return rect
}

// ScrollIntoView attempts to scroll this panel into the current view if it is
// not already there, using ScrollAreas in this panel's hierarchy.
func (p *Panel) ScrollIntoView() {
	p.ScrollRectIntoView(p.ContentRect(true))
}

// ScrollRectIntoView attempts to scroll the rect (in coordinates local to
// this panel) into the current view if it is not already there, using
// ScrollAreas in this panel's hierarchy.
func (p *Panel) ScrollRectIntoView(rect geom.Rect) {
	look := p
	for look != nil {
		if look.ScrollRectIntoViewCallback != nil {
			if look.ScrollRectIntoViewCallback(rect) {
				return
			}
		}
		rect.Point.Add(look.frame.Point)
		look = look.parent
	}
}

// ClientData returns a map of client data for this panel.
func (p *Panel) ClientData() map[string]interface{} {
	if p.data == nil {
		p.data = make(map[string]interface{})
	}
	return p.data
}

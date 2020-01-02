// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package ux

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/clipboard/datatypes"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/keys"
)

// Constants for mouse buttons.
const (
	ButtonLeft  = 0
	ButtonRight = 1
)

// StyleMask controls the look and capabilities of a window.
type StyleMask int

// Possible values for the StyleMask.
const (
	TitledWindowMask StyleMask = 1 << iota
	ClosableWindowMask
	MinimizableWindowMask
	ResizableWindowMask
	NoInternalMenuWindowMask = 1 << 30 // Has no effect on macOS
	BorderlessWindowMask     = 0
	StdWindowMask            = TitledWindowMask | ClosableWindowMask | MinimizableWindowMask | ResizableWindowMask
)

// Window holds window information.
type Window struct {
	id uint64
	// MayCloseCallback is called when the user has requested that the window
	// be closed. Return true to permit it, false to cancel the operation.
	// Defaults to always returning true.
	MayCloseCallback func() bool
	// WillCloseCallback is called just prior to the window closing.
	WillCloseCallback func()
	// GainedFocusCallback is called when the keyboard focus is gained on this
	// window.
	GainedFocusCallback func()
	// LostFocusCallback is called when the keyboard focus is lost from this
	// window.
	LostFocusCallback func()
	// MouseDownCallback is called when the mouse is pressed within this
	// window.
	MouseDownCallback func(where geom.Point, button, clickCount int, mod keys.Modifiers)
	// MouseDragCallback is called when the mouse is dragged after being
	// pressed within this window.
	MouseDragCallback func(where geom.Point, button int, mod keys.Modifiers)
	// MouseUpCallback is called when the mouse is released after being
	// pressed within this window.
	MouseUpCallback func(where geom.Point, button int, mod keys.Modifiers)
	// MouseEnterCallback is called when the mouse enters this window.
	MouseEnterCallback func(where geom.Point, mod keys.Modifiers)
	// MouseMoveCallback is called when the mouse moves within this window.
	MouseMoveCallback func(where geom.Point, mod keys.Modifiers)
	// MouseExitCallback is called when the mouse exits this window.
	MouseExitCallback func()
	// MouseWheelCallback is called when the mouse wheel is rotated over this
	// window.
	MouseWheelCallback func(where, delta geom.Point, mod keys.Modifiers)
	// KeyDownCallback is called when a key is pressed in this window.
	KeyDownCallback func(keyCode int, ch rune, mod keys.Modifiers, repeat bool)
	// KeyDownCallback is called when a key is released in this window.
	KeyUpCallback func(keyCode int, mod keys.Modifiers)
	// DragEnteredCallback is called when a drag enters this window.
	DragEnteredCallback func(di *DragInfo) DragOperation
	// DragUpdatedCallback is called when a drag moves within this window or
	// the drag operation changes.
	DragUpdatedCallback func(di *DragInfo) DragOperation
	// DragExitedCallback is called when a drag exits the window.
	DragExitedCallback func()
	// DragEndedCallback is called when a drag ends.
	DragEndedCallback func()
	// DropIsAcceptableCallback is called when the drag is dropped to
	// determine if it is still valid to drop here.
	DropIsAcceptableCallback func(di *DragInfo) bool
	// DropCallback is called to perform the actual drop. Returns true on
	// success.
	DropCallback func(di *DragInfo) bool
	// DropFinishedCallback is called when the drop concludes.
	DropFinishedCallback func(di *DragInfo)
	data                 map[string]interface{}
	title                string
	root                 *rootPanel
	focus                *Panel
	cursor               *draw.Cursor
	lastMouseDownPanel   *Panel
	lastMouseOverPanel   *Panel
	lastKeyDownPanel     *Panel
	lastDragPanel        *Panel
	lastTooltip          *Panel
	background           draw.Ink
	lastTooltipShownAt   time.Time
	style                StyleMask
	tooltipSequence      int
	diacritics           keys.Diacritics
	wnd                  OSWindow
	valid                bool
}

var windowList []*Window

// WindowCount returns the number of windows that are open.
func WindowCount() int {
	return len(windowList)
}

// Windows returns a slice containing the current set of open windows.
func Windows() []*Window {
	list := make([]*Window, len(windowList))
	copy(list, windowList)
	return list
}

// WindowWithFocus returns the window that currently has the keyboard focus,
// or nil if none of your application's windows has the keyboard focus.
func WindowWithFocus() *Window {
	return osKeyWindow()
}

// AppWindowsToFront attempts to bring all of the application's windows to the
// foreground.
func AppWindowsToFront() {
	osAppWindowsToFront()
}

// WindowContentRectForFrameRect determines the content rect for a window
// based on the given frame rect and window style.
func WindowContentRectForFrameRect(frame geom.Rect, style StyleMask) geom.Rect {
	return osWindowContentRectForFrameRect(frame, style)
}

// WindowFrameRectForContentRect determines the frame rect for a window
// based on the given content rect and window style.
func WindowFrameRectForContentRect(content geom.Rect, style StyleMask) geom.Rect {
	return osWindowFrameRectForContentRect(content, style)
}

// NewWindow creates a new window.
func NewWindow(title string, frame geom.Rect, style StyleMask) (*Window, error) {
	wnd, err := osNewWindow(title, frame, style)
	if err != nil {
		return nil, err
	}
	w := &Window{
		id:         atomic.AddUint64(&nextGlobalID, 1),
		title:      title,
		background: draw.WindowBackgroundColor,
		style:      style,
		wnd:        wnd,
		valid:      true,
	}
	w.GainedFocusCallback = w.focusGained
	w.LostFocusCallback = w.focusLost
	w.MouseDownCallback = w.mouseDown
	w.MouseDragCallback = w.mouseDrag
	w.MouseUpCallback = w.mouseUp
	w.MouseEnterCallback = w.mouseEnter
	w.MouseMoveCallback = w.mouseMove
	w.MouseExitCallback = w.mouseExit
	w.MouseWheelCallback = w.mouseWheel
	w.KeyDownCallback = w.keyDown
	w.KeyUpCallback = w.keyUp
	w.DragEnteredCallback = w.dragEntered
	w.DragUpdatedCallback = w.dragUpdated
	w.DragExitedCallback = w.dragExited
	w.DragEndedCallback = w.dragEnded
	w.DropIsAcceptableCallback = w.dropIsAcceptable
	w.DropCallback = w.drop
	w.DropFinishedCallback = w.dropFinished
	windowList = append(windowList, w)
	w.osAddNativeWindow()
	w.root = newRootPanel(w)
	w.ValidateLayout()
	w.MarkForRedraw()
	return w, nil
}

// RunModal displays and brings this window to the front, the runs a modal
// event loop until StopModal is called.
func (w *Window) RunModal() int {
	return w.osRunModal()
}

// StopModal stops the current modal event loop, closes the window, and
// propagates the provided code as the result to RunModal().
func (w *Window) StopModal(code int) {
	w.osStopModal(code)
	w.Dispose()
}

// IsValid returns true if the window is still valid (i.e. hasn't been
// disposed).
func (w *Window) IsValid() bool {
	return w.valid
}

// OSWindow returns the underlying OS window object.
func (w *Window) OSWindow() OSWindow {
	return w.wnd
}

// ID returns the unique ID for the window.
func (w *Window) ID() uint64 {
	return w.id
}

func (w *Window) String() string {
	return fmt.Sprintf("Window[%s]", w.title)
}

// AttemptClose closes the window if permitted.
func (w *Window) AttemptClose() {
	if w.MayCloseCallback == nil || w.MayCloseCallback() {
		w.Dispose()
	}
}

// Dispose of the window.
func (w *Window) Dispose() {
	if w.WillCloseCallback != nil {
		w.WillCloseCallback()
		w.WillCloseCallback = nil
	}
	if w.root.content != nil {
		w.root.content.RemoveFromParent()
	}
	for i, wnd := range windowList {
		if w != wnd {
			continue
		}
		copy(windowList[i:], windowList[i+1:])
		count := len(windowList) - 1
		windowList[count] = nil
		windowList = windowList[:count]
		break
	}
	w.osRemoveNativeWindow()
	if w.valid {
		w.valid = false
		w.osDispose()
	}
}

// Title returns the title of this window.
func (w *Window) Title() string {
	return w.title
}

// SetTitle sets the title of this window.
func (w *Window) SetTitle(title string) {
	if w.IsValid() && w.title != title {
		w.title = title
		w.osSetTitle(title)
	}
}

// Content returns the content panel for the window.
func (w *Window) Content() *Panel {
	return w.root.content
}

// SetContent sets the content panel for the window.
func (w *Window) SetContent(panel *Panel) {
	w.root.setContent(panel)
	w.ValidateLayout()
	w.MarkForRedraw()
}

// ValidateLayout performs any layout that needs to be run by this window or
// its children.
func (w *Window) ValidateLayout() {
	rect := w.ContentRect()
	rect.X = 0
	rect.Y = 0
	w.root.SetFrameRect(rect)
	w.root.ValidateLayout()
}

// FrameRect returns the boundaries in display coordinates of the frame of
// this window (i.e. the area that includes both the content and its border
// and window controls).
func (w *Window) FrameRect() geom.Rect {
	if w.IsValid() {
		return w.osFrameRect()
	}
	return geom.Rect{}
}

// SetFrameRect sets the boundaries of the frame of this window.
func (w *Window) SetFrameRect(rect geom.Rect) {
	if w.IsValid() {
		current := w.ContentRect()
		w.osSetFrameRect(WindowFrameRectForContentRect(w.adjustContentRectForMinMax(WindowContentRectForFrameRect(rect, w.style)), w.style))
		adjusted := w.ContentRect()
		if current.Size != adjusted.Size {
			w.ValidateLayout()
		}
	}
}

func (w *Window) adjustContentRectForMinMax(rect geom.Rect) geom.Rect {
	min, _, max := w.root.Sizes(geom.Size{})
	if rect.Width < min.Width {
		rect.Width = min.Width
	} else if rect.Width > max.Width {
		rect.Width = max.Width
	}
	if rect.Height < min.Height {
		rect.Height = min.Height
	} else if rect.Height > max.Height {
		rect.Height = max.Height
	}
	return rect
}

// ContentRect returns the boundaries in display coordinates of the window's
// content area.
func (w *Window) ContentRect() geom.Rect {
	if w.IsValid() {
		return w.osContentRect()
	}
	return geom.Rect{}
}

// SetContentRect sets the boundaries of the frame of this window by
// converting the content rect into a suitable frame rect and then applying it
// to the window.
func (w *Window) SetContentRect(rect geom.Rect) {
	w.SetFrameRect(WindowFrameRectForContentRect(rect, w.style))
}

// Pack sets the window's content size to match the preferred size of the
// root panel.
func (w *Window) Pack() {
	_, pref, _ := w.root.Sizes(geom.Size{})
	rect := w.ContentRect()
	rect.Size = pref
	w.SetContentRect(rect)
}

// Focused returns true if the window has the current keyboard focus.
func (w *Window) Focused() bool {
	return w == WindowWithFocus()
}

// Focus returns the panel with the keyboard focus in this window.
func (w *Window) Focus() *Panel {
	if w.focus == nil {
		w.FocusNext()
	}
	return w.focus
}

// SetFocus sets the keyboard focus to the specified target.
func (w *Window) SetFocus(target *Panel) {
	if target != nil {
		tw := target.Window()
		if tw != nil && tw.id == w.id && !target.Is(w.focus) {
			if w.focus != nil && w.focus.LostFocusCallback != nil {
				w.focus.LostFocusCallback()
			}
			w.focus = target
			if w.focus != nil {
				if w.focus.GainedFocusCallback != nil {
					w.focus.GainedFocusCallback()
				}
				w.focus.ScrollIntoView()
			}
		}
	}
}

// FocusNext moves the keyboard focus to the next focusable panel.
func (w *Window) FocusNext() {
	if w.root.content != nil {
		current := w.focus
		if current == nil {
			current = w.root.content
		}
		i, focusables := collectFocusables(w.root.content, current, nil)
		if len(focusables) > 0 {
			i++
			if i >= len(focusables) {
				i = 0
			}
			current = focusables[i]
		}
		w.SetFocus(current)
	}
}

// FocusPrevious moves the keyboard focus to the previous focusable panel.
func (w *Window) FocusPrevious() {
	if w.root.content != nil {
		current := w.focus
		if current == nil {
			current = w.root.content
		}
		i, focusables := collectFocusables(w.root.content, current, nil)
		if len(focusables) > 0 {
			i--
			if i < 0 {
				i = len(focusables) - 1
			}
			current = focusables[i]
		}
		w.SetFocus(current)
	}
}

func collectFocusables(current, target *Panel, focusables []*Panel) (match int, result []*Panel) {
	match = -1
	if current.Focusable() {
		if current.Is(target) {
			match = len(focusables)
		}
		focusables = append(focusables, current)
	}
	for _, child := range current.Children() {
		var m int
		m, focusables = collectFocusables(child, target, focusables)
		if match == -1 && m != -1 {
			match = m
		}
	}
	return match, focusables
}

// ToFront attempts to bring the window to the foreground and give it the
// keyboard focus.
func (w *Window) ToFront() {
	if w.IsValid() {
		w.osToFront()
	}
}

// Minimize performs the minimize function on the window.
func (w *Window) Minimize() {
	if w.IsValid() {
		w.osMinimize()
	}
}

// Zoom performs the zoom function on the window.
func (w *Window) Zoom() {
	if w.IsValid() {
		w.osZoom()
	}
}

// Closable returns true if the window was created with the
// ClosableWindowMask.
func (w *Window) Closable() bool {
	return w.style&ClosableWindowMask != 0
}

// Minimizable returns true if the window was created with the
// MinimizableWindowMask.
func (w *Window) Minimizable() bool {
	return w.style&MinimizableWindowMask != 0
}

// Resizable returns true if the window was created with the
// ResizableWindowMask.
func (w *Window) Resizable() bool {
	return w.style&ResizableWindowMask != 0
}

// MouseLocation returns the current mouse location relative to this window.
func (w *Window) MouseLocation() geom.Point {
	return w.osMouseLocation()
}

// Draw the window contents.
func (w *Window) Draw(gc draw.Context, dirtyRect geom.Rect, inLiveResize bool) {
	if w.root != nil {
		w.root.ValidateLayout()
		gc.Rect(dirtyRect)
		gc.Fill(w.background)
		w.root.Draw(gc, dirtyRect, inLiveResize)
	}
}

// MarkForRedraw marks this window for drawing at the next update.
func (w *Window) MarkForRedraw() {
	if w.IsValid() {
		rect := w.ContentRect()
		rect.X = 0
		rect.Y = 0
		w.osMarkRectForRedraw(rect)
	}
}

// MarkRectForRedraw marks the rect in local coordinates within the window for
// drawing at the next update.
func (w *Window) MarkRectForRedraw(rect geom.Rect) {
	if w.IsValid() {
		cRect := w.ContentRect()
		cRect.X = 0
		cRect.Y = 0
		rect.Intersect(cRect)
		if !rect.IsEmpty() {
			w.osMarkRectForRedraw(rect)
		}
	}
}

// FlushDrawing causes any areas marked for drawing to be drawn now.
func (w *Window) FlushDrawing() {
	if w.IsValid() {
		w.osFlushDrawing()
	}
}

// Background returns the background of this window.
func (w *Window) Background() draw.Ink {
	return w.background
}

// SetBackground sets the background of this window.
func (w *Window) SetBackground(ink draw.Ink) {
	if ink != nil && ink != w.background {
		w.background = ink
		w.MarkForRedraw()
	}
}

// RegisterDragTypes registers the data types the window will accept in a
// drag & drop operation.
func (w *Window) RegisterDragTypes(dt ...datatypes.DataType) {
	w.osRegisterDragTypes(dt...)
}

func (w *Window) updateTooltipAndCursor(target *Panel, where geom.Point) {
	w.updateCursor(target, where)
	w.updateTooltip(target, where)
}

func (w *Window) updateTooltip(target *Panel, where geom.Point) {
	var avoid geom.Rect
	var tip *Panel
	for target != nil {
		avoid = target.ContentRect(true)
		avoid.Point = target.PointToRoot(avoid.Point)
		avoid.Align()
		if target.UpdateTooltipCallback != nil {
			avoid = target.UpdateTooltipCallback(target.PointFromRoot(where), avoid)
		}
		if target.Tooltip != nil {
			tip = target.Tooltip
			break
		}
		target = target.parent
	}
	if !w.lastTooltip.Is(tip) {
		wasShowing := w.root.tooltip != nil
		w.ClearTooltip()
		w.lastTooltip = tip
		if tip != nil {
			ts := &tooltipSequencer{window: w, avoid: avoid, sequence: w.tooltipSequence}
			if wasShowing || time.Since(w.lastTooltipShownAt) < TooltipDismissal {
				ts.show()
			} else {
				InvokeAfter(ts.show, TooltipDelay)
			}
		}
	}
}

// ClearTooltip clears any existing tooltip and resets the timer.
func (w *Window) ClearTooltip() {
	w.tooltipSequence++
	w.lastTooltipShownAt = time.Time{}
	w.root.setTooltip(nil)
}

func (w *Window) updateCursor(target *Panel, where geom.Point) {
	var cursor *draw.Cursor
	for target != nil {
		if target.UpdateCursorCallback == nil {
			target = target.parent
		} else {
			cursor = target.UpdateCursorCallback(target.PointFromRoot(where))
			break
		}
	}
	if cursor == nil {
		cursor = draw.ArrowCursor
	}
	if w.cursor != cursor {
		w.cursor = cursor
		cursor.MakeCurrent()
	}
}

func (w *Window) focusGained() {
	w.ClearTooltip()
	if w.focus == nil {
		w.FocusNext()
	}
	if w.focus != nil {
		w.focus.MarkForRedraw()
	}
}

func (w *Window) focusLost() {
	w.ClearTooltip()
	if w.focus != nil {
		w.focus.MarkForRedraw()
	}
}

func (w *Window) mouseDown(where geom.Point, button, clickCount int, mod keys.Modifiers) {
	if w.Focused() {
		w.ClearTooltip()
		w.lastMouseDownPanel = nil
		panel := w.root.PanelAt(where)
		for panel != nil {
			if panel.Enabled() && panel.MouseDownCallback != nil && panel.MouseDownCallback(panel.PointFromRoot(where), button, clickCount, mod) {
				w.lastMouseDownPanel = panel
				break
			}
			panel = panel.parent
		}
	}
}

func (w *Window) mouseDrag(where geom.Point, button int, mod keys.Modifiers) {
	if w.lastMouseDownPanel != nil && w.lastMouseDownPanel.MouseDragCallback != nil && w.lastMouseDownPanel.Enabled() {
		w.lastMouseDownPanel.MouseDragCallback(w.lastMouseDownPanel.PointFromRoot(where), button, mod)
	}
}

func (w *Window) mouseUp(where geom.Point, button int, mod keys.Modifiers) {
	if w.lastMouseDownPanel != nil && w.lastMouseDownPanel.MouseUpCallback != nil && w.lastMouseDownPanel.Enabled() {
		w.lastMouseDownPanel.MouseUpCallback(w.lastMouseDownPanel.PointFromRoot(where), button, mod)
	}
	if w.MouseExitCallback != nil && w.root != nil && !w.root.PanelAt(where).Is(w.lastMouseOverPanel) {
		w.MouseExitCallback()
	}
	w.updateTooltipAndCursor(w.lastMouseDownPanel, where)
	w.lastMouseDownPanel = nil
}

func (w *Window) mouseEnter(where geom.Point, mod keys.Modifiers) {
	if w.MouseExitCallback != nil {
		w.MouseExitCallback()
	}
	panel := w.root.PanelAt(where)
	if panel.MouseEnterCallback != nil {
		panel.MouseEnterCallback(panel.PointFromRoot(where), mod)
	}
	w.updateTooltipAndCursor(panel, where)
	w.lastMouseOverPanel = panel
}

func (w *Window) mouseMove(where geom.Point, mod keys.Modifiers) {
	panel := w.root.PanelAt(where)
	if panel.Is(w.lastMouseOverPanel) {
		if panel.MouseMoveCallback != nil {
			panel.MouseMoveCallback(panel.PointFromRoot(where), mod)
		}
		w.updateTooltipAndCursor(panel, where)
	} else if w.MouseEnterCallback != nil {
		w.MouseEnterCallback(where, mod)
	}
}

func (w *Window) mouseExit() {
	if w.lastMouseDownPanel == nil && w.lastMouseOverPanel != nil {
		if w.lastMouseOverPanel.MouseExitCallback != nil {
			w.lastMouseOverPanel.MouseExitCallback()
		}
		w.lastMouseOverPanel = nil
		w.cursor = nil
	}
}

func (w *Window) mouseWheel(where, delta geom.Point, mod keys.Modifiers) {
	panel := w.root.PanelAt(where)
	for panel != nil {
		if panel.Enabled() && panel.MouseWheelCallback != nil && panel.MouseWheelCallback(panel.PointFromRoot(where), delta, mod) {
			break
		}
		panel = panel.parent
	}
	if w.lastMouseDownPanel != nil {
		if w.MouseDragCallback != nil {
			w.MouseDragCallback(where, 0, mod)
		}
	} else if w.MouseMoveCallback != nil {
		w.MouseMoveCallback(where, mod)
	}
}

func (w *Window) keyDown(keyCode int, ch rune, mod keys.Modifiers, repeat bool) {
	w.ClearTooltip()
	w.lastKeyDownPanel = nil
	if focus := w.Focus(); focus != nil {
		ch = w.diacritics.ProcessInput(keyCode, ch, mod)
		panel := focus
		for panel != nil {
			if panel.Enabled() && panel.KeyDownCallback != nil && panel.KeyDownCallback(keyCode, ch, mod, repeat) {
				w.lastKeyDownPanel = panel
				return
			}
			panel = panel.parent
		}
		if keyCode == keys.Tab.Code && (mod&(keys.AllModifiers&^keys.ShiftModifier)) == 0 {
			if mod.ShiftDown() {
				w.FocusPrevious()
			} else {
				w.FocusNext()
			}
		}
	}
}

func (w *Window) keyUp(keyCode int, mod keys.Modifiers) {
	if w.lastKeyDownPanel != nil && w.lastKeyDownPanel.KeyUpCallback != nil {
		w.lastKeyDownPanel.KeyUpCallback(keyCode, mod)
	}
}

func (w *Window) dragEntered(di *DragInfo) DragOperation {
	if w.lastDragPanel != nil && w.DragExitedCallback != nil {
		w.DragExitedCallback()
	}
	where := geom.Point{X: di.DragX, Y: di.DragY}
	panel := w.root.PanelAt(where)
	op := DragOperationNone
	if panel.DragEnteredCallback != nil {
		delta := panel.PointFromRoot(where)
		delta.Subtract(where)
		di.ApplyOffset(delta.X, delta.Y)
		op = panel.DragEnteredCallback(di)
		di.ApplyOffset(-delta.X, -delta.Y)
	}
	w.lastDragPanel = panel
	return op
}

func (w *Window) dragUpdated(di *DragInfo) DragOperation {
	where := geom.Point{X: di.DragX, Y: di.DragY}
	panel := w.root.PanelAt(where)
	op := DragOperationNone
	if panel.Is(w.lastDragPanel) {
		if panel.DragUpdatedCallback != nil {
			delta := panel.PointFromRoot(where)
			delta.Subtract(where)
			di.ApplyOffset(delta.X, delta.Y)
			op = panel.DragUpdatedCallback(di)
			di.ApplyOffset(-delta.X, -delta.Y)
		}
	} else if w.DragEnteredCallback != nil {
		op = w.DragEnteredCallback(di)
	}
	return op
}

func (w *Window) dragExited() {
	if w.lastDragPanel != nil {
		if w.lastDragPanel.DragExitedCallback != nil {
			w.lastDragPanel.DragExitedCallback()
		}
		w.lastDragPanel = nil
	}
}

func (w *Window) dragEnded() {
	if w.lastDragPanel != nil {
		if w.lastDragPanel.DragEndedCallback != nil {
			w.lastDragPanel.DragEndedCallback()
		}
		w.lastDragPanel = nil
	}
}

func (w *Window) dropIsAcceptable(di *DragInfo) bool {
	where := geom.Point{X: di.DragX, Y: di.DragY}
	panel := w.root.PanelAt(where)
	var acceptable bool
	if panel.Is(w.lastDragPanel) {
		if panel.DropIsAcceptableCallback != nil {
			delta := panel.PointFromRoot(where)
			delta.Subtract(where)
			di.ApplyOffset(delta.X, delta.Y)
			acceptable = panel.DropIsAcceptableCallback(di)
			di.ApplyOffset(-delta.X, -delta.Y)
		}
	}
	return acceptable
}

func (w *Window) drop(di *DragInfo) bool {
	where := geom.Point{X: di.DragX, Y: di.DragY}
	panel := w.root.PanelAt(where)
	var accepted bool
	if panel.Is(w.lastDragPanel) {
		if panel.DropCallback != nil {
			delta := panel.PointFromRoot(where)
			delta.Subtract(where)
			di.ApplyOffset(delta.X, delta.Y)
			accepted = panel.DropCallback(di)
			di.ApplyOffset(-delta.X, -delta.Y)
		}
	}
	return accepted
}

func (w *Window) dropFinished(di *DragInfo) {
	where := geom.Point{X: di.DragX, Y: di.DragY}
	panel := w.root.PanelAt(where)
	if panel.Is(w.lastDragPanel) {
		if panel.DropFinishedCallback != nil {
			delta := panel.PointFromRoot(where)
			delta.Subtract(where)
			di.ApplyOffset(delta.X, delta.Y)
			panel.DropFinishedCallback(di)
			di.ApplyOffset(-delta.X, -delta.Y)
		}
	}
	w.lastDragPanel = nil
}

// ClientData returns a map of client data for this window.
func (w *Window) ClientData() map[string]interface{} {
	if w.data == nil {
		w.data = make(map[string]interface{})
	}
	return w.data
}

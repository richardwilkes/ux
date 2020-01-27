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
	"github.com/richardwilkes/macos/ca"
	"github.com/richardwilkes/macos/cg"
	"github.com/richardwilkes/macos/ns"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/clipboard"
	"github.com/richardwilkes/ux/clipboard/datatypes"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/keys"
)

var nativeWindowMap = make(map[ns.WindowNative]*Window)

// OSWindow represents the underlying native window.
type OSWindow = *ns.Window

func osKeyWindow() *Window {
	app := ns.SharedApplication()
	wnd := app.KeyWindow()
	if wnd == nil {
		wnd = app.MainWindow()
	}
	if wnd != nil {
		if w, ok := nativeWindowMap[wnd.Native()]; ok && w.IsValid() {
			return w
		}
	}
	return nil
}

func osAppWindowsToFront() {
	ns.RunningApplicationCurrent().ActivateWithOptions(ns.ApplicationActivateAllWindows | ns.ApplicationActivateIgnoringOtherApps)
}

func osWindowContentRectForFrameRect(frame geom.Rect, styleMask StyleMask) geom.Rect {
	_, _, _, screenHeight := ns.MainScreen().Frame()
	x, y, width, height := ns.WindowContentRectForFrameRectStyleMask(frame.X, screenHeight-(frame.Y+frame.Height), frame.Width, frame.Height, styleMaskToMacStyleMask(styleMask))
	return geom.Rect{
		Point: geom.Point{
			X: x,
			Y: screenHeight - (y + height),
		},
		Size: geom.Size{
			Width:  width,
			Height: height,
		},
	}
}

func osWindowFrameRectForContentRect(content geom.Rect, styleMask StyleMask) geom.Rect {
	_, _, _, screenHeight := ns.MainScreen().Frame()
	x, y, width, height := ns.WindowFrameRectForContentRectStyleMask(content.X, screenHeight-(content.Y+content.Height), content.Width, content.Height, styleMaskToMacStyleMask(styleMask))
	return geom.Rect{
		Point: geom.Point{
			X: x,
			Y: screenHeight - (y + height),
		},
		Size: geom.Size{
			Width:  width,
			Height: height,
		},
	}
}

func osNewWindow(title string, frame geom.Rect, styleMask StyleMask) (OSWindow, error) {
	frame = osWindowContentRectForFrameRect(frame, styleMask)
	w := ns.WindowInitWithContentRectStyleMask(frame.X, frame.Y, frame.Width, frame.Height, styleMaskToMacStyleMask(styleMask))
	w.SetIsKeyable(true)
	w.DisableCursorRects()
	w.SetDelegate(&windowDelegate{})
	rootView := ns.NewView(&viewDelegate{})
	w.SetContentView(rootView)
	w.SetTitle(title)
	rootView.AddTrackingArea(ns.TrackingAreaInitWithRectOptionsOwnerUserInfo(0, 0, frame.Width, frame.Height, ns.TrackingMouseEnteredAndExited|ns.TrackingMouseMoved|ns.TrackingActiveInKeyWindow|ns.TrackingInVisibleRect|ns.TrackingCursorUpdate, rootView, 0))
	// By default, register for the generic drag type so that most things just
	// come through without the user of the library having to make a call to
	// RegisterForDraggedTypes.
	w.RegisterForDraggedTypes(datatypes.Generic.UTI)
	return w, nil
}

func (w *Window) osRunModal() int {
	return ns.SharedApplication().RunModalForWindow(w.wnd)
}

func (w *Window) osStopModal(code int) {
	ns.SharedApplication().StopModalWithCode(code)
}

func (w *Window) osAddNativeWindow() {
	nativeWindowMap[w.wnd.Native()] = w
}

func (w *Window) osRemoveNativeWindow() {
	delete(nativeWindowMap, w.wnd.Native())
}

func (w *Window) osDispose() {
	view := w.wnd.ContentView()
	w.wnd.Close()
	view.ReleaseDelegate()
}

func (w *Window) osSetTitle(title string) {
	w.wnd.SetTitle(title)
}

func (w *Window) osFrameRect() geom.Rect {
	_, _, _, screenHeight := ns.MainScreen().Frame()
	x, y, width, height := w.wnd.Frame()
	return geom.Rect{
		Point: geom.Point{
			X: x,
			Y: screenHeight - (y + height),
		},
		Size: geom.Size{
			Width:  width,
			Height: height,
		},
	}
}

func (w *Window) osSetFrameRect(frame geom.Rect) {
	_, _, _, screenHeight := ns.MainScreen().Frame()
	w.wnd.SetFrame(frame.X, screenHeight-(frame.Y+frame.Height), frame.Width, frame.Height, true)
}

func (w *Window) osContentRect() geom.Rect {
	_, _, _, screenHeight := ns.MainScreen().Frame()
	wx, wy, _, _ := w.wnd.Frame()
	x, y, width, height := w.wnd.ContentView().Frame()
	x += wx
	y += wy
	return geom.Rect{
		Point: geom.Point{
			X: x,
			Y: screenHeight - (y + height),
		},
		Size: geom.Size{
			Width:  width,
			Height: height,
		},
	}
}

func (w *Window) osToFront() {
	w.wnd.MakeKeyAndOrderFront()
}

func (w *Window) osMinimize() {
	w.wnd.PerformMiniaturize()
}

func (w *Window) osZoom() {
	w.wnd.PerformZoom()
}

func (w *Window) osMouseLocation() geom.Point {
	x, y := w.wnd.MouseLocationOutsideOfEventStream()
	_, _, _, height := w.wnd.ContentView().Frame()
	return geom.Point{X: x, Y: height - y}
}

func (w *Window) osMarkRectForRedraw(rect geom.Rect) {
	w.wnd.ContentView().SetNeedsDisplayInRect(rect.X, rect.Y, rect.Width, rect.Height)
}

func (w *Window) osFlushDrawing() {
	ca.TransactionFlush()
}

func (w *Window) osRegisterDragTypes(dt ...datatypes.DataType) {
	types := make([]string, len(dt))
	for i := range dt {
		types[i] = dt[i].UTI
	}
	w.wnd.RegisterForDraggedTypes(types...)
}

type windowDelegate struct {
}

func (d *windowDelegate) WindowDidResize(wnd *ns.Window) {
	if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil {
		current := w.ContentRect()
		adjusted := w.adjustContentRectForMinMax(current)
		if adjusted != current {
			w.SetContentRect(adjusted)
		} else {
			w.ValidateLayout()
		}
	}
}

func (d *windowDelegate) WindowDidBecomeKey(wnd *ns.Window) {
	if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil {
		if w.GainedFocusCallback != nil {
			w.GainedFocusCallback()
		}
		if w.MouseEnterCallback != nil {
			w.MouseEnterCallback(w.osMouseLocation(), 0)
		}
	}
}

func (d *windowDelegate) WindowDidResignKey(wnd *ns.Window) {
	if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil && w.LostFocusCallback != nil {
		w.LostFocusCallback()
	}
}

func (d *windowDelegate) WindowShouldClose(wnd *ns.Window) bool {
	if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil && w.MayCloseCallback != nil {
		return w.MayCloseCallback()
	}
	return true
}

func (d *windowDelegate) WindowWillClose(wnd *ns.Window) {
	if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil {
		w.valid = false // Need to set this preemptively to avoid a double dispose
		w.Dispose()
	}
}

func (d *windowDelegate) WindowDragEntered(wnd *ns.Window, di *ns.DraggingInfo) ns.DragOperation {
	return d.windowDragEnteredOrUpdated(wnd, di, true)
}

func (d *windowDelegate) WindowDragUpdated(wnd *ns.Window, di *ns.DraggingInfo) ns.DragOperation {
	return d.windowDragEnteredOrUpdated(wnd, di, false)
}

func (d *windowDelegate) windowDragEnteredOrUpdated(wnd *ns.Window, di *ns.DraggingInfo, entered bool) ns.DragOperation {
	if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil {
		var callback func(dragInfo *DragInfo) DragOperation
		if entered {
			callback = w.DragEnteredCallback
		} else {
			callback = w.DragUpdatedCallback
		}
		if callback != nil {
			info := dragInfoFromNSDraggingInfo(di)
			originalCount := info.ValidItemsForDrop
			result := callback(info)
			if info.ValidItemsForDrop == 0 {
				return DragOperationNone
			}
			if originalCount != info.ValidItemsForDrop {
				di.SetNumberOfValidItemsForDrop(info.ValidItemsForDrop)
			}
			return result & info.SourceOperationMask
		}
	}
	return DragOperationNone
}

func (d *windowDelegate) WindowDragExited(wnd *ns.Window) {
	if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil {
		if w.DragExitedCallback != nil {
			w.DragExitedCallback()
		}
	}
}

func (d *windowDelegate) WindowDragEnded(wnd *ns.Window) {
	if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil {
		if w.DragEndedCallback != nil {
			w.DragEndedCallback()
		}
	}
}

func (d *windowDelegate) WindowDropIsAcceptable(wnd *ns.Window, di *ns.DraggingInfo) bool {
	if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil {
		if w.DropIsAcceptableCallback != nil {
			return w.DropIsAcceptableCallback(dragInfoFromNSDraggingInfo(di))
		}
	}
	return false
}

func (d *windowDelegate) WindowDrop(wnd *ns.Window, di *ns.DraggingInfo) bool {
	if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil {
		if w.DropCallback != nil {
			return w.DropCallback(dragInfoFromNSDraggingInfo(di))
		}
	}
	return false
}

func (d *windowDelegate) WindowDropFinished(wnd *ns.Window, di *ns.DraggingInfo) {
	if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil {
		if w.DropFinishedCallback != nil {
			w.DropFinishedCallback(dragInfoFromNSDraggingInfo(di))
		}
	}
}

func dragInfoFromNSDraggingInfo(di *ns.DraggingInfo) *DragInfo {
	pb := di.Pasteboard()
	info := &DragInfo{
		Sequence:            di.SequenceNumber(),
		SourceOperationMask: di.SourceOperationMask(),
		ValidItemsForDrop:   di.NumberOfValidItemsForDrop(),
		ItemTypes:           clipboard.TypesForPasteboard(pb),
		DataForType:         func(dataType datatypes.DataType) [][]byte { return clipboard.DataFromPasteboard(pb, dataType) },
	}
	info.DragX, info.DragY = di.Location()
	_, _, _, height := di.DestinationWindow().ContentView().Frame()
	info.DragY = height - info.DragY
	info.DragImageX, info.DragImageY = di.ImageLocation()
	info.DragImageY = height - info.DragImageY
	return info
}

type viewDelegate struct {
}

func (d *viewDelegate) ViewDraw(view *ns.View, gc cg.Context, x, y, width, height float64, inLiveResize bool) {
	if wnd := view.Window(); wnd != nil {
		if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil {
			draw.UpdateSystemColors()
			w.Draw(draw.NewContextForOSContext(&gc), geom.Rect{
				Point: geom.Point{X: x, Y: y},
				Size:  geom.Size{Width: width, Height: height},
			}, inLiveResize)
		}
	}
}

func (d *viewDelegate) ViewMouseDownEvent(view *ns.View, x, y float64, button, clickCount, mod int) {
	if wnd := view.Window(); wnd != nil {
		if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil && w.MouseDownCallback != nil {
			w.MouseDownCallback(geom.Point{X: x, Y: y}, button, clickCount, convertModifiers(mod))
		}
	}
}

func (d *viewDelegate) ViewMouseDragEvent(view *ns.View, x, y float64, button, mod int) {
	if wnd := view.Window(); wnd != nil {
		if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil && w.MouseDragCallback != nil {
			w.MouseDragCallback(geom.Point{X: x, Y: y}, button, convertModifiers(mod))
		}
	}
}

func (d *viewDelegate) ViewMouseUpEvent(view *ns.View, x, y float64, button, mod int) {
	if wnd := view.Window(); wnd != nil {
		if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil && w.MouseUpCallback != nil {
			w.MouseUpCallback(geom.Point{X: x, Y: y}, button, convertModifiers(mod))
		}
	}
}

func (d *viewDelegate) ViewMouseEnterEvent(view *ns.View, x, y float64, mod int) {
	if wnd := view.Window(); wnd != nil {
		if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil && w.MouseEnterCallback != nil {
			w.MouseEnterCallback(geom.Point{X: x, Y: y}, convertModifiers(mod))
		}
	}
}

func (d *viewDelegate) ViewMouseMoveEvent(view *ns.View, x, y float64, mod int) {
	if wnd := view.Window(); wnd != nil {
		if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil && w.MouseMoveCallback != nil {
			w.MouseMoveCallback(geom.Point{X: x, Y: y}, convertModifiers(mod))
		}
	}
}

func (d *viewDelegate) ViewMouseExitEvent(view *ns.View) {
	if wnd := view.Window(); wnd != nil {
		if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil && w.MouseExitCallback != nil {
			w.MouseExitCallback()
		}
	}
}

func (d *viewDelegate) ViewMouseWheelEvent(view *ns.View, x, y, dx, dy float64, mod int) {
	if wnd := view.Window(); wnd != nil {
		if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil && w.MouseWheelCallback != nil {
			w.MouseWheelCallback(geom.Point{X: x, Y: y}, geom.Point{X: dx, Y: dy}, convertModifiers(mod))
		}
	}
}

func (d *viewDelegate) ViewCursorUpdateEvent(view *ns.View, x, y float64, mod int) {
	d.ViewMouseEnterEvent(view, x, y, mod)
}

func (d *viewDelegate) ViewKeyDownEvent(view *ns.View, keyCode int, ch rune, mod int, repeat bool) {
	if wnd := view.Window(); wnd != nil {
		if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil && w.KeyDownCallback != nil {
			w.KeyDownCallback(keyCode, ch, convertModifiers(mod), repeat)
		}
	}
}

func (d *viewDelegate) ViewKeyUpEvent(view *ns.View, keyCode, mod int) {
	if wnd := view.Window(); wnd != nil {
		if w, ok := nativeWindowMap[wnd.Native()]; ok && w.wnd != nil && w.KeyUpCallback != nil {
			w.KeyUpCallback(keyCode, convertModifiers(mod))
		}
	}
}

func convertModifiers(mod int) keys.Modifiers {
	return keys.Modifiers((mod & (ns.EventModifierFlagCapsLock |
		ns.EventModifierFlagShift | ns.EventModifierFlagControl |
		ns.EventModifierFlagOption | ns.EventModifierFlagCommand)) >> 16)
}

func styleMaskToMacStyleMask(styleMask StyleMask) ns.WindowStyleMask {
	return ns.WindowStyleMask(styleMask &^ NoInternalMenuWindowMask)
}

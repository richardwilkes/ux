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
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/motif"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xrect"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/clipboard/datatypes"
	"github.com/richardwilkes/ux/globals"
)

var nativeWindowMap = make(map[xproto.Window]*Window)

type OSWindow = *xwindow.Window

func osKeyWindow() *Window {
	return nil // RAW: Implement
}

func osAppWindowsToFront() {
	// RAW: Implement
}

func osWindowContentRectForFrameRect(frame geom.Rect, styleMask StyleMask) geom.Rect {
	return frame // RAW: Implement
}

func osWindowFrameRectForContentRect(content geom.Rect, styleMask StyleMask) geom.Rect {
	return content // RAW: Implement
}

func osNewWindow(title string, frame geom.Rect, styleMask StyleMask) (OSWindow, error) {
	w, err := xwindow.Generate(globals.X11)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	if err = w.CreateChecked(globals.X11.RootWin(), int(frame.X), int(frame.Y), int(frame.Width), int(frame.Height), xproto.CwBackPixel|xproto.CwEventMask, uint32(0xffffff), xproto.EventMaskButtonRelease); err != nil {
		return nil, errs.Wrap(err)
	}
	if err = motif.WmHintsSet(globals.X11, w.Id, &motif.Hints{
		Flags:      motif.HintFunctions | motif.HintDecorations | motif.HintInputMode,
		Function:   functionsForStyleMask(styleMask),
		Decoration: decorationsForStyleMask(styleMask),
		Input:      motif.InputModeless,
	}); err != nil {
		jot.Error(errs.Wrap(err))
	}
	if title != "" {
		if err := ewmh.WmNameSet(globals.X11, w.Id, title); err != nil {
			jot.Error(errs.Wrap(err))
		}
	}
	w.WMGracefulClose(func(me *xwindow.Window) {
		if pw, ok := nativeWindowMap[me.Id]; ok && pw.IsValid() {
			pw.AttemptClose()
		}
	})
	w.Map()
	// For some reason, the initial coordinates are ignored... move it to the
	// asked for position.
	w.Move(int(frame.X), int(frame.Y))
	return w, nil
}

func (w *Window) osRunModal() int {
	return 0 // RAW: Need impl
}

func (w *Window) osStopModal(code int) {
	// RAW: Need impl
}

func (w *Window) osAddNativeWindow() {
	nativeWindowMap[w.wnd.Id] = w
}

func (w *Window) osRemoveNativeWindow() {
	delete(nativeWindowMap, w.wnd.Id)
}

func (w *Window) osDispose() {
	xevent.Detach(globals.X11, w.wnd.Id)
	mousebind.Detach(globals.X11, w.wnd.Id)
	w.wnd.Destroy()
	if len(nativeWindowMap) == 0 && (QuitAfterLastWindowClosedCallback == nil || QuitAfterLastWindowClosedCallback()) {
		AttemptQuit()
	}
}

func (w *Window) osSetTitle(title string) {
	if err := ewmh.WmNameSet(globals.X11, w.wnd.Id, title); err != nil {
		jot.Error(errs.Wrap(err))
	}
}

func (w *Window) osFrameRect() geom.Rect {
	r, err := w.wnd.DecorGeometry()
	if err != nil {
		jot.Error(errs.Wrap(err))
		return geom.Rect{}
	}
	return fromXRectToRect(r)
}

func (w *Window) osSetFrameRect(frame geom.Rect) {
	w.wnd.MoveResize(int(frame.X), int(frame.Y), int(frame.Width), int(frame.Height))
}

func (w *Window) osContentRect() geom.Rect {
	r, err := w.wnd.Geometry()
	if err != nil {
		jot.Error(errs.Wrap(err))
		return geom.Rect{}
	}
	return fromXRectToRect(r)
}

func (w *Window) osToFront() {
	if err := ewmh.RestackWindow(globals.X11, w.wnd.Id); err != nil {
		jot.Error(errs.Wrap(err))
	}
}

func (w *Window) osMinimize() {
	if err := icccm.WmStateSet(globals.X11, w.wnd.Id, &icccm.WmState{State: icccm.StateIconic}); err != nil {
		jot.Error(errs.Wrap(err))
	}
}

func (w *Window) osZoom() {
	if err := icccm.WmStateSet(globals.X11, w.wnd.Id, &icccm.WmState{State: icccm.StateZoomed}); err != nil {
		jot.Error(errs.Wrap(err))
	}
}

func (w *Window) osMouseLocation() geom.Point {
	// RAW: Implement
	return geom.Point{}
}

func (w *Window) osMarkRectForRedraw(rect geom.Rect) {
	// RAW: Implement
}

func (w *Window) osFlushDrawing() {
	// RAW: Implement
}

func (w *Window) osRegisterDragTypes(dt ...datatypes.DataType) {
	// RAW: Implement
}

func functionsForStyleMask(style StyleMask) uint {
	if styleMaskToLinuxStyleMask(style) == BorderlessWindowMask {
		return motif.FunctionNone
	}
	result := uint(motif.FunctionMove)
	if style&ClosableWindowMask == ClosableWindowMask {
		result |= motif.FunctionClose
	}
	if style&MinimizableWindowMask == MinimizableWindowMask {
		result |= motif.FunctionMinimize
	}
	if style&ResizableWindowMask == ResizableWindowMask {
		result |= motif.FunctionResize | motif.FunctionMaximize
	}
	return result
}

func decorationsForStyleMask(style StyleMask) uint {
	if styleMaskToLinuxStyleMask(style) == BorderlessWindowMask {
		return motif.DecorationNone
	}
	result := uint(motif.DecorationBorder | motif.DecorationTitle | motif.DecorationMenu)
	if style&MinimizableWindowMask == MinimizableWindowMask {
		result |= motif.DecorationMinimize
	}
	if style&ResizableWindowMask == ResizableWindowMask {
		result |= motif.DecorationResizeH | motif.DecorationMaximize
	}
	return result
}

func fromXRectToRect(in xrect.Rect) geom.Rect {
	return geom.Rect{
		Point: geom.Point{X: float64(in.X()), Y: float64(in.Y())},
		Size:  geom.Size{Width: float64(in.Width()), Height: float64(in.Height())},
	}
}

func styleMaskToLinuxStyleMask(styleMask StyleMask) StyleMask {
	return styleMask &^ NoInternalMenuWindowMask
}

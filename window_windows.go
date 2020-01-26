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
	"syscall"
	"unsafe"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/clipboard/datatypes"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/globals"
	"github.com/richardwilkes/win32"
	"github.com/richardwilkes/win32/d2d"
)

const windowClassName = "wndClass"

type OSWindow = win32.HWND

type wndInfo struct {
	wnd          *Window
	renderTarget *d2d.HWNDRenderTarget
}

var (
	windowClass win32.ATOM
	// MenuItemSelectionCallback is exposed as an implementation side-effect
	// and is not intended for client use.
	MenuItemSelectionCallback func(id int)
	// MenuValidationCallback is exposed as an implementation side-effect and
	// is not intended for client use.
	MenuValidationCallback func(hmenu win32.HMENU)
	nativeWindowMap        = make(map[win32.HWND]*wndInfo)
	d2dFactory             *d2d.Factory
)

func osKeyWindow() *Window {
	wnd := win32.GetForegroundWindow()
	if wnd != 0 {
		if wi, ok := nativeWindowMap[wnd]; ok && wi.wnd.IsValid() {
			return wi.wnd
		}
	}
	return nil
}

func osAppWindowsToFront() {
	list := make([]*Window, 0)
	win32.EnumWindows(func(wnd win32.HWND, data win32.LPARAM) win32.BOOL {
		if one, ok := nativeWindowMap[wnd]; ok {
			list = append(list, one.wnd)
		}
		return 1
	}, 0)
	for i, one := range list {
		after := win32.HWND_TOP
		flags := uint32(win32.SWP_NOMOVE | win32.SWP_NOSIZE)
		if i != 0 {
			flags |= win32.SWP_NOACTIVATE
			if list[i-1].IsValid() {
				after = list[i-1].wnd
			} else {
				after = 0
			}
		}
		if one.IsValid() {
			win32.SetWindowPos(one.wnd, after, 0, 0, 0, 0, flags)
		}
	}
}

func osWindowContentRectForFrameRect(frame geom.Rect, styleMask StyleMask) geom.Rect {
	rect := win32.RECT{Top: 100, Left: 100, Bottom: 300, Right: 300}
	style, exStyle := styleMaskToWin32Style(styleMask)
	win32.AdjustWindowRectEx(&rect, style, styleMask&NoInternalMenuWindowMask == 0, exStyle)
	frame.Inset(geom.Insets{Top: float64(100 - rect.Top), Left: float64(100 - rect.Left), Bottom: float64(rect.Bottom - 300), Right: float64(rect.Right - 300)})
	return frame
}

func osWindowFrameRectForContentRect(content geom.Rect, styleMask StyleMask) geom.Rect {
	rect := fromRectToWin32Rect(content)
	style, exStyle := styleMaskToWin32Style(styleMask)
	win32.AdjustWindowRectEx(&rect, style, styleMask&NoInternalMenuWindowMask == 0, exStyle)
	return fromWin32RectToRect(rect)
}

func osNewWindow(title string, frame geom.Rect, styleMask StyleMask) (OSWindow, error) {
	style, exStyle := styleMaskToWin32Style(styleMask)
	hwnd := win32.CreateWindowExS(exStyle, windowClassName, title, style, int32(frame.X), int32(frame.Y), int32(frame.Width), int32(frame.Height), win32.NULL, win32.NULL, globals.ModuleInstance, win32.NULL)
	if hwnd == win32.NULL {
		return 0, errs.New("unable to create window")
	}
	return hwnd, nil
}

func (w *Window) osRunModal() int {
	return 0 // RAW: Need impl
}

func (w *Window) osStopModal(code int) {
	// RAW: Need impl
}

func (w *Window) osAddNativeWindow() {
	nativeWindowMap[w.wnd] = &wndInfo{wnd: w}
}

func (w *Window) osRemoveNativeWindow() {
	delete(nativeWindowMap, w.wnd)
}

func (w *Window) osDispose() {
	// RAW: Dispose of any resources, such as d2d
	win32.DestroyWindow(w.wnd)
}

func (w *Window) osSetTitle(title string) {
	win32.SetWindowTextS(w.wnd, title)
}

func (w *Window) osFrameRect() geom.Rect {
	var rect win32.RECT
	win32.GetWindowRect(w.wnd, &rect)
	return fromWin32RectToRect(rect)
}

func (w *Window) osSetFrameRect(frame geom.Rect) {
	win32.MoveWindow(w.wnd, int32(frame.X), int32(frame.Y), int32(frame.Width), int32(frame.Height), true)
}

func (w *Window) osContentRect() geom.Rect {
	var rect win32.RECT
	win32.GetClientRect(w.wnd, &rect)
	win32.MapWindowRect(w.wnd, win32.HWND_DESKTOP, &rect)
	return fromWin32RectToRect(rect)
}

func (w *Window) osToFront() {
	win32.ShowWindow(w.wnd, win32.SW_SHOWNORMAL)
	win32.DrawMenuBar(w.wnd)
	win32.SetActiveWindow(w.wnd)
}

func (w *Window) osMinimize() {
	cmd := int32(win32.SW_MINIMIZE)
	if win32.IsIconic(w.wnd) {
		cmd = win32.SW_RESTORE
	}
	win32.ShowWindow(w.wnd, cmd)
}

func (w *Window) osZoom() {
	cmd := int32(win32.SW_MAXIMIZE)
	if win32.IsZoomed(w.wnd) {
		cmd = win32.SW_RESTORE
	}
	win32.ShowWindow(w.wnd, cmd)
}

func (w *Window) osMouseLocation() geom.Point {
	// RAW: Implement
	return geom.Point{}
}

func (w *Window) osMarkRectForRedraw(rect geom.Rect) {
	nativeRect := fromRectToWin32Rect(rect)
	win32.InvalidateRect(w.wnd, &nativeRect, false)
}

func (w *Window) osFlushDrawing() {
	win32.GdiFlush()
}

func (w *Window) osRegisterDragTypes(dt ...datatypes.DataType) {
	// RAW: Implement
}

// Make this private once the internal platform package is gone
func RegisterWindowClass() {
	wcx := win32.WNDCLASSEX{
		Style:    win32.CS_HREDRAW | win32.CS_VREDRAW,
		WndProc:  syscall.NewCallback(wndProc),
		Instance: globals.ModuleInstance,
		Cursor:   win32.LoadSystemCursor(win32.IDC_ARROW),
	}
	wcx.Size = uint32(unsafe.Sizeof(wcx)) //nolint:gosec
	var err error
	if wcx.ClassName, err = syscall.UTF16PtrFromString(windowClassName); err != nil {
		jot.Error(errs.Wrap(err))
		return
	}
	windowClass = win32.RegisterClassEx(&wcx)

	if d2dFactory = d2d.CreateFactory(false, d2d.DebugLevelNone); d2dFactory == nil {
		jot.Fatal(1, errs.New("unable to create Direct2D factory"))
	}
}

func wndProc(wnd win32.HWND, msg uint32, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
	switch msg {
	case win32.WM_COMMAND:
		if MenuItemSelectionCallback != nil {
			MenuItemSelectionCallback(int(wParam))
		}
		return 0
	case win32.WM_WINDOWPOSCHANGED:
		if wi, ok := nativeWindowMap[wnd]; ok && wi.wnd.IsValid() {
			wp := (*win32.WINDOWPOS)(unsafe.Pointer(lParam))
			if wp.Flags&win32.SWP_NOSIZE == 0 {
				wi.wnd.ValidateLayout()
			}
		}
		return 0
	case win32.WM_PAINT:
		if wi, ok := nativeWindowMap[wnd]; ok && wi.wnd.IsValid() {
			bounds := wi.wnd.ContentRect()
			bounds.X = 0
			bounds.Y = 0
			if wi.renderTarget == nil {
				if wi.renderTarget = d2dFactory.CreateHwndRenderTarget(&d2d.RenderTargetProperties{}, &d2d.HWNDRenderTargetProperties{
					HWND: wnd,
					PixelSize: d2d.SizeU{
						Width:  uint32(bounds.Width),
						Height: uint32(bounds.Height),
					},
				}); wi.renderTarget == nil {
					jot.Error(errs.New("unable to create render target"))
					return 0
				}
			} else {
				size := wi.renderTarget.Size()
				if size.Width != float32(bounds.Width) || size.Height != float32(bounds.Width) {
					wi.renderTarget.Resize(d2d.SizeU{
						Width:  uint32(bounds.Width),
						Height: uint32(bounds.Height),
					})
				}
			}
			gc := draw.NewContextForOSContext(wi.renderTarget)
			wi.wnd.Draw(gc, bounds, false)
			gc.Dispose()
		}
		win32.ValidateRect(wnd, nil)
		return 0
	case win32.WM_CLOSE:
		if wi, ok := nativeWindowMap[wnd]; ok {
			wi.wnd.AttemptClose()
		} else {
			win32.DestroyWindow(wnd)
		}
		if len(nativeWindowMap) == 0 && (QuitAfterLastWindowClosedCallback == nil || QuitAfterLastWindowClosedCallback()) {
			AttemptQuit()
		}
		return 0
	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
		return 0
	case win32.WM_ACTIVATE:
		if wi, ok := nativeWindowMap[wnd]; ok {
			if wParam&(win32.WA_ACTIVE|win32.WA_CLICKACTIVE) != 0 {
				if wi.wnd.GainedFocusCallback != nil {
					wi.wnd.GainedFocusCallback()
				}
				if child := win32.GetWindow(wnd, win32.GW_CHILD); child != win32.NULL {
					win32.SetFocus(child)
				}
				return 0
			}
			if wi.wnd.LostFocusCallback != nil {
				wi.wnd.LostFocusCallback()
			}
		}
	case win32.WM_INITMENUPOPUP:
		if MenuValidationCallback != nil {
			MenuValidationCallback(win32.HMENU(wParam))
		}
	}
	return win32.DefWindowProc(wnd, msg, wParam, lParam)
}

func styleMaskToWin32Style(styleMask StyleMask) (style, exStyle win32.DWORD) {
	return win32.WS_OVERLAPPEDWINDOW | win32.WS_CLIPCHILDREN, 0
}

func fromRectToWin32Rect(in geom.Rect) win32.RECT {
	return win32.RECT{
		Top:    int32(in.Y),
		Left:   int32(in.X),
		Bottom: int32(in.Y + in.Height),
		Right:  int32(in.X + in.Width),
	}
}

func fromWin32RectToRect(in win32.RECT) geom.Rect {
	return geom.Rect{
		Point: geom.Point{X: float64(in.Left), Y: float64(in.Top)},
		Size:  geom.Size{Width: float64(in.Right - in.Left), Height: float64(in.Bottom - in.Top)},
	}
}

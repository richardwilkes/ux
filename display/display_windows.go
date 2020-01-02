// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package display

import (
	"github.com/richardwilkes/win32"
	"unsafe"
)

func osDisplays() []*Display {
	result := make([]*Display, 0)
	win32.EnumDisplayMonitors(0, nil, func(monitor win32.HMONITOR, dc win32.HDC, rect *win32.RECT, param win32.LPARAM) win32.BOOL {
		d := &Display{}
		var info win32.MONITORINFO
		info.Size = win32.DWORD(unsafe.Sizeof(info))
		if win32.GetMonitorInfo(monitor, &info) {
			var dpiX, dpiY uint32
			if !win32.GetDpiForMonitor(monitor, win32.MDT_EFFECTIVE_DPI, &dpiX, &dpiY) {
				// Windows 7 fallback
				overallX := win32.GetDeviceCaps(dc, win32.LOGPIXELSX)
				overallY := win32.GetDeviceCaps(dc, win32.LOGPIXELSY)
				if overallX > 0 && overallY > 0 {
					dpiX = uint32(overallX)
					dpiY = uint32(overallY)
				}
			}
			d.Frame.X = float64(info.MonitorBounds.Left)
			d.Frame.Y = float64(info.MonitorBounds.Top)
			d.Frame.Width = float64(info.MonitorBounds.Right - info.MonitorBounds.Left)
			d.Frame.Height = float64(info.MonitorBounds.Bottom - info.MonitorBounds.Top)
			d.Usable.X = float64(info.WorkBounds.Left)
			d.Usable.Y = float64(info.WorkBounds.Top)
			d.Usable.Width = float64(info.WorkBounds.Right - info.WorkBounds.Left)
			d.Usable.Height = float64(info.WorkBounds.Bottom - info.WorkBounds.Top)
			d.ScalingFactor = float64(dpiX) / 96
			d.Primary = info.Flags&win32.MONITORINFOF_PRIMARY != 0
			result = append(result, d)
		}
		return 1
	}, 0)
	return result
}

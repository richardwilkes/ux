// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package draw

import (
	"github.com/richardwilkes/win32"
	"syscall"
)

// Font holds an instance of a platform font.
type osFont struct {
	ref win32.HFONT
}

func osInitSystemFonts() {
	UserFont = NewFont(FontDescriptor{
		Family: "Segoe UI",
		Size:   12,
	})
	UserMonospacedFont = NewFont(FontDescriptor{
		Family: "Consolas",
		Size:   10,
	})
	SystemFont = NewFont(FontDescriptor{
		Family: "Segoe UI",
		Size:   13,
	})
	EmphasizedSystemFont = NewFont(FontDescriptor{
		Family: "Segoe UI Black",
		Size:   13,
	})
	SmallSystemFont = NewFont(FontDescriptor{
		Family: "Segoe UI",
		Size:   11,
	})
	SmallEmphasizedSystemFont = NewFont(FontDescriptor{
		Family: "Segoe UI Black",
		Size:   11,
	})
	ViewsFont = NewFont(FontDescriptor{
		Family: "Segoe UI",
		Size:   12,
	})
	LabelFont = NewFont(FontDescriptor{
		Family: "Segoe UI",
		Size:   10,
	})
	MenuFont = NewFont(FontDescriptor{
		Family: "Segoe UI",
		Size:   14,
	})
	MenuCmdKeyFont = NewFont(FontDescriptor{
		Family: "Segoe UI",
		Size:   14,
	})
}

func osFontFamilies() []string {
	hdc := win32.GetDC(0)
	defer win32.ReleaseDC(0, hdc)
	all := make(map[string]bool)
	win32.EnumFontFamiliesExW(hdc, &win32.LOGFONT{LfCharSet: win32.DEFAULT_CHARSET}, func(lf *win32.LOGFONT, tm *win32.NEWTEXTMETRIC, fontType win32.DWORD, parm win32.LPARAM) int {
		if fontType == win32.TRUETYPE_FONTTYPE && lf.LfFaceName[0] != '@' {
			all[syscall.UTF16ToString(lf.LfFaceName[:])] = true
		}
		return 1
	}, 0, 0)
	families := make([]string, 0, len(all))
	for k := range all {
		families = append(families, k)
	}
	return families
}

func osNewFont(desc FontDescriptor) *Font {
	hdc := win32.GetDC(0)
	defer win32.ReleaseDC(0, hdc)
	weight := 400
	if desc.Bold {
		weight = 700
	}
	italic := 0
	if desc.Italic {
		italic = 1
	}
	hFont := win32.CreateFont(int(desc.Size*float64(win32.GetDeviceCaps(hdc, win32.LOGPIXELSY))/72),
		0, 0, 0, weight, italic, 0, 0, win32.ANSI_CHARSET, win32.OUT_TT_PRECIS, win32.CLIP_DEFAULT_PRECIS,
		win32.PROOF_QUALITY, win32.DEFAULT_PITCH|win32.FF_DONTCARE, desc.Family)
	win32.SelectObject(hdc, win32.HGDIOBJ(hFont))
	var tm win32.NEWTEXTMETRIC
	win32.GetTextMetrics(hdc, &tm)
	return &Font{
		ascent:     float64(tm.TmAscent),
		descent:    float64(tm.TmDescent),
		leading:    float64(tm.TmExternalLeading),
		desc:       desc,
		osFont:     osFont{ref: hFont},
		monospaced: (tm.TmPitchAndFamily & win32.TMPF_FIXED_PITCH) == 0,
	}
}

func (f *Font) osWidth(str string) float64 {
	var size win32.SIZE
	hdc := win32.GetDC(0)
	win32.SelectObject(hdc, win32.HGDIOBJ(f.ref))
	win32.GetTextExtentPoint(hdc, str, &size)
	win32.ReleaseDC(0, hdc)
	return float64(size.CX)
}

func (f *Font) osIndexForPosition(x float64, str string) int {
	return 0 // RAW: Implement
}

func (f *Font) osPositionForIndex(index int, str string) float64 {
	return 0 // RAW: Implement
}

func (f *Font) osDispose() {
	if f.ref != 0 {
		win32.DeleteObject(win32.HGDIOBJ(f.ref))
		f.ref = 0
	}
}

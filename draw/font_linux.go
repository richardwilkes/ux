// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package draw

type osFont struct {
}

func osInitSystemFonts() {
	UserFont = NewFont(FontDescriptor{
		Family: "Sans",
		Size:   12,
	})
	UserMonospacedFont = NewFont(FontDescriptor{
		Family: "Monospace",
		Size:   10,
	})
	SystemFont = NewFont(FontDescriptor{
		Family: "Sans",
		Size:   13,
	})
	EmphasizedSystemFont = NewFont(FontDescriptor{
		Family: "Sans",
		Size:   13,
		Bold:   true,
	})
	SmallSystemFont = NewFont(FontDescriptor{
		Family: "Sans",
		Size:   11,
	})
	SmallEmphasizedSystemFont = NewFont(FontDescriptor{
		Family: "Sans",
		Size:   11,
		Bold:   true,
	})
	ViewsFont = NewFont(FontDescriptor{
		Family: "Sans",
		Size:   12,
	})
	LabelFont = NewFont(FontDescriptor{
		Family: "Sans",
		Size:   10,
	})
	MenuFont = NewFont(FontDescriptor{
		Family: "Sans",
		Size:   14,
	})
	MenuCmdKeyFont = NewFont(FontDescriptor{
		Family: "Sans",
		Size:   14,
	})
}

func osFontFamilies() []string {
	return []string{"Sans", "Monospaced"} // RAW: Implement
}

func osNewFont(desc FontDescriptor) *Font {
	// RAW: Implement
	return &Font{}
}

func (f *Font) osWidth(str string) float64 {
	return 0 // RAW: Implement
}

func (f *Font) osIndexForPosition(x float64, str string) int {
	return 0 // RAW: Implement
}

func (f *Font) osPositionForIndex(index int, str string) float64 {
	return 0 // RAW: Implement
}

func (f *Font) osDispose() {
	// RAW: Implement
}

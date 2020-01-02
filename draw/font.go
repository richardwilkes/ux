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
	"sort"

	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xmath/geom"
)

var (
	// UserFont is the font used by default for documents and other text under
	// the user’s control (that is, text whose font the user can normally
	// change).
	UserFont *Font
	// UserMonospacedFont is the font used by default for documents and other
	// text under the user’s control when that font is fixed-pitch.
	UserMonospacedFont *Font
	// SystemFont is the system font used for standard user-interface items
	// such as window titles, button labels, etc.
	SystemFont *Font
	// EmphasizedSystemFont is the system font used for emphasis in alerts.
	EmphasizedSystemFont *Font
	// SmallSystemFont is the standard small system font used for informative
	// text in alerts, column headings in lists, help tags, utility window
	// titles, toolbar item labels, tool palettes, tool tips, and small
	// controls.
	SmallSystemFont *Font
	// SmallEmphasizedSystemFont is the small system font used for emphasis.
	SmallEmphasizedSystemFont *Font
	// ViewsFont is the font used as the default font of text in lists and
	// tables.
	ViewsFont *Font
	// LabelFont is the font used for labels.
	LabelFont *Font
	// MenuFont is the font used for menus.
	MenuFont *Font
	// MenuCmdKeyFont is the font used for menu item command key equivalents.
	MenuCmdKeyFont *Font
)

var fontCache = make(map[FontDescriptor]*Font)

// Font holds an instance of a platform font.
type Font struct {
	ascent   float64
	descent  float64
	leading  float64
	desc     FontDescriptor
	refCount uint
	osFont
	monospaced bool
}

// FontFamilies retrieves the names of the installed font families.
func FontFamilies() []string {
	families := osFontFamilies()
	sort.Slice(families, func(i, j int) bool { return txt.NaturalLess(families[i], families[j], true) })
	return families
}

// NewFont creates a new font from the supplied FontDescriptor.
func NewFont(desc FontDescriptor) *Font {
	f, exists := fontCache[desc]
	if !exists {
		f = osNewFont(desc)
		fontCache[desc] = f
	}
	f.refCount++
	return f
}

// Descriptor returns the FontDescriptor for this font.
func (f *Font) Descriptor() FontDescriptor {
	return f.desc
}

// Ascent returns the distance from the baseline to the top of a typical
// capital letter.
func (f *Font) Ascent() float64 {
	return f.ascent
}

// Descent returns the distance from the baseline to the bottom of the typical
// letter that has a descender, such as a lower case 'g'.
func (f *Font) Descent() float64 {
	return f.descent
}

// Leading returns the recommended distance between the bottom of the
// descender line to the top of the next line.
func (f *Font) Leading() float64 {
	return f.leading
}

// Height returns the height of the font, typically Ascent() + Descent().
func (f *Font) Height() float64 {
	return f.ascent + f.descent
}

// Monospaced returns true if the font characters have a fixed width.
func (f *Font) Monospaced() bool {
	return f.monospaced
}

// Width of the string rendered with this font. Note that this does not
// account for any embedded line endings nor tabs.
func (f *Font) Width(str string) float64 {
	return f.osWidth(str)
}

// Extents of the string rendered with this font. Note that this does not
// account for any embedded line endings nor tabs.
func (f *Font) Extents(str string) geom.Size {
	return geom.Size{Width: f.Width(str), Height: f.Height()}
}

// IndexForPosition returns the rune index within the string for the specified
// x-coordinate, where 0 is the start of the string. Note that this does not
// account for any embedded line endings nor tabs.
func (f *Font) IndexForPosition(x float64, str string) int {
	return f.osIndexForPosition(x, str)
}

// PositionForIndex returns the x-coordinate where the specified rune index
// starts. The returned coordinate assumes 0 is the start of the string. Note
// that this does not account for any embedded line endings nor tabs.
func (f *Font) PositionForIndex(index int, str string) float64 {
	return f.osPositionForIndex(index, str)
}

// String implements fmt.Stringer.
func (f *Font) String() string {
	return f.desc.String()
}

// Dispose of the font, releasing any OS resources associated with it.
func (f *Font) Dispose() {
	switch f.refCount {
	case 0:
	case 1:
		f.refCount = 0
		delete(fontCache, f.desc)
		f.osDispose()
	default:
		f.refCount--
	}
}

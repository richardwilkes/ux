// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package layout

import (
	"math"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/border"
)

const (
	// DefaultMaxSize is the default size that should be used for a maximum
	// dimension if the target has no real preference and can be expanded
	// beyond its preferred size. This is intentionally not something very
	// large to allow basic math operations an opportunity to succeed when
	// laying out panels. It is perfectly acceptable to use a larger value
	// than this, however, if that makes sense for your specific target.
	DefaultMaxSize = 10000
	// DefaultHSpacing is used for the default spacing between columns.
	DefaultHSpacing = 4
	// DefaultVSpacing is used for the default spacing between rows.
	DefaultVSpacing = 2
)

// Sizer returns minimum, preferred, and maximum sizes. The hint will contain
// values other than zero for a dimension that has already been determined.
type Sizer func(hint geom.Size) (min, pref, max geom.Size)

// Layout defines methods that all layouts must provide.
type Layout interface {
	Sizes(hint geom.Size) (min, pref, max geom.Size)
	Layout()
}

// Layoutable defines the methods an object that wants to participate in
// layout must implement.
type Layoutable interface {
	SetLayout(layout Layout)
	LayoutData() interface{}
	SetLayoutData(data interface{})
	Sizes(hint geom.Size) (min, pref, max geom.Size)
	Border() border.Border
	FrameRect() geom.Rect
	SetFrameRect(rect geom.Rect)
	ChildrenForLayout() []Layoutable
}

// MaxSize returns the size that is at least as large as DefaultMaxSize in
// both dimensions, but larger if the size that is passed in is larger.
func MaxSize(size geom.Size) geom.Size {
	return geom.Size{
		Width:  math.Max(DefaultMaxSize, size.Width),
		Height: math.Max(DefaultMaxSize, size.Height),
	}
}

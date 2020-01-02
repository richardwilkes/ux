// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package flex

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/layout"
	"github.com/richardwilkes/ux/layout/align"
)

// Data is used to control how an object is laid out by the Flex layout.
type Data struct {
	cacheSize    geom.Size
	minCacheSize geom.Size
	sizeHint     geom.Size
	minSize      geom.Size
	hSpan        int
	vSpan        int
	hAlign       align.Alignment
	vAlign       align.Alignment
	hGrab        bool
	vGrab        bool
}

// NewData creates new flex layout data.
func NewData() *Data {
	return &Data{
		hSpan:  1,
		vSpan:  1,
		hAlign: align.Start,
		vAlign: align.Middle,
	}
}

// SizeHint sets a hint requesting a particular size for the target. Defaults
// to no hint (0, 0).
func (f *Data) SizeHint(sizeHint geom.Size) *Data {
	f.sizeHint = sizeHint
	return f
}

// MinSize overrides the minimum size of the target. Defaults to no override
// (0, 0).
func (f *Data) MinSize(minSize geom.Size) *Data {
	f.minSize = minSize
	return f
}

// HSpan sets the number of columns the target should span. Defaults to 1.
func (f *Data) HSpan(hSpan int) *Data {
	f.hSpan = hSpan
	return f
}

// VSpan sets the number of rows the target should span. Defaults to 1.
func (f *Data) VSpan(vSpan int) *Data {
	f.vSpan = vSpan
	return f
}

// HAlign sets the horizontal alignment of the target within its available
// space. Defaults to Start.
func (f *Data) HAlign(hAlign align.Alignment) *Data {
	f.hAlign = hAlign
	return f
}

// VAlign sets the vertical alignment of the target within its available
// space. Defaults to Middle.
func (f *Data) VAlign(vAlign align.Alignment) *Data {
	f.vAlign = vAlign
	return f
}

// HGrab sets whether excess horizontal space should be grabbed. Defaults to
// false.
func (f *Data) HGrab(hGrab bool) *Data {
	f.hGrab = hGrab
	return f
}

// VGrab sets whether excess vertical space should be grabbed. Defaults to
// false.
func (f *Data) VGrab(vGrab bool) *Data {
	f.vGrab = vGrab
	return f
}

// Apply the layout data to the target. A copy is made of this data and that
// is applied to the target, so this data may be applied to other targets.
func (f *Data) Apply(target layout.Layoutable) {
	flexData := *f
	flexData.normalizeAndResetCache()
	target.SetLayoutData(&flexData)
}

func (f *Data) normalizeAndResetCache() {
	f.cacheSize.Width = 0
	f.cacheSize.Height = 0
	f.minCacheSize.Width = 0
	f.minCacheSize.Height = 0
	if f.sizeHint.Width < 0 {
		f.sizeHint.Width = 0
	}
	if f.sizeHint.Height < 0 {
		f.sizeHint.Height = 0
	}
	if f.minSize.Width < 0 {
		f.minSize.Width = 0
	}
	if f.minSize.Height < 0 {
		f.minSize.Height = 0
	}
	if f.hSpan < 1 {
		f.hSpan = 1
	}
	if f.vSpan < 1 {
		f.vSpan = 1
	}
}

func (f *Data) computeCacheSize(sizer layout.Sizer, hint geom.Size, useMinimumSize bool) {
	f.normalizeAndResetCache()
	min, pref, max := sizer(hint)
	if hint.Width > 0 || hint.Height > 0 {
		if f.minSize.Width > 0 {
			f.minCacheSize.Width = f.minSize.Width
		} else {
			f.minCacheSize.Width = min.Width
		}
		if hint.Width > 0 && hint.Width < f.minCacheSize.Width {
			hint.Width = f.minCacheSize.Width
		}
		if hint.Width > 0 && hint.Width > max.Width {
			hint.Width = max.Width
		}
		if f.minSize.Height > 0 {
			f.minCacheSize.Height = f.minSize.Height
		} else {
			f.minCacheSize.Height = min.Height
		}
		if hint.Height > 0 && hint.Height < f.minCacheSize.Height {
			hint.Height = f.minCacheSize.Height
		}
		if hint.Height > 0 && hint.Height > max.Height {
			hint.Height = max.Height
		}
	}
	if useMinimumSize {
		f.cacheSize = min
		if f.minSize.Width > 0 {
			f.minCacheSize.Width = f.minSize.Width
		} else {
			f.minCacheSize.Width = min.Width
		}
		if f.minSize.Height > 0 {
			f.minCacheSize.Height = f.minSize.Height
		} else {
			f.minCacheSize.Height = min.Height
		}
	} else {
		f.cacheSize = pref
	}
	if hint.Width > 0 {
		f.cacheSize.Width = hint.Width
	}
	if f.minSize.Width > 0 && f.cacheSize.Width < f.minSize.Width {
		f.cacheSize.Width = f.minSize.Width
	}
	if f.sizeHint.Width > 0 {
		f.cacheSize.Width = f.sizeHint.Width
	}
	if hint.Height > 0 {
		f.cacheSize.Height = hint.Height
	}
	if f.minSize.Height > 0 && f.cacheSize.Height < f.minSize.Height {
		f.cacheSize.Height = f.minSize.Height
	}
	if f.sizeHint.Height > 0 {
		f.cacheSize.Height = f.sizeHint.Height
	}
}

package layout

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/layout/align"
)

// FlexData is used to control how an object is laid out by the Flex layout.
type FlexData struct {
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

// NewFlexData creates new flex layout data.
func NewFlexData() *FlexData {
	return &FlexData{
		hSpan:  1,
		vSpan:  1,
		hAlign: align.Start,
		vAlign: align.Middle,
	}
}

// SizeHint sets a hint requesting a particular size for the target. Defaults
// to no hint (0, 0).
func (f *FlexData) SizeHint(sizeHint geom.Size) *FlexData {
	f.sizeHint = sizeHint
	return f
}

// MinSize overrides the minimum size of the target. Defaults to no override
// (0, 0).
func (f *FlexData) MinSize(minSize geom.Size) *FlexData {
	f.minSize = minSize
	return f
}

// HSpan sets the number of columns the target should span. Defaults to 1.
func (f *FlexData) HSpan(hSpan int) *FlexData {
	f.hSpan = hSpan
	return f
}

// VSpan sets the number of rows the target should span. Defaults to 1.
func (f *FlexData) VSpan(vSpan int) *FlexData {
	f.vSpan = vSpan
	return f
}

// HAlign sets the horizontal alignment of the target within its available
// space. Defaults to Start.
func (f *FlexData) HAlign(hAlign align.Alignment) *FlexData {
	f.hAlign = hAlign
	return f
}

// VAlign sets the vertical alignment of the target within its available
// space. Defaults to Middle.
func (f *FlexData) VAlign(vAlign align.Alignment) *FlexData {
	f.vAlign = vAlign
	return f
}

// HGrab sets whether excess horizontal space should be grabbed. Defaults to
// false.
func (f *FlexData) HGrab(hGrab bool) *FlexData {
	f.hGrab = hGrab
	return f
}

// VGrab sets whether excess vertical space should be grabbed. Defaults to
// false.
func (f *FlexData) VGrab(vGrab bool) *FlexData {
	f.vGrab = vGrab
	return f
}

// Apply the layout data to the target. A copy is made of this data and that
// is applied to the target, so this data may be applied to other targets.
func (f *FlexData) Apply(target Layoutable) {
	flexData := *f
	flexData.normalizeAndResetCache()
	target.SetLayoutData(&flexData)
}

func (f *FlexData) normalizeAndResetCache() {
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

func (f *FlexData) computeCacheSize(sizer Sizer, hint geom.Size, useMinimumSize bool) {
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

package layout

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/layout/align"
)

// FlexData is used to control how an object is laid out by the Flex layout.
type FlexData struct {
	cacheSize    geom.Size
	minCacheSize geom.Size
	SizeHint     geom.Size       // Hint requesting a particular size of the Layoutable
	MinSize      geom.Size       // Override for the minimum size of the Layoutable
	HSpan        int             // Number of columns the Layoutable should span
	VSpan        int             // Number of rows the Layoutable should span
	HAlign       align.Alignment // Horizontal alignment of the Layoutable within its space
	VAlign       align.Alignment // Vertical alignment of the Layoutable within its space
	HGrab        bool            // Grab excess horizontal space if true
	VGrab        bool            // Grab excess vertical space if true
}

// NewFlexData creates new flex layout data.
func NewFlexData() *FlexData {
	return &FlexData{
		HSpan:  1,
		VSpan:  1,
		HAlign: align.Start,
		VAlign: align.Middle,
	}
}

// Clone the Data.
func (data *FlexData) Clone() *FlexData {
	d := *data
	d.normalizeAndResetCache()
	return &d
}

func (data *FlexData) normalizeAndResetCache() {
	data.cacheSize.Width = 0
	data.cacheSize.Height = 0
	data.minCacheSize.Width = 0
	data.minCacheSize.Height = 0
	if data.SizeHint.Width < 0 {
		data.SizeHint.Width = 0
	}
	if data.SizeHint.Height < 0 {
		data.SizeHint.Height = 0
	}
	if data.MinSize.Width < 0 {
		data.MinSize.Width = 0
	}
	if data.MinSize.Height < 0 {
		data.MinSize.Height = 0
	}
	if data.HSpan < 1 {
		data.HSpan = 1
	}
	if data.VSpan < 1 {
		data.VSpan = 1
	}
}

func (data *FlexData) computeCacheSize(sizer Sizer, hint geom.Size, useMinimumSize bool) {
	data.normalizeAndResetCache()
	min, pref, max := sizer(hint)
	if hint.Width > 0 || hint.Height > 0 {
		if data.MinSize.Width > 0 {
			data.minCacheSize.Width = data.MinSize.Width
		} else {
			data.minCacheSize.Width = min.Width
		}
		if hint.Width > 0 && hint.Width < data.minCacheSize.Width {
			hint.Width = data.minCacheSize.Width
		}
		if hint.Width > 0 && hint.Width > max.Width {
			hint.Width = max.Width
		}
		if data.MinSize.Height > 0 {
			data.minCacheSize.Height = data.MinSize.Height
		} else {
			data.minCacheSize.Height = min.Height
		}
		if hint.Height > 0 && hint.Height < data.minCacheSize.Height {
			hint.Height = data.minCacheSize.Height
		}
		if hint.Height > 0 && hint.Height > max.Height {
			hint.Height = max.Height
		}
	}
	if useMinimumSize {
		data.cacheSize = min
		if data.MinSize.Width > 0 {
			data.minCacheSize.Width = data.MinSize.Width
		} else {
			data.minCacheSize.Width = min.Width
		}
		if data.MinSize.Height > 0 {
			data.minCacheSize.Height = data.MinSize.Height
		} else {
			data.minCacheSize.Height = min.Height
		}
	} else {
		data.cacheSize = pref
	}
	if hint.Width > 0 {
		data.cacheSize.Width = hint.Width
	}
	if data.MinSize.Width > 0 && data.cacheSize.Width < data.MinSize.Width {
		data.cacheSize.Width = data.MinSize.Width
	}
	if data.SizeHint.Width > 0 {
		data.cacheSize.Width = data.SizeHint.Width
	}
	if hint.Height > 0 {
		data.cacheSize.Height = hint.Height
	}
	if data.MinSize.Height > 0 && data.cacheSize.Height < data.MinSize.Height {
		data.cacheSize.Height = data.MinSize.Height
	}
	if data.SizeHint.Height > 0 {
		data.cacheSize.Height = data.SizeHint.Height
	}
}

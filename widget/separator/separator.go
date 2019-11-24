package separator

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/layout"
)

// Separator provides a simple vertical or horizontal separator line.
type Separator struct {
	ux.Panel
	FillInk    draw.Ink
	horizontal bool
}

// NewHorizontal creates a new horizontal separator line.
func NewHorizontal() *Separator {
	return New(true)
}

// NewVertical creates a new vertical separator line.
func NewVertical() *Separator {
	return New(false)
}

// New creates a new horizontal or vertical separator line.
func New(horizontal bool) *Separator {
	s := &Separator{
		FillInk:    draw.SeparatorColor,
		horizontal: horizontal,
	}
	s.InitTypeAndID(s)
	s.SetSizer(s.DefaultSizes)
	s.DrawCallback = s.DefaultDraw
	return s
}

// DefaultSizes provides the default sizing.
func (s *Separator) DefaultSizes(hint geom.Size) (min, pref, max geom.Size) {
	if s.horizontal {
		if hint.Width < 1 {
			pref.Width = 1
		} else {
			pref.Width = hint.Width
		}
		min.Width = 1
		max.Width = layout.DefaultMaxSize
		min.Height = 1
		pref.Height = 1
		max.Height = 1
	} else {
		if hint.Height < 1 {
			pref.Height = 1
		} else {
			pref.Height = hint.Height
		}
		min.Height = 1
		max.Height = layout.DefaultMaxSize
		min.Width = 1
		pref.Width = 1
		max.Width = 1
	}
	if border := s.Border(); border != nil {
		insets := border.Insets()
		min.AddInsets(insets)
		pref.AddInsets(insets)
		max.AddInsets(insets)
	}
	return min, pref, max
}

// DefaultDraw provides the default drawing.
func (s *Separator) DefaultDraw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	rect := s.ContentRect(false)
	if s.horizontal {
		if rect.Height > 1 {
			rect.Y += (rect.Height - 1) / 2
			rect.Height = 1
		}
	} else if rect.Width > 1 {
		rect.X += (rect.Width - 1) / 2
		rect.Width = 1
	}
	gc.Rect(rect)
	gc.Fill(s.FillInk)
}

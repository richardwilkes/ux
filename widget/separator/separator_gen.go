// Code created from "widget.go.tmpl" - don't edit by hand

package separator

import (
	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/draw"
)

type managed struct {
	fillInk draw.Ink
}

func (m *managed) initialize() {
	m.fillInk = draw.SeparatorColor
}

// FillInk returns the ink that will be used for the separator line.
func (s *Separator) FillInk() draw.Ink {
	return s.fillInk
}

// SetFillInk sets the ink that will be used for the separator line. Pass in
// nil to use the default.
func (s *Separator) SetFillInk(value draw.Ink) *Separator {
	if value == nil {
		value = draw.SeparatorColor
	}
	if s.fillInk != value {
		s.fillInk = value
		s.MarkForRedraw()
	}
	return s
}

// SetBorder sets the border. May be nil.
func (s *Separator) SetBorder(value border.Border) *Separator {
	s.Panel.SetBorder(value)
	return s
}

// SetEnabled sets enabled state.
func (s *Separator) SetEnabled(enabled bool) *Separator {
	s.Panel.SetEnabled(enabled)
	return s
}

// SetFocusable whether it can have the keyboard focus.
func (s *Separator) SetFocusable(focusable bool) *Separator {
	s.Panel.SetFocusable(focusable)
	return s
}

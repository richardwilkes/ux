package draw

import (
	"fmt"

	"github.com/richardwilkes/toolbox/xmath/geom"
)

var _ Ink = &Gradient{}

// Stop provides information about the color and position of one 'color stop'
// in a gradient.
type Stop struct {
	Color    *DynamicColor
	Location float64
}

func (s Stop) String() string {
	return fmt.Sprintf("%v:%v", s.Color.Color, s.Location)
}

// Gradient defines a smooth transition between colors across an area. Start
// and End should hold values from 0 to 1. These will be be used to set a
// relative starting and ending position for the gradient. If StartRadius and
// EndRadius are both greater than 0, then the gradient will be a radial one
// instead of a linear one.
type Gradient struct {
	Start       geom.Point
	StartRadius float64
	End         geom.Point
	EndRadius   float64
	Stops       []Stop
}

// NewHorizontalEvenlySpacedGradient creates a new gradient with the specified
// colors evenly spread across the whole range.
func NewHorizontalEvenlySpacedGradient(colors ...*DynamicColor) *Gradient {
	return NewEvenlySpacedGradient(geom.Point{}, geom.Point{X: 1}, 0, 0, colors...)
}

// NewVerticalEvenlySpacedGradient creates a new gradient with the specified
// colors evenly spread across the whole range.
func NewVerticalEvenlySpacedGradient(colors ...*DynamicColor) *Gradient {
	return NewEvenlySpacedGradient(geom.Point{}, geom.Point{Y: 1}, 0, 0, colors...)
}

// NewEvenlySpacedGradient creates a new gradient with the specified colors
// evenly spread across the whole range. start and end should hold values from
// 0 to 1, representing the percentage position within the area that will be
// filled.
func NewEvenlySpacedGradient(start, end geom.Point, startRadius, endRadius float64, colors ...*DynamicColor) *Gradient {
	gradient := &Gradient{
		Start:       start,
		StartRadius: startRadius,
		End:         end,
		EndRadius:   endRadius,
		Stops:       make([]Stop, len(colors)),
	}
	switch len(colors) {
	case 0:
	case 1:
		gradient.Stops[0].Color = colors[0]
	case 2:
		gradient.Stops[0].Color = colors[0]
		gradient.Stops[1].Color = colors[1]
		gradient.Stops[1].Location = 1
	default:
		step := 1 / float64(len(colors)-1)
		var location float64
		for i, color := range colors {
			gradient.Stops[i].Color = color
			gradient.Stops[i].Location = location
			if i < len(colors)-1 {
				location += step
			} else {
				location = 1
			}
		}
	}
	return gradient
}

// Reversed creates a copy of the current Gradient and inverts the locations
// of each color stop in that copy.
func (g *Gradient) Reversed() *Gradient {
	other := *g
	other.Stops = make([]Stop, len(g.Stops))
	for i, stop := range g.Stops {
		stop.Location = 1 - stop.Location
		other.Stops[i] = stop
	}
	return &other
}

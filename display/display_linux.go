package display

import (
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/ux/globals"
)

func osDisplays() []*Display {
	screen := globals.X11.Screen()
	display := &Display{
		ScalingFactor: (float64(screen.WidthInPixels) * 25.4 / float64(screen.WidthInMillimeters)) / 96,
		Primary:       true,
	}
	display.Frame.Width = 1024
	display.Frame.Height = 768
	display.Usable = display.Frame
	result := []*Display{display}
	size, err := ewmh.DesktopGeometryGet(globals.X11)
	if err != nil {
		jot.Error(errs.Wrap(err))
		return result
	}
	display.Frame.Width = float64(size.Width)
	display.Frame.Height = float64(size.Height)
	display.Usable.Size = display.Frame.Size
	cur, err := ewmh.CurrentDesktopGet(globals.X11)
	if err != nil {
		jot.Error(errs.Wrap(err))
		return result
	}
	pos, err := ewmh.DesktopViewportGet(globals.X11)
	if err != nil {
		jot.Error(errs.Wrap(err))
		return result
	}
	display.Frame.X = float64(pos[cur].X)
	display.Frame.Y = float64(pos[cur].Y)
	display.Usable.Point = display.Frame.Point
	workAreas, err := ewmh.WorkareaGet(globals.X11)
	if err != nil {
		jot.Error(errs.Wrap(err))
		return result
	}
	display.Usable.X = float64(workAreas[cur].X)
	display.Usable.Y = float64(workAreas[cur].Y)
	display.Usable.Width = float64(workAreas[cur].Width)
	display.Usable.Height = float64(workAreas[cur].Height)
	return result
}

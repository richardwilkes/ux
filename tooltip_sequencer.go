package ux

import (
	"time"

	"github.com/richardwilkes/toolbox/xmath/geom"
)

var (
	// TooltipDelay holds the delay before a tooltip will be shown.
	TooltipDelay = 1500 * time.Millisecond
	// TooltipDismissal holds the delay before a tooltip will be dismissed.
	TooltipDismissal = 3 * time.Second
)

type tooltipSequencer struct {
	window   *Window
	avoid    geom.Rect
	sequence int
}

func (ts *tooltipSequencer) show() {
	if ts.window.tooltipSequence == ts.sequence {
		tip := ts.window.lastTooltip
		_, pref, _ := tip.Sizes(geom.Size{})
		rect := geom.Rect{Point: geom.Point{X: ts.avoid.X, Y: ts.avoid.Y + ts.avoid.Height + 1}, Size: pref}
		if rect.X < 0 {
			rect.X = 0
		}
		if rect.Y < 0 {
			rect.Y = 0
		}
		viewSize := ts.window.root.ContentRect(true).Size
		if viewSize.Width < rect.Width {
			_, pref, _ = tip.Sizes(geom.Size{Width: viewSize.Width})
			if viewSize.Width < pref.Width {
				rect.X = 0
				rect.Width = viewSize.Width
			} else {
				rect.Width = pref.Width
			}
			rect.Height = pref.Height
		}
		if viewSize.Width < rect.X+rect.Width {
			rect.X = viewSize.Width - rect.Width
		}
		if viewSize.Height < rect.Y+rect.Height {
			rect.Y = ts.avoid.Y - (rect.Height + 1)
			if rect.Y < 0 {
				rect.Y = 0
			}
		}
		tip.SetFrameRect(rect)
		ts.window.root.setTooltip(tip)
		ts.window.lastTooltipShownAt = time.Now()
		InvokeAfter(ts.close, TooltipDismissal)
	}
}

func (ts *tooltipSequencer) close() {
	if ts.window.tooltipSequence == ts.sequence {
		ts.window.root.setTooltip(nil)
	}
}

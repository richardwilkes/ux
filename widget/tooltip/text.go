package tooltip

import (
	"strings"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/layout"
	"github.com/richardwilkes/ux/widget/label"
)

// NewBase returns the base for a tooltip.
func NewBase() *ux.Panel {
	tip := ux.NewPanel()
	tip.SetBorder(border.NewCompound(border.NewLine(draw.UnemphasizedSelectedContentBackgroundColor, geom.NewUniformInsets(1), false), border.NewEmpty(geom.Insets{Top: 2, Left: 4, Bottom: 2, Right: 4})))
	tip.DrawCallback = func(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
		gc.Rect(dirty)
		gc.Fill(draw.ControlBackgroundColor)
	}
	return tip
}

// NewWithText creates a standard text tooltip panel.
func NewWithText(text string) *ux.Panel {
	tip := NewBase()
	layout.NewFlex(tip)
	for _, str := range strings.Split(text, "\n") {
		l := label.NewWithText(str)
		l.Font = draw.SystemFont
		tip.AddChild(l.AsPanel())
	}
	return tip
}

// NewWithSecondaryText creates a text tooltip panel containing a primary
// piece of text along with a secondary piece of text in a slightly smaller
// font.
func NewWithSecondaryText(primary, secondary string) *ux.Panel {
	tip := NewBase()
	layout.NewFlex(tip)
	for _, str := range strings.Split(primary, "\n") {
		l := label.NewWithText(str)
		l.Font = draw.SystemFont
		tip.AddChild(l.AsPanel())
	}
	if secondary != "" {
		for _, str := range strings.Split(secondary, "\n") {
			l := label.NewWithText(str)
			l.Font = draw.SmallSystemFont
			tip.AddChild(l.AsPanel())
		}
	}
	return tip
}

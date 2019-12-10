package label

import (
	"fmt"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/draw"
)

// CellFactory provides a simple implementation of a CellFactory that uses
// Labels for its cells.
type CellFactory struct {
	Height float64
}

// CellHeight implements widget.CellFactory.
func (f *CellFactory) CellHeight() float64 {
	return f.Height
}

// CreateCell implements widget.CellFactory.
func (f *CellFactory) CreateCell(owner *ux.Panel, element interface{}, index int, selected, focused bool) *ux.Panel {
	txtLabel := New().SetText(fmt.Sprintf("%v", element)).SetFont(draw.ViewsFont)
	txtLabel.SetBorder(border.NewEmpty(geom.Insets{Left: 4, Right: 4}))
	if selected {
		txtLabel.SetInk(draw.AlternateSelectedControlTextColor)
	}
	return txtLabel.AsPanel()
}

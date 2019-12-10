package widget

import (
	"fmt"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/widget/label"
)

// CellFactory defines methods all cell factories must implement.
type CellFactory interface {
	// CellHeight returns the height to use for the cells. A value less than 1
	// indicates that each cell's height may be different.
	CellHeight() float64

	// CreateCell creates a new cell for 'owner' using 'element' as the
	// content. 'index' indicates which row the element came from. 'selected'
	// indicates the cell should be created in its selected state. 'focused'
	// indicates the cell should be created in its focused state.
	CreateCell(owner *ux.Panel, element interface{}, index int, selected, focused bool) *ux.Panel
}

// LabelCellFactory provides a simple implementation of a CellFactory that
// uses Labels for its cells.
type LabelCellFactory struct {
	Height float64
}

// CellHeight implements the CellFactory interface.
func (f *LabelCellFactory) CellHeight() float64 {
	return f.Height
}

// CreateCell implements the CellFactory interface.
func (f *LabelCellFactory) CreateCell(owner *ux.Panel, element interface{}, index int, selected, focused bool) *ux.Panel {
	txtLabel := label.New().SetText(fmt.Sprintf("%v", element)).SetFont(draw.ViewsFont)
	txtLabel.SetBorder(border.NewEmpty(geom.Insets{Left: 4, Right: 4}))
	if selected {
		txtLabel.SetInk(draw.AlternateSelectedControlTextColor)
	}
	return txtLabel.AsPanel()
}

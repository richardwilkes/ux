// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package list

import (
	"math"

	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/ids"
	"github.com/richardwilkes/ux/keys"
	"github.com/richardwilkes/ux/layout"
)

// List provides a control that allows the user to select from a list of
// items, represented by cells.
type List struct {
	ux.Panel
	managed
	DoubleClickCallback  func()
	NewSelectionCallback func()
	rows                 []interface{}
	Selection            *xmath.BitSet
	savedSelection       *xmath.BitSet
	anchor               int
	pressed              bool
}

// New creates a new List control.
func New() *List {
	l := &List{
		Selection:      &xmath.BitSet{},
		savedSelection: &xmath.BitSet{},
		anchor:         -1,
	}
	l.managed.initialize()
	l.InitTypeAndID(l)
	l.SetFocusable(true)
	l.SetSizer(l.DefaultSizes)
	l.DrawCallback = l.DefaultDraw
	l.MouseDownCallback = l.DefaultMouseDown
	l.MouseDragCallback = l.DefaultMouseDrag
	l.MouseUpCallback = l.DefaultMouseUp
	l.KeyDownCallback = l.DefaultKeyDown
	l.CanPerformCmdCallback = l.DefaultCanPerformCmd
	l.PerformCmdCallback = l.DefaultPerformCmd
	return l
}

// Append values to the list of items.
func (l *List) Append(values ...interface{}) {
	l.rows = append(l.rows, values...)
	l.MarkForLayoutAndRedraw()
}

// Insert values at the specified index.
func (l *List) Insert(index int, values ...interface{}) {
	l.rows = append(l.rows[:index], append(values, l.rows[index:]...)...)
	l.MarkForLayoutAndRedraw()
}

// Remove the item at the specified index.
func (l *List) Remove(index int) {
	copy(l.rows[index:], l.rows[index+1:])
	size := len(l.rows) - 1
	l.rows[size] = nil
	l.rows = l.rows[:size]
	l.MarkForLayoutAndRedraw()
}

// DefaultSizes provides the default sizing.
func (l *List) DefaultSizes(hint geom.Size) (min, pref, max geom.Size) {
	max = layout.MaxSize(max)
	height := math.Ceil(l.factory.CellHeight())
	if height < 1 {
		height = 0
	}
	size := geom.Size{Width: hint.Width, Height: height}
	for i, row := range l.rows {
		cell := l.factory.CreateCell(l.AsPanel(), row, i, false, false)
		_, cpref, cmax := cell.Sizes(size)
		cpref.GrowToInteger()
		cmax.GrowToInteger()
		if pref.Width < cpref.Width {
			pref.Width = cpref.Width
		}
		if max.Width < cmax.Width {
			max.Width = cmax.Width
		}
		if height < 1 {
			pref.Height += cpref.Height
			max.Height += cmax.Height
		}
	}
	if height >= 1 {
		count := float64(len(l.rows))
		if count < 1 {
			count = 1
		}
		pref.Height = count * height
		max.Height = count * height
		if max.Height < layout.DefaultMaxSize {
			max.Height = layout.DefaultMaxSize
		}
	}
	if border := l.Border(); border != nil {
		insets := border.Insets()
		pref.AddInsets(insets)
		max.AddInsets(insets)
	}
	pref.GrowToInteger()
	max.GrowToInteger()
	return pref, pref, max
}

// DefaultDraw provides the default drawing.
func (l *List) DefaultDraw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	index, y := l.rowAt(dirty.Y)
	if index >= 0 {
		cellHeight := math.Ceil(l.factory.CellHeight())
		count := len(l.rows)
		ymax := dirty.Y + dirty.Height
		focused := l.Focused()
		selCount := l.Selection.Count()
		rect := l.ContentRect(false)
		for index < count && y < ymax {
			selected := l.Selection.State(index)
			cell := l.factory.CreateCell(l.AsPanel(), l.rows[index], index, selected, focused && selected && selCount == 1)
			cellRect := geom.Rect{Point: geom.Point{X: rect.X, Y: y}, Size: geom.Size{Width: rect.Width, Height: cellHeight}}
			if cellHeight < 1 {
				_, pref, _ := cell.Sizes(geom.Size{})
				pref.GrowToInteger()
				cellRect.Height = pref.Height
			}
			cell.SetFrameRect(cellRect)
			y += cellRect.Height
			var ink draw.Ink
			switch {
			case selected:
				ink = l.selectedBackgroundInk
			case index%2 == 0:
				ink = l.backgroundInk
			default:
				ink = l.alternateBackgroundInk
			}
			gc.Rect(geom.Rect{Point: geom.Point{X: rect.X, Y: cellRect.Y}, Size: geom.Size{Width: rect.Width, Height: cellRect.Height}})
			gc.Fill(ink)
			gc.Save()
			tl := cellRect.Point
			dirty.Point.Subtract(tl)
			gc.Translate(cellRect.X, cellRect.Y)
			cellRect.X = 0
			cellRect.Y = 0
			cell.Draw(gc, dirty, inLiveResize)
			dirty.Point.Add(tl)
			gc.Restore()
			index++
		}
	}
}

// DefaultMouseDown provides the default mouse down handling.
func (l *List) DefaultMouseDown(where geom.Point, button, clickCount int, mod keys.Modifiers) bool {
	l.RequestFocus()
	l.savedSelection = l.Selection.Clone()
	if index, _ := l.rowAt(where.Y); index >= 0 {
		switch {
		case mod.CommandDown():
			l.Selection.Flip(index)
			l.anchor = index
		case mod.ShiftDown():
			if l.anchor != -1 {
				l.Selection.SetRange(l.anchor, index)
			} else {
				l.Selection.Set(index)
				l.anchor = index
			}
		case l.Selection.State(index):
			l.anchor = index
			if clickCount == 2 && l.DoubleClickCallback != nil {
				l.DoubleClickCallback()
				return true
			}
		default:
			l.Selection.Reset()
			l.Selection.Set(index)
			l.anchor = index
		}
		if !l.Selection.Equal(l.savedSelection) {
			l.MarkForRedraw()
		}
	}
	l.pressed = true
	return true
}

// DefaultMouseDrag provides the default mouse drag handling.
func (l *List) DefaultMouseDrag(where geom.Point, button int, mod keys.Modifiers) {
	if l.pressed {
		l.Selection.Copy(l.savedSelection)
		if index, _ := l.rowAt(where.Y); index >= 0 {
			if l.anchor == -1 {
				l.anchor = index
			}
			switch {
			case mod.CommandDown():
				l.Selection.FlipRange(l.anchor, index)
			case mod.ShiftDown():
				l.Selection.SetRange(l.anchor, index)
			default:
				l.Selection.Reset()
				l.Selection.SetRange(l.anchor, index)
			}
			if !l.Selection.Equal(l.savedSelection) {
				l.MarkForRedraw()
			}
		}
	}
}

// DefaultMouseUp provides the default mouse up handling.
func (l *List) DefaultMouseUp(where geom.Point, button int, mod keys.Modifiers) {
	if l.pressed {
		l.pressed = false
		if l.NewSelectionCallback != nil && !l.Selection.Equal(l.savedSelection) {
			l.NewSelectionCallback()
		}
	}
	l.savedSelection = nil
}

// DefaultKeyDown provides the default key down handling.
func (l *List) DefaultKeyDown(keyCode int, ch rune, mod keys.Modifiers, repeat bool) bool {
	if keys.IsControlAction(keyCode) {
		if l.DoubleClickCallback != nil && l.Selection.Count() > 0 {
			l.DoubleClickCallback()
		}
	} else {
		switch keyCode {
		case keys.Up.Code, keys.NumpadUp.Code:
			var first int
			if l.Selection.Count() == 0 {
				first = len(l.rows) - 1
			} else {
				first = l.Selection.FirstSet() - 1
				if first < 0 {
					first = 0
				}
			}
			l.Select(mod.ShiftDown(), first)
			if l.NewSelectionCallback != nil {
				l.NewSelectionCallback()
			}
		case keys.Down.Code, keys.NumpadDown.Code:
			last := l.Selection.LastSet() + 1
			if last >= len(l.rows) {
				last = len(l.rows) - 1
			}
			l.Select(mod.ShiftDown(), last)
			if l.NewSelectionCallback != nil {
				l.NewSelectionCallback()
			}
		case keys.Home.Code, keys.NumpadHome.Code:
			l.Select(mod.ShiftDown(), 0)
			if l.NewSelectionCallback != nil {
				l.NewSelectionCallback()
			}
		case keys.End.Code, keys.NumpadEnd.Code:
			l.Select(mod.ShiftDown(), len(l.rows)-1)
			if l.NewSelectionCallback != nil {
				l.NewSelectionCallback()
			}
		default:
			return false
		}
	}
	return true
}

// DefaultCanPerformCmd provides the default can perform cmd handling.
func (l *List) DefaultCanPerformCmd(source interface{}, id int) bool {
	return id == ids.SelectAllItemID && l.Selection.Count() < len(l.rows)
}

// DefaultPerformCmd provides the default perform cmd handling.
func (l *List) DefaultPerformCmd(source interface{}, id int) {
	if id == ids.SelectAllItemID {
		l.SelectRange(0, len(l.rows)-1, false)
	}
}

// SelectRange selects items from 'start' to 'end', inclusive. If 'add' is
// true, then any existing selection is added to rather than replaced.
func (l *List) SelectRange(start, end int, add bool) {
	if !add {
		l.Selection.Reset()
		l.anchor = -1
	}
	max := len(l.rows) - 1
	start = xmath.MaxInt(xmath.MinInt(start, max), 0)
	end = xmath.MaxInt(xmath.MinInt(end, max), 0)
	l.Selection.SetRange(start, end)
	if l.anchor == -1 {
		l.anchor = start
	}
	l.MarkForRedraw()
}

// Select items at the specified indexes. If 'add' is true, then any existing
// selection is added to rather than replaced.
func (l *List) Select(add bool, index ...int) {
	if !add {
		l.Selection.Reset()
		l.anchor = -1
	}
	max := len(l.rows)
	for _, v := range index {
		if v >= 0 && v < max {
			l.Selection.Set(v)
			if l.anchor == -1 {
				l.anchor = v
			}
		}
	}
	l.MarkForRedraw()
}

func (l *List) rowAt(y float64) (index int, top float64) {
	count := len(l.rows)
	top = l.ContentRect(false).Y
	cellHeight := math.Ceil(l.factory.CellHeight())
	if cellHeight < 1 {
		for index < count {
			cell := l.factory.CreateCell(l.AsPanel(), l.rows[index], index, false, false)
			_, pref, _ := cell.Sizes(geom.Size{})
			pref.GrowToInteger()
			if top+pref.Height >= y {
				break
			}
			top += pref.Height
			index++
		}
	} else {
		index = int(math.Floor((y - top) / cellHeight))
		top += float64(index) * cellHeight
	}
	if index >= count {
		index = -1
		top = 0
	}
	return index, top
}

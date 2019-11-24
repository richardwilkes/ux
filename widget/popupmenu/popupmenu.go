package popupmenu

import (
	"fmt"
	"math"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/ids"
	"github.com/richardwilkes/ux/keys"
	"github.com/richardwilkes/ux/layout"
	"github.com/richardwilkes/ux/layout/align"
	"github.com/richardwilkes/ux/menu"
	"github.com/richardwilkes/ux/widget"
	"github.com/richardwilkes/ux/widget/label"
)

// PopupMenu represents a clickable button that displays a menu of choices.
type PopupMenu struct {
	ux.Panel
	SelectionCallback    func()
	items                []interface{}
	Font                 *draw.Font // The font to use
	BackgroundInk        draw.Ink   // The background ink when enabled but not pressed or focused
	FocusedBackgroundInk draw.Ink   // The background ink when enabled and focused
	PressedBackgroundInk draw.Ink   // The background ink when enabled and pressed
	EdgeInk              draw.Ink   // The ink to use on the edges
	TextInk              draw.Ink   // The text ink to use
	PressedTextInk       draw.Ink   // The text ink when enabled and pressed
	CornerRadius         float64    // The amount of rounding to use on the corners
	HMargin              float64    // The margin on the left and right side of the content
	VMargin              float64    // The margin on the top and bottom of the content
	selectedIndex        int
	Pressed              bool
}

type separationMarker struct {
}

// New creates a new PopupMenu.
func New() *PopupMenu {
	p := &PopupMenu{
		Font:                 draw.SystemFont,
		BackgroundInk:        draw.ControlBackgroundInk,
		FocusedBackgroundInk: draw.ControlFocusedBackgroundInk,
		PressedBackgroundInk: draw.ControlPressedBackgroundInk,
		EdgeInk:              draw.ControlEdgeAdjColor,
		TextInk:              draw.ControlTextColor,
		PressedTextInk:       draw.AlternateSelectedControlTextColor,
		CornerRadius:         4,
		HMargin:              8,
		VMargin:              1,
	}
	p.InitTypeAndID(p)
	p.SetFocusable(true)
	p.SetSizer(p.DefaultSizes)
	p.DrawCallback = p.DefaultDraw
	p.GainedFocusCallback = p.MarkForRedraw
	p.LostFocusCallback = p.MarkForRedraw
	p.MouseDownCallback = p.DefaultMouseDown
	p.KeyDownCallback = p.DefaultKeyDown
	return p
}

// DefaultSizes provides the default sizing.
func (p *PopupMenu) DefaultSizes(hint geom.Size) (min, pref, max geom.Size) {
	pref = label.Size("M", p.Font, nil, 0, 0)
	for _, one := range p.items {
		switch one.(type) {
		case *separationMarker:
		default:
			size := label.Size(fmt.Sprintf("%v", one), p.Font, nil, 0, 0)
			if pref.Width < size.Width {
				pref.Width = size.Width
			}
			if pref.Height < size.Height {
				pref.Height = size.Height
			}
		}
	}
	if border := p.Border(); border != nil {
		pref.AddInsets(border.Insets())
	}
	pref.Height += p.VMargin*2 + 2
	pref.Width += p.HMargin*2 + 2 + pref.Height*0.75
	pref.GrowToInteger()
	pref.ConstrainForHint(hint)
	max.Width = math.Max(layout.DefaultMaxSize, pref.Width)
	max.Height = pref.Height
	return pref, pref, max
}

// DefaultDraw provides the default drawing.
func (p *PopupMenu) DefaultDraw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	if !p.Enabled() {
		gc.SetOpacity(0.33)
	}
	rect := p.ContentRect(false)
	widget.DrawRoundedRectBase(gc, rect, p.CornerRadius, p.currentBackgroundInk(), p.EdgeInk)
	rect.InsetUniform(1.5)
	rect.X += p.HMargin
	rect.Y += p.VMargin
	rect.Width -= p.HMargin * 2
	rect.Height -= p.VMargin * 2
	triWidth := rect.Height * 0.75
	triHeight := triWidth / 2
	rect.Width -= triWidth
	label.Draw(gc, rect, align.Start, align.Middle, p.Text(), p.Font, p.currentTextInk(), nil, 0, 0, p.Enabled())
	rect.Width += triWidth + p.HMargin/2
	gc.MoveTo(rect.X+rect.Width, rect.Y+(rect.Height-triHeight)/2)
	gc.LineTo(rect.X+rect.Width-triWidth, rect.Y+(rect.Height-triHeight)/2)
	gc.LineTo(rect.X+rect.Width-triWidth/2, rect.Y+(rect.Height-triHeight)/2+triHeight)
	gc.ClosePath()
	gc.Fill(p.currentTextInk())
}

// Text the currently shown text.
func (p *PopupMenu) Text() string {
	if p.selectedIndex >= 0 && p.selectedIndex < len(p.items) {
		item := p.items[p.selectedIndex]
		switch item.(type) {
		case *separationMarker:
		default:
			return fmt.Sprintf("%v", item)
		}
	}
	return ""
}

func (p *PopupMenu) currentBackgroundInk() draw.Ink {
	switch {
	case p.Pressed:
		return p.PressedBackgroundInk
	case p.Focused():
		return p.FocusedBackgroundInk
	default:
		return p.BackgroundInk
	}
}

func (p *PopupMenu) currentTextInk() draw.Ink {
	if p.Pressed || p.Focused() {
		return p.PressedTextInk
	}
	return p.TextInk
}

// Click performs any animation associated with a click and triggers the
// popup menu to appear.
func (p *PopupMenu) Click() {
	hasItem := false
	m := menu.New(ids.PopupMenuTemporaryBaseID, "", nil)
	defer m.Dispose()
	for i := range p.items {
		if p.addItemToMenu(m, i) {
			hasItem = true
		}
	}
	if hasItem {
		m.Popup(p.Window(), p.ToRoot(p.ContentRect(true).Point), p.selectedIndex)
	}
}

func (p *PopupMenu) addItemToMenu(m *menu.Menu, index int) bool {
	one := p.items[index]
	switch one.(type) {
	case *separationMarker:
		m.InsertSeparator(-1)
		return false
	default:
		m.InsertItem(-1, ids.PopupMenuTemporaryBaseID+1+index, fmt.Sprintf("%v", one), nil, 0, nil, func() {
			if index != p.SelectedIndex() {
				p.SelectIndex(index)
			}
		})
		return true
	}
}

// AddItem appends an item to the end of the PopupMenu.
func (p *PopupMenu) AddItem(item interface{}) {
	p.items = append(p.items, item)
}

// AddSeparator adds a separator to the end of the PopupMenu.
func (p *PopupMenu) AddSeparator() {
	p.items = append(p.items, &separationMarker{})
}

// IndexOfItem returns the index of the specified item. -1 will be returned if
// the item isn't present.
func (p *PopupMenu) IndexOfItem(item interface{}) int {
	for i, one := range p.items {
		if one == item {
			return i
		}
	}
	return -1
}

// RemoveItem from the PopupMenu.
func (p *PopupMenu) RemoveItem(item interface{}) {
	p.RemoveItemAt(p.IndexOfItem(item))
}

// RemoveItemAt the specified index from the PopupMenu.
func (p *PopupMenu) RemoveItemAt(index int) {
	if index >= 0 {
		length := len(p.items)
		if index < length {
			if p.selectedIndex == index {
				if p.selectedIndex > length-2 {
					p.selectedIndex = length - 2
					if p.selectedIndex < 0 {
						p.selectedIndex = 0
					}
				}
				p.MarkForRedraw()
			} else if p.selectedIndex > index {
				p.selectedIndex--
			}
			copy(p.items[index:], p.items[index+1:])
			length--
			p.items[length] = nil
			p.items = p.items[:length]
		}
	}
}

// Selected returns the currently selected item or nil.
func (p *PopupMenu) Selected() interface{} {
	if p.selectedIndex >= 0 && p.selectedIndex < len(p.items) {
		return p.items[p.selectedIndex]
	}
	return nil
}

// SelectedIndex returns the currently selected item index.
func (p *PopupMenu) SelectedIndex() int {
	return p.selectedIndex
}

// Select an item.
func (p *PopupMenu) Select(item interface{}) {
	p.SelectIndex(p.IndexOfItem(item))
}

// SelectIndex selects an item by its index.
func (p *PopupMenu) SelectIndex(index int) {
	if index != p.selectedIndex && index >= 0 && index < len(p.items) {
		p.selectedIndex = index
		p.MarkForRedraw()
		if p.SelectionCallback != nil {
			p.SelectionCallback()
		}
	}
}

// DefaultMouseDown provides the default mouse down handling.
func (p *PopupMenu) DefaultMouseDown(where geom.Point, button, clickCount int, mod keys.Modifiers) bool {
	p.Click()
	return true
}

// DefaultKeyDown provides the default key down handling.
func (p *PopupMenu) DefaultKeyDown(keyCode int, ch rune, mod keys.Modifiers, repeat bool) bool {
	if keys.IsControlAction(keyCode) {
		p.Click()
		return true
	}
	return false
}

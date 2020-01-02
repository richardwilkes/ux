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
)

// PopupMenu represents a clickable button that displays a menu of choices.
type PopupMenu struct {
	ux.Panel
	managed
	SelectionCallback func()
	items             []interface{}
	selectedIndex     int
	Pressed           bool
}

type separationMarker struct {
}

// New creates a new PopupMenu.
func New() *PopupMenu {
	p := &PopupMenu{}
	p.managed.initialize()
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
	pref = widget.LabelSize("M", p.font, nil, 0, 0)
	for _, one := range p.items {
		switch one.(type) {
		case *separationMarker:
		default:
			size := widget.LabelSize(fmt.Sprintf("%v", one), p.font, nil, 0, 0)
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
	pref.Height += p.vMargin*2 + 2
	pref.Width += p.hMargin*2 + 2 + pref.Height*0.75
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
	widget.DrawRoundedRectBase(gc, rect, p.cornerRadius, p.currentBackgroundInk(), p.edgeInk)
	rect.InsetUniform(1.5)
	rect.X += p.hMargin
	rect.Y += p.vMargin
	rect.Width -= p.hMargin * 2
	rect.Height -= p.vMargin * 2
	triWidth := rect.Height * 0.75
	triHeight := triWidth / 2
	rect.Width -= triWidth
	widget.DrawLabel(gc, rect, align.Start, align.Middle, p.Text(), p.font, p.currentTextInk(), nil, 0, 0, p.Enabled())
	rect.Width += triWidth + p.hMargin/2
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
		return p.pressedBackgroundInk
	case p.Focused():
		return p.focusedBackgroundInk
	default:
		return p.backgroundInk
	}
}

func (p *PopupMenu) currentTextInk() draw.Ink {
	if p.Pressed || p.Focused() {
		return p.pressedTextInk
	}
	return p.textInk
}

// Click performs any animation associated with a click and triggers the
// popup menu to appear.
func (p *PopupMenu) Click() {
	hasItem := false
	m := menu.New("", nil)
	defer m.Dispose()
	for i, item := range p.items {
		if _, ok := item.(*separationMarker); ok {
			m.InsertSeparator(-1)
		} else {
			hasItem = true
			index := i
			m.InsertItem(-1, ids.PopupMenuTemporaryBaseID+index, fmt.Sprintf("%v", item), nil, 0, nil, func(*menu.Item) {
				if index != p.SelectedIndex() {
					p.SelectIndex(index)
				}
			})
		}
	}
	if hasItem {
		m.Popup(p.Window(), p.RectToRoot(p.ContentRect(true)), p.selectedIndex)
	}
}

// AddItem appends an item to the end of the PopupMenu.
func (p *PopupMenu) AddItem(item interface{}) *PopupMenu {
	p.items = append(p.items, item)
	return p
}

// AddSeparator adds a separator to the end of the PopupMenu.
func (p *PopupMenu) AddSeparator() *PopupMenu {
	p.items = append(p.items, &separationMarker{})
	return p
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
func (p *PopupMenu) RemoveItem(item interface{}) *PopupMenu {
	p.RemoveItemAt(p.IndexOfItem(item))
	return p
}

// RemoveItemAt the specified index from the PopupMenu.
func (p *PopupMenu) RemoveItemAt(index int) *PopupMenu {
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
	return p
}

// ItemCount returns the number of items in this PopupMenu.
func (p *PopupMenu) ItemCount() int {
	return len(p.items)
}

// ItemAt returns the item at the specified index or nil.
func (p *PopupMenu) ItemAt(index int) interface{} {
	if index >= 0 && index < len(p.items) {
		return p.items[index]
	}
	return nil
}

// Selected returns the currently selected item or nil.
func (p *PopupMenu) Selected() interface{} {
	return p.ItemAt(p.selectedIndex)
}

// SelectedIndex returns the currently selected item index.
func (p *PopupMenu) SelectedIndex() int {
	return p.selectedIndex
}

// Select an item.
func (p *PopupMenu) Select(item interface{}) *PopupMenu {
	p.SelectIndex(p.IndexOfItem(item))
	return p
}

// SelectIndex selects an item by its index.
func (p *PopupMenu) SelectIndex(index int) *PopupMenu {
	if index != p.selectedIndex && index >= 0 && index < len(p.items) {
		p.selectedIndex = index
		p.MarkForRedraw()
		if p.SelectionCallback != nil {
			p.SelectionCallback()
		}
	}
	return p
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

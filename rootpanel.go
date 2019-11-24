package ux

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/layout"
)

type rootPanel struct {
	Panel
	window  *Window
	menubar *Panel
	content *Panel
	tooltip *Panel
}

func newRootPanel(wnd *Window) *rootPanel {
	p := &rootPanel{}
	p.InitTypeAndID(p)
	p.SetLayout(&rootLayout{root: p})
	p.window = wnd
	content := NewPanel()
	layout.NewFlow(content, 4, 2)
	p.setContent(content)
	return p
}

func (p *rootPanel) setMenuBar(bar *Panel) { //nolint:unused
	if p.menubar != nil {
		p.RemoveChild(p.menubar)
	}
	p.menubar = bar
	if bar != nil {
		p.AddChildAtIndex(bar, 0)
	}
	p.NeedsLayout = true
	p.MarkForRedraw()
}

func (p *rootPanel) setContent(content *Panel) {
	if p.content != nil {
		p.RemoveChild(p.content)
	}
	p.content = content
	if content != nil {
		index := 0
		if p.menubar != nil {
			index = 1
		}
		p.AddChildAtIndex(content, index)
	}
	p.NeedsLayout = true
	p.MarkForRedraw()
}

func (p *rootPanel) setTooltip(tip *Panel) {
	if p.tooltip != nil {
		p.tooltip.MarkForRedraw()
		p.RemoveChild(p.tooltip)
	}
	p.tooltip = tip
	if tip != nil {
		p.AddChild(tip)
		tip.MarkForRedraw()
	}
}

type rootLayout struct {
	root *rootPanel
}

func (l *rootLayout) Sizes(hint geom.Size) (min, pref, max geom.Size) {
	min, pref, max = l.root.content.Sizes(hint)
	if l.root.menubar != nil {
		_, barSize, _ := l.root.menubar.Sizes(geom.Size{})
		for _, size := range []*geom.Size{&min, &pref, &max} {
			size.Height += barSize.Height
			if size.Width < barSize.Width {
				size.Width = barSize.Width
			}
		}
	}
	return
}

func (l *rootLayout) Layout() {
	rect := l.root.frame
	rect.X = 0
	rect.Y = 0
	if l.root.menubar != nil {
		_, size, _ := l.root.menubar.Sizes(geom.Size{})
		l.root.menubar.SetFrameRect(geom.Rect{Size: geom.Size{Width: rect.Width, Height: size.Height}})
		rect.Y += size.Height
		rect.Height -= size.Height
	}
	l.root.content.SetFrameRect(rect)
}

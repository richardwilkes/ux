// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package scrollarea

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/layout"
	"github.com/richardwilkes/ux/widget/scrollarea/behavior"
)

type scrollAreaLayout struct {
	scrollArea *ScrollArea
}

func (l *scrollAreaLayout) Sizes(hint geom.Size) (min, pref, max geom.Size) {
	_, hBarSize, _ := l.scrollArea.ScrollBar(true).Sizes(geom.Size{})
	_, vBarSize, _ := l.scrollArea.ScrollBar(false).Sizes(geom.Size{})
	min.Width = vBarSize.Width * 2
	min.Height = hBarSize.Height * 2
	if l.scrollArea.content != nil {
		_, pref, _ = l.scrollArea.content.Sizes(hint)
	}
	if pref.Width < min.Width {
		pref.Width = min.Width
	}
	if pref.Height < min.Height {
		pref.Height = min.Height
	}
	if b := l.scrollArea.Border(); b != nil {
		insets := b.Insets()
		min.AddInsets(insets)
		pref.AddInsets(insets)
		max.AddInsets(insets)
	}
	return min, pref, layout.MaxSize(pref)
}

func (l *scrollAreaLayout) Layout() {
	hBar := l.scrollArea.ScrollBar(true)
	_, hBarSize, _ := hBar.Sizes(geom.Size{})
	vBar := l.scrollArea.ScrollBar(false)
	_, vBarSize, _ := vBar.Sizes(geom.Size{})
	needHBar := false
	needVBar := false
	var insets geom.Insets
	if b := l.scrollArea.Border(); b != nil {
		insets = b.Insets()
	}
	area := l.scrollArea.ContentRect(false)
	visibleSize := area.Size
	var contentSize geom.Size
	var prefContentSize geom.Size
	if l.scrollArea.content != nil {
		var hint geom.Size
		if l.scrollArea.behavior == behavior.FollowsWidth {
			hint.Width = area.Width
		} else if l.scrollArea.behavior == behavior.FollowsHeight {
			hint.Height = area.Height
		}
		_, prefContentSize, _ = l.scrollArea.content.Sizes(hint)
		contentSize = prefContentSize
		switch l.scrollArea.behavior {
		case behavior.FillWidth:
			if visibleSize.Width > contentSize.Width {
				contentSize.Width = visibleSize.Width
			}
		case behavior.FillHeight:
			if visibleSize.Height > contentSize.Height {
				contentSize.Height = visibleSize.Height
			}
		case behavior.Fill:
			if visibleSize.Width > contentSize.Width {
				contentSize.Width = visibleSize.Width
			}
			if visibleSize.Height > contentSize.Height {
				contentSize.Height = visibleSize.Height
			}
		case behavior.FollowsWidth:
			prefContentSize.Width = visibleSize.Width
			contentSize.Width = visibleSize.Width
		case behavior.FollowsHeight:
			prefContentSize.Height = visibleSize.Height
			contentSize.Height = visibleSize.Height
		default:
		}
	}
	if visibleSize.Width < contentSize.Width {
		visibleSize.Height -= hBarSize.Height
		if insets.Bottom >= 1 {
			visibleSize.Height++
		}
		if l.scrollArea.behavior == behavior.FillHeight || l.scrollArea.behavior == behavior.Fill {
			if visibleSize.Height > prefContentSize.Height {
				contentSize.Height = visibleSize.Height
			}
		} else if l.scrollArea.behavior == behavior.FollowsHeight {
			contentSize.Height = visibleSize.Height
		}
		needHBar = true
	}
	if visibleSize.Height < contentSize.Height {
		visibleSize.Width -= vBarSize.Width
		if insets.Right >= 1 {
			visibleSize.Width++
		}
		if l.scrollArea.behavior == behavior.FillWidth || l.scrollArea.behavior == behavior.Fill {
			if visibleSize.Width > prefContentSize.Width {
				contentSize.Width = visibleSize.Width
			}
		} else if l.scrollArea.behavior == behavior.FollowsWidth {
			contentSize.Width = visibleSize.Width
		}
		needVBar = true
	}
	if !needHBar && visibleSize.Width < contentSize.Width {
		visibleSize.Height -= hBarSize.Height
		if insets.Bottom >= 1 {
			visibleSize.Height++
		}
		if l.scrollArea.behavior == behavior.FillHeight || l.scrollArea.behavior == behavior.Fill {
			if visibleSize.Height > prefContentSize.Height {
				contentSize.Height = visibleSize.Height
			}
		} else if l.scrollArea.behavior == behavior.FollowsHeight {
			contentSize.Height = visibleSize.Height
		}
		needHBar = true
	}
	if needHBar {
		if hBar.Parent() == nil {
			l.scrollArea.AddChild(hBar.AsPanel())
		}
	} else {
		hBar.RemoveFromParent()
	}
	if needVBar {
		if vBar.Parent() == nil {
			l.scrollArea.AddChild(vBar.AsPanel())
		}
	} else {
		vBar.RemoveFromParent()
	}
	l.scrollArea.view.SetFrameRect(geom.Rect{Point: area.Point, Size: visibleSize})
	if needHBar {
		hBarSize.Width = visibleSize.Width
		barRect := geom.Rect{Point: geom.Point{X: area.X, Y: area.Y + visibleSize.Height}, Size: hBarSize}
		if insets.Left >= 1 {
			barRect.X--
			barRect.Width++
		}
		if insets.Right >= 1 {
			barRect.Width++
		}
		hBar.SetFrameRect(barRect)
	}
	if needVBar {
		vBarSize.Height = visibleSize.Height
		barRect := geom.Rect{Point: geom.Point{X: area.X + visibleSize.Width, Y: area.Y}, Size: vBarSize}
		if insets.Top >= 1 {
			barRect.Y--
			barRect.Height++
		}
		if insets.Bottom >= 1 {
			barRect.Height++
		}
		vBar.SetFrameRect(barRect)
	}
	if l.scrollArea.content != nil {
		contentRect := l.scrollArea.content.FrameRect()
		contentRect.Size = contentSize
		l.scrollArea.content.SetFrameRect(contentRect)
	}
}

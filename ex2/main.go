// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package main

import (
	"runtime"

	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/log/jotrotate"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/layout/align"
	"github.com/richardwilkes/ux/layout/flex"
	"github.com/richardwilkes/ux/menu"
)

func main() {
	runtime.LockOSThread()

	cmdline.AppName = "Example"
	cmdline.AppCmdName = "example"
	cmdline.AppVersion = "0.1"
	cmdline.CopyrightYears = "2019"
	cmdline.CopyrightHolder = "Richard A. Wilkes"
	cmdline.AppIdentifier = "com.trollworks.ux.example"

	cl := cmdline.New(true)
	jotrotate.ParseAndSetup(cl)

	ux.WillFinishStartupCallback = finishStartup
	ux.QuittingCallback = func() { jot.Info("Quitting") }
	ux.Start() // Never returns
}

func finishStartup() {
	wnd, err := ux.NewWindow("Test", geom.Rect{}, ux.StdWindowMask)
	jot.FatalIfErr(err)
	if bar, global, first := menu.BarForWindow(wnd, nil); !global || first {
		bar.InsertStdMenus(nil, nil, nil)
	}
	content := wnd.Content()
	content.SetBorder(border.NewEmpty(geom.NewUniformInsets(10)))
	flex.New().VSpacing(10).Apply(content)
	panel := ux.NewPanel()
	panel.DrawCallback = func(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
		rect := panel.ContentRect(false)
		// gc.Rect(rect)
		gc.RoundedRect(rect, 8)
		// gc.Ellipse(rect)
		gc.Fill(draw.Orange)
	}
	flex.NewData().HGrab(true).HAlign(align.Fill).MinSize(geom.Size{Width: 200, Height: 30}).Apply(panel)
	// flow.New().HSpacing(5).VSpacing(5).Apply(panel)
	// btn := button.New().SetText("Press Me")
	// btn.ClickCallback = func() { jot.Infof("%v was clicked.", btn) }
	// btn.Tooltip = tooltip.NewWithText(fmt.Sprintf("This is the tooltip for %v", btn))
	// btn.SetLayoutData(align.Middle)
	// panel.AddChild(btn.AsPanel())
	content.AddChild(panel)

	wnd.Pack()
	wndRect := wnd.FrameRect()
	wndRect.X = 100
	wndRect.Y = 100
	wnd.SetFrameRect(wndRect)
	wnd.ToFront()
}

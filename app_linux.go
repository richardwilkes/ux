// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package ux

import (
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/richardwilkes/toolbox/atexit"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/globals"
)

var awaitingQuitDecision bool

func osStart() {
	mousebind.Initialize(globals.X11)
	draw.UpdateSystemColors()
	draw.Initialize()
	if WillFinishStartupCallback != nil {
		WillFinishStartupCallback()
	}
	if DidFinishStartupCallback != nil {
		DidFinishStartupCallback()
	}
	xevent.Main(globals.X11)
}

func osAttemptQuit() {
	response := Now
	if CheckQuitCallback == nil {
		response = CheckQuitCallback()
	}
	switch response {
	case Cancel:
		return
	case Now:
		atexit.Exit(0)
	case Later:
		awaitingQuitDecision = true
	}
}

func osMayQuitNow(quit bool) {
	if awaitingQuitDecision {
		awaitingQuitDecision = false
		if quit {
			atexit.Exit(0)
		}
	} else {
		jot.Error("call to MayQuitNow without AttemptQuit")
	}
}

func osHideApp() {
	// Not supported
}

func osHideOtherApps() {
	// Not supported
}

func osShowAllApps() {
	// Not supported
}

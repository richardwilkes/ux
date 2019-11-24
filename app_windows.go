package ux

import (
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/globals"
	"github.com/richardwilkes/win32"
)

var awaitingQuitDecision bool

func osStart() {
	RegisterWindowClass()
	draw.UpdateSystemColors()
	draw.Initialize()
	if WillFinishStartupCallback != nil {
		WillFinishStartupCallback()
	}
	if DidFinishStartupCallback != nil {
		DidFinishStartupCallback()
	}
	for {
		var msg win32.MSG
		quit, err := win32.GetMessage(&msg, win32.NULL, 0, 0)
		if err != nil {
			jot.Error(errs.Wrap(err))
			continue
		}
		if quit {
			break
		}
		if msg.Message == globals.InvokeMsgID {
			dispatchTask(uint64(uint32(msg.WParam))<<32 | uint64(uint32(msg.LParam)))
		} else {
			win32.TranslateMessage(&msg)
			win32.DispatchMessage(&msg)
		}
	}
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
		win32.PostQuitMessage(0)
	case Later:
		awaitingQuitDecision = true
	}
}

func osMayQuitNow(quit bool) {
	if awaitingQuitDecision {
		awaitingQuitDecision = false
		if quit {
			win32.PostQuitMessage(0)
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

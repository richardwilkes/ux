package ux

import (
	"runtime"
	"sync"

	"github.com/richardwilkes/toolbox/atexit"
	"github.com/richardwilkes/ux/globals"
)

// Possible termination responses
const (
	Cancel QuitResponse = iota
	Now
	Later // Must make a call to MayQuitNow() at some point in the future.
)

var (
	// WillFinishStartupCallback is called right before application startup
	// has completed. This is a good point to create any windows your app
	// wants to display.
	WillFinishStartupCallback func()

	// DidFinishStartupCallback is called once application startup has
	// completed and it is about to start servicing the event loop.
	DidFinishStartupCallback func()

	// WillActivateCallback is called right before the application is
	// activated.
	WillActivateCallback func()

	// DidActivateCallback is called once the application is activated.
	DidActivateCallback func()

	// WillDeactivateCallback is called right before the application is
	// deactivated.
	WillDeactivateCallback func()

	// DidDeactivateCallback is called once the application is deactivated.
	DidDeactivateCallback func()

	// OpenURLsCallback is called when the application is asked to open one or
	// more URLs by the OS or external application.
	OpenURLsCallback func(urls []string)

	// OSThemeChangedCallback is called when the OS theme is changed.
	OSThemeChangedCallback func()

	// QuitAfterLastWindowClosedCallback is called when the last window is
	// closed to determine if the application should quit as a result.
	QuitAfterLastWindowClosedCallback func() bool

	// CheckQuitCallback is called when termination has been requested.
	CheckQuitCallback func() QuitResponse

	// QuittingCallback is called when the app will in fact terminate.
	QuittingCallback func()
	quitOnce         sync.Once
)

// QuitResponse is used to respond to requests for app termination.
type QuitResponse int

// Start the application. This function does NOT return. No calls to anything
// in the ux package tree should be made before this call, although setting
// the various callbacks above can and should be done before this call.
func Start() {
	runtime.LockOSThread()
	atexit.Register(runQuitCallback)
	globals.Initialize()
	osStart()
	atexit.Exit(0)
}

// HideApp will hide the application on platforms where that is supported.
func HideApp() {
	osHideApp()
}

// HideOtherApps will hide all other applications on platforms where that is
// supported.
func HideOtherApps() {
	osHideOtherApps()
}

// ShowAllApps will unhide any hidden applications on platforms where that is
// supported.
func ShowAllApps() {
	osShowAllApps()
}

// AttemptQuit initiates the termination sequence.
func AttemptQuit() {
	osAttemptQuit()
}

// MayQuitNow resumes the termination sequence that was delayed when Later is
// returned from the CheckQuitCallback. Passing in false for the quit
// parameter will cancel the termination sequence while true will allow it to
// proceed.
func MayQuitNow(quit bool) {
	osMayQuitNow(quit)
}

func runQuitCallback() {
	if QuittingCallback != nil {
		quitOnce.Do(QuittingCallback)
	}
}

package ux

import (
	"github.com/richardwilkes/macos/ns"
	"github.com/richardwilkes/toolbox/atexit"
	"github.com/richardwilkes/ux/draw"
)

func osStart() {
	pool := ns.NewAutoreleasePool()
	defer pool.Release()
	app := ns.SharedApplication()
	app.SetDelegate(&delegate{})
	// Required for apps without bundle & Info.plist
	app.SetActivationPolicy(ns.ApplicationActivationPolicyRegular)
	// Required to use 'NSApplicationActivateIgnoringOtherApps' otherwise our windows end up in the background.
	ns.RunningApplicationCurrent().ActivateWithOptions(ns.ApplicationActivateAllWindows | ns.ApplicationActivateIgnoringOtherApps)
	app.Run()
}

func osAttemptQuit() {
	ns.SharedApplication().Terminate(nil)
}

func osMayQuitNow(quit bool) {
	ns.SharedApplication().ReplyToApplicationShouldTerminate(quit)
}

func osHideApp() {
	ns.RunningApplicationCurrent().Hide()
}

func osHideOtherApps() {
	ns.SharedApplication().HideOtherApplications(nil)
}

func osShowAllApps() {
	ns.SharedApplication().UnhideAllApplications(nil)
}

type delegate struct {
}

func (d *delegate) ApplicationWillFinishLaunching(notification *ns.Notification) {
	draw.UpdateSystemColors()
	draw.Initialize()
	if WillFinishStartupCallback != nil {
		WillFinishStartupCallback()
	}
}

func (d *delegate) ApplicationDidFinishLaunching(notification *ns.Notification) {
	if DidFinishStartupCallback != nil {
		DidFinishStartupCallback()
	}
}

func (d *delegate) ApplicationShouldTerminate(sender *ns.Application) ns.ApplicationTerminateReply {
	if CheckQuitCallback != nil {
		return ns.ApplicationTerminateReply(CheckQuitCallback())
	}
	return ns.TerminateNow
}

func (d *delegate) ApplicationShouldTerminateAfterLastWindowClosed(app *ns.Application) bool {
	if QuitAfterLastWindowClosedCallback != nil {
		return QuitAfterLastWindowClosedCallback()
	}
	return true
}

func (d *delegate) ApplicationWillTerminate(notification *ns.Notification) {
	atexit.Exit(0)
}

func (d *delegate) ApplicationWillBecomeActive(notification *ns.Notification) {
	if WillActivateCallback != nil {
		WillActivateCallback()
	}
}

func (d *delegate) ApplicationDidBecomeActive(notification *ns.Notification) {
	if DidActivateCallback != nil {
		DidActivateCallback()
	}
}

func (d *delegate) ApplicationWillResignActive(notification *ns.Notification) {
	if WillDeactivateCallback != nil {
		WillDeactivateCallback()
	}
}

func (d *delegate) ApplicationDidResignActive(notification *ns.Notification) {
	if DidDeactivateCallback != nil {
		DidDeactivateCallback()
	}
}

func (d *delegate) ThemeChanged(notification *ns.Notification) {
	draw.MarkSystemColorsForUpdate()
	if OSThemeChangedCallback != nil {
		OSThemeChangedCallback()
	}
	for _, wnd := range Windows() {
		wnd.MarkForRedraw()
	}
}

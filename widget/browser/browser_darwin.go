// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package browser

import (
	"github.com/richardwilkes/macos/ns"
	"github.com/richardwilkes/macos/wk"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
)

var browserCallbackLog = false

type osBrowser = *wk.WebView

func osNewBrowser(wnd *ux.Window) (osBrowser, error) {
	webView := wk.WebViewInitWithFrameConfiguration(0, 0, 0, 0, wk.NewWebViewConfiguration(), &webViewDelegate{})
	wnd.OSWindow().ContentView().AddSubview(&webView.View)
	return webView, nil
}

func (b *Browser) osSetFrame(rect geom.Rect) {
	b.browser.SetFrame(rect.X, rect.Y, rect.Width, rect.Height)
	b.browser.SetNeedsLayout(true)
	b.browser.SetNeedsDisplay(true)
}

func (b *Browser) osLoadURL(url string) {
	b.browser.LoadRequest(ns.URLRequestWithURL(ns.URLWithString(url)))
}

func (b *Browser) osDispose() {
	b.browser.RemoveFromSuperview()
	b.browser.Release()
}

type webViewDelegate struct {
}

func (d *webViewDelegate) WebViewDidCommitNavigation(webView *wk.WebView, nav *wk.Navigation) {
	if browserCallbackLog {
		// RAW: Do something?
		jot.Infof("browserDidCommitNavigation: %v", nav.Native())
	}
}

func (d *webViewDelegate) WebViewDidStartProvisionalNavigation(webView *wk.WebView, nav *wk.Navigation) {
	if browserCallbackLog {
		// RAW: Do something?
		jot.Infof("browserDidStartProvisionalNavigation: %v", nav.Native())
	}
}

func (d *webViewDelegate) WebViewDidReceiveServerRedirectForProvisionNavigation(webView *wk.WebView, nav *wk.Navigation) {
	if browserCallbackLog {
		// RAW: Do something?
		jot.Infof("browserDidReceiveServerRedirectForProvisionalNavigation: %v", nav.Native())
	}
}

func (d *webViewDelegate) WebViewDidReceiveAuthenticationChallenge(webView *wk.WebView, challenge *ns.URLAuthenticationChallenge) (disposition ns.URLSessionAuthChallengeDisposition, credential *ns.URLCredential) {
	if browserCallbackLog {
		jot.Info("browserDidReceiveAuthenticationChallenge")
	}
	trust := challenge.ProtectionSpace().ServerTrust()
	exceptions := trust.CopyExceptions()
	trust.SetExceptions(exceptions)
	exceptions.Release()
	return ns.URLSessionAuthChallengeUseCredential, ns.URLCredentialForTrust(trust)
}

func (d *webViewDelegate) WebViewDidFailNavigationWithError(webView *wk.WebView, nav *wk.Navigation, errorMsg string) {
	if browserCallbackLog {
		// RAW: Do something?
		jot.Infof("browserDidFailNavigation: %v: %s", nav.Native(), errorMsg)
	}
}

func (d *webViewDelegate) WebViewDidFailProvisionalNavigationWithError(webView *wk.WebView, nav *wk.Navigation, errorMsg string) {
	if browserCallbackLog {
		// RAW: Do something?
		jot.Infof("browserDidFailProvisionalNavigation: %v: %s", nav.Native(), errorMsg)
	}
}

func (d *webViewDelegate) WebViewDidFinishNavigation(webView *wk.WebView, nav *wk.Navigation) {
	if browserCallbackLog {
		// RAW: Do something?
		jot.Infof("browserDidFinishNavigation: %v", nav.Native())
	}
}

func (d *webViewDelegate) WebViewWebContentProcessDidTerminate(webView *wk.WebView) {
	if browserCallbackLog {
		// RAW: Do something?
		jot.Info("browserWebContentProcessDidTerminate")
	}
}

func (d *webViewDelegate) WebViewDecidePolicyForNavigationAction(webView *wk.WebView, action *wk.NavigationAction) wk.NavigationActionPolicy {
	if browserCallbackLog {
		// RAW: Do something?
		jot.Infof("browserDecidePolicyForNavigationAction: %v", action)
	}
	return wk.NavigationActionPolicyAllow
}

func (d *webViewDelegate) WebViewDecidePolicyForNavigationResponse(webView *wk.WebView, response *wk.NavigationResponse) wk.NavigationResponsePolicy {
	if browserCallbackLog {
		// RAW: Do something?
		jot.Infof("browserDecidePolicyForNavigationResponse: %v", response)
	}
	return wk.NavigationResponsePolicyAllow
}

package browser

import (
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/ux"
)

// Browser represents a native web view. Since it is a native component, it
// does not respect the panel hierarchy and effectively draws on top of all
// other panels.
type Browser struct {
	ux.Panel
	browser osBrowser
	valid   bool
}

// New creates a new, empty browser. Note that this panel behaves differently
// than other panels because it is backed by a native widget. In particular,
// you must pass in a valid Window at construction time and you must manually
// dispose of the browser when no longer needed. Not all platforms currently
// provide browser support. Those that don't will return an error.
func New(wnd *ux.Window) (*Browser, error) {
	if !wnd.IsValid() {
		return nil, errs.New("invalid window")
	}
	browser, err := osNewBrowser(wnd)
	if err != nil {
		return nil, err
	}
	b := &Browser{
		browser: browser,
		valid:   true,
	}
	b.InitTypeAndID(b)
	b.SetFocusable(true)
	b.FrameChangeCallback = b.DefaultFrameChange
	return b, nil
}

// IsValid returns true if the browser is still valid (i.e. hasn't been
// disposed).
func (b *Browser) IsValid() bool {
	return b.valid
}

// LoadURL loads the specified URL into the browser.
func (b *Browser) LoadURL(url string) {
	if b.IsValid() && url != "" {
		b.osLoadURL(url)
	}
}

// DefaultFrameChange adjusts the native component rect to match the panel.
func (b *Browser) DefaultFrameChange() {
	if b.IsValid() {
		b.osSetFrame(b.RectToRoot(b.ContentRect(false)))
	}
}

// Dispose of the browser, releasing any system resources associated with it.
func (b *Browser) Dispose() {
	b.RemoveFromParent()
	if b.valid {
		b.valid = false
		b.osDispose()
	}
}

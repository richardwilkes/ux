package browser

import (
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
)

type osBrowser = int

func osNewBrowser(wnd *ux.Window) (osBrowser, error) {
	// RAW: Implement
	return 0, errs.New("browser panel not supported")
}

func (b *Browser) osSetFrame(rect geom.Rect) {
	// RAW: Implement
}

func (b *Browser) osLoadURL(url string) {
	// RAW: Implement
}

func (b *Browser) osDispose() {
	// RAW: Implement
}

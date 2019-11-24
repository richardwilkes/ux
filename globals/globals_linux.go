package globals

import (
	"github.com/BurntSushi/xgbutil"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
)

// Server holds the connection to the X11 server.
var X11 *xgbutil.XUtil

func osInitialize() {
	var err error
	X11, err = xgbutil.NewConn()
	jot.FatalIfErr(errs.Wrap(err))
}

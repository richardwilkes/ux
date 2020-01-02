// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

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

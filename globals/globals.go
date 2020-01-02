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
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/softref"
)

var (
	// Pool is the soft reference pool used by all soft references allocated
	// by the ux packages.
	Pool = softref.NewPool(&jot.Logger{})
)

// Initialize the globals.
func Initialize() {
	osInitialize()
}

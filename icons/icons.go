// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package icons

import (
	"sync"

	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/ux/draw"
)

var (
	iconCacheLock sync.Mutex
	iconCache     = make(map[string]*draw.Image)
)

// Error returns the standard error alert image.
func Error() *draw.Image {
	return retrieveImage("/error.png")
}

// Question returns the standard question alert image.
func Question() *draw.Image {
	return retrieveImage("/question.png")
}

func retrieveImage(name string) *draw.Image {
	iconCacheLock.Lock()
	icon, ok := iconCache[name]
	if !ok {
		var err error
		icon, err = draw.NewImageFromBytes(EFS.PrimaryFileSystem().MustContentAsBytes(name), 0.5)
		jot.FatalIfErr(err)
		iconCache[name] = icon
	}
	iconCacheLock.Unlock()
	return icon
}

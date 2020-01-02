// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package draw

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/draw/quality"
)

type osImage = int

func osNewImageFromBytes(buffer []byte) (img osImage, width, height int, err error) {
	// RAW: Implement
	return 0, 0, 0, errs.New(errUnableToCreateImage)
}

func osNewImageFromData(data *ImageData) (osImage, error) {
	// RAW: Implement
	return 0, errs.New(errUnableToCreateImage)
}

func (img *imageRef) osNewScaledImage(width, height int, q quality.Quality) (osImage, error) {
	// RAW: Implement
	return 0, errs.New(errUnableToCreateImage)
}

func (img *imageRef) osIsValid() bool {
	// RAW: Implement
	return img.osImg != 0
}

func (img *imageRef) osImagePixels(pixels []Color) {
	// RAW: Implement
}

func (img *imageRef) osDrawInRect(gc Context, rect geom.Rect) {
	// RAW: Implement
}

func (img *imageRef) osDispose() {
	// RAW: Implement
}

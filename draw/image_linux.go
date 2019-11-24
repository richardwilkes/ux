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

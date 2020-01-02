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
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"unsafe"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/draw/quality"
	"github.com/richardwilkes/win32"
)

type osImage struct {
	hBitmap win32.HBITMAP
	lpBits  unsafe.Pointer
}

func osNewImageFromBytes(buffer []byte) (img osImage, width, height int, err error) {
	decoded, _, err := image.Decode(bytes.NewBuffer(buffer))
	if err != nil {
		return osImage{}, 0, 0, errs.NewWithCause(errUnableToCreateImage, err)
	}
	bounds := decoded.Bounds()
	hBitmap, lpBits, err := newWinImage(int32(bounds.Dx()), int32(bounds.Dy()))
	if err != nil {
		return osImage{}, 0, 0, err
	}
	bitmap_array := (*[1 << 30]byte)(lpBits)
	i := 0
	for y := bounds.Min.Y; y != bounds.Max.Y; y++ {
		for x := bounds.Min.X; x != bounds.Max.X; x++ {
			r, g, b, a := decoded.At(x, y).RGBA()
			bitmap_array[i] = byte(b >> 8)
			bitmap_array[i+1] = byte(g >> 8)
			bitmap_array[i+2] = byte(r >> 8)
			bitmap_array[i+3] = byte(a >> 8)
			i += 4
		}
	}
	return osImage{
		hBitmap: hBitmap,
		lpBits:  lpBits,
	}, bounds.Dx(), bounds.Dy(), nil
}

func osNewImageFromData(data *ImageData) (osImage, error) {
	hBitmap, lpBits, err := newWinImage(int32(data.Width), int32(data.Height))
	if err != nil {
		return osImage{}, err
	}
	bitmap_array := (*[1 << 30]byte)(lpBits)
	i := 0
	for _, pixel := range data.Pixels {
		bitmap_array[i] = byte(pixel.Blue())
		bitmap_array[i+1] = byte(pixel.Green())
		bitmap_array[i+2] = byte(pixel.Red())
		bitmap_array[i+3] = byte(pixel.Alpha())
		i += 4
	}
	return osImage{
		hBitmap: hBitmap,
		lpBits:  lpBits,
	}, nil
}

func newWinImage(width, height int32) (hBitmap win32.HBITMAP, lpBits unsafe.Pointer, err error) {
	var bi win32.BITMAPV5HEADER
	bi.BiSize = uint32(unsafe.Sizeof(bi))
	bi.BiWidth = width
	bi.BiHeight = -height
	bi.BiPlanes = 1
	bi.BiBitCount = 32
	bi.BiCompression = win32.BI_BITFIELDS
	bi.BV4RedMask = 0x00FF0000
	bi.BV4GreenMask = 0x0000FF00
	bi.BV4BlueMask = 0x000000FF
	bi.BV4AlphaMask = 0xFF000000
	hdc := win32.GetDC(0)
	defer win32.ReleaseDC(0, hdc)
	if hBitmap = win32.CreateDIBSection(hdc, &bi.BITMAPINFOHEADER, win32.DIB_RGB_COLORS, &lpBits, 0, 0); hBitmap == 0 {
		return 0, nil, errs.New(errUnableToCreateImage)
	}
	return hBitmap, lpBits, nil
}

func (img *imageRef) osNewScaledImage(width, height int, q quality.Quality) (osImage, error) {
	// RAW: Implement
	return osImage{}, errs.New(errUnableToCreateImage)
}

func (img *imageRef) osIsValid() bool {
	return img.osImg.hBitmap != 0
}

func (img *imageRef) osImagePixels(pixels []Color) {
	win32.GdiFlush()
	bitmap_array := (*[1 << 30]byte)(img.osImg.lpBits)
	i := 0
	for pi := range pixels {
		pixels[pi] = ARGB(float64(bitmap_array[i+3])/255, int(bitmap_array[i+2]), int(bitmap_array[i+1]), int(bitmap_array[i]))
		i += 4
	}
}

func (img *imageRef) osDrawInRect(gc Context, rect geom.Rect) {
	// RAW: Implement
}

func (img *imageRef) osDispose() {
	if img.osImg.hBitmap != 0 {
		win32.DeleteObject(win32.HGDIOBJ(img.osImg.hBitmap))
		img.osImg.hBitmap = 0
		img.osImg.lpBits = nil
	}
}

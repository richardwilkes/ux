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
	"unsafe"

	"github.com/richardwilkes/macos/cf"
	"github.com/richardwilkes/macos/cg"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/draw/quality"
)

type osImage = cg.Image

type imgContext struct {
	osImg     osImage
	osContext OSContext
}

var (
	imageLayerCache     = make(map[imgContext]cg.Layer)
	imageToContextCache = make(map[osImage][]OSContext)
	contextToImageCache = make(map[OSContext][]osImage)
)

func osNewImageFromBytes(buffer []byte) (img osImage, width, height int, err error) {
	data := cf.DataCreate(buffer)
	imgSrc := cg.ImageSourceCreateWithData(data, 0)
	data.Release()
	if imgSrc != 0 {
		defer imgSrc.Release()
		if img = imgSrc.CreateImageAtIndex(0, 0); img != 0 {
			return img, img.GetWidth(), img.GetHeight(), nil
		}
	}
	return 0, 0, 0, errs.New(errUnableToCreateImage)
}

func osNewImageFromData(data *ImageData) (osImage, error) {
	pd := make([]Color, len(data.Pixels))
	for i := range data.Pixels {
		pd[i] = data.Pixels[i].Premultiply()
	}
	colorspace := cg.ColorSpaceCreateSRGB()
	defer colorspace.Release()
	pixels := cf.DataCreate(((*[1 << 30]byte)(unsafe.Pointer(&pd[0])))[:len(pd)*4]) //nolint:gosec
	defer pixels.Release()
	dataProvider := cg.DataProviderCreateWithCFData(pixels)
	defer dataProvider.Release()
	if imgRef := cg.ImageCreate(data.Width, data.Height, 8, 32, data.Width*4, colorspace, cg.BitmapAlphaPremultipliedFirst, dataProvider, nil, false, cg.RenderingIntentDefault); imgRef != 0 {
		return imgRef, nil
	}
	return 0, errs.New(errUnableToCreateImage)
}

func (img *imageRef) osNewSubImage(x, y, width, height int) (osImage, error) {
	if imgRef := img.osImg.CreateWithImageInRect(float64(x), float64(y), float64(width), float64(height)); imgRef != 0 {
		return imgRef, nil
	}
	return 0, errs.New(errUnableToCreateImage)
}

func (img *imageRef) osNewScaledImage(width, height int, q quality.Quality) (osImage, error) {
	colorspace := cg.ColorSpaceCreateSRGB()
	defer colorspace.Release()
	if ctx := cg.BitmapContextCreate(nil, width, height, 8, 0, colorspace, cg.BitmapAlphaPremultipliedFirst); ctx != 0 {
		defer ctx.Release()
		ctx.SetInterpolationQuality(cg.InterpolationQuality(q))
		ctx.DrawImage(0, 0, float64(width), float64(height), img.osImg)
		if imgRef := ctx.BitmapContextCreateImage(); imgRef != 0 {
			return imgRef, nil
		}
	}
	return 0, errs.New(errUnableToCreateImage)
}

func (img *imageRef) osIsValid() bool {
	return img.osImg != 0
}

func (img *imageRef) osImagePixels(pixels []Color) {
	colorspace := cg.ColorSpaceCreateSRGB()
	if ctx := cg.BitmapContextCreate(unsafe.Pointer(&pixels[0]), img.width, img.height, 8, img.width*4, colorspace, cg.BitmapAlphaPremultipliedFirst); ctx != 0 { //nolint:gosec
		ctx.SetInterpolationQuality(cg.InterpolationQualityNone)
		ctx.DrawImage(0, 0, float64(img.width), float64(img.height), img.osImg)
		ctx.Release()
	}
	colorspace.Release()
	for i := range pixels {
		pixels[i] = pixels[i].Unpremultiply()
	}
}

func (img *imageRef) osDrawInRect(gc Context, rect geom.Rect) {
	ic := imgContext{
		osImg:     img.osImg,
		osContext: gc.OSContext(),
	}
	layer, ok := imageLayerCache[ic]
	if !ok {
		w := float64(img.width)
		h := float64(img.height)
		layer = cg.LayerCreateWithContext(gc.OSContext(), w, h)
		layer.Retain()
		ctx := layer.Context()
		ctx.DrawImage(0, 0, w, h, img.osImg)
		imageLayerCache[ic] = layer
		var contexts []OSContext
		if contexts, ok = imageToContextCache[ic.osImg]; !ok {
			imageToContextCache[ic.osImg] = []OSContext{ic.osContext}
		} else {
			contexts = append(contexts, ic.osContext)
			imageToContextCache[ic.osImg] = contexts
		}
		var images []osImage
		if images, ok = contextToImageCache[ic.osContext]; !ok {
			contextToImageCache[ic.osContext] = []osImage{ic.osImg}
		} else {
			images = append(images, ic.osImg)
			contextToImageCache[ic.osContext] = images
		}
	}
	gc.Save()
	gc.Translate(0, rect.Y+rect.Height)
	gc.Scale(1, -1)
	gc.OSContext().DrawLayer(rect.X, 0, rect.Width, rect.Height, layer)
	gc.Restore()
}

func (img *imageRef) osDispose() {
	if img.osImg != 0 {
		if contexts, ok := imageToContextCache[img.osImg]; ok {
			delete(imageToContextCache, img.osImg)
			for _, ctx := range contexts {
				var images []osImage
				if images, ok = contextToImageCache[ctx]; ok {
					for i := range images {
						if img.osImg == images[i] {
							if len(images) == 1 {
								delete(contextToImageCache, ctx)
							} else {
								images[i] = images[len(images)-1]
								images[len(images)-1] = 0
								images = images[:len(images)-1]
								contextToImageCache[ctx] = images
							}
							break
						}
					}
				}
				ic := imgContext{
					osImg:     img.osImg,
					osContext: ctx,
				}
				var layer cg.Layer
				if layer, ok = imageLayerCache[ic]; ok {
					layer.Release()
					delete(imageLayerCache, ic)
				}
			}
		}
		img.osImg.Release()
		img.osImg = 0
	}
}

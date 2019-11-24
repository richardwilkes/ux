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

var imageLayerCache = make(map[cg.Image]cg.Layer)

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
	colorspace := cg.ColorSpaceCreateDeviceRGB()
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
	colorspace := cg.ColorSpaceCreateDeviceRGB()
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
	colorspace := cg.ColorSpaceCreateDeviceRGB()
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
	layer, ok := imageLayerCache[img.osImg]
	if !ok {
		w := float64(img.width)
		h := float64(img.height)
		layer = cg.LayerCreateWithContext(gc.OSContext(), w, h)
		layer.Retain()
		ctx := layer.Context()
		ctx.DrawImage(0, 0, w, h, img.osImg)
		imageLayerCache[img.osImg] = layer
	}
	gc.Save()
	gc.Translate(0, rect.Y+rect.Height)
	gc.Scale(1, -1)
	gc.OSContext().DrawLayer(rect.X, 0, rect.Width, rect.Height, layer)
	gc.Restore()
}

func (img *imageRef) osDispose() {
	if img.osImg != 0 {
		if layer, ok := imageLayerCache[img.osImg]; ok {
			layer.Release()
			delete(imageLayerCache, img.osImg)
		}
		img.osImg.Release()
		img.osImg = 0
	}
}

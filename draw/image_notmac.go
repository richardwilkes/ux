// +build !darwin

package draw

func (img *imageRef) osNewSubImage(x, y, width, height int) (osImage, error) {
	origPixels := make([]Color, width*height)
	img.osImagePixels(origPixels)
	xe := x + width
	ye := y + height
	pixels := make([]Color, width*height)
	for yy := y; yy < ye; yy++ {
		for xx := x; xx < xe; xx++ {
			pixels[(yy-y)*width+(xx-x)] = origPixels[yy*img.width+xx]
		}
	}
	return osNewImageFromData(&ImageData{
		Pixels: pixels,
		Width:  width,
		Height: height,
		Scale:  img.scale,
	})
}

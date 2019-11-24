package draw

func (c Color) osPrepareForFill(gc Context) {
	gc.OSContext().SetRGBFillColor(c.RedIntensity(), c.GreenIntensity(), c.BlueIntensity(), c.AlphaIntensity())
}

func (c Color) osFill(gc Context) {
	c.osPrepareForFill(gc)
	gc.OSContext().FillPath()
}

func (c Color) osFillEvenOdd(gc Context) {
	c.osPrepareForFill(gc)
	gc.OSContext().EOFillPath()
}

func (c Color) osStroke(gc Context) {
	g := gc.OSContext()
	g.SetRGBStrokeColor(c.RedIntensity(), c.GreenIntensity(), c.BlueIntensity(), c.AlphaIntensity())
	g.StrokePath()
}

package draw

// Ink holds a color/pattern/gradient to draw with.
type Ink interface {
	osPrepareForFill(gc Context)
	osFill(gc Context)
	osFillEvenOdd(gc Context)
	osStroke(gc Context)
}

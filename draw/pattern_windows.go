package draw

type osPattern = int

func osNewPattern(img *Image) osPattern {
	return 0 // RAW: Implement
}

func (p *Pattern) osPrepareForFill(gc Context) {
	// RAW: Implement
}

func (p *Pattern) osFill(gc Context) {
	// RAW: Implement
}

func (p *Pattern) osFillEvenOdd(gc Context) {
	// RAW: Implement
}

func (p *Pattern) osStroke(gc Context) {
	// RAW: Implement
}

func (r *patternRef) osIsValid() bool {
	return false // RAW: Implement
}

func (r *patternRef) osDispose() {
	// RAW: Implement
}

package linejoin

// LineJoin specifies how to draw the junction between connected line
// segments.
type LineJoin int

const (
	// Miter is a join with a sharp, angled, corner. The line is drawn with
	// the outer sides of the line beyond the endpoint of the path, until they
	// meet. If the length of the miter divided by the line width is greater
	// than the miter limit, a bevel join is used instead. This is the
	// default.
	Miter LineJoin = iota
	// Round is a join with a rounded end. The line is drawn to extend beyond
	// the endpoint of the path. The line ends with a semicircular arc with a
	// radius of 1/2 the line’s width, centered on the endpoint.
	Round
	// Bevel is a join with a squared-off end. The line is drawn to extend
	// beyond the endpoint of the path, for a distance of 1/2 the line’s
	// width.
	Bevel
)

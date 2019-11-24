package linecap

// LineCap defines styles for rendering the endpoint of a stroked line.
type LineCap int

const (
	// Butt is a line with a squared-off end. Lines are drawn to extend only
	// to the exact endpoint of the path. This is the default.
	Butt LineCap = iota
	// Round is a line with a rounded end. Lines are drawn to extend beyond
	// the endpoint of the path. The line ends with a semicircular arc with a
	// radius of half the line's width, centered on the endpoint.
	Round
	// Square is a line with a squared-off end. Lines are drawn beyond the
	// endpoint of the path for a distance equal to half the line width.
	Square
)

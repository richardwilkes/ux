// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

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

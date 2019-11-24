package draw

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/draw/linecap"
	"github.com/richardwilkes/ux/draw/linejoin"
	"github.com/richardwilkes/ux/draw/quality"
)

// Context defines the methods required of all drawing contexts.
type Context interface {
	Pather
	// OSContext returns the underlying OS graphics context.
	OSContext() OSContext
	// Save pushes a copy of the current graphics state onto the stack for the
	// context. The current path is not saved as part of this call.
	Save()
	// Restore sets the current graphics state to the state most recently
	// saved by a call to Save().
	Restore()
	// SetOpacity sets the opacity level used when drawing. Values can range
	// from 0 (transparent) to 1 (opaque).
	SetOpacity(opacity float64)
	// SetPatternOffset sets the offset to start the ink pattern at.
	SetPatternOffset(x, y float64)
	// SetStrokeWidth sets the current stroke width.
	SetStrokeWidth(width float64)
	// SetLineCap sets the current line cap style.
	SetLineCap(lineCap linecap.LineCap)
	// SetLineJoin sets the current line join style.
	SetLineJoin(lineJoin linejoin.LineJoin)
	// SetMiterLimit sets the current miter limit for joins of connected
	// lines. Only relevant when LineJoinMiter is used.
	SetMiterLimit(limit float64)
	// SetLineDash sets the pattern for dashed lines. The phase specifies how
	// far into the dash pattern the line starts. The segments specify the
	// length of painted and unpainted segments of the dash pattern.
	SetLineDash(phase float64, segments ...float64)
	// SetInterpolationQualityHint sets the interpolation quality hint for
	// image rendering.
	SetInterpolationQualityHint(q quality.Quality)
	// Translate the coordinate system.
	Translate(x, y float64)
	// Scale the coordinate system.
	Scale(x, y float64)
	// Rotate the coordinate system.
	Rotate(angleInRadians float64)
	// Fill the current path using the non-zero winding rule, then clear the
	// current path state. https://en.wikipedia.org/wiki/Nonzero-rule
	Fill(ink Ink)
	// FillEvenOdd the current path using the even-odd rule, then clear the
	// current path state. https://en.wikipedia.org/wiki/Even-odd_rule
	FillEvenOdd(ink Ink)
	// Stroke the current path, then clear the current path state.
	Stroke(ink Ink)
	// GetClipRect returns the rectangle that encompasses the current clipping
	// path.
	GetClipRect() geom.Rect
	// Clip sets the clipping path to the intersection of the current clipping
	// path and the current path using the non-zero winding rule, then clears
	// the current path state. https://en.wikipedia.org/wiki/Nonzero-rule
	Clip()
	// ClipEvenOdd sets the clipping path to the intersection of the current
	// clipping path and the current path using the even-odd rule, then clears
	// the current path state. https://en.wikipedia.org/wiki/Even-odd_rule
	ClipEvenOdd()
	// DrawString draws a string at the specified location using the specified
	// font and ink.
	DrawString(x, y float64, font *Font, ink Ink, str string)
	// Dispose of the context, releasing any OS resources associated with it.
	Dispose()
}

// Initialize is not intended for client code to call.
func Initialize() {
	osInitSystemCursors()
	osInitSystemFonts()
}

// NewContextForOSContext creates a new graphics context for the given OS
// graphic context.
func NewContextForOSContext(gc OSContext) Context {
	return osNewContextForOSContext(gc)
}

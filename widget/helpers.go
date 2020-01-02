// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package widget

import (
	"math"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/layout/align"
	"github.com/richardwilkes/ux/layout/side"
)

// DrawRectBase fills and strokes a rectangle.
func DrawRectBase(gc draw.Context, rect geom.Rect, fillInk, strokeInk draw.Ink) {
	gc.Rect(rect)
	gc.Fill(fillInk)
	rect.InsetUniform(0.5)
	gc.Rect(rect)
	gc.Stroke(strokeInk)
}

// DrawRoundedRectBase fills and strokes a rounded rectangle.
func DrawRoundedRectBase(gc draw.Context, rect geom.Rect, cornerRadius float64, fillInk, strokeInk draw.Ink) {
	gc.RoundedRect(rect, cornerRadius)
	gc.Fill(fillInk)
	rect.InsetUniform(0.5)
	gc.RoundedRect(rect, math.Max(cornerRadius-0.5, 0))
	gc.Stroke(strokeInk)
}

// DrawEllipseBase fills and strokes an ellipse.
func DrawEllipseBase(gc draw.Context, rect geom.Rect, fillInk, strokeInk draw.Ink) {
	gc.Ellipse(rect)
	gc.Fill(fillInk)
	rect.InsetUniform(0.5)
	gc.Ellipse(rect)
	gc.Stroke(strokeInk)
}

// LabelSize returns the preferred size of a label. Provided as a standalone
// function so that other types of panels can make use of it.
func LabelSize(text string, font *draw.Font, image *draw.Image, imgSide side.Side, imgGap float64) geom.Size {
	var size geom.Size
	if text != "" {
		size = font.Extents(text)
		size.GrowToInteger()
	}
	adjustLabelSizeForImage(text, image, imgSide, imgGap, &size)
	size.GrowToInteger()
	return size
}

func adjustLabelSizeForImage(text string, image *draw.Image, imgSide side.Side, imgGap float64, size *geom.Size) {
	if image != nil {
		logicalSize := image.LogicalGeomSize()
		switch {
		case text == "":
			*size = logicalSize
		case imgSide.Horizontal():
			size.Width += logicalSize.Width + imgGap
			if size.Height < logicalSize.Height {
				size.Height = logicalSize.Height
			}
		default:
			size.Height += logicalSize.Height + imgGap
			if size.Width < logicalSize.Width {
				size.Width = logicalSize.Width
			}
		}
	}
}

// DrawLabel draws a label. Provided as a standalone function so that other
// types of panels can make use of it.
func DrawLabel(gc draw.Context, rect geom.Rect, hAlign, vAlign align.Alignment, text string, font *draw.Font, textInk draw.Ink, image *draw.Image, imgSide side.Side, imgGap float64, enabled bool) {
	// Determine overall size of content
	var size, txtSize geom.Size
	if text != "" {
		txtSize = font.Extents(text)
		size = txtSize
	}
	adjustLabelSizeForImage(text, image, imgSide, imgGap, &size)

	// Adjust the working area for the content size
	switch hAlign {
	case align.Middle, align.Fill:
		rect.X = math.Floor(rect.X + (rect.Width-size.Width)/2)
	case align.End:
		rect.X += rect.Width - size.Width
	default: // Start
	}
	switch vAlign {
	case align.Middle, align.Fill:
		rect.Y = math.Floor(rect.Y + (rect.Height-size.Height)/2)
	case align.End:
		rect.Y += rect.Height - size.Height
	default: // Start
	}
	rect.Size = size

	// Determine image and text areas
	imgX := rect.X
	imgY := rect.Y
	txtX := rect.X
	txtY := rect.Y
	if text != "" && image != nil {
		logicalSize := image.LogicalGeomSize()
		switch imgSide {
		case side.Top:
			txtY += logicalSize.Height + imgGap
			if logicalSize.Width > txtSize.Width {
				txtX = math.Floor(txtX + (logicalSize.Width-txtSize.Width)/2)
			} else {
				imgX = math.Floor(imgX + (txtSize.Width-logicalSize.Width)/2)
			}
		case side.Left:
			txtX += logicalSize.Width + imgGap
			if logicalSize.Height > txtSize.Height {
				txtY = math.Floor(txtY + (logicalSize.Height-txtSize.Height)/2)
			} else {
				imgY = math.Floor(imgY + (txtSize.Height-logicalSize.Height)/2)
			}
		case side.Bottom:
			imgY += rect.Height - logicalSize.Height
			txtY = imgY - (imgGap + txtSize.Height)
			if logicalSize.Width > txtSize.Width {
				txtX = math.Floor(txtX + (logicalSize.Width-txtSize.Width)/2)
			} else {
				imgX = math.Floor(imgX + (txtSize.Width-logicalSize.Width)/2)
			}
		case side.Right:
			imgX += rect.Width - logicalSize.Width
			txtX = imgX - (imgGap + txtSize.Width)
			if logicalSize.Height > txtSize.Height {
				txtY = math.Floor(txtY + (logicalSize.Height-txtSize.Height)/2)
			} else {
				imgY = math.Floor(imgY + (txtSize.Height-logicalSize.Height)/2)
			}
		}
	}

	if !enabled {
		gc.SetOpacity(0.33)
	}

	// Draw the image
	if image != nil {
		rect.X = imgX
		rect.Y = imgY
		rect.Size = image.LogicalGeomSize()
		image.DrawInRect(gc, rect)
	}

	// Draw the text
	if text != "" {
		gc.DrawString(txtX, txtY, font, textInk, text)
	}
}

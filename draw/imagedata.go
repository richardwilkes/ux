// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package draw

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"image"
	"image/color"
	"io"

	"github.com/richardwilkes/toolbox/errs"
)

// ImageData is the raw information that makes up an Image.
type ImageData struct {
	Pixels []Color // The pixel data
	Width  int     // The width of the image, in pixels
	Height int     // The height of the image, in pixels
	Scale  float64 // The scale to apply to the image size to obtain the device-independent dimensions
}

// LogicalWidth returns the logical (device-independent) width.
func (imgData *ImageData) LogicalWidth() int {
	return int(float64(imgData.Width) * imgData.Scale)
}

// LogicalHeight returns the logical (device-independent) height.
func (imgData *ImageData) LogicalHeight() int {
	return int(float64(imgData.Height) * imgData.Scale)
}

// LogicalSize returns the logical (device-independent) size.
func (imgData *ImageData) LogicalSize() (width, height int) {
	return imgData.LogicalWidth(), imgData.LogicalHeight()
}

// Validate checks the ImageData and reports any errors found.
func (imgData *ImageData) Validate() error {
	if imgData == nil {
		return errs.New("image data may not be nil")
	}
	if imgData.Width < 1 {
		return errs.New("image data has invalid width")
	}
	if imgData.Height < 1 {
		return errs.New("image data has invalid height")
	}
	if imgData.LogicalWidth() < 1 || imgData.LogicalHeight() < 1 {
		return errs.New("image data scale results in zero-sized image")
	}
	if len(imgData.Pixels) != imgData.Width*imgData.Height {
		return errs.New("image data size does not match pixel data provided")
	}
	return nil
}

type jsonImageData struct {
	Pixels string  `json:"pixels"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Scale  float64 `json:"scale"`
}

// MarshalJSON implements json.Marshaler.
func (imgData *ImageData) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer
	w := gzip.NewWriter(&buffer)
	pixelBuffer := make([]byte, 4)
	for _, pixel := range imgData.Pixels {
		binary.LittleEndian.PutUint32(pixelBuffer, uint32(pixel))
		if _, err := w.Write(pixelBuffer); err != nil {
			return nil, err
		}
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return json.Marshal(&jsonImageData{
		Pixels: base64.RawStdEncoding.EncodeToString(buffer.Bytes()),
		Width:  imgData.Width,
		Height: imgData.Height,
		Scale:  imgData.Scale,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (imgData *ImageData) UnmarshalJSON(data []byte) error {
	var tmp jsonImageData
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	pixels, err := base64.RawStdEncoding.DecodeString(tmp.Pixels)
	if err != nil {
		return err
	}
	r, err := gzip.NewReader(bytes.NewBuffer(pixels))
	if err != nil {
		return err
	}
	imgData.Pixels = make([]Color, tmp.Width*tmp.Height)
	pixelBuffer := make([]byte, 4)
	for i := range imgData.Pixels {
		var n int
		if n, err = r.Read(pixelBuffer); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if n != 4 {
			return errs.New("not enough pixel data")
		}
		imgData.Pixels[i] = Color(binary.LittleEndian.Uint32(pixelBuffer))
	}
	imgData.Width = tmp.Width
	imgData.Height = tmp.Height
	imgData.Scale = tmp.Scale
	return imgData.Validate()
}

// ColorModel implements image.Image.
func (imgData *ImageData) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds implements image.Image.
func (imgData *ImageData) Bounds() image.Rectangle {
	return image.Rect(0, 0, imgData.Width, imgData.Height)
}

// At implements image.Image.
func (imgData *ImageData) At(x, y int) color.Color {
	return imgData.Pixels[y*imgData.Width+x]
}

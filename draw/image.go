package draw

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/softref"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/draw/quality"
	"github.com/richardwilkes/ux/globals"
)

// Image file extensions.
const (
	PNGExt  = ".png"
	JPGExt  = ".jpg"
	JPEGExt = ".jpeg"
	GIFExt  = ".gif"
)

const (
	errUnableToCreateImage = "unable to create image"
	errInvalidImage        = "invalid image"
	imgURLKey              = "imgurl"
	revisionLatest         = "/revision/latest"
	fileScheme             = "file"
	httpScheme             = "http"
	httpsScheme            = "https"
)

// Image holds a reference to an image.
type Image softref.SoftRef

// NewImageFromBytes creates a new image from raw bytes.
func NewImageFromBytes(buffer []byte, scale float64) (*Image, error) {
	if len(buffer) < 1 {
		return nil, errs.New("no data in input buffer")
	}
	img, width, height, err := osNewImageFromBytes(buffer)
	if err != nil {
		return nil, err
	}
	return newImage(width, height, scale, img)
}

// NewImageFromData creates a new image from the ImageData.
func NewImageFromData(data *ImageData) (*Image, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}
	img, err := osNewImageFromData(data)
	if err != nil {
		return nil, err
	}
	return newImage(data.Width, data.Height, data.Scale, img)
}

// NewImageFromURL creates a new image from data retrieved from the URL. The
// http.DefaultClient will be used if the data URL is remote.
func NewImageFromURL(urlStr string, scale float64) (*Image, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, errs.NewWithCause(urlStr, err)
	}
	var data []byte
	switch u.Scheme {
	case fileScheme:
		if data, err = ioutil.ReadFile(u.Path); err != nil {
			return nil, errs.NewWithCause(urlStr, err)
		}
	case httpScheme, httpsScheme:
		if data, err = retrieveDataFromURL(urlStr); err != nil {
			return nil, errs.NewWithCause(urlStr, err)
		}
	default:
		return nil, errs.Newf("invalid url: %s", urlStr)
	}
	return NewImageFromBytes(data, scale)
}

func (r *Image) osImage() osImage {
	return r.Resource.(*imageRef).osImg
}

// IsValid returns true if the image is still valid (i.e. hasn't been
// disposed).
func (r *Image) IsValid() bool {
	return r.Resource.(*imageRef).osIsValid()
}

// NewSubImage creates a new image from a portion of this image.
func (r *Image) NewSubImage(x, y, width, height int) (*Image, error) {
	if !r.IsValid() {
		return nil, errs.New(errInvalidImage)
	}
	img := r.Resource.(*imageRef) //nolint:errcheck
	if x < 0 || width < 1 || y < 0 || height < 1 || x+width > img.width || y+height > img.height {
		return nil, errs.New("rect must be within the image")
	}
	sub, err := img.osNewSubImage(x, y, width, height)
	if err != nil {
		return nil, err
	}
	return newImage(width, height, img.scale, sub)
}

// NewScaledImage creates a new image by scaling this image.
func (r *Image) NewScaledImage(width, height int, q quality.Quality) (*Image, error) {
	if !r.IsValid() {
		return nil, errs.New(errInvalidImage)
	}
	if width < 1 || height < 1 {
		return nil, errs.New("invalid size")
	}
	img := r.Resource.(*imageRef) //nolint:errcheck
	scaled, err := img.osNewScaledImage(width, height, q)
	if err != nil {
		return nil, err
	}
	return newImage(width, height, img.scale, scaled)
}

// LogicalWidth returns the logical (device-independent) width.
func (r *Image) LogicalWidth() int {
	img := r.Resource.(*imageRef) //nolint:errcheck
	return int(float64(img.width) * img.scale)
}

// LogicalHeight returns the logical (device-independent) height.
func (r *Image) LogicalHeight() int {
	img := r.Resource.(*imageRef) //nolint:errcheck
	return int(float64(img.height) * img.scale)
}

// LogicalSize returns the logical (device-independent) size.
func (r *Image) LogicalSize() (width, height int) {
	return r.LogicalWidth(), r.LogicalHeight()
}

// LogicalGeomSize returns the logical (device-independent) size.
func (r *Image) LogicalGeomSize() geom.Size {
	return geom.Size{
		Width:  float64(r.LogicalWidth()),
		Height: float64(r.LogicalHeight()),
	}
}

// PixelSize returns the pixel size of the image, i.e. the actual size of the
// image determined by the raw pixels.
func (r *Image) PixelSize() (width, height int) {
	img := r.Resource.(*imageRef) //nolint:errcheck
	return img.width, img.height
}

// Data extracts the raw pixel data.
func (r *Image) Data() *ImageData {
	img := r.Resource.(*imageRef) //nolint:errcheck
	pixels := make([]Color, img.width*img.height)
	if r.IsValid() {
		img.osImagePixels(pixels)
	}
	return &ImageData{
		Pixels: pixels,
		Width:  img.width,
		Height: img.height,
		Scale:  img.scale,
	}
}

// Draw the image at the specified location using its logical size.
func (r *Image) Draw(gc Context, where geom.Point) {
	r.DrawInRect(gc, geom.Rect{Point: where, Size: r.LogicalGeomSize()})
}

// DrawInRect draws the image into the area specified by the rect, scaling if
// necessary.
func (r *Image) DrawInRect(gc Context, rect geom.Rect) {
	if r.IsValid() {
		r.Resource.(*imageRef).osDrawInRect(gc, rect)
	}
}

type imageRef struct {
	key    string
	width  int
	height int
	scale  float64
	osImg  osImage
}

func newImage(width, height int, scale float64, img osImage) (*Image, error) {
	imgRef := &imageRef{
		width:  width,
		height: height,
		scale:  scale,
		osImg:  img,
	}
	pixels := make([]Color, width*height)
	imgRef.osImagePixels(pixels)
	s := sha256.New224()
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, math.Float64bits(scale))
	if _, err := s.Write(buffer); err != nil {
		return nil, errs.Wrap(err)
	}
	buffer = buffer[:4]
	for _, pixel := range pixels {
		binary.LittleEndian.PutUint32(buffer, uint32(pixel))
		if _, err := s.Write(buffer); err != nil {
			return nil, errs.Wrap(err)
		}
	}
	imgRef.key = base64.RawURLEncoding.EncodeToString(s.Sum(nil)[:sha256.Size224])
	ref, existedPreviously := globals.Pool.NewSoftRef(imgRef)
	if existedPreviously {
		imgRef.Release()
	}
	return (*Image)(ref), nil
}

func (r *imageRef) Key() string {
	return r.key
}

func (r *imageRef) Release() {
	r.osDispose()
}

// DistillImageURL distills a URL string into either a URL string that likely
// has an image we can use, or an empty string.
func DistillImageURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		jot.Warn(err)
		return ""
	}
	switch u.Scheme {
	case fileScheme:
		if HasImageSuffix(urlStr) {
			return urlStr
		}
	case httpScheme, httpsScheme:
		if HasImageSuffix(u.Path) {
			return urlStr
		}
		if alt, ok := u.Query()[imgURLKey]; ok && len(alt) > 0 {
			return DistillImageURL(alt[0])
		}
		if strings.HasSuffix(u.Path, revisionLatest) {
			u.RawPath = ""
			u.Path = u.Path[:len(u.Path)-len(revisionLatest)]
			return DistillImageURL(u.String())
		}
	default:
		jot.Warnf("unhandled url scheme: %s", urlStr)
	}
	return ""
}

// HasImageSuffix returns true if the string has one of the image suffixes
// that we can process.
func HasImageSuffix(p string) bool {
	p = strings.ToLower(p)
	return strings.HasSuffix(p, PNGExt) || strings.HasSuffix(p, JPGExt) ||
		strings.HasSuffix(p, JPEGExt) || strings.HasSuffix(p, GIFExt)
}

func retrieveDataFromURL(urlStr string) ([]byte, error) {
	rsp, err := http.Get(urlStr) //nolint:gosec
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer xio.CloseIgnoringErrors(rsp.Body)
	if rsp.StatusCode < 200 || rsp.StatusCode > 299 {
		return nil, errs.Newf("received status %d (%s)", rsp.StatusCode, rsp.Status)
	}
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return data, nil
}

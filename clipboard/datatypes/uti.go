package datatypes

// Well-known data types. The UTIs should always be unique.
var (
	None       = DataType{}                                                           // The empty type; not a valid type
	Generic    = DataType{UTI: "public.item"}                                         // Generic base type
	Data       = DataType{UTI: "public.data", Mime: "application/octet-stream"}       // Bytes
	URL        = DataType{UTI: "public.url", Mime: "application/octet-stream"}        // URL
	FileURL    = DataType{UTI: "public.file-url", Mime: "application/octet-stream"}   // File URL
	PlainText  = DataType{UTI: "public.utf8-plain-text", Mime: "text/plain"}          // UTF-8 encoded text
	RTFText    = DataType{UTI: "public.rtf", Mime: "text/rtf"}                        // Rich text
	HTMLText   = DataType{UTI: "public.html", Mime: "text/html"}                      // HTML text
	XMLText    = DataType{UTI: "public.xml", Mime: "application/xml"}                 // XML text
	JPEG       = DataType{UTI: "public.jpeg", Mime: "image/jpeg"}                     // JPEG image
	TIFF       = DataType{UTI: "public.tiff", Mime: "image/tiff"}                     // TIFF image
	PNG        = DataType{UTI: "public.png", Mime: "image/png"}                       // PNG image
	XBM        = DataType{UTI: "public.xbitmap-image", Mime: "image/x-xbitmap"}       // X bitmap image
	BMP        = DataType{UTI: "com.microsoft.bmp", Mime: "image/bmp"}                // BMP image
	ICO        = DataType{UTI: "com.microsoft.ico", Mime: "image/vnd.microsoft.icon"} // ICO image
	GIF        = DataType{UTI: "com.compuserve.gif", Mime: "image/gif"}               // GIF image
	MP3        = DataType{UTI: "public.mp3", Mime: "audio/mpeg"}                      // MPEG-3 audio
	MPEG4Audio = DataType{UTI: "public.mpeg-4-audio", Mime: "audio/mp4"}              // MPEG-4 audio
	AIFF       = DataType{UTI: "public.aiff-audio", Mime: "audio/aiff"}               // AIFF audio
	AVI        = DataType{UTI: "public.avi", Mime: "video/avi"}                       // AVI movie
	MPEG       = DataType{UTI: "public.mpeg", Mime: "video/mpeg"}                     // MPEG-1 or MPEG-2 movie
	MPEG4      = DataType{UTI: "public.mpeg-4", Mime: "video/mp4"}                    // MPEG-4 movie
	PDF        = DataType{UTI: "com.adobe.pdf", Mime: "application/pdf"}              // PDF
	ByUTI      = make(map[string]DataType)
	ByMime     = make(map[string]DataType)
)

// DataType holds Uniform Type Identifiers and their corresponding Mime Types
type DataType struct {
	UTI  string
	Mime string
}

func init() {
	for _, one := range []DataType{
		Generic, Data, URL, FileURL, PlainText, RTFText, HTMLText, XMLText,
		JPEG, TIFF, PNG, XBM, BMP, ICO, GIF, MP3, MPEG4Audio, AIFF, AVI, MPEG,
		MPEG4, PDF,
	} {
		ByUTI[one.UTI] = one
		// The MimeTypes may not be unique, so the first entry will be the one
		// found by ByMime[].
		if _, exists := ByMime[one.Mime]; !exists {
			ByMime[one.Mime] = one
		}
	}
}

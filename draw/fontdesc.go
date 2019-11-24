package draw

import (
	"bytes"
	"fmt"
)

// FontDescriptor holds information necessary to construct a Font.
type FontDescriptor struct {
	Family string  // Family name of the Font, such as "Times New Roman".
	Size   float64 // Size of the font, in points (1/72 of an inch).
	Bold   bool
	Italic bool
}

// String implements the fmt.Stringer interface.
func (d *FontDescriptor) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(d.Family)
	fmt.Fprintf(&buffer, " %v", d.Size)
	if d.Bold {
		buffer.WriteString(" bold")
	}
	if d.Italic {
		buffer.WriteString(" italic")
	}
	return buffer.String()
}

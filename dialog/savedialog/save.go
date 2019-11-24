package savedialog

import "strings"

// SaveDialog represents a system save dialog.
type SaveDialog struct {
	dialog osSaveDialog
}

// New creates a new save dialog.
func New() *SaveDialog {
	return &SaveDialog{dialog: osNew()}
}

// DirectoryURL returns a URL pointing to the directory the dialog will open
// up in.
func (d *SaveDialog) DirectoryURL() string {
	return d.osDirectoryURL()
}

// SetDirectoryURL sets the directory the dialog will open up in.
func (d *SaveDialog) SetDirectoryURL(dirURL string) {
	d.osSetDirectoryURL(dirURL)
}

// AllowedFileTypes returns the set of permitted file types. nil will be
// returned if all files are allowed.
func (d *SaveDialog) AllowedFileTypes() []string {
	return d.osAllowedFileTypes()
}

// SetAllowedFileTypes sets the permitted file types that may be selected for
// saving. Pass in nil to allow all files.
func (d *SaveDialog) SetAllowedFileTypes(allowedExtensions []string) {
	var actual []string
	for _, ext := range allowedExtensions {
		for strings.HasPrefix(ext, ".") {
			ext = ext[1:]
		}
		if ext != "" {
			actual = append(actual, ext)
		}
	}
	d.osSetAllowedFileTypes(actual)
}

// URL returns the URL that was chosen.
func (d *SaveDialog) URL() string {
	return d.osURL()
}

// RunModal displays the dialog, allowing the user to make a selection.
// Returns true if successful or false if canceled.
func (d *SaveDialog) RunModal() bool {
	return d.osRunModal()
}

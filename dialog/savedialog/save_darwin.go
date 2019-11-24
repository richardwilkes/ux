package savedialog

import (
	"github.com/richardwilkes/macos/ns"
)

type osSaveDialog = *ns.SavePanel

func osNew() osSaveDialog {
	return ns.NewSavePanel()
}

func (d *SaveDialog) osDirectoryURL() string {
	return d.dialog.DirectoryURL()
}

func (d *SaveDialog) osSetDirectoryURL(dirURL string) {
	d.dialog.SetDirectoryURL(dirURL)
}

func (d *SaveDialog) osAllowedFileTypes() []string {
	return d.dialog.AllowedFileTypes()
}

func (d *SaveDialog) osSetAllowedFileTypes(allowed []string) {
	d.dialog.SetAllowedFileTypes(allowed)
}

func (d *SaveDialog) osURL() string {
	return d.dialog.URL()
}

func (d *SaveDialog) osRunModal() bool {
	return d.dialog.RunModal()
}

package opendialog

import (
	"github.com/richardwilkes/macos/ns"
)

type osOpenDialog = *ns.OpenPanel

func osNew() osOpenDialog {
	return ns.NewOpenPanel()
}

func (d *OpenDialog) osDirectoryURL() string {
	return d.dialog.DirectoryURL()
}

func (d *OpenDialog) osSetDirectoryURL(dirURL string) {
	d.dialog.SetDirectoryURL(dirURL)
}

func (d *OpenDialog) osAllowedFileTypes() []string {
	return d.dialog.AllowedFileTypes()
}

func (d *OpenDialog) osSetAllowedFileTypes(allowed []string) {
	d.dialog.SetAllowedFileTypes(allowed)
}

func (d *OpenDialog) osCanChooseFiles() bool {
	return d.dialog.CanChooseFiles()
}

func (d *OpenDialog) osSetCanChooseFiles(canChoose bool) {
	d.dialog.SetCanChooseFiles(canChoose)
}

func (d *OpenDialog) osCanChooseDirectories() bool {
	return d.dialog.CanChooseDirectories()
}

func (d *OpenDialog) osSetCanChooseDirectories(canChoose bool) {
	d.dialog.SetCanChooseDirectories(canChoose)
}

func (d *OpenDialog) osResolvesAliases() bool {
	return d.dialog.ResolvesAliases()
}

func (d *OpenDialog) osSetResolvesAliases(resolves bool) {
	d.dialog.SetResolvesAliases(resolves)
}

func (d *OpenDialog) osAllowsMultipleSelection() bool {
	return d.dialog.AllowsMultipleSelection()
}

func (d *OpenDialog) osSetAllowsMultipleSelection(allow bool) {
	d.dialog.SetAllowsMultipleSelection(allow)
}

func (d *OpenDialog) osURLs() []string {
	return d.dialog.URLs()
}

func (d *OpenDialog) osRunModal() bool {
	return d.dialog.RunModal()
}

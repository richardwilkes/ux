package opendialog

import "strings"

// OpenDialog represents a system open dialog.
type OpenDialog struct {
	dialog osOpenDialog
}

// New creates a new open dialog.
func New() *OpenDialog {
	return &OpenDialog{dialog: osNew()}
}

// DirectoryURL returns a URL pointing to the directory the dialog will open
// up in.
func (d *OpenDialog) DirectoryURL() string {
	return d.osDirectoryURL()
}

// SetDirectoryURL sets the directory the dialog will open up in.
func (d *OpenDialog) SetDirectoryURL(dirURL string) {
	d.osSetDirectoryURL(dirURL)
}

// AllowedFileTypes returns the set of permitted file types. nil will be
// returned if all files are allowed.
func (d *OpenDialog) AllowedFileTypes() []string {
	return d.osAllowedFileTypes()
}

// SetAllowedFileTypes sets the permitted file types that may be selected for
// opening. Pass in nil to allow all files.
func (d *OpenDialog) SetAllowedFileTypes(allowedExtensions []string) {
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

// CanChooseFiles returns true if the open dialog is permitted to select
// files.
func (d *OpenDialog) CanChooseFiles() bool {
	return d.osCanChooseFiles()
}

// SetCanChooseFiles sets whether the open dialog is permitted to select
// files.
func (d *OpenDialog) SetCanChooseFiles(canChoose bool) {
	d.osSetCanChooseFiles(canChoose)
}

// CanChooseDirectories returns true if the open dialog is permitted to select
// directories.
func (d *OpenDialog) CanChooseDirectories() bool {
	return d.osCanChooseDirectories()
}

// SetCanChooseDirectories sets whether the open dialog is permitted to select
// directories.
func (d *OpenDialog) SetCanChooseDirectories(canChoose bool) {
	d.osSetCanChooseDirectories(canChoose)
}

// ResolvesAliases returns whether the returned URLs have been resolved in the
// case where the selection was an alias.
func (d *OpenDialog) ResolvesAliases() bool {
	return d.osResolvesAliases()
}

// SetResolvesAliases sets whether the returned URLs will be resolved in the
// case where the selection was an alias.
func (d *OpenDialog) SetResolvesAliases(resolves bool) {
	d.osSetResolvesAliases(resolves)
}

// AllowsMultipleSelection returns true if more than one item can be selected.
func (d *OpenDialog) AllowsMultipleSelection() bool {
	return d.osAllowsMultipleSelection()
}

// SetAllowsMultipleSelection sets whether more than one item can be selected.
func (d *OpenDialog) SetAllowsMultipleSelection(allow bool) {
	d.osSetAllowsMultipleSelection(allow)
}

// URLs returns the URLs that were chosen.
func (d *OpenDialog) URLs() []string {
	return d.osURLs()
}

// RunModal displays the dialog, allowing the user to make a selection.
// Returns true if successful or false if canceled.
func (d *OpenDialog) RunModal() bool {
	return d.osRunModal()
}

// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

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
func (d *SaveDialog) SetDirectoryURL(dirURL string) *SaveDialog {
	d.osSetDirectoryURL(dirURL)
	return d
}

// AllowedFileTypes returns the set of permitted file types. nil will be
// returned if all files are allowed.
func (d *SaveDialog) AllowedFileTypes() []string {
	return d.osAllowedFileTypes()
}

// SetAllowedFileTypes sets the permitted file types that may be selected for
// saving. Pass in nil to allow all files.
func (d *SaveDialog) SetAllowedFileTypes(allowedExtensions []string) *SaveDialog {
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
	return d
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

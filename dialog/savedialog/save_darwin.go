// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

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

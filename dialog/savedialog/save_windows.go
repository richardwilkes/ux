// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package savedialog

type osSaveDialog = int

func osNew() osSaveDialog {
	return 0 // TODO: Implement
}

func (d *SaveDialog) osDirectoryURL() string {
	return "" // TODO: Implement
}

func (d *SaveDialog) osSetDirectoryURL(dirURL string) {
	// TODO: Implement
}

func (d *SaveDialog) osAllowedFileTypes() []string {
	return nil // TODO: Implement
}

func (d *SaveDialog) osSetAllowedFileTypes(allowed []string) {
	// TODO: Implement
}

func (d *SaveDialog) osURL() string {
	return "" // TODO: Implement
}

func (d *SaveDialog) osRunModal() bool {
	return false // TODO: Implement
}

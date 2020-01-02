// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package clipboard

import "github.com/richardwilkes/ux/clipboard/datatypes"

func osClear() {
	// RAW: Implement
}

func osChangeCount() int {
	return 0
}

func osLoadTypes() {
	// RAW: Implement
}

func osGetData(dataType datatypes.DataType) [][]byte {
	// RAW: Implement
	return nil
}

func osSetData(data []map[datatypes.DataType][]byte) {
	// RAW: Implement
}

func osBytesToURL(in []byte) string {
	return string(in)
}

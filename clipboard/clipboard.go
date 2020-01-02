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

var (
	clipboardLastChangeCount = -1
	clipboardDataTypes       []datatypes.DataType
)

// Clear clears the clipboard contents.
func Clear() {
	clipboardLastChangeCount = -1
	clipboardDataTypes = nil
	osClear()
}

// HasType returns true if the specified data type exists on the clipboard.
func HasType(dataType datatypes.DataType) bool {
	for _, one := range Types() {
		if dataType == one || (dataType.UTI != "" && dataType.UTI == one.UTI) || (dataType.Mime != "" && dataType.Mime == one.Mime) {
			return true
		}
	}
	return false
}

// Types returns the types of data currently on the clipboard.
func Types() []datatypes.DataType {
	changeCount := osChangeCount()
	if changeCount != clipboardLastChangeCount {
		clipboardLastChangeCount = changeCount
		clipboardDataTypes = nil
		osLoadTypes()
	}
	return clipboardDataTypes
}

// GetFirstData returns the bytes for the first item associated with the
// specified data type on the clipboard. An empty slice will be returned if no
// such data type is present.
func GetFirstData(dataType datatypes.DataType) []byte {
	data := osGetData(dataType)
	if len(data) > 0 {
		return data[0]
	}
	return nil
}

// GetData returns a slice holding slices of bytes associated with the
// specified data type on the clipboard. An empty slice will be returned if no
// such data type is present.
func GetData(dataType datatypes.DataType) [][]byte {
	return osGetData(dataType)
}

// SetDataWithType sets the data into the system clipboard.
func SetDataWithType(data []byte, dataType datatypes.DataType) {
	osSetData([]map[datatypes.DataType][]byte{{dataType: data}})
}

// SetDataWithMultipleTypes sets the data into the system clipboard.
func SetDataWithMultipleTypes(data map[datatypes.DataType][]byte) {
	osSetData([]map[datatypes.DataType][]byte{data})
}

// SetData sets the data into the system clipboard.
func SetData(data []map[datatypes.DataType][]byte) {
	osSetData(data)
}

// BytesToURL converts bytes into a URL. On most platforms, this is just a
// simple string() cast. However, macOS has a file reference URL type that
// needs special handling to resolve properly.
func BytesToURL(in []byte) string {
	return osBytesToURL(in)
}

// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package clipboard

import (
	"github.com/richardwilkes/macos/cf"
	"github.com/richardwilkes/macos/ns"
	"github.com/richardwilkes/macos/ut"
	"github.com/richardwilkes/ux/clipboard/datatypes"
)

func osClear() {
	ns.PasteboardGeneral().ClearContents()
}

func osChangeCount() int {
	return ns.PasteboardGeneral().ChangeCount()
}

func osLoadTypes() {
	clipboardDataTypes = TypesForPasteboard(ns.PasteboardGeneral())
}

func osGetData(dataType datatypes.DataType) [][]byte {
	return DataFromPasteboard(ns.PasteboardGeneral(), dataType)
}

func osSetData(data []map[datatypes.DataType][]byte) {
	SetDataToPasteboard(ns.PasteboardGeneral(), data)
}

func osBytesToURL(in []byte) string {
	resolved := ns.URLWithString(string(in)).ResolveFilePath()
	if resolved != "" {
		return resolved
	}
	return string(in)
}

// TypesForPasteboard is exposed for use by drag and drop.
func TypesForPasteboard(pb *ns.Pasteboard) []datatypes.DataType {
	clipTypes := pb.Types()
	count := clipTypes.GetCount()
	dataTypes := make([]datatypes.DataType, 0, count)
	if count > 0 {
		set := make(map[datatypes.DataType]bool)
		for i := 0; i < count; i++ {
			uti := cf.String(clipTypes.GetValueAtIndex(i)).String()
			dt, ok := datatypes.ByUTI[uti]
			if !ok {
				dt = datatypes.DataType{UTI: uti, Mime: ut.TypeCopyPreferredTagWithClassMimeType(uti)}
			}
			if !set[dt] {
				set[dt] = true
				dataTypes = append(dataTypes, dt)
			}
		}
	}
	return dataTypes
}

// DataFromPasteboard is exposed for use by drag and drop.
func DataFromPasteboard(pb *ns.Pasteboard, dataType datatypes.DataType) [][]byte {
	var result [][]byte
	items := pb.Items()
	for _, item := range items {
		itemTypes := item.Types()
		count := itemTypes.GetCount()
		for i := 0; i < count; i++ {
			if dataType.UTI == cf.String(itemTypes.GetValueAtIndex(i)).String() {
				data := item.DataForType(dataType.UTI)
				result = append(result, data.GetBytes(0, data.GetLength()))
				break
			}
		}
	}
	return result
}

// SetDataToPasteboard is exposed for use by drag and drop.
func SetDataToPasteboard(pb *ns.Pasteboard, data []map[datatypes.DataType][]byte) {
	pb.ClearContents()
	items := make([]*ns.PasteboardItem, len(data))
	for i, m := range data {
		item := ns.NewPasteboardItem()
		for k, v := range m {
			item.SetDataForType(cf.DataCreate(v), k.UTI)
		}
		items[i] = item
	}
	pb.SetItems(items)
}

// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package ux

import (
	"math"

	"github.com/richardwilkes/ux/clipboard/datatypes"
)

// DragOperation holds the type of drag operation.
type DragOperation = uint32

// Possible DragOperation values.
const (
	DragOperationCopy DragOperation = 1 << iota
	DragOperationLink
	DragOperationGeneric
	DragOperationPrivate
	DragOperationMove
	DragOperationDelete
	DragOperationNone  DragOperation = 0
	DragOperationEvery DragOperation = math.MaxUint32
)

// DragInfo holds information about the current drag.
type DragInfo struct {
	Sequence            int
	SourceOperationMask DragOperation
	DragX               float64
	DragY               float64
	DragImageX          float64
	DragImageY          float64
	ValidItemsForDrop   int
	ItemTypes           []datatypes.DataType
	DataForType         func(dataType datatypes.DataType) [][]byte
}

// HasType returns true if the specified data type is present.
func (di *DragInfo) HasType(dataType datatypes.DataType) bool {
	for _, dt := range di.ItemTypes {
		if dt.UTI == dataType.UTI {
			return true
		}
	}
	return false
}

// FirstTypePresent returns the first data type that matches the available
// data types, or returns datatypes.None.
func (di *DragInfo) FirstTypePresent(dataType ...datatypes.DataType) datatypes.DataType {
	for _, one := range dataType {
		for _, dt := range di.ItemTypes {
			if dt.UTI == one.UTI {
				return one
			}
		}
	}
	return datatypes.None
}

// ApplyOffset applies the delta to the drag and drag image positions.
func (di *DragInfo) ApplyOffset(dx, dy float64) {
	di.DragX += dx
	di.DragY += dy
	di.DragImageX += dx
	di.DragImageY += dy
}

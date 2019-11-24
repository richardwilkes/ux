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

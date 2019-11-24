package clipboard

import (
	"syscall"
	"unsafe"

	"github.com/richardwilkes/ux/clipboard/datatypes"
	"github.com/richardwilkes/win32"
)

var (
	clipboard   map[datatypes.DataType][]byte
	clipboardID = uint64(1)
)

func osClear() {
	if win32.OpenClipboard(0) {
		win32.EmptyClipboard()
		win32.CloseClipboard()
	}
	clipboard = nil
	clipboardID++
}

func osChangeCount() int {
	return win32.GetClipboardSequenceNumber()
}

func osLoadTypes() {
	if win32.OpenClipboard(0) {
		defer win32.CloseClipboard()
		format := uint(0)
		haveText := false
		for {
			format = win32.EnumClipboardFormats(format)
			switch format {
			case 0:
				return
			case win32.CF_TEXT, win32.CF_OEMTEXT, win32.CF_UNICODETEXT:
				if !haveText {
					clipboardDataTypes = append(clipboardDataTypes, datatypes.PlainText)
					haveText = true
				}
			case win32.CF_PRIVATEFIRST:
				if data := win32.GetClipboardData(format); data == win32.HANDLE(clipboardID) {
					for k := range clipboard {
						clipboardDataTypes = append(clipboardDataTypes, k)
					}
				}
			default:
			}
		}
	}
}

func osGetData(dataType datatypes.DataType) [][]byte {
	// FIXME: Only retrieves one item for the given type
	if win32.OpenClipboard(0) {
		defer win32.CloseClipboard()
		switch dataType.UTI {
		case datatypes.PlainText.UTI:
			if data := win32.GetClipboardData(win32.CF_UNICODETEXT); data != 0 {
				ptr := (*[1 << 30]byte)(unsafe.Pointer(data))
				i := 0
				for {
					if ptr[i] == 0 {
						break
					}
					i++
				}
				str := make([]byte, i)
				copy(str, ptr[:])
				return [][]byte{str}
			}
		default:
			if data := win32.GetClipboardData(win32.CF_PRIVATEFIRST); data == win32.HANDLE(clipboardID) {
				if v, ok := clipboard[dataType]; ok {
					return [][]byte{v}
				}
			}
		}
	}
	return nil
}

func osSetData(data []map[datatypes.DataType][]byte) {
	// FIXME: Only stores the first item
	if hwnd := win32.GetActiveWindow(); hwnd != 0 && win32.OpenClipboard(hwnd) {
		defer win32.CloseClipboard()
		win32.EmptyClipboard()
		clipboard = make(map[datatypes.DataType][]byte)
		clipboardID++
		setPrivate := false
		if len(data) > 0 {
			for k, v := range data[0] {
				switch k.UTI {
				case datatypes.PlainText.UTI:
					if str, err := syscall.UTF16FromString(string(v)); err != nil {
						if mem := win32.GlobalAlloc(win32.GMEM_MOVEABLE, len(str)*2); mem != 0 {
							if p := win32.GlobalLock(mem); p != nil {
								win32.MoveMemory(p, unsafe.Pointer(&str[0]), len(str)*2)
								win32.GlobalUnlock(mem)
							}
							if !win32.SetClipboardData(win32.CF_UNICODETEXT, win32.HANDLE(mem)) {
								win32.GlobalFree(mem)
							}
						}
					}
				default:
					if !setPrivate {
						setPrivate = true
						win32.SetClipboardData(win32.CF_PRIVATEFIRST, win32.HANDLE(clipboardID))
					}
					clipboard[k] = v
				}
			}
		}
	}
}

func osBytesToURL(in []byte) string {
	return string(in)
}

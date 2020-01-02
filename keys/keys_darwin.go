// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package keys

// Known key codes.
var (
	A              = &Key{Code: 0x00, Name: "A", Rune: 'a'}
	B              = &Key{Code: 0x0b, Name: "B", Rune: 'b'}
	C              = &Key{Code: 0x08, Name: "C", Rune: 'c'}
	D              = &Key{Code: 0x02, Name: "D", Rune: 'd'}
	E              = &Key{Code: 0x0e, Name: "E", Rune: 'e'}
	F              = &Key{Code: 0x03, Name: "F", Rune: 'f'}
	G              = &Key{Code: 0x05, Name: "G", Rune: 'g'}
	H              = &Key{Code: 0x04, Name: "H", Rune: 'h'}
	I              = &Key{Code: 0x22, Name: "I", Rune: 'i'}
	J              = &Key{Code: 0x26, Name: "J", Rune: 'j'}
	K              = &Key{Code: 0x28, Name: "K", Rune: 'k'}
	L              = &Key{Code: 0x25, Name: "L", Rune: 'l'}
	M              = &Key{Code: 0x2e, Name: "M", Rune: 'm'}
	N              = &Key{Code: 0x2d, Name: "N", Rune: 'n'}
	O              = &Key{Code: 0x1f, Name: "O", Rune: 'o'}
	P              = &Key{Code: 0x23, Name: "P", Rune: 'p'}
	Q              = &Key{Code: 0x0c, Name: "Q", Rune: 'q'}
	R              = &Key{Code: 0x0f, Name: "R", Rune: 'r'}
	S              = &Key{Code: 0x01, Name: "S", Rune: 's'}
	T              = &Key{Code: 0x11, Name: "T", Rune: 't'}
	U              = &Key{Code: 0x20, Name: "U", Rune: 'u'}
	V              = &Key{Code: 0x09, Name: "V", Rune: 'v'}
	W              = &Key{Code: 0x0d, Name: "W", Rune: 'w'}
	X              = &Key{Code: 0x07, Name: "X", Rune: 'x'}
	Y              = &Key{Code: 0x10, Name: "Y", Rune: 'y'}
	Z              = &Key{Code: 0x06, Name: "Z", Rune: 'z'}
	One            = &Key{Code: 0x12, Name: "1"}
	Two            = &Key{Code: 0x13, Name: "2"}
	Three          = &Key{Code: 0x14, Name: "3"}
	Four           = &Key{Code: 0x15, Name: "4"}
	Five           = &Key{Code: 0x17, Name: "5"}
	Six            = &Key{Code: 0x16, Name: "6"}
	Seven          = &Key{Code: 0x1a, Name: "7"}
	Eight          = &Key{Code: 0x1c, Name: "8"}
	Nine           = &Key{Code: 0x19, Name: "9"}
	Zero           = &Key{Code: 0x1d, Name: "0"}
	Return         = &Key{Code: 0x24, Name: "Return", Rune: '\x0d'}
	Escape         = &Key{Code: 0x35, Name: "Escape", Rune: '\x1b'}
	Backspace      = &Key{Code: 0x33, Name: "Backspace", Rune: '\x08'}
	Tab            = &Key{Code: 0x30, Name: "Tab", Rune: '\x09'}
	Space          = &Key{Code: 0x31, Name: "Space", Rune: ' '}
	Minus          = &Key{Code: 0x1b, Name: "Minus", Rune: '-'}
	Equal          = &Key{Code: 0x18, Name: "="}
	LeftBracket    = &Key{Code: 0x21, Name: "["}
	RightBracket   = &Key{Code: 0x1e, Name: "]"}
	Backslash      = &Key{Code: 0x2a, Name: `\`}
	Semicolon      = &Key{Code: 0x29, Name: ";"}
	Quote          = &Key{Code: 0x27, Name: "'"}
	Backquote      = &Key{Code: 0x32, Name: "`"}
	Comma          = &Key{Code: 0x2b, Name: ","}
	Period         = &Key{Code: 0x2f, Name: "."}
	Slash          = &Key{Code: 0x2c, Name: "/"}
	F1             = &Key{Code: 0x7a, Name: "F1", Rune: '\uf704'}
	F2             = &Key{Code: 0x78, Name: "F2", Rune: '\uf705'}
	F3             = &Key{Code: 0x63, Name: "F3", Rune: '\uf706'}
	F4             = &Key{Code: 0x76, Name: "F4", Rune: '\uf707'}
	F5             = &Key{Code: 0x60, Name: "F5", Rune: '\uf708'}
	F6             = &Key{Code: 0x61, Name: "F6", Rune: '\uf709'}
	F7             = &Key{Code: 0x62, Name: "F7", Rune: '\uf70a'}
	F8             = &Key{Code: 0x64, Name: "F8", Rune: '\uf70b'}
	F9             = &Key{Code: 0x65, Name: "F9", Rune: '\uf70c'}
	F10            = &Key{Code: 0x6d, Name: "F10", Rune: '\uf70d'}
	F11            = &Key{Code: 0x67, Name: "F11", Rune: '\uf70e'}
	F12            = &Key{Code: 0x6f, Name: "F12", Rune: '\uf70f'}
	F13            = &Key{Code: 0x69, Name: "F13", Rune: '\uf710'}
	F14            = &Key{Code: 0x6b, Name: "F14", Rune: '\uf711'}
	F15            = &Key{Code: 0x71, Name: "F15", Rune: '\uf712'}
	Delete         = &Key{Code: 0x75, Name: "Delete", Rune: '\uf728'}
	Home           = &Key{Code: 0x73, Name: "Home", Rune: '\uf729'}
	End            = &Key{Code: 0x77, Name: "End", Rune: '\uf72b'}
	PageUp         = &Key{Code: 0x74, Name: "PageUp", Rune: '\uf72c'}
	PageDown       = &Key{Code: 0x79, Name: "PageDown", Rune: '\uf72d'}
	Left           = &Key{Code: 0x7b, Name: "Left", Rune: '\uf702'}
	Up             = &Key{Code: 0x7e, Name: "Up", Rune: '\uf700'}
	Right          = &Key{Code: 0x7c, Name: "Right", Rune: '\uf703'}
	Down           = &Key{Code: 0x7d, Name: "Down", Rune: '\uf701'}
	Clear          = &Key{Code: 0x47, Name: "Clear"}
	NumpadDivide   = &Key{Code: 0x4b, Name: "/"}
	NumpadMultiply = &Key{Code: 0x43, Name: "*"}
	NumpadAdd      = &Key{Code: 0x45, Name: "+"}
	NumpadSubtract = &Key{Code: 0x4e, Name: "-"}
	NumpadDecimal  = &Key{Code: 0x41, Name: "."}
	NumpadEnter    = &Key{Code: 0x4c, Name: "Enter", Rune: '\x0d'}
	Numpad1        = &Key{Code: 0x53, Name: "1"}
	Numpad2        = &Key{Code: 0x54, Name: "2"}
	Numpad3        = &Key{Code: 0x55, Name: "3"}
	Numpad4        = &Key{Code: 0x56, Name: "4"}
	Numpad5        = &Key{Code: 0x57, Name: "5"}
	Numpad6        = &Key{Code: 0x58, Name: "6"}
	Numpad7        = &Key{Code: 0x59, Name: "7"}
	Numpad8        = &Key{Code: 0x5b, Name: "8"}
	Numpad9        = &Key{Code: 0x5c, Name: "9"}
	Numpad0        = &Key{Code: 0x52, Name: "0"}
)

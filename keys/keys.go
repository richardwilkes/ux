package keys

// ByCode is a map of all known key codes.
var ByCode = map[int]*Key{
	A.Code:              A,
	B.Code:              B,
	C.Code:              C,
	D.Code:              D,
	E.Code:              E,
	F.Code:              F,
	G.Code:              G,
	H.Code:              H,
	I.Code:              I,
	J.Code:              J,
	K.Code:              K,
	L.Code:              L,
	M.Code:              M,
	N.Code:              N,
	O.Code:              O,
	P.Code:              P,
	Q.Code:              Q,
	R.Code:              R,
	S.Code:              S,
	T.Code:              T,
	U.Code:              U,
	V.Code:              V,
	W.Code:              W,
	X.Code:              X,
	Y.Code:              Y,
	Z.Code:              Z,
	One.Code:            One,
	Two.Code:            Two,
	Three.Code:          Three,
	Four.Code:           Four,
	Five.Code:           Five,
	Six.Code:            Six,
	Seven.Code:          Seven,
	Eight.Code:          Eight,
	Nine.Code:           Nine,
	Zero.Code:           Zero,
	Return.Code:         Return,
	Escape.Code:         Escape,
	Backspace.Code:      Backspace,
	Tab.Code:            Tab,
	Space.Code:          Space,
	Minus.Code:          Minus,
	Equal.Code:          Equal,
	LeftBracket.Code:    LeftBracket,
	RightBracket.Code:   RightBracket,
	Backslash.Code:      Backslash,
	Semicolon.Code:      Semicolon,
	Quote.Code:          Quote,
	Backquote.Code:      Backquote,
	Comma.Code:          Comma,
	Period.Code:         Period,
	Slash.Code:          Slash,
	F1.Code:             F1,
	F2.Code:             F2,
	F3.Code:             F3,
	F4.Code:             F4,
	F5.Code:             F5,
	F6.Code:             F6,
	F7.Code:             F7,
	F8.Code:             F8,
	F9.Code:             F9,
	F10.Code:            F10,
	F11.Code:            F11,
	F12.Code:            F12,
	F13.Code:            F13,
	F14.Code:            F14,
	F15.Code:            F15,
	Delete.Code:         Delete,
	Home.Code:           Home,
	End.Code:            End,
	PageUp.Code:         PageUp,
	PageDown.Code:       PageDown,
	Left.Code:           Left,
	Up.Code:             Up,
	Right.Code:          Right,
	Down.Code:           Down,
	Clear.Code:          Clear,
	NumpadDivide.Code:   NumpadDivide,
	NumpadMultiply.Code: NumpadMultiply,
	NumpadAdd.Code:      NumpadAdd,
	NumpadSubtract.Code: NumpadSubtract,
	NumpadDecimal.Code:  NumpadDecimal,
	NumpadEnter.Code:    NumpadEnter,
	Numpad1.Code:        Numpad1,
	Numpad2.Code:        Numpad2,
	Numpad3.Code:        Numpad3,
	Numpad4.Code:        Numpad4,
	Numpad5.Code:        Numpad5,
	Numpad6.Code:        Numpad6,
	Numpad7.Code:        Numpad7,
	Numpad8.Code:        Numpad8,
	Numpad9.Code:        Numpad9,
	Numpad0.Code:        Numpad0,
}

// Some common aliases
var (
	NumpadDelete   = NumpadDecimal
	NumpadEnd      = Numpad1
	NumpadDown     = Numpad2
	NumpadPageDown = Numpad3
	NumpadLeft     = Numpad4
	NumpadRight    = Numpad6
	NumpadHome     = Numpad7
	NumpadUp       = Numpad8
	NumpadPageUp   = Numpad9
)

// Key holds information about a key.
type Key struct {
	Code int
	Name string
	Rune rune // Only used by macOS
}

// RuneStr returns the key code string needed by macOS for menu processing.
func (k *Key) RuneStr() string {
	if k.Rune != 0 {
		return string([]rune{k.Rune})
	}
	return k.Name
}

// IsControlAction returns true if the keyCode should trigger a control, such
// as a button, that is focused.
func IsControlAction(keyCode int) bool {
	return keyCode == Return.Code || keyCode == NumpadEnter.Code || keyCode == Space.Code
}

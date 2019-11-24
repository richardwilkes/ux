package keys

// Diacritics provides handling of diacritic input.
type Diacritics int

// ProcessInput processes keyboard input and returns the resulting rune.
func (d *Diacritics) ProcessInput(keyCode int, ch rune, mod Modifiers) rune {
	value := int(*d)
	if value != 0 {
		if mod&^ShiftModifier == 0 {
			switch ch {
			case 'a':
				switch value {
				case E.Code:
					ch = 'á'
				case I.Code:
					ch = 'â'
				case Backquote.Code:
					ch = 'à'
				case N.Code:
					ch = 'ã'
				case U.Code:
					ch = 'ä'
				}
			case 'A':
				switch value {
				case E.Code:
					ch = 'Á'
				case I.Code:
					ch = 'Â'
				case Backquote.Code:
					ch = 'À'
				case N.Code:
					ch = 'Ã'
				case U.Code:
					ch = 'Ä'
				}
			case 'e':
				switch value {
				case E.Code:
					ch = 'é'
				case I.Code:
					ch = 'ê'
				case Backquote.Code:
					ch = 'è'
				case U.Code:
					ch = 'ë'
				}
			case 'E':
				switch value {
				case E.Code:
					ch = 'É'
				case I.Code:
					ch = 'Ê'
				case Backquote.Code:
					ch = 'È'
				case U.Code:
					ch = 'Ë'
				}
			case 'i':
				switch value {
				case E.Code:
					ch = 'í'
				case I.Code:
					ch = 'î'
				case Backquote.Code:
					ch = 'ì'
				case U.Code:
					ch = 'ï'
				}
			case 'I':
				switch value {
				case E.Code:
					ch = 'Í'
				case I.Code:
					ch = 'Î'
				case Backquote.Code:
					ch = 'Ì'
				case U.Code:
					ch = 'Ï'
				}
			case 'o':
				switch value {
				case E.Code:
					ch = 'ó'
				case I.Code:
					ch = 'ô'
				case Backquote.Code:
					ch = 'ò'
				case N.Code:
					ch = 'õ'
				case U.Code:
					ch = 'ö'
				}
			case 'O':
				switch value {
				case E.Code:
					ch = 'Ó'
				case I.Code:
					ch = 'Ô'
				case Backquote.Code:
					ch = 'Ò'
				case N.Code:
					ch = 'Õ'
				case U.Code:
					ch = 'Ö'
				}
			case 'u':
				switch value {
				case E.Code:
					ch = 'ú'
				case I.Code:
					ch = 'û'
				case Backquote.Code:
					ch = 'ù'
				case U.Code:
					ch = 'ü'
				}
			case 'U':
				switch value {
				case E.Code:
					ch = 'Ú'
				case I.Code:
					ch = 'Û'
				case Backquote.Code:
					ch = 'Ù'
				case U.Code:
					ch = 'Ü'
				}
			}
		}
		value = 0
	}
	if mod&^ShiftModifier == OptionModifier {
		switch keyCode {
		case E.Code, I.Code, Backquote.Code, N.Code, U.Code:
			value = keyCode
		default:
			value = 0
		}
	}
	if value != 0 {
		ch = 0
	}
	*d = Diacritics(value)
	return ch
}

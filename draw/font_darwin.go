package draw

import (
	"github.com/richardwilkes/macos/cf"
	"github.com/richardwilkes/macos/ct"
)

type osFont struct {
	ref ct.Font
}

var osSystemFontMapping = make(map[string]ct.FontUIFontType)

func osInitSystemFonts() {
	for _, one := range []ct.FontUIFontType{
		ct.FontUIFontUser,
		ct.FontUIFontUserFixedPitch,
		ct.FontUIFontSystem,
		ct.FontUIFontEmphasizedSystem,
		ct.FontUIFontSmallSystem,
		ct.FontUIFontSmallEmphasizedSystem,
		ct.FontUIFontMiniSystem,
		ct.FontUIFontMiniEmphasizedSystem,
		ct.FontUIFontViews,
		ct.FontUIFontApplication,
		ct.FontUIFontLabel,
		ct.FontUIFontMenuTitle,
		ct.FontUIFontMenuItem,
		ct.FontUIFontMenuItemMark,
		ct.FontUIFontMenuItemCmdKey,
		ct.FontUIFontWindowTitle,
		ct.FontUIFontPushButton,
		ct.FontUIFontUtilityWindowTitle,
		ct.FontUIFontAlertHeader,
		ct.FontUIFontSystemDetail,
		ct.FontUIFontEmphasizedSystemDetail,
		ct.FontUIFontToolbar,
		ct.FontUIFontSmallToolbar,
		ct.FontUIFontMessage,
		ct.FontUIFontPalette,
		ct.FontUIFontToolTip,
		ct.FontUIFontControlContent,
	} {
		osSystemFontMapping[fontDesc(one).Family] = one
	}
	UserFont = NewFont(fontDesc(ct.FontUIFontUser))
	UserMonospacedFont = NewFont(fontDesc(ct.FontUIFontUserFixedPitch))
	SystemFont = NewFont(fontDesc(ct.FontUIFontSystem))
	EmphasizedSystemFont = NewFont(fontDesc(ct.FontUIFontEmphasizedSystem))
	SmallSystemFont = NewFont(fontDesc(ct.FontUIFontSmallSystem))
	SmallEmphasizedSystemFont = NewFont(fontDesc(ct.FontUIFontSmallEmphasizedSystem))
	ViewsFont = NewFont(fontDesc(ct.FontUIFontViews))
	LabelFont = NewFont(fontDesc(ct.FontUIFontLabel))
	MenuFont = NewFont(fontDesc(ct.FontUIFontMenuItem))
	MenuCmdKeyFont = NewFont(fontDesc(ct.FontUIFontMenuItemCmdKey))
}

func osFontFamilies() []string {
	var families []string
	unique := make(map[string]bool)
	allFontsCollection := ct.FontCollectionCreateFromAvailableFonts(0)
	allFonts := allFontsCollection.CreateMatchingFontDescriptors()
	for i := allFonts.GetCount() - 1; i >= 0; i-- {
		str := cf.String(ct.FontDescriptor(allFonts.GetValueAtIndex(i)).CopyAttribute(ct.FontFamilyNameAttribute))
		family := str.String()
		str.Release()
		if _, exists := unique[family]; !exists {
			unique[family] = true
			families = append(families, family)
		}
	}
	allFonts.Release()
	allFontsCollection.Release()
	return families
}

func osNewFont(desc FontDescriptor) *Font {
	var ref ct.Font
	if fontType, exists := osSystemFontMapping[desc.Family]; exists {
		ref = ct.FontCreateUIFontForLanguage(fontType, desc.Size, "")
	} else {
		ref = ct.FontCreateWithName(desc.Family, desc.Size, nil)
	}
	bold := desc.Bold
	italic := desc.Italic
	if bold || italic {
		var traits ct.FontSymbolicTraits
	tryAgain:
		if bold {
			traits |= ct.FontBoldTrait
		}
		if italic {
			traits |= ct.FontItalicTrait
		}
		adjustedFont := ref.CreateCopyWithSymbolicTraits(desc.Size, nil, traits, ct.FontItalicTrait|ct.FontBoldTrait)
		if adjustedFont == 0 {
			traits = 0
			if italic {
				italic = false
			} else if bold {
				bold = false
			}
			if bold || italic {
				goto tryAgain
			}
			adjustedFont = ref
		}
		ref = adjustedFont
	}
	ref.Retain()
	return &Font{
		ascent:     ref.GetAscent(),
		descent:    ref.GetDescent(),
		leading:    ref.GetLeading(),
		desc:       desc,
		osFont:     osFont{ref: ref},
		monospaced: ref.GetSymbolicTraits()&ct.FontMonoSpaceTrait == ct.FontMonoSpaceTrait,
	}
}

func (f *Font) osWidth(str string) float64 {
	as := f.createAttributedString(str)
	line := ct.LineCreateWithAttributedString(cf.AttributedString(as))
	width := line.GetTypographicBounds(nil, nil, nil)
	line.Release()
	as.Release()
	return width
}

func (f *Font) osIndexForPosition(x float64, str string) int {
	as := f.createAttributedString(str)
	line := ct.LineCreateWithAttributedString(cf.AttributedString(as))
	i := line.GetStringIndexForPosition(x, 0)
	line.Release()
	as.Release()
	return i
}

func (f *Font) osPositionForIndex(index int, str string) float64 {
	as := f.createAttributedString(str)
	line := ct.LineCreateWithAttributedString(cf.AttributedString(as))
	x := line.GetOffsetForStringIndex(index, nil)
	line.Release()
	as.Release()
	return x
}

// Dispose of the font, releasing any OS resources associated with it.
func (f *Font) osDispose() {
	if f.ref != 0 {
		f.ref.Release()
		f.ref = 0
	}
}

func (f *Font) createAttributedString(str string) cf.MutableAttributedString {
	as := cf.AttributedStringCreateMutable(0)
	as.BeginEditing()
	s := cf.StringCreateWithString(str)
	as.ReplaceString(0, 0, s)
	as.SetAttribute(0, s.GetLength(), ct.FontAttributeName, cf.Type(f.ref))
	as.EndEditing()
	return as
}

func fontDesc(fontType ct.FontUIFontType) FontDescriptor {
	f := ct.FontCreateUIFontForLanguage(fontType, 0, "")
	traits := f.GetSymbolicTraits()
	return FontDescriptor{
		Family: f.FamilyName(),
		Size:   f.GetSize(),
		Bold:   traits&ct.FontBoldTrait == ct.FontBoldTrait,
		Italic: traits&ct.FontItalicTrait == ct.FontItalicTrait,
	}
}

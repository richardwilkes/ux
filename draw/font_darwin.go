// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package draw

import (
	"github.com/richardwilkes/macos/cf"
	"github.com/richardwilkes/macos/ct"
	"github.com/richardwilkes/toolbox/errs"
)

type osFont struct {
	ref ct.Font
}

var (
	osSystemFontMapping = make(map[string]ct.FontUIFontType)
	osLoadedFontMapping = make(map[string]ct.FontDescriptor)
)

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
		osSystemFontMapping[osiFontDesc(one).Family] = one
	}
	UserFont = NewFont(osiFontDesc(ct.FontUIFontUser))
	UserMonospacedFont = NewFont(osiFontDesc(ct.FontUIFontUserFixedPitch))
	SystemFont = NewFont(osiFontDesc(ct.FontUIFontSystem))
	EmphasizedSystemFont = NewFont(osiFontDesc(ct.FontUIFontEmphasizedSystem))
	SmallSystemFont = NewFont(osiFontDesc(ct.FontUIFontSmallSystem))
	SmallEmphasizedSystemFont = NewFont(osiFontDesc(ct.FontUIFontSmallEmphasizedSystem))
	ViewsFont = NewFont(osiFontDesc(ct.FontUIFontViews))
	LabelFont = NewFont(osiFontDesc(ct.FontUIFontLabel))
	MenuFont = NewFont(osiFontDesc(ct.FontUIFontMenuItem))
	MenuCmdKeyFont = NewFont(osiFontDesc(ct.FontUIFontMenuItemCmdKey))
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
	var exists bool
	var fontType ct.FontUIFontType
	var fd ct.FontDescriptor
	if fontType, exists = osSystemFontMapping[desc.Family]; exists {
		ref = ct.FontCreateUIFontForLanguage(fontType, desc.Size, "")
	} else if fd, exists = osLoadedFontMapping[desc.Family]; exists {
		ref = ct.FontCreateWithFontDescriptor(fd, desc.Size, nil)
	} else {
		ref = ct.FontCreateWithName(desc.Family, desc.Size, nil)
	}
	return osiNewFont(ref, desc)
}

func osiNewFont(ref ct.Font, desc FontDescriptor) *Font {
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

func osNewFontFromData(data []byte) (*Font, error) {
	cfData := cf.DataCreate(data)
	defer cfData.Release()
	desc := ct.FontManagerCreateFontDescriptorFromData(cfData)
	if desc == 0 {
		return nil, errs.New("unable to create font descriptor from data")
	}
	ref := ct.FontCreateWithFontDescriptor(desc, 0, nil)
	if ref == 0 {
		return nil, errs.New("unable to create font from data")
	}
	fd := osiFontDescFromFont(ref)
	osLoadedFontMapping[fd.Family] = desc
	return osiNewFont(ref, osiFontDescFromFont(ref)), nil
}

func (f *Font) osWidth(str string) float64 {
	as := f.osiCreateAttributedString(str)
	line := ct.LineCreateWithAttributedString(cf.AttributedString(as))
	width := line.GetTypographicBounds(nil, nil, nil)
	line.Release()
	as.Release()
	return width
}

func (f *Font) osIndexForPosition(x float64, str string) int {
	as := f.osiCreateAttributedString(str)
	line := ct.LineCreateWithAttributedString(cf.AttributedString(as))
	i := line.GetStringIndexForPosition(x, 0)
	line.Release()
	as.Release()
	return i
}

func (f *Font) osPositionForIndex(index int, str string) float64 {
	as := f.osiCreateAttributedString(str)
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

func (f *Font) osiCreateAttributedString(str string) cf.MutableAttributedString {
	as := cf.AttributedStringCreateMutable(0)
	as.BeginEditing()
	s := cf.StringCreateWithString(str)
	as.ReplaceString(0, 0, s)
	as.SetAttribute(0, s.GetLength(), ct.FontAttributeName, cf.Type(f.ref))
	as.EndEditing()
	return as
}

func osiFontDesc(fontType ct.FontUIFontType) FontDescriptor {
	return osiFontDescFromFont(ct.FontCreateUIFontForLanguage(fontType, 0, ""))
}

func osiFontDescFromFont(f ct.Font) FontDescriptor {
	traits := f.GetSymbolicTraits()
	return FontDescriptor{
		Family: f.FamilyName(),
		Size:   f.GetSize(),
		Bold:   traits&ct.FontBoldTrait == ct.FontBoldTrait,
		Italic: traits&ct.FontItalicTrait == ct.FontItalicTrait,
	}
}

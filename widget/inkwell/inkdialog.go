// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package inkwell

import (
	"strconv"

	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/dialog"
	"github.com/richardwilkes/ux/dialog/opendialog"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/ids"
	"github.com/richardwilkes/ux/layout/align"
	"github.com/richardwilkes/ux/layout/flex"
	"github.com/richardwilkes/ux/widget/button"
	"github.com/richardwilkes/ux/widget/label"
	"github.com/richardwilkes/ux/widget/textfield"
)

type inkDialog struct {
	well             *InkWell
	ink              draw.Ink
	dialog           *dialog.Dialog
	panel            *ux.Panel
	preview          *ux.Panel
	right            *ux.Panel
	redField         *textfield.TextField
	greenField       *textfield.TextField
	blueField        *textfield.TextField
	alphaField       *textfield.TextField
	validationFields []*textfield.TextField
}

// TODO: Implement gradient selection

func showDialog(well *InkWell) {
	d := &inkDialog{
		well:    well,
		ink:     well.Ink(),
		panel:   ux.NewPanel(),
		preview: ux.NewPanel().SetBorder(border.NewCompound(border.NewLine(draw.Black, 0, geom.NewUniformInsets(1), false), border.NewLine(draw.White, 0, geom.NewUniformInsets(1), false))),
		right:   ux.NewPanel(),
	}
	flex.New().Columns(2).HAlign(align.Fill).Apply(d.panel)
	flex.NewData().VAlign(align.Start).SizeHint(geom.Size{Width: 64, Height: 64}).Apply(d.preview)
	d.preview.DrawCallback = func(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
		if p, ok := d.ink.(*draw.Pattern); ok {
			gc.Clip()
			p.Image().DrawInRect(gc, d.preview.ContentRect(false))
		} else {
			gc.Rect(dirty)
			gc.Fill(d.ink)
		}
	}
	d.panel.AddChild(d.preview)
	flex.New().Columns(2).Apply(d.right)
	flex.NewData().HGrab(true).HAlign(align.Fill).Apply(d.right)
	d.panel.AddChild(d.right)
	allowed := well.AllowedTypes()
	if allowed&ColorInkWellMask != 0 {
		color := draw.Black
		switch inkColor := d.ink.(type) {
		case draw.Color:
			color = inkColor
		case *draw.Color:
			color = *inkColor
		default:
		}
		d.redField = d.addEntryField(i18n.Text("Red:"), color.Red())
		d.greenField = d.addEntryField(i18n.Text("Green:"), color.Green())
		d.blueField = d.addEntryField(i18n.Text("Blue:"), color.Blue())
		d.alphaField = d.addEntryField(i18n.Text("Alpha:"), color.Alpha())
	}
	if allowed&PatternInkWellMask != 0 {
		b := button.New().SetText(i18n.Text("Select Image…"))
		flex.NewData().HSpan(2).HAlign(align.Middle).Apply(b)
		b.ClickCallback = func() {
			openDialog := opendialog.New().SetAllowedFileTypes([]string{"png", "jpg", "jpeg", "gif"})
			if openDialog.RunModal() {
				unable := i18n.Text("Unable to load image")
				if urlStr := draw.DistillImageURL(openDialog.URLs()[0]); urlStr == "" {
					dialog.ErrorDialogWithMessage(unable, "Invalid URL")
				} else {
					if img, err := draw.NewImageFromURL(urlStr, d.well.imageScale); err != nil {
						dialog.ErrorDialogWithMessage(unable, err.Error())
					} else {
						if d.well.ValidateImageCallback != nil {
							img = d.well.ValidateImageCallback(img)
						}
						if img == nil {
							dialog.ErrorDialogWithMessage(unable, "")
						} else {
							d.ink = draw.NewPattern(img)
							d.preview.MarkForRedraw()
						}
					}
				}
			}
		}
		if len(d.right.Children()) > 0 {
			b.SetBorder(border.NewEmpty(geom.Insets{Top: 10}))
		}
		d.right.AddChild(b.AsPanel())
	}
	var err error
	d.dialog, err = dialog.NewDialog(nil, d.panel, []*dialog.ButtonInfo{dialog.NewCancelButtonInfo(), dialog.NewOKButtonInfo()})
	if err != nil {
		jot.Error(err)
		return
	}
	d.dialog.Window().SetTitle(i18n.Text("Choose an ink"))
	if d.dialog.RunModal() == ids.ModalResponseOK {
		well.SetInk(d.ink)
	}
}

func (d *inkDialog) addEntryField(title string, value int) *textfield.TextField {
	l := label.New().SetText(title).SetHAlign(align.End)
	flex.NewData().HAlign(align.End).Apply(l)
	d.right.AddChild(l.AsPanel())
	field := textfield.New().SetText(strconv.Itoa(value)).SetMinimumTextWidth(50)
	flex.NewData().HGrab(true).HAlign(align.Fill).Apply(field)
	field.ValidateCallback = func() bool {
		v, err := strconv.Atoi(field.Text())
		valid := err == nil && v >= 0 && v <= 255
		if valid {
			var r int
			if r, err = strconv.Atoi(d.redField.Text()); err == nil {
				var g int
				if g, err = strconv.Atoi(d.greenField.Text()); err == nil {
					var b int
					if b, err = strconv.Atoi(d.blueField.Text()); err == nil {
						var a int
						if a, err = strconv.Atoi(d.alphaField.Text()); err == nil {
							d.ink = draw.ARGB(float64(a)/255, r, g, b)
							d.preview.MarkForRedraw()
						}
					}
				}
			}
		}
		d.adjustOKButton(field, valid)
		return valid
	}
	d.validationFields = append(d.validationFields, field)
	d.right.AddChild(field.AsPanel())
	return field
}

func (d *inkDialog) adjustOKButton(field *textfield.TextField, valid bool) {
	if d.dialog != nil {
		enabled := valid
		if enabled {
			for _, f := range d.validationFields {
				if f != field && f.Invalid() {
					enabled = false
					break
				}
			}
		}
		d.dialog.Button(ids.ModalResponseOK).SetEnabled(enabled)
	}
}

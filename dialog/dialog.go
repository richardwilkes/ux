package dialog

import (
	"strings"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/display"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/icons"
	"github.com/richardwilkes/ux/ids"
	"github.com/richardwilkes/ux/keys"
	"github.com/richardwilkes/ux/layout"
	"github.com/richardwilkes/ux/layout/align"
	"github.com/richardwilkes/ux/layout/flex"
	"github.com/richardwilkes/ux/widget/button"
	"github.com/richardwilkes/ux/widget/label"
)

type buttonData struct {
	info   *ButtonInfo
	button *button.Button
}

// Dialog holds information about a dialog.
type Dialog struct {
	wnd     *ux.Window
	buttons map[int]*buttonData
	err     error
}

// Window returns the underlying window.
func (d *Dialog) Window() *ux.Window {
	return d.wnd
}

// Button returns the button mapped to the given response code.
func (d *Dialog) Button(responseCode int) *button.Button {
	if bd, ok := d.buttons[responseCode]; ok {
		return bd.button
	}
	return nil
}

// RunModal displays and brings this dialog to the front, the runs a modal
// event loop until StopModal is called.
func (d *Dialog) RunModal() int {
	return d.wnd.RunModal()
}

// StopModal stops the current modal event loop, closes the dialog, and
// propagates the provided code as the result to RunModal().
func (d *Dialog) StopModal(code int) {
	d.wnd.StopModal(code)
}

// NewDialog creates a new standard dialog. To show the dialog you must call
// .RunModal() on the returned window. If the dialog window cannot be created,
// nil will be returned.
func NewDialog(img *draw.Image, msgPanel *ux.Panel, buttonInfo []*ButtonInfo) (*Dialog, error) {
	d := &Dialog{buttons: make(map[int]*buttonData)}
	var frame geom.Rect
	if focused := ux.WindowWithFocus(); focused != nil {
		frame = focused.FrameRect()
	} else {
		frame = display.Primary().Usable
	}
	d.wnd, d.err = ux.NewWindow("", frame, ux.TitledWindowMask|ux.ResizableWindowMask)
	if d.err != nil {
		return nil, errs.NewWithCause("unable to create dialog", d.err)
	}
	content := d.wnd.Content()
	content.SetBorder(border.NewEmpty(geom.NewUniformInsets(16)))
	columns := 1
	if img != nil {
		columns++
		icon := label.New().SetImage(img).SetBorder(border.NewEmpty(geom.Insets{Bottom: 16, Right: 8}))
		flex.NewData().VAlign(align.Start).Apply(icon)
		content.AddChild(icon.AsPanel())
	}
	flex.New().Columns(columns).Apply(content)
	if b := msgPanel.Border(); b != nil {
		msgPanel.SetBorder(border.NewCompound(border.NewEmpty(geom.Insets{Bottom: 16}), b))
	} else {
		msgPanel.SetBorder(border.NewEmpty(geom.Insets{Bottom: 16}))
	}
	flex.NewData().HGrab(true).VGrab(true).HAlign(align.Fill).VAlign(align.Start).Apply(msgPanel)
	content.AddChild(msgPanel)
	buttonPanel := ux.NewPanel()
	flex.New().Columns(len(buttonInfo) + 1).EqualColumns(true).HSpacing(layout.DefaultHSpacing * 2).Apply(buttonPanel)
	buttonPanel.AddChild(ux.NewPanel())
	for _, bi := range buttonInfo {
		b := bi.NewButton(d)
		d.buttons[bi.ResponseCode] = &buttonData{
			info:   bi,
			button: b,
		}
		buttonPanel.AddChild(b.AsPanel())
	}
	flex.NewData().HAlign(align.End).HSpan(columns).Apply(buttonPanel)
	content.AddChild(buttonPanel)
	originalKeyDownCallback := content.KeyDownCallback
	content.KeyDownCallback = func(keyCode int, ch rune, mod keys.Modifiers, repeat bool) bool {
		if originalKeyDownCallback == nil || !originalKeyDownCallback(keyCode, ch, mod, repeat) {
			if mod&keys.NonStickyModifiers == 0 {
				for _, one := range d.buttons {
					for _, kc := range one.info.KeyCode {
						if kc.Code == keyCode {
							if one.button.Enabled() {
								one.button.Click()
							}
							return true
						}
					}
				}
			}
			return false
		}
		return true
	}
	d.wnd.Pack()
	wndFrame := d.wnd.FrameRect()
	frame.Y += (frame.Height - wndFrame.Height) / 3
	frame.Height = wndFrame.Height
	frame.X += (frame.Width - wndFrame.Width) / 2
	frame.Width = wndFrame.Width
	frame.Align()
	d.wnd.SetFrameRect(frame)
	return d, nil
}

// NewMessagePanel creates a new panel containing the given primary and detail
// messages. Embedded line feeds are OK.
func NewMessagePanel(primary, detail string) *ux.Panel {
	panel := ux.NewPanel()
	flex.New().Apply(panel)
	breakTextIntoLabels(panel, primary, draw.EmphasizedSystemFont)
	breakTextIntoLabels(panel, detail, draw.SystemFont)
	flex.NewData().MinSize(geom.Size{Width: 200}).Apply(panel)
	return panel
}

func breakTextIntoLabels(panel *ux.Panel, text string, font *draw.Font) {
	if text != "" {
		returns := 0
		for {
			if i := strings.Index(text, "\n"); i != -1 {
				if i == 0 {
					returns++
					text = text[1:]
				} else {
					part := text[:i]
					l := label.New().SetText(part).SetFont(font)
					if returns > 1 {
						l.SetBorder(border.NewEmpty(geom.Insets{Top: 8}))
					}
					panel.AddChild(l.AsPanel())
					text = text[i+1:]
					returns = 1
				}
			} else {
				if text != "" {
					l := label.New().SetText(text).SetFont(font)
					if returns > 1 {
						l.SetBorder(border.NewEmpty(geom.Insets{Top: 8}))
					}
					panel.AddChild(l.AsPanel())
				}
				break
			}
		}
	}
}

// ErrorDialogWithMessage displays a standard error dialog with the specified
// primary and detail messages. Embedded line feeds are OK.
func ErrorDialogWithMessage(primary, detail string) {
	ErrorDialogWithPanel(NewMessagePanel(primary, detail))
}

// ErrorDialogWithPanel displays a standard error dialog with the specified
// panel.
func ErrorDialogWithPanel(msgPanel *ux.Panel) {
	if dialog, err := NewDialog(icons.Error(), msgPanel, []*ButtonInfo{NewOKButtonInfo()}); err != nil {
		jot.Error(err)
	} else {
		dialog.RunModal()
	}
}

// QuestionDialog displays a standard question dialog with the specified
// primary and detail messages. Embedded line feeds are OK. This function
// returns ids.ModalResponseOK if the OK button was pressed and
// ids.ModalResponseCancel if the Cancel button was pressed.
func QuestionDialog(primary, detail string) int {
	return QuestionDialogWithPanel(NewMessagePanel(primary, detail))
}

// QuestionDialogWithPanel displays a standard question dialog with the
// specified panel. This function returns ids.ModalResponseOK if the OK button
// was pressed and ids.ModalResponseCancel if the Cancel button was pressed.
func QuestionDialogWithPanel(msgPanel *ux.Panel) int {
	buttonPanel := ux.NewPanel()
	flex.New().Columns(2).EqualColumns(true).HSpacing(layout.DefaultHSpacing * 2).Apply(buttonPanel)
	cancelButton := button.New().SetText(i18n.Text("Cancel"))
	buttonPanel.AddChild(cancelButton.AsPanel())
	okButton := button.New().SetText(i18n.Text("OK"))
	buttonPanel.AddChild(okButton.AsPanel())
	for _, p := range buttonPanel.Children() {
		flex.NewData().HAlign(align.Fill).Apply(p)
	}
	if dialog, err := NewDialog(icons.Question(), msgPanel, []*ButtonInfo{NewCancelButtonInfo(), NewOKButtonInfo()}); err != nil {
		jot.Error(err)
	} else {
		return dialog.RunModal()
	}
	return ids.ModalResponseCancel
}

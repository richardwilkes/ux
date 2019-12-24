package dialog

import (
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/ux/ids"
	"github.com/richardwilkes/ux/keys"
	"github.com/richardwilkes/ux/layout/align"
	"github.com/richardwilkes/ux/layout/flex"
	"github.com/richardwilkes/ux/widget/button"
)

// ButtonInfo holds information for constructing the dialog button panel.
type ButtonInfo struct {
	Title        string
	ResponseCode int
	KeyCode      []*keys.Key
}

// NewButton creates a new button for the dialog.
func (bi *ButtonInfo) NewButton(d *Dialog) *button.Button {
	b := button.New().SetText(bi.Title)
	b.ClickCallback = func() { d.StopModal(bi.ResponseCode) }
	flex.NewData().HAlign(align.Fill).Apply(b)
	return b
}

// NewCancelButtonInfo creates a standard cancel button.
func NewCancelButtonInfo() *ButtonInfo {
	return &ButtonInfo{
		Title:        i18n.Text("Cancel"),
		ResponseCode: ids.ModalResponseCancel,
		KeyCode:      []*keys.Key{keys.Escape},
	}
}

// NewOKButtonInfo creates a standard OK button.
func NewOKButtonInfo() *ButtonInfo {
	return NewOKButtonInfoWithTitle(i18n.Text("OK"))
}

// NewOKButtonInfoWithTitle creates a standard OK button with a specific
// title.
func NewOKButtonInfoWithTitle(title string) *ButtonInfo {
	return &ButtonInfo{
		Title:        title,
		ResponseCode: ids.ModalResponseOK,
		KeyCode:      []*keys.Key{keys.Return, keys.NumpadEnter},
	}
}

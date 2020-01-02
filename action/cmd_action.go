// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package action

import (
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/ids"
	"github.com/richardwilkes/ux/keys"
)

var (
	// Cut removes the selection and places it on the clipboard.
	Cut = &CmdAction{
		ActionID:        ids.CutItemID,
		ActionTitle:     i18n.Text("Cut"),
		ActionHotKey:    keys.X,
		ActionModifiers: keys.OSMenuCmdModifier(),
	}
	// Copy the selection and place it on the clipboard.
	Copy = &CmdAction{
		ActionID:        ids.CopyItemID,
		ActionTitle:     i18n.Text("Copy"),
		ActionHotKey:    keys.C,
		ActionModifiers: keys.OSMenuCmdModifier(),
	}
	// Paste the contents of the clipboard, replacing the selection.
	Paste = &CmdAction{
		ActionID:        ids.PasteItemID,
		ActionTitle:     i18n.Text("Paste"),
		ActionHotKey:    keys.V,
		ActionModifiers: keys.OSMenuCmdModifier(),
	}
	// Delete the selection.
	Delete = &CmdAction{
		ActionID:     ids.DeleteItemID,
		ActionTitle:  i18n.Text("Delete"),
		ActionHotKey: keys.Backspace,
	}
	// SelectAll selects everything in the current focus.
	SelectAll = &CmdAction{
		ActionID:        ids.SelectAllItemID,
		ActionTitle:     i18n.Text("Select All"),
		ActionHotKey:    keys.A,
		ActionModifiers: keys.OSMenuCmdModifier(),
	}
)

// CmdAction provides a standardized way to issue commands to focused UI
// elements.
type CmdAction struct {
	ActionID        int
	ActionTitle     string
	ActionHotKey    *keys.Key
	ActionModifiers keys.Modifiers
}

var _ Action = &CmdAction{}

// ID implements action.Action.
func (a *CmdAction) ID() int {
	return a.ActionID
}

// Title implements action.Action.
func (a *CmdAction) Title() string {
	return a.ActionTitle
}

// HotKey implements action.Action.
func (a *CmdAction) HotKey() *keys.Key {
	return a.ActionHotKey
}

// HotKeyModifiers implements action.Action.
func (a *CmdAction) HotKeyModifiers() keys.Modifiers {
	return a.ActionModifiers
}

// Enabled implements action.Action.
func (a *CmdAction) Enabled(source interface{}) bool {
	if wnd := ux.WindowWithFocus(); wnd != nil {
		focus := wnd.Focus()
		return focus != nil && focus.CanPerformCmdCallback != nil && focus.CanPerformCmdCallback(source, a.ActionID)
	}
	return false
}

// Execute implements action.Action.
func (a *CmdAction) Execute(source interface{}) {
	if wnd := ux.WindowWithFocus(); wnd != nil {
		if focus := wnd.Focus(); focus != nil && focus.PerformCmdCallback != nil {
			focus.PerformCmdCallback(source, a.ActionID)
		}
	}
}

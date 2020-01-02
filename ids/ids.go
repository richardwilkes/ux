// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package ids

// Pre-defined menu IDs. Apps should start their IDs at UserBaseID.
const (
	AppMenuID = 1 + iota
	FileMenuID
	EditMenuID
	WindowMenuID
	HelpMenuID
	ServicesMenuID
	AboutItemID
	PreferencesItemID
	QuitItemID
	CutItemID
	CopyItemID
	PasteItemID
	DeleteItemID
	SelectAllItemID
	MinimizeItemID
	ZoomItemID
	BringAllWindowsToFrontItemID
	CloseItemID
	HideItemID
	HideOthersItemID
	ShowAllItemID
	PopupMenuTemporaryBaseID
	UserBaseID        = 1000
	MaxUserBaseID     = 1<<30 - 1
	ContextMenuIDFlag = 1 << 30 // Should be or'd into IDs for context menus
)

// Pre-defined modal response codes. Apps should start their codes at
// ModalResponseUserBase.
const (
	ModalResponseCancel = iota
	ModalResponseOK
	ModalResponseDiscard
	ModalResponseUserBase = 1000
)

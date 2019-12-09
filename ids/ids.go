package ids

// Pre-defined menu IDs. Apps should start their IDs at UserBaseID.
const (
	BarID = 1 + iota
	AppMenuID
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

package savedialog

type osSaveDialog = int

func osNew() osSaveDialog {
	return 0 // TODO: Implement
}

func (d *SaveDialog) osDirectoryURL() string {
	return "" // TODO: Implement
}

func (d *SaveDialog) osSetDirectoryURL(dirURL string) {
	// TODO: Implement
}

func (d *SaveDialog) osAllowedFileTypes() []string {
	return nil // TODO: Implement
}

func (d *SaveDialog) osSetAllowedFileTypes(allowed []string) {
	// TODO: Implement
}

func (d *SaveDialog) osURL() string {
	return "" // TODO: Implement
}

func (d *SaveDialog) osRunModal() bool {
	return false // TODO: Implement
}

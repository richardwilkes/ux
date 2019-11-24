package opendialog

type osOpenDialog = int

func osNew() osOpenDialog {
	return 0 // TODO: Implement
}

func (d *OpenDialog) osDirectoryURL() string {
	// TODO: Implement
	return ""
}

func (d *OpenDialog) osSetDirectoryURL(dirURL string) {
	// TODO: Implement
}

func (d *OpenDialog) osAllowedFileTypes() []string {
	return nil // TODO: Implement
}

func (d *OpenDialog) osSetAllowedFileTypes(allowed []string) {
	// TODO: Implement
}

func (d *OpenDialog) osCanChooseFiles() bool {
	return true // TODO: Implement
}

func (d *OpenDialog) osSetCanChooseFiles(canChoose bool) {
	// TODO: Implement
}

func (d *OpenDialog) osCanChooseDirectories() bool {
	return false // TODO: Implement
}

func (d *OpenDialog) osSetCanChooseDirectories(canChoose bool) {
	// TODO: Implement
}

func (d *OpenDialog) osResolvesAliases() bool {
	return true // TODO: Implement
}

func (d *OpenDialog) osSetResolvesAliases(resolves bool) {
	// TODO: Implement
}

func (d *OpenDialog) osAllowsMultipleSelection() bool {
	return false // TODO: Implement
}

func (d *OpenDialog) osSetAllowsMultipleSelection(allow bool) {
	// TODO: Implement
}

func (d *OpenDialog) osURLs() []string {
	return nil // TODO: Implement
}

func (d *OpenDialog) osRunModal() bool {
	return false // TODO: Implement
}

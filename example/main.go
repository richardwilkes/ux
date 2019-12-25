package main

import (
	"fmt"
	"runtime"
	"strings"
	"unicode"

	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/log/jotrotate"
	"github.com/richardwilkes/toolbox/xio/fs/embedded"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/display"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/layout/align"
	"github.com/richardwilkes/ux/layout/flex"
	"github.com/richardwilkes/ux/layout/flow"
	"github.com/richardwilkes/ux/menu"
	"github.com/richardwilkes/ux/widget/browser"
	"github.com/richardwilkes/ux/widget/button"
	"github.com/richardwilkes/ux/widget/checkbox"
	"github.com/richardwilkes/ux/widget/checkbox/state"
	"github.com/richardwilkes/ux/widget/label"
	"github.com/richardwilkes/ux/widget/list"
	"github.com/richardwilkes/ux/widget/popupmenu"
	"github.com/richardwilkes/ux/widget/radiobutton"
	"github.com/richardwilkes/ux/widget/scrollarea"
	"github.com/richardwilkes/ux/widget/scrollarea/behavior"
	"github.com/richardwilkes/ux/widget/selectable"
	"github.com/richardwilkes/ux/widget/separator"
	"github.com/richardwilkes/ux/widget/textfield"
	"github.com/richardwilkes/ux/widget/tooltip"
)

var (
	aboutWindow         *ux.Window
	appleCursor         *draw.Cursor
	homeImg             *draw.Image
	classicAppleLogoImg *draw.Image
	mountainsImg        *draw.Image
)

func main() {
	runtime.LockOSThread()

	cmdline.AppName = "Example"
	cmdline.AppCmdName = "example"
	cmdline.AppVersion = "0.1"
	cmdline.CopyrightYears = "2019"
	cmdline.CopyrightHolder = "Richard A. Wilkes"
	cmdline.AppIdentifier = "com.trollworks.ux.example"

	cl := cmdline.New(true)
	jotrotate.ParseAndSetup(cl)

	ux.WillFinishStartupCallback = finishStartup
	ux.QuittingCallback = func() { jot.Info("Quitting") }
	ux.Start() // Never returns
}

func finishStartup() {
	fs := embedded.NewFileSystemFromEmbeddedZip("example/images")
	var err error
	homeImg, err = draw.NewImageFromBytes(fs.MustContentAsBytes("home.png"), 0.5)
	jot.FatalIfErr(err)
	classicAppleLogoImg, err = draw.NewImageFromBytes(fs.MustContentAsBytes("classic-apple-logo.png"), 0.5)
	jot.FatalIfErr(err)
	mountainsImg, err = draw.NewImageFromBytes(fs.MustContentAsBytes("mountains.jpg"), 0.5)
	jot.FatalIfErr(err)
	size := classicAppleLogoImg.LogicalGeomSize()
	appleCursor = draw.NewCursor(classicAppleLogoImg, geom.Point{
		X: size.Width / 2,
		Y: size.Height / 2,
	})

	usable := display.Primary().Usable
	w1 := createButtonsWindow("Demo #1", usable.Point)
	frame1 := w1.FrameRect()
	createButtonsWindow("Demo #2", geom.Point{X: frame1.X + frame1.Width, Y: frame1.Y})
}

func createButtonsWindow(title string, where geom.Point) *ux.Window {
	wnd, err := ux.NewWindow(title, geom.Rect{}, ux.StdWindowMask)
	jot.FatalIfErr(err)
	if bar, global, first := menu.BarForWindow(wnd, nil); !global || first {
		bar.InsertStdMenus(createAboutWindow, createPreferencesWindow, nil)
	}

	content := wnd.Content()
	content.SetBorder(border.NewEmpty(geom.NewUniformInsets(10)))
	flex.New().VSpacing(10).Apply(content)

	buttonsPanel := createButtonsPanel()
	flexData := flex.NewData().HGrab(true)
	flexData.Apply(buttonsPanel)
	content.AddChild(buttonsPanel)

	addSeparator(content)

	checkBoxPanel := createCheckBoxPanel()
	flexData.Apply(checkBoxPanel)
	content.AddChild(checkBoxPanel)

	addSeparator(content)

	toggleButtonsPanel := createToggleButtonsPanel()
	flexData.Apply(toggleButtonsPanel)
	content.AddChild(toggleButtonsPanel)

	addSeparator(content)

	radioButtonsPanel := createRadioButtonsPanel()
	flexData.Apply(radioButtonsPanel)
	content.AddChild(radioButtonsPanel)

	addSeparator(content)

	popupMenusPanel := createPopupMenusPanel()
	flexData.Apply(popupMenusPanel)
	content.AddChild(popupMenusPanel)

	addSeparator(content)

	wrapper := ux.NewPanel()
	flex.New().Columns(2).EqualColumns(true).HSpacing(10).Apply(wrapper)
	flexData.HAlign(align.Fill).Apply(wrapper)
	wrapper.SetLayoutData(flexData)
	textFieldsPanel := createTextFieldsPanel()
	flexData.Apply(textFieldsPanel)
	wrapper.AddChild(textFieldsPanel)
	wrapper.AddChild(createListPanel())
	content.AddChild(wrapper)

	addSeparator(content)

	if title == "Demo #1" {
		var b *browser.Browser
		if b, err = browser.New(wnd); err == nil {
			flex.NewData().HAlign(align.Fill).VAlign(align.Fill).HGrab(true).VGrab(true).SizeHint(geom.Size{Width: 1024, Height: 768}).Apply(b)
			b.LoadURL("https://gurpscharactersheet.com")
			content.AddChild(b.AsPanel())
		}
	} else {
		imgPanel := label.New().SetImage(mountainsImg)
		imgPanel.SetFocusable(true)
		_, prefSize, _ := imgPanel.Sizes(geom.Size{})
		imgPanel.SetFrameRect(geom.Rect{Size: prefSize})
		imgPanel.UpdateCursorCallback = func(where geom.Point) *draw.Cursor {
			return appleCursor
		}
		imgPanel.UpdateTooltipCallback = func(where geom.Point, avoid geom.Rect) geom.Rect {
			imgPanel.Tooltip = tooltip.NewWithText(where.String())
			avoid.X = where.X - 16
			avoid.Y = where.Y - 16
			avoid.Point = imgPanel.PointToRoot(avoid.Point)
			avoid.Width = 32
			avoid.Height = 32
			return avoid
		}
		scrollArea := scrollarea.New().SetContent(imgPanel.AsPanel(), behavior.Unmodified)
		flex.NewData().HAlign(align.Fill).VAlign(align.Fill).HGrab(true).VGrab(true).Apply(scrollArea)
		content.AddChild(scrollArea.AsPanel())
	}

	wnd.SetFocus(textFieldsPanel.Children()[0])
	wnd.Pack()
	rect := wnd.FrameRect()
	rect.Point = where
	wnd.SetFrameRect(rect)
	wnd.ToFront()
	return wnd
}

func createListPanel() *ux.Panel {
	lst := list.New()
	lst.Append(
		"One",
		"Two",
		"Three with some long text to make it interesting",
		"Four",
		"Five",
	)
	lst.NewSelectionCallback = func() {
		var buffer strings.Builder
		fmt.Fprintf(&buffer, "Selection changed in %v. Now:", lst)
		index := -1
		first := true
		for {
			index = lst.Selection.NextSet(index + 1)
			if index == -1 {
				break
			}
			if first {
				first = false
			} else {
				buffer.WriteString(",")
			}
			fmt.Fprintf(&buffer, " %d", index)
		}
		jot.Info(buffer.String())
	}
	lst.DoubleClickCallback = func() {
		jot.Infof("Double-clicked on %v", lst)
	}
	_, prefSize, _ := lst.Sizes(geom.Size{})
	lst.SetFrameRect(geom.Rect{Size: prefSize})
	scroller := scrollarea.New().SetContent(lst.AsPanel(), behavior.Fill)
	flex.NewData().HAlign(align.Fill).VAlign(align.Fill).HGrab(true).VGrab(true).Apply(scroller)
	return scroller.AsPanel()
}

func addSeparator(parent *ux.Panel) {
	sep := separator.NewHorizontal()
	flex.NewData().HAlign(align.Fill).Apply(sep)
	parent.AddChild(sep.AsPanel())
}

func createButtonsPanel() *ux.Panel {
	panel := ux.NewPanel()
	flow.New().HSpacing(5).VSpacing(5).Apply(panel)
	btn := createButton("Press Me", panel)
	btn.SetLayoutData(align.Middle)
	btn = createButton("Disabled", panel)
	btn.SetLayoutData(align.Middle)
	btn.SetEnabled(false)
	btn = createImageButton(homeImg, panel)
	btn.SetLayoutData(align.Middle)
	btn = createImageButton(homeImg, panel)
	btn.SetLayoutData(align.Middle)
	btn.SetEnabled(false)
	btn = createImageButton(classicAppleLogoImg, panel)
	btn.SetLayoutData(align.Middle)
	btn = createImageButton(classicAppleLogoImg, panel)
	btn.SetLayoutData(align.Middle)
	btn.SetEnabled(false)
	return panel
}

func createButton(title string, panel *ux.Panel) *button.Button {
	btn := button.New().SetText(title)
	btn.ClickCallback = func() { jot.Infof("%v was clicked.", btn) }
	btn.Tooltip = tooltip.NewWithText(fmt.Sprintf("This is the tooltip for %v", btn))
	panel.AddChild(btn.AsPanel())
	return btn
}

func createImageButton(img *draw.Image, panel *ux.Panel) *button.Button {
	btn := button.New().SetImage(img)
	btn.ClickCallback = func() { jot.Infof("%v was clicked.", btn) }
	btn.Tooltip = tooltip.NewWithText(fmt.Sprintf("This is the tooltip for %v", btn))
	panel.AddChild(btn.AsPanel())
	return btn
}

func createCheckBoxPanel() *ux.Panel {
	panel := ux.NewPanel()
	flex.New().Apply(panel)
	createCheckBox("Press Me", panel)
	createCheckBox("Initially Mixed", panel).SetState(state.Mixed)
	createCheckBox("Disabled", panel).SetEnabled(false)
	createCheckBox("Disabled w/Check", panel).SetState(state.Checked).SetEnabled(false)
	return panel
}

func createCheckBox(title string, panel *ux.Panel) *checkbox.CheckBox {
	check := checkbox.New().SetText(title)
	check.ClickCallback = func() { jot.Infof("%v was clicked.", check) }
	check.Tooltip = tooltip.NewWithText(fmt.Sprintf("This is the tooltip for %v", check))
	panel.AddChild(check.AsPanel())
	return check
}

func createToggleButtonsPanel() *ux.Panel {
	panel := ux.NewPanel()
	flow.New().HSpacing(5).VSpacing(5).Apply(panel)
	group := selectable.NewGroup()
	first := createToggleButton(homeImg, panel, group)
	createToggleButton(classicAppleLogoImg, panel, group)
	group.Select(first.AsSelectable())
	return panel
}

func createToggleButton(img *draw.Image, panel *ux.Panel, group *selectable.Group) *button.Button {
	btn := createImageButton(img, panel).SetSticky(true)
	btn.SetLayoutData(align.Middle)
	group.Add(btn.AsSelectable())
	return btn
}

func createRadioButtonsPanel() *ux.Panel {
	panel := ux.NewPanel()
	flex.New().Apply(panel)
	group := selectable.NewGroup()
	first := createRadioButton("First", panel, group)
	createRadioButton("Second", panel, group)
	createRadioButton("Third (disabled)", panel, group).SetEnabled(false)
	createRadioButton("Fourth", panel, group)
	group.Select(first.AsSelectable())
	return panel
}

func createRadioButton(title string, panel *ux.Panel, group *selectable.Group) *radiobutton.RadioButton {
	rb := radiobutton.New().SetText(title)
	rb.ClickCallback = func() { jot.Infof("%v was clicked.", rb) }
	rb.Tooltip = tooltip.NewWithText(fmt.Sprintf("This is the tooltip for %v", rb))
	panel.AddChild(rb.AsPanel())
	group.Add(rb.AsSelectable())
	return rb
}

func createPopupMenusPanel() *ux.Panel {
	panel := ux.NewPanel()
	flex.New().Apply(panel)
	createPopupMenu(panel, 1, "One", "Two", "Three", "", "Four", "Five", "Six")
	createPopupMenu(panel, 2, "Red", "Blue", "Green").SetEnabled(false)
	return panel
}

func createPopupMenu(panel *ux.Panel, selection int, titles ...string) *popupmenu.PopupMenu {
	p := popupmenu.New()
	p.Tooltip = tooltip.NewWithText(fmt.Sprintf("This is the tooltip for %v", p))
	for _, title := range titles {
		if title == "" {
			p.AddSeparator()
		} else {
			p.AddItem(title)
		}
	}
	p.SelectIndex(selection)
	p.SelectionCallback = func() { jot.Infof("The '%v' item was selected from the PopupMenu.", p.Selected()) }
	panel.AddChild(p.AsPanel())
	return p
}

func createTextFieldsPanel() *ux.Panel {
	panel := ux.NewPanel()
	flex.New().Apply(panel)
	field := createTextField("First Text Field", panel)
	createTextField("Second Text Field (disabled)", panel).SetEnabled(false)
	createTextField("", panel).SetWatermark("Watermarked")
	field = createTextField("", panel).SetWatermark("Enter only numbers")
	field.ValidateCallback = func() bool {
		for _, r := range field.Text() {
			if !unicode.IsDigit(r) {
				return false
			}
		}
		return true
	}
	return panel
}

func createTextField(text string, panel *ux.Panel) *textfield.TextField {
	field := textfield.New()
	field.SetText(text)
	flex.NewData().HAlign(align.Fill).HGrab(true).Apply(field)
	field.Tooltip = tooltip.NewWithText(fmt.Sprintf("This is the tooltip for %v", field))
	panel.AddChild(field.AsPanel())
	return field
}

func createAboutWindow() {
	if aboutWindow == nil {
		var err error
		aboutWindow, err = ux.NewWindow("About "+cmdline.AppName, geom.Rect{}, ux.TitledWindowMask|ux.ClosableWindowMask)
		if err != nil {
			jot.Error(err)
			return
		}
		aboutWindow.WillCloseCallback = func() { aboutWindow = nil }
		content := aboutWindow.Content()
		content.SetBorder(border.NewEmpty(geom.NewUniformInsets(10)))
		flex.New().Apply(content)
		title := label.New().SetText(cmdline.AppName).SetFont(draw.EmphasizedSystemFont)
		flexData := flex.NewData().HAlign(align.Fill).HGrab(true)
		flexData.Apply(title)
		content.AddChild(title.AsPanel())
		desc := label.New().SetText("Simple app to demonstrate the\ncapabilities of the ui framework.")
		flexData.Apply(desc)
		content.AddChild(desc.AsPanel())
		aboutWindow.Pack()
	}
	aboutWindow.ToFront()
}

func createPreferencesWindow() {
	jot.Info("Preferences...")
}

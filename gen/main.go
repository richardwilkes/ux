package main

//go:generate go run main.go

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/txt"
)

const (
	typeBool      = "bool"
	typeFloat64   = "float64"
	typeString    = "string"
	typeImage     = "*draw.Image"
	typeFont      = "*draw.Font"
	typeInk       = "draw.Ink"
	typeAlignment = "align.Alignment"
	typeSide      = "side.Side"
	typeDuration  = "time.Duration"
)

// Var holds information about a variable we want to manage on a widget.
type Var struct {
	Name            string
	Type            string
	Default         string
	Comment         string
	UseDefaultIfNil bool
	MinOfZero       bool
	Redraw          bool
	Layout          bool
}

type widgetVars struct {
	Name       string
	Instance   string
	Vars       []*Var
	Selectable bool
}

var widgetList = []*widgetVars{
	{
		Name:       "Button",
		Instance:   "b",
		Selectable: true,
		Vars: []*Var{
			{
				Name:    "image",
				Type:    typeImage,
				Comment: "the image. May be nil",
				Redraw:  true,
				Layout:  true,
			},
			{
				Name:    "text",
				Type:    typeString,
				Comment: "the text content",
				Redraw:  true,
				Layout:  true,
			},
			{
				Name:            "font",
				Type:            typeFont,
				Default:         "draw.SystemFont",
				Comment:         "the font that will be used when drawing text content",
				UseDefaultIfNil: true,
				Redraw:          true,
				Layout:          true,
			},
			{
				Name:            "backgroundInk",
				Type:            typeInk,
				Default:         "draw.ControlBackgroundInk",
				Comment:         "the ink that will be used for the background when enabled but not selected, pressed or focused",
				UseDefaultIfNil: true,
				Redraw:          true,
			},
			{
				Name:            "selectedBackgroundInk",
				Type:            typeInk,
				Default:         "draw.ControlSelectedBackgroundInk",
				Comment:         "the ink that will be used for the background when enabled and selected, but not pressed or focused",
				UseDefaultIfNil: true,
				Redraw:          true,
			},
			{
				Name:            "focusedBackgroundInk",
				Type:            typeInk,
				Default:         "draw.ControlFocusedBackgroundInk",
				Comment:         "the ink that will be used for the background when enabled and focused",
				UseDefaultIfNil: true,
				Redraw:          true,
			},
			{
				Name:            "pressedBackgroundInk",
				Type:            typeInk,
				Default:         "draw.ControlPressedBackgroundInk",
				Comment:         "the ink that will be used for the background when enabled and pressed",
				UseDefaultIfNil: true,
				Redraw:          true,
			},
			{
				Name:            "edgeInk",
				Type:            typeInk,
				Default:         "draw.ControlEdgeAdjColor",
				Comment:         "the ink that will be used for the edges",
				UseDefaultIfNil: true,
				Redraw:          true,
			},
			{
				Name:            "textInk",
				Type:            typeInk,
				Default:         "draw.ControlTextColor",
				Comment:         "the ink that will be used for the text when disabled or not pressed",
				UseDefaultIfNil: true,
				Redraw:          true,
			},
			{
				Name:            "pressedTextInk",
				Type:            typeInk,
				Default:         "draw.AlternateSelectedControlTextColor",
				Comment:         "the ink that will be used for the text when enabled and pressed",
				UseDefaultIfNil: true,
				Redraw:          true,
			},
			{
				Name:      "gap",
				Type:      typeFloat64,
				Default:   "3",
				Comment:   "the gap to put between the image and text",
				MinOfZero: true,
				Redraw:    true,
				Layout:    true,
			},
			{
				Name:      "cornerRadius",
				Type:      typeFloat64,
				Default:   "4",
				Comment:   "the amount of rounding to use on the corners",
				MinOfZero: true,
				Redraw:    true,
			},
			{
				Name:      "hMargin",
				Type:      typeFloat64,
				Default:   "8",
				Comment:   "the margin on the left and right side of the content",
				MinOfZero: true,
				Redraw:    true,
			},
			{
				Name:      "vMargin",
				Type:      typeFloat64,
				Default:   "1",
				Comment:   "the margin on the top and bottom side of the content",
				MinOfZero: true,
				Redraw:    true,
			},
			{
				Name:      "imageOnlyHMargin",
				Type:      typeFloat64,
				Default:   "3",
				Comment:   "the margin on the left and right side of the content when only an image is present",
				MinOfZero: true,
				Redraw:    true,
			},
			{
				Name:      "imageOnlyVMargin",
				Type:      typeFloat64,
				Default:   "3",
				Comment:   "the margin on the top and bottom side of the content when only an image is present",
				MinOfZero: true,
				Redraw:    true,
			},
			{
				Name:      "clickAnimationTime",
				Type:      typeDuration,
				Default:   "time.Millisecond * 100",
				Comment:   "the amount of time to spend animating the click action",
				MinOfZero: true,
			},
			{
				Name:    "hAlign",
				Type:    typeAlignment,
				Default: "align.Middle",
				Comment: "the horizontal alignment",
				Redraw:  true,
			},
			{
				Name:    "vAlign",
				Type:    typeAlignment,
				Default: "align.Middle",
				Comment: "the vertical alignment",
				Redraw:  true,
			},
			{
				Name:    "side",
				Type:    typeSide,
				Default: "side.Left",
				Comment: "the side of the text the image should be on",
				Redraw:  true,
			},
			{
				Name:    "sticky",
				Type:    typeBool,
				Comment: "whether the button will visually retain its selected state",
				Redraw:  true,
			},
		},
	},
	{
		Name:     "CheckBox",
		Instance: "c",
		Vars: []*Var{
			{
				Name:    "image",
				Type:    typeImage,
				Comment: "the image. May be nil",
				Redraw:  true,
				Layout:  true,
			},
			{
				Name:    "text",
				Type:    typeString,
				Comment: "the text content",
				Redraw:  true,
				Layout:  true,
			},
			{
				Name:            "font",
				Type:            typeFont,
				Default:         "draw.SystemFont",
				Comment:         "the font that will be used when drawing text content",
				UseDefaultIfNil: true,
				Redraw:          true,
				Layout:          true,
			},
			{
				Name:            "backgroundInk",
				Type:            typeInk,
				Default:         "draw.ControlBackgroundInk",
				Comment:         "the ink that will be used for the background when enabled but not pressed or focused",
				UseDefaultIfNil: true,
				Redraw:          true,
			},
			{
				Name:            "focusedBackgroundInk",
				Type:            typeInk,
				Default:         "draw.ControlFocusedBackgroundInk",
				Comment:         "the ink that will be used for the background when enabled and focused",
				UseDefaultIfNil: true,
				Redraw:          true,
			},
			{
				Name:            "pressedBackgroundInk",
				Type:            typeInk,
				Default:         "draw.ControlPressedBackgroundInk",
				Comment:         "the ink that will be used for the background when enabled and pressed",
				UseDefaultIfNil: true,
				Redraw:          true,
			},
			{
				Name:            "edgeInk",
				Type:            typeInk,
				Default:         "draw.ControlEdgeAdjColor",
				Comment:         "the ink that will be used for the edges",
				UseDefaultIfNil: true,
				Redraw:          true,
			},
			{
				Name:            "textInk",
				Type:            typeInk,
				Default:         "draw.ControlTextColor",
				Comment:         "the ink that will be used for the text when disabled or not pressed",
				UseDefaultIfNil: true,
				Redraw:          true,
			},
			{
				Name:            "pressedTextInk",
				Type:            typeInk,
				Default:         "draw.AlternateSelectedControlTextColor",
				Comment:         "the ink that will be used for the text when enabled and pressed",
				UseDefaultIfNil: true,
				Redraw:          true,
			},
			{
				Name:      "gap",
				Type:      typeFloat64,
				Default:   "3",
				Comment:   "the gap to put between the checkbox, image and text",
				MinOfZero: true,
				Redraw:    true,
				Layout:    true,
			},
			{
				Name:      "cornerRadius",
				Type:      typeFloat64,
				Default:   "4",
				Comment:   "the amount of rounding to use on the corners",
				MinOfZero: true,
				Redraw:    true,
			},
			{
				Name:      "clickAnimationTime",
				Type:      typeDuration,
				Default:   "time.Millisecond * 100",
				Comment:   "the amount of time to spend animating the click action",
				MinOfZero: true,
			},
			{
				Name:    "hAlign",
				Type:    typeAlignment,
				Default: "align.Start",
				Comment: "the horizontal alignment",
				Redraw:  true,
			},
			{
				Name:    "vAlign",
				Type:    typeAlignment,
				Default: "align.Middle",
				Comment: "the vertical alignment",
				Redraw:  true,
			},
			{
				Name:    "side",
				Type:    typeSide,
				Default: "side.Left",
				Comment: "the side of the text the image should be on",
				Redraw:  true,
			},
		},
	},
	{
		Name:     "InkWell",
		Instance: "well",
		Vars: []*Var{
			{
				Name:            "backgroundInk",
				Type:            typeInk,
				Default:         "draw.ControlBackgroundInk",
				Comment:         "the ink that will be used for the background when enabled but not pressed or focused",
				UseDefaultIfNil: true,
				Redraw:          true,
			},
			{
				Name:            "focusedBackgroundInk",
				Type:            typeInk,
				Default:         "draw.ControlFocusedBackgroundInk",
				Comment:         "the ink that will be used for the background when enabled and focused",
				UseDefaultIfNil: true,
				Redraw:          true,
			},
			{
				Name:            "pressedBackgroundInk",
				Type:            typeInk,
				Default:         "draw.ControlPressedBackgroundInk",
				Comment:         "the ink that will be used for the background when enabled and pressed",
				UseDefaultIfNil: true,
				Redraw:          true,
			},
			{
				Name:            "edgeInk",
				Type:            typeInk,
				Default:         "draw.ControlEdgeAdjColor",
				Comment:         "the ink that will be used for the edges",
				UseDefaultIfNil: true,
				Redraw:          true,
			},
			{
				Name:            "edgeHighlightInk",
				Type:            typeInk,
				Default:         "draw.ControlEdgeHighlightAdjColor",
				Comment:         "the ink that will be used just inside the edges",
				UseDefaultIfNil: true,
				Redraw:          true,
			},
			{
				Name:      "imageScale",
				Type:      typeFloat64,
				Default:   "0.5",
				Comment:   "the image scale to use for images dropped onto the well. Defaults to 0.5 to support retina displays",
				MinOfZero: true,
				Redraw:    true,
			},
			{
				Name:      "contentSize",
				Type:      typeFloat64,
				Default:   "20",
				Comment:   "the content width and height",
				MinOfZero: true,
				Redraw:    true,
				Layout:    true,
			},
			{
				Name:      "cornerRadius",
				Type:      typeFloat64,
				Default:   "4",
				Comment:   "the amount of rounding to use on the corners",
				MinOfZero: true,
				Redraw:    true,
			},
			{
				Name:      "clickAnimationTime",
				Type:      typeDuration,
				Default:   "time.Millisecond * 100",
				Comment:   "the amount of time to spend animating the click action",
				MinOfZero: true,
			},
		},
	},
	{
		Name:     "Label",
		Instance: "l",
		Vars: []*Var{
			{
				Name:    "image",
				Type:    typeImage,
				Comment: "the image. May be nil",
				Redraw:  true,
				Layout:  true,
			},
			{
				Name:    "text",
				Type:    typeString,
				Comment: "the text content",
				Redraw:  true,
				Layout:  true,
			},
			{
				Name:            "font",
				Type:            typeFont,
				Default:         "draw.LabelFont",
				Comment:         "the font that will be used when drawing text content",
				UseDefaultIfNil: true,
				Redraw:          true,
				Layout:          true,
			},
			{
				Name:            "ink",
				Type:            typeInk,
				Default:         "draw.LabelColor",
				Comment:         "the ink that will be used when drawing text content",
				UseDefaultIfNil: true,
				Redraw:          true,
			},
			{
				Name:      "gap",
				Type:      typeFloat64,
				Default:   "3",
				Comment:   "the gap to put between the image and text",
				MinOfZero: true,
				Redraw:    true,
				Layout:    true,
			},
			{
				Name:    "hAlign",
				Type:    typeAlignment,
				Default: "align.Start",
				Comment: "the horizontal alignment",
				Redraw:  true,
			},
			{
				Name:    "vAlign",
				Type:    typeAlignment,
				Default: "align.Middle",
				Comment: "the vertical alignment",
				Redraw:  true,
			},
			{
				Name:    "side",
				Type:    typeSide,
				Default: "side.Left",
				Comment: "the side of the text the image should be on",
				Redraw:  true,
			},
		},
	},
}

func main() {
	for _, w := range widgetList {
		name := strings.ToLower(w.Name)
		processTemplate("widget", filepath.Join("..", "widget", name, name+"_gen.go"), w)
	}
}

func processTemplate(name, dstPath string, arg interface{}) {
	var buffer bytes.Buffer
	baseName := name + ".go.tmpl"
	fmt.Fprintf(&buffer, "// Code created from \"%s\" - don't edit by hand\n\n", baseName)
	tmpl, err := template.New(baseName).Funcs(template.FuncMap{
		"firstToLower": firstToLower,
		"firstToUpper": firstToUpper,
		"imports":      imports,
		"package":      pkg,
		"comment":      comment,
	}).ParseFiles(filepath.Join("tmpl", baseName))
	fatalIfErr(err)
	fatalIfErr(tmpl.Execute(&buffer, arg))
	var data []byte
	if data, err = format.Source(buffer.Bytes()); err != nil {
		jot.Warn(errs.NewWithCause(dstPath, err))
		data = buffer.Bytes()
	}
	fatalIfErr(ioutil.WriteFile(dstPath, data, 0644))
}

func pkg(w *widgetVars) string {
	return strings.ToLower(w.Name)
}

func imports(w *widgetVars) []string {
	sys := make(map[string]bool)
	usr := make(map[string]bool)
	usr["github.com/richardwilkes/ux/border"] = true
	for _, v := range w.Vars {
		switch v.Type {
		case typeImage, typeFont, typeInk:
			usr["github.com/richardwilkes/ux/draw"] = true
		case typeAlignment:
			usr["github.com/richardwilkes/ux/layout/align"] = true
		case typeSide:
			usr["github.com/richardwilkes/ux/layout/side"] = true
		case typeDuration:
			sys["time"] = true
		}
	}
	all := extractImports(sys)
	if len(all) > 0 && len(usr) > 0 {
		all = append(all, "")
	}
	return append(all, extractImports(usr)...)
}

func extractImports(m map[string]bool) []string {
	list := make([]string, 0, len(m))
	for k := range m {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}

func comment(in string, length int) string {
	return txt.Wrap("// ", in, length)
}

func firstToLower(in string) string {
	return string(unicode.ToLower(rune(in[0]))) + in[1:]
}

func firstToUpper(in string) string {
	return string(unicode.ToUpper(rune(in[0]))) + in[1:]
}

func fatalIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package draw

var _ Ink = &DynamicColor{}

// System provided dynamic colors. The defaults are based upon macOS Mojave's
// (10.14.x) standard light theme. Each OS implementation should update them
// appropriately at app launch as well as whenever the system theme changes.
var (
	needSystemColorUpdate                      = true
	AlternateSelectedControlTextColor          = &DynamicColor{Color: White}
	ControlAccentColor                         = &DynamicColor{Color: RGB(0, 122, 255)}
	ControlBackgroundColor                     = &DynamicColor{Color: White}
	ControlColor                               = &DynamicColor{Color: White}
	ControlEdgeAdjColor                        = &DynamicColor{Color: ThemeBlack().SetAlphaIntensity(0.35)}
	ControlEdgeHighlightAdjColor               = &DynamicColor{Color: ThemeWhite().SetAlphaIntensity(0.35)}
	ControlTextColor                           = &DynamicColor{Color: ARGB(0.847059, 0, 0, 0)}
	DisabledControlTextColor                   = &DynamicColor{Color: ARGB(0.247059, 0, 0, 0)}
	FindHighlightColor                         = &DynamicColor{Color: Yellow}
	GridColor                                  = &DynamicColor{Color: RGB(204, 204, 204)}
	HeaderTextColor                            = &DynamicColor{Color: ARGB(0.847059, 0, 0, 0)}
	HighlightColor                             = &DynamicColor{Color: White}
	InvalidBackgroundColor                     = &DynamicColor{Color: RGB(255, 204, 204)}
	KeyboardFocusIndicatorColor                = &DynamicColor{Color: ARGB(0.247059, 0, 103, 244)}
	LabelColor                                 = &DynamicColor{Color: ARGB(0.847059, 0, 0, 0)}
	LinkColor                                  = &DynamicColor{Color: RGB(0, 104, 218)}
	PlaceholderTextColor                       = &DynamicColor{Color: ARGB(0.247059, 0, 0, 0)}
	QuaternaryLabelColor                       = &DynamicColor{Color: ARGB(0.098039, 0, 0, 0)}
	SecondaryLabelColor                        = &DynamicColor{Color: ARGB(0.498039, 0, 0, 0)}
	SelectedContentBackgroundColor             = &DynamicColor{Color: RGB(0, 99, 225)}
	SelectedControlColor                       = &DynamicColor{Color: RGB(179, 215, 255)}
	SelectedControlTextColor                   = &DynamicColor{Color: ARGB(0.847059, 0, 0, 0)}
	SelectedMenuItemTextColor                  = &DynamicColor{Color: White}
	SelectedTextBackgroundColor                = &DynamicColor{Color: RGB(179, 215, 255)}
	SelectedTextColor                          = &DynamicColor{Color: Black}
	SeparatorColor                             = &DynamicColor{Color: ARGB(0.098039, 0, 0, 0)}
	ShadowColor                                = &DynamicColor{Color: Black}
	SystemBlueColor                            = &DynamicColor{Color: RGB(0, 122, 255)}
	SystemBrownColor                           = &DynamicColor{Color: RGB(162, 132, 94)}
	SystemGrayColor                            = &DynamicColor{Color: RGB(142, 142, 147)}
	SystemGreenColor                           = &DynamicColor{Color: RGB(40, 205, 65)}
	SystemOrangeColor                          = &DynamicColor{Color: RGB(255, 149, 0)}
	SystemPinkColor                            = &DynamicColor{Color: RGB(255, 45, 85)}
	SystemPurpleColor                          = &DynamicColor{Color: RGB(175, 82, 222)}
	SystemRedColor                             = &DynamicColor{Color: RGB(255, 59, 48)}
	SystemYellowColor                          = &DynamicColor{Color: RGB(255, 204, 0)}
	TertiaryLabelColor                         = &DynamicColor{Color: ARGB(0.247059, 0, 0, 0)}
	TextAlternateBackgroundColor               = &DynamicColor{Color: RGB(244, 245, 245)}
	TextBackgroundColor                        = &DynamicColor{Color: White}
	TextColor                                  = &DynamicColor{Color: Black}
	UnderPageBackgroundColor                   = &DynamicColor{Color: ARGB(0.898039, 150, 150, 150)}
	UnemphasizedSelectedContentBackgroundColor = &DynamicColor{Color: Gainsboro}
	UnemphasizedSelectedTextBackgroundColor    = &DynamicColor{Color: Gainsboro}
	UnemphasizedSelectedTextColor              = &DynamicColor{Color: Black}
	WindowBackgroundColor                      = &DynamicColor{Color: RGB(236, 236, 236)}
	WindowFrameTextColor                       = &DynamicColor{Color: ARGB(0.847059, 0, 0, 0)}
)

var (
	controlBackgroundGradient         = initBackgroundGradient(ControlColor.Color)
	controlSelectedBackgroundGradient = initBackgroundGradient(ControlAccentColor.Color.Blend(ControlColor.Color, 0.5))
	controlFocusedBackgroundGradient  = initBackgroundGradient(ControlAccentColor.Color)
	controlPressedBackgroundGradient  = initBackgroundGradient(ControlAccentColor.Color.Blend(ThemeBlack(), 0.2))
)

// System provided dynamic inks.
var (
	ControlBackgroundInk         Ink = controlBackgroundGradient
	ControlSelectedBackgroundInk Ink = controlSelectedBackgroundGradient
	ControlFocusedBackgroundInk  Ink = controlFocusedBackgroundGradient
	ControlPressedBackgroundInk  Ink = controlPressedBackgroundGradient
)

func initBackgroundGradient(c Color) *Gradient {
	return NewVerticalEvenlySpacedGradient(&DynamicColor{Color: c.AdjustBrightness(0.2)}, &DynamicColor{Color: c.AdjustBrightness(-0.2)})
}

// DynamicColor holds a color that may be changed.
type DynamicColor struct {
	Color Color
}

func (c *DynamicColor) osPrepareForFill(gc Context) {
	c.Color.osPrepareForFill(gc)
}

func (c *DynamicColor) osFill(gc Context) {
	c.Color.osFill(gc)
}

func (c *DynamicColor) osFillEvenOdd(gc Context) {
	c.Color.osFillEvenOdd(gc)
}

func (c *DynamicColor) osStroke(gc Context) {
	c.Color.osStroke(gc)
}

// MarkSystemColorsForUpdate marks the system colors to be updated the next
// time UpdateSystemColors() is called.
func MarkSystemColorsForUpdate() {
	needSystemColorUpdate = true
}

// UpdateSystemColors updates the system colors, but only if a call to
// MarkSystemColorsForUpdate() was made since the last time this function was
// called.
func UpdateSystemColors() {
	osUpdateSystemColors()
}

// ThemeIsDark returns true if the current theme is a dark one.
func ThemeIsDark() bool {
	return WindowBackgroundColor.Color.Luminance() <= 0.5
}

// ThemeWhite returns White for a light theme and Black for a dark theme.
func ThemeWhite() Color {
	if ThemeIsDark() {
		return Black
	}
	return White
}

// ThemeBlack returns Black for a light theme and White for a dark theme.
func ThemeBlack() Color {
	if ThemeIsDark() {
		return White
	}
	return Black
}

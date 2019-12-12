package draw

import (
	"github.com/richardwilkes/macos/ns"
)

func osUpdateSystemColors() {
	if needSystemColorUpdate {
		AlternateSelectedControlTextColor.Color = convertColor(ns.AlternateSelectedControlTextColor())
		ControlAccentColor.Color = convertColor(ns.ControlAccentColor())
		ControlBackgroundColor.Color = convertColor(ns.ControlBackgroundColor())
		ControlColor.Color = convertColor(ns.ControlColor())
		ControlTextColor.Color = convertColor(ns.ControlTextColor())
		DisabledControlTextColor.Color = convertColor(ns.DisabledControlTextColor())
		FindHighlightColor.Color = convertColor(ns.FindHighlightColor())
		GridColor.Color = convertColor(ns.GridColor())
		HeaderTextColor.Color = convertColor(ns.HeaderTextColor())
		HighlightColor.Color = convertColor(ns.HighlightColor())
		KeyboardFocusIndicatorColor.Color = convertColor(ns.KeyboardFocusIndicatorColor())
		LabelColor.Color = convertColor(ns.LabelColor())
		LinkColor.Color = convertColor(ns.LinkColor())
		PlaceholderTextColor.Color = convertColor(ns.PlaceholderTextColor())
		QuaternaryLabelColor.Color = convertColor(ns.QuaternaryLabelColor())
		SecondaryLabelColor.Color = convertColor(ns.SecondaryLabelColor())
		SelectedContentBackgroundColor.Color = convertColor(ns.SelectedContentBackgroundColor())
		SelectedControlColor.Color = convertColor(ns.SelectedControlColor())
		SelectedControlTextColor.Color = convertColor(ns.SelectedControlTextColor())
		SelectedMenuItemTextColor.Color = convertColor(ns.SelectedMenuItemTextColor())
		SelectedTextBackgroundColor.Color = convertColor(ns.SelectedTextBackgroundColor())
		SelectedTextColor.Color = convertColor(ns.SelectedTextColor())
		SeparatorColor.Color = convertColor(ns.SeparatorColor())
		ShadowColor.Color = convertColor(ns.ShadowColor())
		SystemBlueColor.Color = convertColor(ns.SystemBlueColor())
		SystemBrownColor.Color = convertColor(ns.SystemBrownColor())
		SystemGrayColor.Color = convertColor(ns.SystemGrayColor())
		SystemGreenColor.Color = convertColor(ns.SystemGreenColor())
		SystemOrangeColor.Color = convertColor(ns.SystemOrangeColor())
		SystemPinkColor.Color = convertColor(ns.SystemPinkColor())
		SystemPurpleColor.Color = convertColor(ns.SystemPurpleColor())
		SystemRedColor.Color = convertColor(ns.SystemRedColor())
		SystemYellowColor.Color = convertColor(ns.SystemYellowColor())
		TertiaryLabelColor.Color = convertColor(ns.TertiaryLabelColor())
		TextBackgroundColor.Color = convertColor(ns.TextBackgroundColor())
		TextColor.Color = convertColor(ns.TextColor())
		UnderPageBackgroundColor.Color = convertColor(ns.UnderPageBackgroundColor())
		UnemphasizedSelectedContentBackgroundColor.Color = convertColor(ns.UnemphasizedSelectedContentBackgroundColor())
		UnemphasizedSelectedTextBackgroundColor.Color = convertColor(ns.UnemphasizedSelectedTextBackgroundColor())
		UnemphasizedSelectedTextColor.Color = convertColor(ns.UnemphasizedSelectedTextColor())
		WindowBackgroundColor.Color = convertColor(ns.WindowBackgroundColor())
		WindowFrameTextColor.Color = convertColor(ns.WindowFrameTextColor())

		if colors := ns.AlternatingContentBackgroundColors(); len(colors) > 1 {
			// We currently ignore the first color, since we want to reuse TextBackgroundColor for it
			TextAlternateBackgroundColor.Color = convertColor(colors[1])
		}

		controlBackgroundGradient.Stops[0].Color.Color = ControlColor.Color.AdjustBrightness(0.2)
		controlBackgroundGradient.Stops[1].Color.Color = ControlColor.Color.AdjustBrightness(-0.2)
		c := ControlAccentColor.Color.Blend(ControlColor.Color, 0.5)
		controlSelectedBackgroundGradient.Stops[0].Color.Color = c.AdjustBrightness(0.2)
		controlSelectedBackgroundGradient.Stops[1].Color.Color = c.AdjustBrightness(-0.2)
		controlFocusedBackgroundGradient.Stops[0].Color.Color = ControlAccentColor.Color.AdjustBrightness(0.2)
		controlFocusedBackgroundGradient.Stops[1].Color.Color = ControlAccentColor.Color.AdjustBrightness(-0.2)
		c = ControlAccentColor.Color.Blend(ThemeBlack(), 0.2)
		controlPressedBackgroundGradient.Stops[0].Color.Color = c.AdjustBrightness(0.2)
		controlPressedBackgroundGradient.Stops[1].Color.Color = c.AdjustBrightness(-0.2)

		ControlEdgeAdjColor.Color = ThemeBlack().SetAlphaIntensity(0.35)
		ControlEdgeHighlightAdjColor.Color = ThemeWhite().SetAlphaIntensity(0.35)

		needSystemColorUpdate = false
	}
}

func convertColor(color *ns.Color) Color {
	c := color.ColorUsingColorSpace(ns.ColorSpaceDeviceRGBColorSpace())
	r, g, b, a := c.GetRedGreenBlueAlpha()
	return ARGBfloat(a, r, g, b)
}

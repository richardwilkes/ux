// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package behavior

// Possible ways to handle auto-sizing of the scroll content's preferred size.
const (
	Unmodified Behavior = iota
	FillWidth
	FillHeight
	Fill
	FollowsWidth
	FollowsHeight
)

// Behavior controls how auto-sizing of the scroll content's preferred size is
// handled.
type Behavior uint8

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

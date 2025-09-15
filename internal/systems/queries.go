package systems

import (
	"github.com/gfeyer/ebit/internal/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var (
	QSelectable = donburi.NewQuery(filter.Contains(components.Position, components.Sprite, components.SelectableRes))
	QDrag       = donburi.NewQuery(filter.Contains(components.DragRes))
)

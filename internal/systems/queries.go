package systems

import (
	"github.com/gfeyer/ebit/internal/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var (
	QSelectable = donburi.NewQuery(filter.Contains(components.Position, components.Sprite, components.SelectableRes, components.UnitRes))
	QDrag       = donburi.NewQuery(filter.Contains(components.DragRes))
	QSpice      = donburi.NewQuery(filter.Contains(components.SpiceRes, components.Position, components.Sprite))
)

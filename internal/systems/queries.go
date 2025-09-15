package systems

import (
	"github.com/gfeyer/ebit/internal/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var (
		QSelectable = donburi.NewQuery(filter.And(
		filter.Contains(components.Position, components.SelectableRes),
		filter.Or(filter.Contains(components.UnitRes), filter.Contains(components.RefineryRes)),
	))
	QDrag       = donburi.NewQuery(filter.Contains(components.DragRes))
	QSpice      = donburi.NewQuery(filter.Contains(components.SpiceRes, components.Position, components.Sprite))
)

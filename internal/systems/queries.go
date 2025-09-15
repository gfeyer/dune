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
	QPlayer     = donburi.NewQuery(filter.Contains(components.PlayerRes))

	// Query for the unit build menu
	UnitMenuQuery = donburi.NewQuery(filter.Contains(components.UnitInfoRes))

	// Query for selected buildings
	SelectedBuildingQuery = donburi.NewQuery(filter.Contains(components.SelectableRes))

	// Query for selectable units
	SelectableUnitQuery = donburi.NewQuery(filter.And(
		filter.Contains(components.Position, components.SelectableRes, components.UnitRes),
	))
)

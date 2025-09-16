package systems

import (
	"github.com/gfeyer/ebit/internal/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var (
	// QSelectable retrieves all entities that can be selected by the player, including units and buildings.
	QSelectable = donburi.NewQuery(filter.And(
		filter.Contains(components.Position, components.SelectableRes),
		filter.Or(filter.Contains(components.UnitRes), filter.Contains(components.RefineryRes)),
	))
	// QDrag retrieves the entity that manages the state of the drag-selection box.
	QDrag = donburi.NewQuery(filter.Contains(components.DragRes))
	// QSpice retrieves all spice fields on the map.
	QSpice = donburi.NewQuery(filter.Contains(components.SpiceRes, components.Position, components.Sprite))
	// QPlayer retrieves the player's entity, used for accessing resources like money.
	QPlayer = donburi.NewQuery(filter.Contains(components.PlayerRes))

	// UnitMenuQuery retrieves all entities that represent an option in the unit training menu.
	UnitMenuQuery = donburi.NewQuery(filter.Contains(components.UnitInfoRes))

	// SelectedBuildingQuery retrieves all entities that are currently selected.
	SelectedBuildingQuery = donburi.NewQuery(filter.Contains(components.SelectableRes))

	// SelectableUnitQuery retrieves all unit entities that can be selected by the player.
	SelectableUnitQuery = donburi.NewQuery(filter.And(
		filter.Contains(components.Position, components.SelectableRes, components.UnitRes),
	))
)

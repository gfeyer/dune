package systems

import (
	"github.com/yohamta/donburi/ecs"
)

const (
	// LayerSpice is the rendering layer for spice fields.
	LayerSpice ecs.LayerID = iota
	// LayerBuildings is the rendering layer for buildings.
	LayerBuildings
	// LayerUnits is the rendering layer for units.
	LayerUnits
	// LayerUI is the rendering layer for general UI elements.
	LayerUI
	// LayerMinimap is the rendering layer for the minimap.
	LayerMinimap
	// LayerMenus is the rendering layer for menus.
	LayerMenus
	// LayerBuildMenuUI is the rendering layer for the build menu UI.
	LayerBuildMenuUI
	// LayerPlacement is the rendering layer for the building placement preview.
	LayerPlacement
	// LayerFog is the rendering layer for the fog of war.
	LayerFog
)


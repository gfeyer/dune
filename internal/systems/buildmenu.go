package systems

import (
	"github.com/gfeyer/ebit/internal/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var (
	// BuildMenuQuery retrieves all entities that represent an option in the building construction menu.
	BuildMenuQuery = donburi.NewQuery(filter.Contains(components.BuildInfoRes))
)

// DrawBuildMenu renders the appropriate build menu on the screen.
// If a building is selected, it draws the menu for training units.
// Otherwise, it draws the menu for constructing new buildings.
func DrawBuildMenu(ecs *ecs.ECS, screen *ebiten.Image) {
	// Find minimap to position the build menu below it
	minimapEntry, ok := MinimapQuery.First(ecs.World)
	if !ok {
		return // No minimap, no build menu
	}
	minimap := components.MinimapRes.Get(minimapEntry)

	// Calculate the layout for the menu based on the minimap's position and size.
	menuX := minimap.X
	menuY := minimap.Y + minimap.Height + 10 // 10px padding
	padding := 5
	iconWidth := (minimap.Width - padding) / 2
	rowHeight := 64 + padding // from game.go

	// Determine which menu to draw based on whether a building is currently selected.
	var selectedBuilding *donburi.Entry
	SelectedBuildingQuery.Each(ecs.World, func(entry *donburi.Entry) {
		if components.SelectableRes.Get(entry).Selected {
			// Ensure the selected entity is a building before changing the menu.
			if entry.HasComponent(components.RefineryRes) || entry.HasComponent(components.BarracksRes) {
				selectedBuilding = entry
			}
		}
	})

	i := 0
	if selectedBuilding != nil {
		// A building is selected, so draw the menu for available units.
		// A building is selected, draw the unit menu
		var buildingType components.BuildingType
		if selectedBuilding.HasComponent(components.RefineryRes) {
			buildingType = components.BuildingRefinery
		} else if selectedBuilding.HasComponent(components.BarracksRes) {
			buildingType = components.BuildingBarracks
		}

		// Iterate through all available units and draw the ones that can be built from the selected building.
		UnitMenuQuery.Each(ecs.World, func(entry *donburi.Entry) {
			unitInfo := components.UnitInfoRes.Get(entry)
			if unitInfo.RequiredBuilding != buildingType {
				return
			}

			// Position the icon in a grid layout.
			col := i % 2
			row := i / 2
			iconX := menuX + col*(iconWidth+padding)
			iconY := menuY + row*rowHeight

			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(iconX), float64(iconY))
			screen.DrawImage(unitInfo.Icon, opts)

			i++
		})
	} else {
		// No building is selected, so draw the menu for constructing buildings.
		// Iterate through the available buildings and draw their icons.
		BuildMenuQuery.Each(ecs.World, func(entry *donburi.Entry) {
			buildInfo := components.BuildInfoRes.Get(entry)

			// Position the icon in a grid layout.
			col := i % 2
			row := i / 2
			iconX := menuX + col*(iconWidth+padding)
			iconY := menuY + row*rowHeight

			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(iconX), float64(iconY))
			screen.DrawImage(buildInfo.Icon, opts)

			i++
		})
	}
}

package systems

import (
	"github.com/gfeyer/ebit/internal/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var (
	BuildMenuQuery = donburi.NewQuery(filter.Contains(components.BuildInfoRes))
)

func DrawBuildMenu(ecs *ecs.ECS, screen *ebiten.Image) {
	// Find minimap to position the build menu below it
	minimapEntry, ok := MinimapQuery.First(ecs.World)
	if !ok {
		return // No minimap, no build menu
	}
	minimap := components.MinimapRes.Get(minimapEntry)

	// Menu layout
	menuX := minimap.X
	menuY := minimap.Y + minimap.Height + 10 // 10px padding
	padding := 5
	iconWidth := (minimap.Width - padding) / 2
	rowHeight := 64 + padding // from game.go

	// Determine which menu to draw based on selection
	var selectedBuilding *donburi.Entry
	SelectedBuildingQuery.Each(ecs.World, func(entry *donburi.Entry) {
		if components.SelectableRes.Get(entry).Selected {
			// Ensure the selected entity is a building before changing the menu
			if entry.HasComponent(components.RefineryRes) || entry.HasComponent(components.BarracksRes) {
				selectedBuilding = entry
			}
		}
	})

	i := 0
	if selectedBuilding != nil {
		// A building is selected, draw the unit menu
		var buildingType components.BuildingType
		if selectedBuilding.HasComponent(components.RefineryRes) {
			buildingType = components.BuildingRefinery
		} else if selectedBuilding.HasComponent(components.BarracksRes) {
			buildingType = components.BuildingBarracks
		}

		UnitMenuQuery.Each(ecs.World, func(entry *donburi.Entry) {
			unitInfo := components.UnitInfoRes.Get(entry)
			if unitInfo.RequiredBuilding != buildingType {
				return
			}

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
		// No building is selected, draw the building menu
		BuildMenuQuery.Each(ecs.World, func(entry *donburi.Entry) {
			buildInfo := components.BuildInfoRes.Get(entry)

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

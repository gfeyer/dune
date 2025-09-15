package systems

import (
	"github.com/gfeyer/ebit/internal/camera"
	"github.com/gfeyer/ebit/internal/components"
	"github.com/gfeyer/ebit/internal/factory"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var (
	PlacementQuery = donburi.NewQuery(filter.Contains(components.PlacementRes))
	PlayerQuery    = donburi.NewQuery(filter.Contains(components.PlayerRes))
)

func UpdateBuildInput(ecs *ecs.ECS) {
	placementEntry, ok := PlacementQuery.First(ecs.World)
	if !ok {
		return
	}
	placement := components.PlacementRes.Get(placementEntry)

	if placement.IsPlacing {
		// Handle cancellation
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			placement.IsPlacing = false
			return
		}

		// Handle placement
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			playerEntry, ok := PlayerQuery.First(ecs.World)
			if !ok {
				return
			}
			player := components.PlayerRes.Get(playerEntry)

			// Check for sufficient funds
			if player.Money < placement.Cost {
				// Not enough money, exit placement mode
				placement.IsPlacing = false
				return
			}

			// Deduct cost and place building
			player.Money -= placement.Cost

			cameraEntry, ok := camera.CameraQuery.First(ecs.World)
			if !ok {
				return
			}
			cam := camera.CameraRes.Get(cameraEntry)
			mx, my := ebiten.CursorPosition()
			wx, wy := cam.ScreenToWorld(float64(mx), float64(my))

			// Create the building
			switch placement.BuildingType {
			case components.BuildingRefinery:
				factory.CreateRefinery(ecs.World, wx, wy)
			case components.BuildingBarracks:
				factory.CreateBarracks(ecs.World, wx, wy)
			}

			// Exit placement mode
			placement.IsPlacing = false
		}
		return
	}

	// If not in placement mode, check for clicks on the build menu
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()

		minimapEntry, ok := MinimapQuery.First(ecs.World)
		if !ok {
			return
		}
		minimap := components.MinimapRes.Get(minimapEntry)

		menuX := minimap.X
		menuY := minimap.Y + minimap.Height + 10
		iconWidth := 32
		iconHeight := 32
		padding := 5
		colWidth := iconWidth + padding
		rowHeight := iconHeight + padding + 15

		i := 0
		BuildMenuQuery.Each(ecs.World, func(entry *donburi.Entry) {
			buildInfo := components.BuildInfoRes.Get(entry)

			col := i % 2
			row := i / 2

			iconX := menuX + col*colWidth
			iconY := menuY + row*rowHeight

			if mx >= iconX && mx < iconX+iconWidth && my >= iconY && my < iconY+iconHeight {
				placement.IsPlacing = true
				placement.BuildingType = buildInfo.Type
				placement.Icon = buildInfo.Icon
				placement.Cost = buildInfo.Cost
			}
			i++
		})
	}
}

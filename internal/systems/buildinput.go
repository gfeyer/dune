// Package systems contains the logic for updating the game state in response to user input and game events.
package systems

import (
	"math/rand"

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
	// PlacementQuery retrieves the entity that manages the state of placing a new building.
	PlacementQuery = donburi.NewQuery(filter.Contains(components.PlacementRes))
	// PlayerQuery retrieves the player's entity, used for accessing resources like money.
	PlayerQuery = donburi.NewQuery(filter.Contains(components.PlayerRes))

	// SelectableBuildingQuery retrieves all building entities that can be selected by the player.
	SelectableBuildingQuery = donburi.NewQuery(
		filter.And(
			filter.Contains(components.SelectableRes),
			filter.Or(
				filter.Contains(components.RefineryRes),
				filter.Contains(components.BarracksRes),
			),
		),
	)
)

// UpdateBuildInput handles all user input related to building placement and unit creation.
// It's called every frame to check for mouse clicks on the world or the build menu.
func UpdateBuildInput(ecs *ecs.ECS) {
	placementEntry, ok := PlacementQuery.First(ecs.World)
	if !ok {
		return
	}
	placement := components.PlacementRes.Get(placementEntry)

	// If the player is currently in the process of placing a building.
	if placement.IsPlacing {
		// Cancel placement with a right-click or by pressing the Escape key.
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			placement.IsPlacing = false
			return
		}

		// Confirm building placement with a left-click.
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

	// If not in placement mode, check for clicks on the build menu or for selecting buildings in the world.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()

		// First, check for clicks on the build menu
		if checkBuildMenuClick(ecs, mx, my) {
			return // Click was handled by the menu
		}

		// If not a menu click, handle world selection
		cameraEntry, _ := camera.CameraQuery.First(ecs.World)
		cam := camera.CameraRes.Get(cameraEntry)
		wx, wy := cam.ScreenToWorld(float64(mx), float64(my))

		// Deselect any currently selected entity
		SelectedBuildingQuery.Each(ecs.World, func(entry *donburi.Entry) {
			selectable := components.SelectableRes.Get(entry)
			if selectable.Selected {
				selectable.Selected = false
			}
		})

		// Check if a selectable building was clicked
		var clickedBuilding *donburi.Entry
		SelectableBuildingQuery.Each(ecs.World, func(entry *donburi.Entry) {
			p := components.Position.Get(entry)
			sprite := components.Sprite.Get(entry)
			bounds := (*sprite).Bounds()

			if wx >= p.X && wx < p.X+float64(bounds.Dx()) && wy >= p.Y && wy < p.Y+float64(bounds.Dy()) {
				clickedBuilding = entry
			}
		})

		// If a building was clicked, select it
		if clickedBuilding != nil {
			selectable := components.SelectableRes.Get(clickedBuilding)
			selectable.Selected = true
		}
	}
}

// checkBuildMenuClick determines if a mouse click at screen coordinates (mx, my) has occurred on a build menu icon.
// It returns true if a menu item was clicked, handling the corresponding action, and false otherwise.
func checkBuildMenuClick(ecs *ecs.ECS, mx, my int) bool {
	placementEntry, ok := PlacementQuery.First(ecs.World)
	if !ok {
		return false
	}
	placement := components.PlacementRes.Get(placementEntry)

	minimapEntry, ok := MinimapQuery.First(ecs.World)
	if !ok {
		return false
	}
	minimap := components.MinimapRes.Get(minimapEntry)

	menuX := minimap.X
	menuY := minimap.Y + minimap.Height + 10
	padding := 5
	iconWidth := (minimap.Width - padding) / 2
	rowHeight := 64 + padding // from game.go

	clickedOnMenu := false

	// Check if a building is currently selected to determine whether to show the unit menu or the building menu.
	var selectedBuilding *donburi.Entry
	SelectedBuildingQuery.Each(ecs.World, func(entry *donburi.Entry) {
		if components.SelectableRes.Get(entry).Selected {
			selectedBuilding = entry
		}
	})

	// If a building is selected, display the menu for units that can be built from it.
	if selectedBuilding != nil {
		var buildingType components.BuildingType
		if selectedBuilding.HasComponent(components.RefineryRes) {
			buildingType = components.BuildingRefinery
		} else if selectedBuilding.HasComponent(components.BarracksRes) {
			buildingType = components.BuildingBarracks
		}

		i := 0
		UnitMenuQuery.Each(ecs.World, func(entry *donburi.Entry) {
			unitInfo := components.UnitInfoRes.Get(entry)
			if unitInfo.RequiredBuilding != buildingType {
				return
			}

			col := i % 2
			row := i / 2
			iconX := menuX + col*(iconWidth+padding)
			iconY := menuY + row*rowHeight
			actualIconWidth := unitInfo.Icon.Bounds().Dx()
			actualIconHeight := unitInfo.Icon.Bounds().Dy()

			if mx >= iconX && mx < iconX+actualIconWidth && my >= iconY && my < iconY+actualIconHeight {
				// Handle unit creation
				playerEntry, ok := PlayerQuery.First(ecs.World)
				if !ok {
					return
				}
				player := components.PlayerRes.Get(playerEntry)

				if player.Money >= unitInfo.Cost {
					player.Money -= unitInfo.Cost
					buildingPos := components.Position.Get(selectedBuilding)
					spawnX := buildingPos.X + float64(rand.Intn(64)) + 32 // Spawn to the right of the building
					spawnY := buildingPos.Y + float64(rand.Intn(64)) + 32

					switch unitInfo.Type {
					case components.Harvester:
						factory.CreateHarvester(ecs.World, spawnX, spawnY)
					case components.Trike:
						factory.CreateTrike(ecs.World, spawnX, spawnY)
					case components.Quad:
						factory.CreateQuad(ecs.World, spawnX, spawnY)
					}
				}
				clickedOnMenu = true
			}
			i++
		})
	} else {
		// If no building is selected, display the menu for constructing new buildings.
		i := 0
		BuildMenuQuery.Each(ecs.World, func(entry *donburi.Entry) {
			buildInfo := components.BuildInfoRes.Get(entry)
			col := i % 2
			row := i / 2
			iconX := menuX + col*(iconWidth+padding)
			iconY := menuY + row*rowHeight
			actualIconWidth := buildInfo.Icon.Bounds().Dx()
			actualIconHeight := buildInfo.Icon.Bounds().Dy()

			if mx >= iconX && mx < iconX+actualIconWidth && my >= iconY && my < iconY+actualIconHeight {
				// Clicked on this build option
				placement.IsPlacing = true
				placement.BuildingType = buildInfo.Type
				placement.Icon = buildInfo.Icon
				placement.Cost = buildInfo.Cost
				clickedOnMenu = true
			}
			i++
		})
	}

	return clickedOnMenu
}

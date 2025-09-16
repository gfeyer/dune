package systems

import (
	"image/color"

	"github.com/gfeyer/ebit/internal/camera"
	"github.com/gfeyer/ebit/internal/components"
	"github.com/gfeyer/ebit/internal/fog"
	"github.com/gfeyer/ebit/internal/settings"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var (
	// qPlayerUnits retrieves all player-controlled units and buildings that provide vision.
	qPlayerUnits = donburi.NewQuery(filter.And(
		filter.Or(filter.Contains(components.UnitRes), filter.Contains(components.RefineryRes), filter.Contains(components.BarracksRes)),
		filter.Contains(components.Position),
	))
)

// UpdateFog manages the fog of war logic. It first shrouds all previously visible areas
// and then reveals the areas around the player's current units and buildings.
func UpdateFog(ecs *ecs.ECS) {
	fogRes := fog.GetFog(ecs.World)

	// 1. Set all currently visible tiles to shrouded. This creates the effect of areas becoming
	//    unexplored again when units move away.
	for y := 0; y < fogRes.Height; y++ {
		for x := 0; x < fogRes.Width; x++ {
			if fogRes.Grid[y][x] == fog.Visible {
				fogRes.Grid[y][x] = fog.Shroud
			}
		}
	}

	// 2. Reveal the fog around each of the player's units and buildings.
	qPlayerUnits.Each(ecs.World, func(entry *donburi.Entry) {
		p := components.Position.Get(entry)
		visionRadius := 16 // in tiles

		tileX := int(p.X) / fogRes.TileSize
		tileY := int(p.Y) / fogRes.TileSize

		for y := tileY - visionRadius; y <= tileY+visionRadius; y++ {
			for x := tileX - visionRadius; x <= tileX+visionRadius; x++ {
				if x >= 0 && x < fogRes.Width && y >= 0 && y < fogRes.Height {
					dx := x - tileX
					dy := y - tileY
					if dx*dx+dy*dy <= visionRadius*visionRadius {
						fogRes.Grid[y][x] = fog.Visible
					}
				}
			}
		}
	})
}

// DrawFog renders the fog of war onto the screen. It draws a black rectangle for hidden areas
// and a semi-transparent one for shrouded areas.
func DrawFog(ecs *ecs.ECS, screen *ebiten.Image) {
	cameraEntry, _ := camera.CameraQuery.First(ecs.World)
	cam := camera.CameraRes.Get(cameraEntry)
	fogRes := fog.GetFog(ecs.World)
	s := settings.GetSettings(ecs.World)

	for y := 0; y < fogRes.Height; y++ {
		for x := 0; x < fogRes.Width; x++ {
			worldX := float32(x * fogRes.TileSize)
			worldY := float32(y * fogRes.TileSize)

			screenX := worldX - float32(cam.X)
			screenY := worldY - float32(cam.Y)

			// Culling: Don't draw fog tiles that are outside the camera's view.
			if screenX+float32(fogRes.TileSize) < 0 || screenX > float32(s.ScreenWidth) || screenY+float32(fogRes.TileSize) < 0 || screenY > float32(s.ScreenHeight) {
				continue
			}

			// Render the fog tile based on its state (Hidden, Shroud, or Visible).
			switch fogRes.Grid[y][x] {
			case fog.Hidden:
				vector.DrawFilledRect(screen, screenX, screenY, float32(fogRes.TileSize), float32(fogRes.TileSize), color.Black, false)
			case fog.Shroud:
				vector.DrawFilledRect(screen, screenX, screenY, float32(fogRes.TileSize), float32(fogRes.TileSize), color.RGBA{0, 0, 0, 180}, false)
			}
		}
	}
}

package systems

import (
	"image/color"
	"math"

	"github.com/gfeyer/ebit/internal/camera"
	"github.com/gfeyer/ebit/internal/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/font/basicfont"
)

const (
	LayerSprites ecs.LayerID = iota // 0
	LayerUI                         // 1 (in-world UI like health bars)
	LayerMinimap                    // 2
	LayerMenus                      // 3 (for future build menus etc.)
)

var (
	qSprites = donburi.NewQuery(filter.Contains(components.Position, components.Sprite))
	qTrikeUI = donburi.NewQuery(filter.And(
		filter.Contains(components.Position, components.Sprite, components.UnitRes, components.HealthRes),
		filter.Not(filter.Contains(components.HarvesterRes)),
	))
	qHarvesterUI = donburi.NewQuery(filter.Contains(components.Position, components.Sprite, components.UnitRes, components.HealthRes, components.HarvesterRes))
	qRefineryUI  = donburi.NewQuery(filter.Contains(components.Position, components.Sprite, components.RefineryRes))
)

func DrawSprites(ecs *ecs.ECS, screen *ebiten.Image) {
	cameraEntry, _ := camera.CameraQuery.First(ecs.World)
	cam := camera.CameraRes.Get(cameraEntry)

	qSprites.Each(ecs.World, func(entry *donburi.Entry) {
		p := components.Position.Get(entry)
		img := components.Sprite.Get(entry)
		op := &ebiten.DrawImageOptions{}

		// Handle rotation for entities that have velocity
		if entry.HasComponent(components.Velocity) {
			v := components.Velocity.Get(entry)
			if v.X != 0 || v.Y != 0 {
				bounds := (*img).Bounds()
				centerX, centerY := float64(bounds.Dx())/2, float64(bounds.Dy())/2
				op.GeoM.Translate(-centerX, -centerY)
				op.GeoM.Rotate(math.Atan2(v.Y, v.X) + math.Pi/2)
				op.GeoM.Translate(p.X-cam.X+centerX, p.Y-cam.Y+centerY)
			} else {
				op.GeoM.Translate(p.X-cam.X, p.Y-cam.Y)
			}
		} else {
			op.GeoM.Translate(p.X-cam.X, p.Y-cam.Y)
		}

		// Handle selection tint for selectable entities
		if entry.HasComponent(components.SelectableRes) && components.SelectableRes.Get(entry).Selected {
			cm := colorm.ColorM{}
			cm.Scale(0, 0, 0, 1)
			cm.Translate(0, 1, 0, 0)
			colorm.DrawImage(screen, *img, cm, &colorm.DrawImageOptions{
				GeoM: op.GeoM,
			})
		} else {
			screen.DrawImage(*img, op)
		}
	})
}

func DrawUI(ecs *ecs.ECS, screen *ebiten.Image) {
	cameraEntry, _ := camera.CameraQuery.First(ecs.World)
	cam := camera.CameraRes.Get(cameraEntry)

	// Draw UI elements on top of Trikes
	qTrikeUI.Each(ecs.World, func(entry *donburi.Entry) {
		p := components.Position.Get(entry)
		img := components.Sprite.Get(entry)
		health := components.HealthRes.Get(entry)
		barWidth := float32((*img).Bounds().Dx())

		// Health bar
		healthBarY := float32(p.Y-cam.Y) - 8
		healthPercentage := float32(health.Current) / float32(health.Max)
		vector.DrawFilledRect(screen, float32(p.X-cam.X), healthBarY, barWidth, 4, color.RGBA{R: 255, A: 255}, false)
		vector.DrawFilledRect(screen, float32(p.X-cam.X), healthBarY, barWidth*healthPercentage, 4, color.RGBA{G: 255, A: 255}, false)

		// Unit label
		labelY := int(healthBarY) - 2
		text.Draw(screen, "Trike", basicfont.Face7x13, int(p.X-cam.X), labelY, color.White)
	})

	// Draw UI elements on top of Harvesters
	qHarvesterUI.Each(ecs.World, func(entry *donburi.Entry) {
		p := components.Position.Get(entry)
		img := components.Sprite.Get(entry)
		health := components.HealthRes.Get(entry)
		barWidth := float32((*img).Bounds().Dx())

		// Health bar
		healthBarY := float32(p.Y-cam.Y) - 12
		healthPercentage := float32(health.Current) / float32(health.Max)
		vector.DrawFilledRect(screen, float32(p.X-cam.X), healthBarY, barWidth, 4, color.RGBA{R: 255, A: 255}, false)
		vector.DrawFilledRect(screen, float32(p.X-cam.X), healthBarY, barWidth*healthPercentage, 4, color.RGBA{G: 255, A: 255}, false)

		// Spice bar for Harvester
		harvester := components.HarvesterRes.Get(entry)
		spiceBarY := healthBarY + 5
		spicePercentage := float32(harvester.CarriedAmount) / float32(harvester.Capacity)
		vector.DrawFilledRect(screen, float32(p.X-cam.X), spiceBarY, barWidth, 4, color.RGBA{R: 64, G: 64, B: 64, A: 255}, false)
		vector.DrawFilledRect(screen, float32(p.X-cam.X), spiceBarY, barWidth*spicePercentage, 4, color.RGBA{R: 255, G: 140, B: 0, A: 255}, false)

		// Unit label
		labelY := int(healthBarY) - 2
		text.Draw(screen, "Harvester", basicfont.Face7x13, int(p.X-cam.X), labelY, color.White)
	})

	// Draw UI elements on top of Refineries
	qRefineryUI.Each(ecs.World, func(entry *donburi.Entry) {
		p := components.Position.Get(entry)
		labelY := int(p.Y-cam.Y) - 2
		text.Draw(screen, "Refinery", basicfont.Face7x13, int(p.X-cam.X), labelY, color.White)
	})

	// Draw drag selection
	dragEntry, ok := QDrag.First(ecs.World)
	if ok {
		drag := components.DragRes.Get(dragEntry)
		if drag.IsDragging {
			vector.StrokeRect(screen, float32(drag.StartX), float32(drag.StartY), float32(drag.EndX-drag.StartX), float32(drag.EndY-drag.StartY), 1, color.White, false)
		}
	}
}

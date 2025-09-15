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
	LayerDefault ecs.LayerID = iota
)

var (
	qSprites = donburi.NewQuery(filter.Contains(components.Position, components.Sprite))
	qUnitUI  = donburi.NewQuery(filter.Contains(components.Position, components.Sprite, components.UnitRes, components.HealthRes, components.Velocity, components.SelectableRes))
)

func Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	screen.Fill(color.RGBA{210, 180, 140, 255}) // sand color

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

	// Draw UI elements on top of units
	qUnitUI.Each(ecs.World, func(entry *donburi.Entry) {
		p := components.Position.Get(entry)
		img := components.Sprite.Get(entry)
		health := components.HealthRes.Get(entry)
		unit := components.UnitRes.Get(entry)
		barWidth := float32((*img).Bounds().Dx())

		// Health bar
		healthBarY := float32(p.Y-cam.Y) - 8
		healthPercentage := float32(health.Current) / float32(health.Max)
		vector.DrawFilledRect(screen, float32(p.X-cam.X), healthBarY, barWidth, 4, color.RGBA{R: 255, A: 255}, false)
		vector.DrawFilledRect(screen, float32(p.X-cam.X), healthBarY, barWidth*healthPercentage, 4, color.RGBA{G: 255, A: 255}, false)

		// Spice bar for Harvester
		if unit.Type == components.Harvester {
			harvester := components.HarvesterRes.Get(entry)
			spiceBarY := healthBarY + 5
			spicePercentage := float32(harvester.CarriedAmount) / float32(harvester.Capacity)
			vector.DrawFilledRect(screen, float32(p.X-cam.X), spiceBarY, barWidth, 2, color.RGBA{R: 255, G: 165, A: 255}, false)
			vector.DrawFilledRect(screen, float32(p.X-cam.X), spiceBarY, barWidth*spicePercentage, 2, color.RGBA{R: 255, G: 140, B: 0, A: 255}, false)
		}

		// Unit label
		labelY := int(healthBarY) - 2
		var label string
		if unit.Type == components.Harvester {
			label = "Harvester"
		} else {
			label = "Trike"
		}
		text.Draw(screen, label, basicfont.Face7x13, int(p.X-cam.X), labelY, color.White)
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

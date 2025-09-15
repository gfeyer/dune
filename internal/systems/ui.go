package systems

import (
	"fmt"
	"image/color"

	"github.com/gfeyer/ebit/internal/camera"
	"github.com/gfeyer/ebit/internal/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/font/basicfont"
)

var (
	qTrikeUI = donburi.NewQuery(filter.And(
		filter.Contains(components.Position, components.Sprite, components.UnitRes, components.HealthRes),
		filter.Not(filter.Contains(components.HarvesterRes)),
	))
	qHarvesterUI = donburi.NewQuery(filter.Contains(components.Position, components.Sprite, components.UnitRes, components.HealthRes, components.HarvesterRes))
	qRefineryUI  = donburi.NewQuery(filter.Contains(components.Position, components.Sprite, components.RefineryRes))
	qBarracksUI  = donburi.NewQuery(filter.Contains(components.Position, components.Sprite, components.BarracksRes))
)

func DrawUI(ecs *ecs.ECS, screen *ebiten.Image) {
	cameraEntry, _ := camera.CameraQuery.First(ecs.World)
	cam := camera.CameraRes.Get(cameraEntry)

	// Draw player money
	playerEntry, ok := QPlayer.First(ecs.World)
	if ok {
		player := components.PlayerRes.Get(playerEntry)
		moneyText := fmt.Sprintf("$%d", player.Money)
		text.Draw(screen, moneyText, basicfont.Face7x13, 10, 20, color.White)
	}

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

	// Draw UI elements on top of Barracks
	qBarracksUI.Each(ecs.World, func(entry *donburi.Entry) {
		p := components.Position.Get(entry)
		labelY := int(p.Y-cam.Y) - 2
		text.Draw(screen, "Barracks", basicfont.Face7x13, int(p.X-cam.X), labelY, color.White)
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

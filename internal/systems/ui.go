package systems

import (
	"fmt"
	"image/color"

	"github.com/gfeyer/ebit/internal/camera"
	"github.com/gfeyer/ebit/internal/components"
	"github.com/gfeyer/ebit/internal/settings"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/font/basicfont"
)

var (
	// qTrikeUI retrieves all Trike units to draw their specific UI elements.
	qTrikeUI = donburi.NewQuery(filter.And(
		filter.Contains(components.Position, components.Sprite, components.UnitRes, components.HealthRes),
		filter.Not(filter.Contains(components.HarvesterRes)),
	))
	// qHarvesterUI retrieves all Harvester units for their UI.
	qHarvesterUI = donburi.NewQuery(filter.Contains(components.Position, components.Sprite, components.UnitRes, components.HealthRes, components.HarvesterRes))
	// qRefineryUI retrieves all Refinery buildings for their UI.
	qRefineryUI = donburi.NewQuery(filter.Contains(components.Position, components.Sprite, components.RefineryRes))
	// qBarracksUI retrieves all Barracks buildings for their UI.
	qBarracksUI = donburi.NewQuery(filter.Contains(components.Position, components.Sprite, components.BarracksRes))
)

// DrawUI renders all the in-game user interface elements, such as health bars, unit labels, and resource counters.
func DrawUI(ecs *ecs.ECS, screen *ebiten.Image) {
	cameraEntry, _ := camera.CameraQuery.First(ecs.World)
	cam := camera.CameraRes.Get(cameraEntry)

	// Display the player's current amount of money at the top-left of the screen.
	playerEntry, ok := QPlayer.First(ecs.World)
	if ok {
		player := components.PlayerRes.Get(playerEntry)
		moneyText := fmt.Sprintf("$%d", player.Money)
		text.Draw(screen, moneyText, basicfont.Face7x13, 10, 20, color.White)
	}

	// Draw health bars and labels for all Trike units.
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

	// Draw health bars, spice capacity bars, and labels for all Harvester units.
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

	// Draw labels for all Refinery buildings.
	qRefineryUI.Each(ecs.World, func(entry *donburi.Entry) {
		p := components.Position.Get(entry)
		labelY := int(p.Y-cam.Y) - 2
		text.Draw(screen, "Refinery", basicfont.Face7x13, int(p.X-cam.X), labelY, color.White)
	})

	// Draw labels for all Barracks buildings.
	qBarracksUI.Each(ecs.World, func(entry *donburi.Entry) {
		p := components.Position.Get(entry)
		labelY := int(p.Y-cam.Y) - 2
		text.Draw(screen, "Barracks", basicfont.Face7x13, int(p.X-cam.X), labelY, color.White)
	})

	// If the player is drag-selecting, draw the selection rectangle.
	dragEntry, ok := QDrag.First(ecs.World)
	if ok {
		drag := components.DragRes.Get(dragEntry)
		if drag.IsDragging {
			vector.StrokeRect(screen, float32(drag.StartX), float32(drag.StartY), float32(drag.EndX-drag.StartX), float32(drag.EndY-drag.StartY), 1, color.White, false)
		}
	}

	// Display the current frames per second at the bottom-right of the screen.
	s := settings.GetSettings(ecs.World)
	fpsText := fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS())
	text.Draw(screen, fpsText, basicfont.Face7x13, s.ScreenWidth-80, s.ScreenHeight-10, color.White)
}

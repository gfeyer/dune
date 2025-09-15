package systems

import (
	"image/color"

	"github.com/gfeyer/ebit/internal/camera"
	"github.com/gfeyer/ebit/internal/components"
	"github.com/gfeyer/ebit/internal/settings"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var MinimapQuery = donburi.NewQuery(filter.Contains(components.MinimapRes))

func UpdateMinimap(ecs *ecs.ECS) {
	minimapEntry, ok := MinimapQuery.First(ecs.World)
	if !ok {
		return
	}
	minimap := components.MinimapRes.Get(minimapEntry)

	// Left-click to move camera
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		if mx >= minimap.X && mx < minimap.X+minimap.Width && my >= minimap.Y && my < minimap.Y+minimap.Height {
			settings := settings.GetSettings(ecs.World)
			cameraEntry, _ := camera.CameraQuery.First(ecs.World)
			cam := camera.CameraRes.Get(cameraEntry)

			scaleX := float64(minimap.Width) / float64(settings.MapWidth)
			scaleY := float64(minimap.Height) / float64(settings.MapHeight)

			cam.X = (float64(mx - minimap.X) / scaleX) - float64(settings.ScreenWidth)/2
			cam.Y = (float64(my - minimap.Y) / scaleY) - float64(settings.ScreenHeight)/2
		}
	}

	// Right-click to command units
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		mx, my := ebiten.CursorPosition()
		if mx >= minimap.X && mx < minimap.X+minimap.Width && my >= minimap.Y && my < minimap.Y+minimap.Height {
			settings := settings.GetSettings(ecs.World)
			scaleX := float64(minimap.Width) / float64(settings.MapWidth)
			scaleY := float64(minimap.Height) / float64(settings.MapHeight)

			wx := (float64(mx - minimap.X) / scaleX)
			wy := (float64(my - minimap.Y) / scaleY)

			QSelectable.Each(ecs.World, func(entry *donburi.Entry) {
				if components.SelectableRes.Get(entry).Selected {
					*components.TargetRes.Get(entry) = components.Target{X: wx, Y: wy}
				}
			})
		}
	}
}

func DrawMinimap(ecs *ecs.ECS, screen *ebiten.Image) {
	minimapEntry, ok := MinimapQuery.First(ecs.World)
	if !ok {
		return
	}
	minimap := components.MinimapRes.Get(minimapEntry)

	settings := settings.GetSettings(ecs.World)
	cameraEntry, _ := camera.CameraQuery.First(ecs.World)
	cam := camera.CameraRes.Get(cameraEntry)

	// Minimap background
	vector.DrawFilledRect(screen, float32(minimap.X), float32(minimap.Y), float32(minimap.Width), float32(minimap.Height), color.RGBA{50, 50, 50, 150}, false)

	// Scale factors
	scaleX := float64(minimap.Width) / float64(settings.MapWidth)
	scaleY := float64(minimap.Height) / float64(settings.MapHeight)

	// Draw units (green)
	QSelectable.Each(ecs.World, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)
		unitX := float32(minimap.X + int(pos.X*scaleX))
		unitY := float32(minimap.Y + int(pos.Y*scaleY))
		vector.DrawFilledRect(screen, unitX-1, unitY-1, 3, 3, color.RGBA{G: 255, A: 255}, false)
	})

	// Draw spice (orange)
	QSpice.Each(ecs.World, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)
		spiceX := float32(minimap.X + int(pos.X*scaleX))
		spiceY := float32(minimap.Y + int(pos.Y*scaleY))
		vector.DrawFilledRect(screen, spiceX, spiceY, 2, 2, color.RGBA{R: 255, G: 140, A: 255}, false)
	})

	// Draw camera view
	camX := float32(minimap.X + int(cam.X*scaleX))
	camY := float32(minimap.Y + int(cam.Y*scaleY))
	camW := float32(float64(settings.ScreenWidth) * scaleX)
	camH := float32(float64(settings.ScreenHeight) * scaleY)
	vector.StrokeRect(screen, camX, camY, camW, camH, 1, color.White, false)
}

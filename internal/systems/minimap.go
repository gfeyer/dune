package systems

import (
	"image/color"

	"github.com/gfeyer/ebit/internal/camera"
	"github.com/gfeyer/ebit/internal/components"
	"github.com/gfeyer/ebit/internal/fog"
	"github.com/gfeyer/ebit/internal/settings"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var (
	MinimapQuery = donburi.NewQuery(filter.Contains(components.MinimapRes))

	minimapFogImage *ebiten.Image
)

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

			cam.X = (float64(mx-minimap.X) / scaleX) - float64(settings.ScreenWidth)/2
			cam.Y = (float64(my-minimap.Y) / scaleY) - float64(settings.ScreenHeight)/2
		}
	}

	// Right-click to command units
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		mx, my := ebiten.CursorPosition()
		if mx >= minimap.X && mx < minimap.X+minimap.Width && my >= minimap.Y && my < minimap.Y+minimap.Height {
			settings := settings.GetSettings(ecs.World)
			scaleX := float64(minimap.Width) / float64(settings.MapWidth)
			scaleY := float64(minimap.Height) / float64(settings.MapHeight)

			wx := (float64(mx-minimap.X) / scaleX)
			wy := (float64(my-minimap.Y) / scaleY)

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

	// Minimap background (sand color)
	vector.DrawFilledRect(screen, float32(minimap.X), float32(minimap.Y), float32(minimap.Width), float32(minimap.Height), color.RGBA{210, 180, 140, 255}, false)

	// Scale factors
	scaleX := float64(minimap.Width) / float64(settings.MapWidth)
	scaleY := float64(minimap.Height) / float64(settings.MapHeight)

	fogRes := fog.GetFog(ecs.World)

	// Draw units (green)
	QSelectable.Each(ecs.World, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)

		tileX := int(pos.X) / fogRes.TileSize
		tileY := int(pos.Y) / fogRes.TileSize

		if tileX >= 0 && tileX < fogRes.Width && tileY >= 0 && tileY < fogRes.Height {
			if fogRes.Grid[tileY][tileX] != fog.Hidden {
				unitX := float32(minimap.X + int(pos.X*scaleX))
				unitY := float32(minimap.Y + int(pos.Y*scaleY))
				vector.DrawFilledRect(screen, unitX-1, unitY-1, 3, 3, color.RGBA{G: 255, A: 255}, false)
			}
		}
	})

	// Draw spice (orange)
	QSpice.Each(ecs.World, func(entry *donburi.Entry) {
		pos := components.Position.Get(entry)

		tileX := int(pos.X) / fogRes.TileSize
		tileY := int(pos.Y) / fogRes.TileSize

		if tileX >= 0 && tileX < fogRes.Width && tileY >= 0 && tileY < fogRes.Height {
			if fogRes.Grid[tileY][tileX] != fog.Hidden {
				spiceX := float32(minimap.X + int(pos.X*scaleX))
				spiceY := float32(minimap.Y + int(pos.Y*scaleY))
				vector.DrawFilledRect(screen, spiceX, spiceY, 2, 2, color.RGBA{R: 255, G: 140, A: 255}, false)
			}
		}
	})

	// Lazily initialize the minimap fog image
	if minimapFogImage == nil || minimapFogImage.Bounds().Dx() != fogRes.Width || minimapFogImage.Bounds().Dy() != fogRes.Height {
		minimapFogImage = ebiten.NewImage(fogRes.Width, fogRes.Height)
	}

	// Update the fog image pixels
	pixels := make([]byte, fogRes.Width*fogRes.Height*4)
	for y := 0; y < fogRes.Height; y++ {
		for x := 0; x < fogRes.Width; x++ {
			idx := (y*fogRes.Width + x) * 4
			switch fogRes.Grid[y][x] {
			case fog.Hidden:
				pixels[idx+3] = 255 // Black
			case fog.Shroud:
				pixels[idx+3] = 180 // Semi-transparent black
			}
		}
	}
	minimapFogImage.WritePixels(pixels)

	// Draw the fog image over the minimap
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(minimap.Width)/float64(fogRes.Width), float64(minimap.Height)/float64(fogRes.Height))
	op.GeoM.Translate(float64(minimap.X), float64(minimap.Y))
	screen.DrawImage(minimapFogImage, op)

	// Draw camera view
	camX := float32(minimap.X + int(cam.X*scaleX))
	camY := float32(minimap.Y + int(cam.Y*scaleY))
	camW := float32(float64(settings.ScreenWidth) * scaleX)
	camH := float32(float64(settings.ScreenHeight) * scaleY)
	vector.StrokeRect(screen, camX, camY, camW, camH, 1, color.White, false)

}

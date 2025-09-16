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
	// MinimapQuery retrieves the minimap entity.
	MinimapQuery = donburi.NewQuery(filter.Contains(components.MinimapRes))

	// minimapFogImage is a pre-rendered image of the fog of war for the minimap.
	minimapFogImage *ebiten.Image
)

// UpdateMinimap handles user input on the minimap, such as moving the camera or commanding units.
func UpdateMinimap(ecs *ecs.ECS) {
	minimapEntry, ok := MinimapQuery.First(ecs.World)
	if !ok {
		return
	}
	minimap := components.MinimapRes.Get(minimapEntry)

	// A left-click on the minimap moves the camera to the corresponding world position.
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

	// A right-click on the minimap commands all selected units to move to the corresponding world position.
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

// DrawMinimap renders the minimap, including the terrain, units, spice, fog of war, and camera view.
func DrawMinimap(ecs *ecs.ECS, screen *ebiten.Image) {
	minimapEntry, ok := MinimapQuery.First(ecs.World)
	if !ok {
		return
	}
	minimap := components.MinimapRes.Get(minimapEntry)

	settings := settings.GetSettings(ecs.World)
	cameraEntry, _ := camera.CameraQuery.First(ecs.World)
	cam := camera.CameraRes.Get(cameraEntry)

	// Draw the minimap's sand-colored background.
	vector.DrawFilledRect(screen, float32(minimap.X), float32(minimap.Y), float32(minimap.Width), float32(minimap.Height), color.RGBA{210, 180, 140, 255}, false)

	// Calculate the scaling factors to convert world coordinates to minimap coordinates.
	scaleX := float64(minimap.Width) / float64(settings.MapWidth)
	scaleY := float64(minimap.Height) / float64(settings.MapHeight)

	fogRes := fog.GetFog(ecs.World)

	// Draw all visible units on the minimap as green dots.
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

	// Draw all visible spice fields on the minimap as orange dots.
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

	// Lazily initialize the pre-rendered fog image for the minimap if it doesn't exist or its size is incorrect.
	if minimapFogImage == nil || minimapFogImage.Bounds().Dx() != fogRes.Width || minimapFogImage.Bounds().Dy() != fogRes.Height {
		minimapFogImage = ebiten.NewImage(fogRes.Width, fogRes.Height)
	}

	// Update the fog image's pixels based on the current state of the fog of war.
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

	// Draw the pre-rendered fog image over the minimap.
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(minimap.Width)/float64(fogRes.Width), float64(minimap.Height)/float64(fogRes.Height))
	op.GeoM.Translate(float64(minimap.X), float64(minimap.Y))
	screen.DrawImage(minimapFogImage, op)

	// Draw a rectangle on the minimap to represent the camera's current view.
	camX := float32(minimap.X + int(cam.X*scaleX))
	camY := float32(minimap.Y + int(cam.Y*scaleY))
	camW := float32(float64(settings.ScreenWidth) * scaleX)
	camH := float32(float64(settings.ScreenHeight) * scaleY)
	vector.StrokeRect(screen, camX, camY, camW, camH, 1, color.White, false)

	// Draw a border around the minimap.
	vector.StrokeRect(screen, float32(minimap.X), float32(minimap.Y), float32(minimap.Width), float32(minimap.Height), 1, color.White, false)

}

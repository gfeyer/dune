package systems

import (
	"image/color"
	"math"

	"github.com/gfeyer/ebit/internal/camera"
	"github.com/gfeyer/ebit/internal/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

const (
	LayerDefault ecs.LayerID = iota
)

var (
	qSprites = donburi.NewQuery(filter.Contains(components.Position, components.Sprite, components.Velocity))
)

func Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	screen.Fill(color.RGBA{210, 180, 140, 255}) // sand color

	cameraEntry, _ := camera.CameraQuery.First(ecs.World)
	cam := camera.CameraRes.Get(cameraEntry)

	qSprites.Each(ecs.World, func(entry *donburi.Entry) {
		p := components.Position.Get(entry)
		img := components.Sprite.Get(entry)

		v := components.Velocity.Get(entry)

		// Draw sprite
		op := &ebiten.DrawImageOptions{}

		// Rotate sprite if it's moving
		if v.X != 0 || v.Y != 0 {
			bounds := (*img).Bounds()
			centerX, centerY := float64(bounds.Dx())/2, float64(bounds.Dy())/2

			// Translate to center for rotation
			op.GeoM.Translate(-centerX, -centerY)
			// Rotate
			op.GeoM.Rotate(math.Atan2(v.Y, v.X) + math.Pi/2)
			// Translate to final position
			op.GeoM.Translate(p.X-cam.X+centerX, p.Y-cam.Y+centerY)
		} else {
			op.GeoM.Translate(p.X-cam.X, p.Y-cam.Y)
		}

		screen.DrawImage(*img, op)

		// Draw selection indicator
		if components.SelectableRes.Get(entry).Selected {
			bounds := (*img).Bounds()
			vector.StrokeRect(screen, float32(p.X-cam.X), float32(p.Y-cam.Y), float32(bounds.Dx()), float32(bounds.Dy()), 1, color.White, false)
		}
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

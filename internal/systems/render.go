package systems

import (
	"image/color"

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
	qSprites = donburi.NewQuery(filter.Contains(components.Position, components.Sprite))
)

func Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 12, 17, 255}) // clear

	cameraEntry, _ := camera.CameraQuery.First(ecs.World)
	cam := camera.CameraRes.Get(cameraEntry)

	qSprites.Each(ecs.World, func(entry *donburi.Entry) {
		p := components.Position.Get(entry)
		img := components.Sprite.Get(entry)

		// Draw sprite
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(p.X-cam.X, p.Y-cam.Y)
		screen.DrawImage(*img, op)

		// Draw selection indicator
		if components.SelectableRes.Get(entry).Selected {
			bounds := (*img).Bounds()
			vector.StrokeRect(screen, float32(p.X-cam.X), float32(p.Y-cam.Y), float32(bounds.Dx()), float32(bounds.Dy()), 1, color.White, false)
		}
	})
}

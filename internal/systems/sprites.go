package systems

import (
	"math"

	"github.com/gfeyer/ebit/internal/camera"
	"github.com/gfeyer/ebit/internal/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var (
	qUnits = donburi.NewQuery(filter.And(
		filter.Contains(components.Position, components.Sprite, components.UnitRes),
	))
	qBuildings = donburi.NewQuery(filter.And(
		filter.Contains(components.Position, components.Sprite),
		filter.Or(
			filter.Contains(components.RefineryRes),
			filter.Contains(components.BarracksRes),
		),
	))
)

func DrawBuildings(ecs *ecs.ECS, screen *ebiten.Image) {
	cameraEntry, _ := camera.CameraQuery.First(ecs.World)
	cam := camera.CameraRes.Get(cameraEntry)

	qBuildings.Each(ecs.World, func(entry *donburi.Entry) {
		p := components.Position.Get(entry)
		img := components.Sprite.Get(entry)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(p.X-cam.X, p.Y-cam.Y)

		// Handle selection tint
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

func DrawSpice(ecs *ecs.ECS, screen *ebiten.Image) {
	cameraEntry, _ := camera.CameraQuery.First(ecs.World)
	cam := camera.CameraRes.Get(cameraEntry)

	qSpice.Each(ecs.World, func(entry *donburi.Entry) {
		p := components.Position.Get(entry)
		img := components.Sprite.Get(entry)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(p.X-cam.X, p.Y-cam.Y)
		screen.DrawImage(*img, op)
	})
}

func DrawUnits(ecs *ecs.ECS, screen *ebiten.Image) {
	cameraEntry, _ := camera.CameraQuery.First(ecs.World)
	cam := camera.CameraRes.Get(cameraEntry)

	qUnits.Each(ecs.World, func(entry *donburi.Entry) {
		p := components.Position.Get(entry)
		img := components.Sprite.Get(entry)
		op := &ebiten.DrawImageOptions{}

		// Handle rotation for entities that have velocity
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

		// Handle selection tint
		if components.SelectableRes.Get(entry).Selected {
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

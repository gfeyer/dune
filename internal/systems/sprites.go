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

		if !isSpriteInView(p, *img, cam, screen) {
			return
		}

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

		if !isSpriteInView(p, *img, cam, screen) {
			return
		}

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

		if !isSpriteInView(p, *img, cam, screen) {
			return
		}

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

func isSpriteInView(p *components.Pos, img *ebiten.Image, cam *camera.Camera, screen *ebiten.Image) bool {
	screenW, screenH := screen.Bounds().Dx(), screen.Bounds().Dy()
	spriteW, spriteH := img.Bounds().Dx(), img.Bounds().Dy()

	// Bounding box of the object in world coordinates.
	objLeft := p.X
	objRight := p.X + float64(spriteW)
	objTop := p.Y
	objBottom := p.Y + float64(spriteH)

	// Bounding box of the camera in world coordinates.
	camLeft := cam.X
	camRight := cam.X + float64(screenW)
	camTop := cam.Y
	camBottom := cam.Y + float64(screenH)

	// Check if the object's bounding box is outside the camera's bounding box.
	if objRight < camLeft || objLeft > camRight || objBottom < camTop || objTop > camBottom {
		return false
	}

	return true
}

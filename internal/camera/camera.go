package camera

import (
	"github.com/gfeyer/ebit/internal/settings"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

// Camera represents the game's camera.
type Camera struct {
	X, Y float64
}

// ScreenToWorld converts screen coordinates to world coordinates.
func (c *Camera) ScreenToWorld(x, y float64) (float64, float64) {
	return x + c.X, y + c.Y
}

var CameraRes = donburi.NewComponentType[Camera]()

var (
	CameraQuery = donburi.NewQuery(filter.Contains(CameraRes))
)

// Update handles camera movement.
func Update(ecs *ecs.ECS) {
	cameraEntry, _ := CameraQuery.First(ecs.World)
	cam := CameraRes.Get(cameraEntry)

	settings := settings.GetSettings(ecs.World)

	// Pan with arrow keys
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		cam.X -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		cam.X += 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		cam.Y -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		cam.Y += 5
	}

	// Clamp camera to map boundaries
	if cam.X < 0 {
		cam.X = 0
	}
	if cam.Y < 0 {
		cam.Y = 0
	}
	if cam.X > float64(settings.MapWidth-settings.ScreenWidth) {
		cam.X = float64(settings.MapWidth - settings.ScreenWidth)
	}
	if cam.Y > float64(settings.MapHeight-settings.ScreenHeight) {
		cam.Y = float64(settings.MapHeight - settings.ScreenHeight)
	}
}

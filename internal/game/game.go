package game

import (
	"github.com/gfeyer/ebit/internal/camera"
	"github.com/gfeyer/ebit/internal/factory"
	"github.com/gfeyer/ebit/internal/settings"
	"github.com/gfeyer/ebit/internal/systems"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type Game struct {
	ecs *ecs.ECS
}

func NewGame(w, h int) *Game {
	world := donburi.NewWorld()
	ecs := ecs.NewECS(world)

	// Register settings
	e := world.Create(settings.SettingsRes)
	entry := world.Entry(e)
	*settings.SettingsRes.Get(entry) = settings.Settings{
		ScreenWidth:  w,
		ScreenHeight: h,
		MapWidth:     w * 2,
		MapHeight:    h * 2,
	}

	// Register camera
	ce := world.Create(camera.CameraRes)
	centry := world.Entry(ce)
	*camera.CameraRes.Get(centry) = camera.Camera{}

	// Register systems
	ecs.AddSystem(systems.UpdateMovement)
	ecs.AddSystem(systems.UpdateInput)
	ecs.AddSystem(camera.Update)

	// Register renderers
	ecs.AddRenderer(systems.LayerDefault, systems.Draw)

	// Spawn initial units
	factory.CreateTrike(world, 100, 100)
	factory.CreateHarvester(world, 200, 200)

	return &Game{ecs: ecs}
}

func (g *Game) Update() error {
	g.ecs.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.ecs.Draw(screen)
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	entry, _ := settings.SettingsQuery.First(g.ecs.World)
	s := settings.SettingsRes.Get(entry)
	return s.ScreenWidth, s.ScreenHeight
}

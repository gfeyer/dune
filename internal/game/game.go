package game

import (
	"github.com/gfeyer/ebit/internal/camera"
	"github.com/gfeyer/ebit/internal/factory"
	"github.com/gfeyer/ebit/internal/settings"
	"github.com/gfeyer/ebit/internal/systems"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)


type Game struct {
	world donburi.World
}

func NewGame(w, h int) *Game {
	world := donburi.NewWorld()

	e := world.Create(settings.SettingsRes)
	entry := world.Entry(e)
	*settings.SettingsRes.Get(entry) = settings.Settings{
		ScreenWidth:  w,
		ScreenHeight: h,
		MapWidth:     w * 2,
		MapHeight:    h * 2,
	}

	// Create camera
	ce := world.Create(camera.CameraRes)
	centry := world.Entry(ce)
	*camera.CameraRes.Get(centry) = camera.Camera{}

	factory.CreateTrike(world, 100, 100)
	factory.CreateHarvester(world, 200, 200)

	g := &Game{
		world: world,
	}

	return g
}

func (g *Game) Update() error {
	systems.UpdateMovement(g.world)
	systems.UpdateInput(g.world)
	camera.Update(g.world)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	systems.Draw(g.world, screen)
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	entry, _ := settings.SettingsQuery.First(g.world)
	s := settings.SettingsRes.Get(entry)
	return s.ScreenWidth, s.ScreenHeight
}

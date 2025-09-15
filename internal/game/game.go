package game

import (
	"github.com/gfeyer/ebit/internal/settings"
	"github.com/gfeyer/ebit/internal/systems"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var (
	SettingsQuery = donburi.NewQuery(filter.Contains(settings.SettingsRes))
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
	}

	g := &Game{
		world: world,
	}

	return g
}

func (g *Game) Update() error {
	systems.UpdateMovement(g.world)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	systems.Draw(g.world, screen)
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	entry, _ := SettingsQuery.First(g.world)
	s := settings.SettingsRes.Get(entry)
	return s.ScreenWidth, s.ScreenHeight
}

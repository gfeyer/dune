package game

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/gfeyer/ebit/internal/components"
	"github.com/gfeyer/ebit/internal/systems"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type Game struct {
	world            donburi.World
	screenW, screenH int
}

func NewGame(w, h int) *Game {
	g := &Game{
		world:   donburi.NewWorld(),
		screenW: w, screenH: h,
	}
	spawnRandom(g.world, w, h, 32)
	return g
}

func (g *Game) Update() error {
	systems.UpdateMovement(g.world, g.screenW, g.screenH)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	systems.Draw(g.world, screen)
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return g.screenW, g.screenH
}

func spawnRandom(w donburi.World, screenW, screenH, n int) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		e := w.Create(components.Position, components.Velocity, components.Sprite)
		entry := w.Entry(e)

		// random position/velocity
		px := rng.Float64() * float64(screenW-24)
		py := rng.Float64() * float64(screenH-24)
		vx := (rng.Float64()*2 - 1) * 120 // px/sec
		vy := (rng.Float64()*2 - 1) * 120 // px/sec

		*components.Position.Get(entry) = components.Pos{X: px, Y: py}
		*components.Velocity.Get(entry) = components.Vel{X: vx, Y: vy}

		// tiny colored square sprite
		size := 16 + rng.Intn(12)
		img := ebiten.NewImage(size, size)
		img.Fill(color.RGBA{
			R: uint8(rng.Intn(80) + 120),
			G: uint8(rng.Intn(80) + 120),
			B: uint8(rng.Intn(80) + 120),
			A: 255,
		})
		*components.Sprite.Get(entry) = img
	}
}

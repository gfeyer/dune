package main

import (
	"math/rand"
	"time"

	"github.com/gfeyer/ebit/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	const W, H = 1280, 720
	rand.Seed(time.Now().UnixNano()) // seed once
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Dune II")
	ebiten.SetTPS(60)

	g := game.NewGame(W, H)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}

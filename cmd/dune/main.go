package main

import (
	"github.com/gfeyer/ebit/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	const W, H = 1280, 720

	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Dune II")

	g := game.NewGame(W, H)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}

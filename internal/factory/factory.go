package factory

import (
	"image/color"

	"github.com/gfeyer/ebit/internal/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

func CreateHarvester(w donburi.World, x, y float64) {
	e := w.Create(components.Position, components.Sprite, components.UnitRes, components.SelectableRes, components.TargetRes, components.Velocity)
	entry := w.Entry(e)

	// Harvester is a blue square
	img := ebiten.NewImage(24, 24)
	img.Fill(color.RGBA{R: 0, G: 0, B: 255, A: 255})

	*components.Position.Get(entry) = components.Pos{X: x, Y: y}
	*components.Sprite.Get(entry) = img
	*components.UnitRes.Get(entry) = components.Unit{Type: components.Harvester}
	*components.SelectableRes.Get(entry) = components.Selectable{Selected: false}
	*components.Velocity.Get(entry) = components.Vel{}
}

func CreateTrike(w donburi.World, x, y float64) {
	e := w.Create(components.Position, components.Sprite, components.UnitRes, components.SelectableRes, components.TargetRes, components.Velocity)
	entry := w.Entry(e)

	// Trike is a blue triangle
	img := ebiten.NewImage(24, 24)
	triangle := []ebiten.Vertex{
		{DstX: 12, DstY: 2, SrcX: 1, SrcY: 1, ColorR: 0, ColorG: 0, ColorB: 1, ColorA: 1},
		{DstX: 2, DstY: 22, SrcX: 1, SrcY: 1, ColorR: 0, ColorG: 0, ColorB: 1, ColorA: 1},
		{DstX: 22, DstY: 22, SrcX: 1, SrcY: 1, ColorR: 0, ColorG: 0, ColorB: 1, ColorA: 1},
	}
	whiteSubimage := ebiten.NewImage(1, 1)
	whiteSubimage.Fill(color.White)
	img.DrawTriangles(triangle, []uint16{0, 1, 2}, whiteSubimage, &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.FillAll,
	})

	*components.Position.Get(entry) = components.Pos{X: x, Y: y}
	*components.Sprite.Get(entry) = img
	*components.UnitRes.Get(entry) = components.Unit{Type: components.Trike}
	*components.SelectableRes.Get(entry) = components.Selectable{Selected: false}
	*components.Velocity.Get(entry) = components.Vel{}
}

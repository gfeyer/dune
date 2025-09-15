package factory

import (
	"fmt"
	"image/color"

	"github.com/gfeyer/ebit/internal/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/yohamta/donburi"
	"golang.org/x/image/font/basicfont"
)

func CreateHarvester(w donburi.World, x, y float64) {
	e := w.Create(components.Position, components.Sprite, components.UnitRes, components.SelectableRes, components.TargetRes, components.Velocity, components.HarvesterRes, components.HealthRes)
	entry := w.Entry(e)

	// Harvester is a blue square
	img := ebiten.NewImage(16, 16)
	img.Fill(color.RGBA{R: 0, G: 0, B: 255, A: 255})

	*components.Position.Get(entry) = components.Pos{X: x, Y: y}
	*components.Sprite.Get(entry) = img
	*components.UnitRes.Get(entry) = components.Unit{Type: components.Harvester}
	*components.SelectableRes.Get(entry) = components.Selectable{Selected: false}
	*components.Velocity.Get(entry) = components.Vel{}
	*components.HarvesterRes.Get(entry) = components.HarvesterData{Capacity: 100}
	*components.HealthRes.Get(entry) = components.Health{Current: 100, Max: 100}
}

func CreateTrike(w donburi.World, x, y float64) {
	e := w.Create(components.Position, components.Sprite, components.UnitRes, components.SelectableRes, components.TargetRes, components.Velocity, components.HealthRes)
	entry := w.Entry(e)

	// Trike is a blue triangle
	img := ebiten.NewImage(24, 24)
	r, g, b, a := color.RGBA{R: 0, G: 0, B: 255, A: 255}.RGBA()
	triangle := []ebiten.Vertex{
		{DstX: 12, DstY: 2, SrcX: 0, SrcY: 0, ColorR: float32(r) / 0xffff, ColorG: float32(g) / 0xffff, ColorB: float32(b) / 0xffff, ColorA: float32(a) / 0xffff},
		{DstX: 2, DstY: 22, SrcX: 0, SrcY: 0, ColorR: float32(r) / 0xffff, ColorG: float32(g) / 0xffff, ColorB: float32(b) / 0xffff, ColorA: float32(a) / 0xffff},
		{DstX: 22, DstY: 22, SrcX: 0, SrcY: 0, ColorR: float32(r) / 0xffff, ColorG: float32(g) / 0xffff, ColorB: float32(b) / 0xffff, ColorA: float32(a) / 0xffff},
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
	*components.HealthRes.Get(entry) = components.Health{Current: 50, Max: 50}
}

func CreateSpice(w donburi.World, x, y float64) {
	e := w.Create(components.Position, components.Sprite, components.SpiceRes, components.Velocity, components.SelectableRes, components.SpiceAmountRes)
	entry := w.Entry(e)

	// Spice is an orange square
	img := ebiten.NewImage(32, 32)
	img.Fill(color.RGBA{R: 210, G: 105, B: 30, A: 255})

	*components.Position.Get(entry) = components.Pos{X: x, Y: y}
	*components.Sprite.Get(entry) = img
	*components.SpiceRes.Get(entry) = components.Spice{}
	*components.Velocity.Get(entry) = components.Vel{}
	*components.SelectableRes.Get(entry) = components.Selectable{Selected: false}
	*components.SpiceAmountRes.Get(entry) = components.SpiceAmount{Amount: 1000}
}

func CreateBuildOption(w donburi.World, btype components.BuildingType, name string, cost int, width, height int) {
	e := w.Create(components.BuildInfoRes)
	entry := w.Entry(e)

	icon := ebiten.NewImage(width, height)
	bgColor := color.RGBA{R: 128, G: 128, B: 128, A: 255} // Gray background
	icon.Fill(bgColor)

	// Draw text on the icon
	nameText := name
	costText := fmt.Sprintf("Cost: %d", cost)

	// Center the name text
	nameBounds := text.BoundString(basicfont.Face7x13, nameText)
	nameX := (width - nameBounds.Dx()) / 2
	text.Draw(icon, nameText, basicfont.Face7x13, nameX, 15, color.White)

	// Center the cost text
	costBounds := text.BoundString(basicfont.Face7x13, costText)
	costX := (width - costBounds.Dx()) / 2
	text.Draw(icon, costText, basicfont.Face7x13, costX, 30, color.White)

	*components.BuildInfoRes.Get(entry) = components.BuildInfo{
		Type: btype,
		Name: name,
		Cost: cost,
		Icon: icon,
	}
}

func CreateBarracks(w donburi.World, x, y float64) {
	e := w.Create(components.Position, components.Sprite, components.BarracksRes)
	entry := w.Entry(e)

	// Barracks is a red square
	img := ebiten.NewImage(64, 64)
	img.Fill(color.RGBA{R: 255, G: 0, B: 0, A: 255})

	*components.Position.Get(entry) = components.Pos{X: x, Y: y}
	*components.Sprite.Get(entry) = img
	*components.BarracksRes.Get(entry) = components.Barracks{}
}

func CreateRefinery(w donburi.World, x, y float64) {
	e := w.Create(components.Position, components.Sprite, components.RefineryRes)
	entry := w.Entry(e)

	// Refinery is a gray square
	img := ebiten.NewImage(64, 64)
	img.Fill(color.RGBA{R: 128, G: 128, B: 128, A: 255})

	*components.Position.Get(entry) = components.Pos{X: x, Y: y}
	*components.Sprite.Get(entry) = img
	*components.RefineryRes.Get(entry) = components.Refinery{}
}

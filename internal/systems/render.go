package systems

import (
	"image/color"

	"github.com/gfeyer/ebit/internal/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var (
	qSprites = donburi.NewQuery(filter.Contains(components.Position, components.Sprite))
)

func Draw(w donburi.World, screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 12, 17, 255}) // clear

	qSprites.Each(w, func(entry *donburi.Entry) {
		p := components.Position.Get(entry)
		img := components.Sprite.Get(entry)

		var op ebiten.DrawImageOptions
		op.GeoM.Translate(p.X, p.Y)
		screen.DrawImage(*img, &op)
	})
}

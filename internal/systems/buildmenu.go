package systems

import (
	"github.com/gfeyer/ebit/internal/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var (
	BuildMenuQuery = donburi.NewQuery(filter.Contains(components.BuildInfoRes))
)

func DrawBuildMenu(ecs *ecs.ECS, screen *ebiten.Image) {
	// Find minimap to position the build menu below it
	minimapEntry, ok := MinimapQuery.First(ecs.World)
	if !ok {
		return // No minimap, no build menu
	}
	minimap := components.MinimapRes.Get(minimapEntry)

	// Menu layout
	menuX := minimap.X
	menuY := minimap.Y + minimap.Height + 10 // 10px padding
	padding := 5
	colWidth := (minimap.Width - padding) / 2
	rowHeight := colWidth

	i := 0
	BuildMenuQuery.Each(ecs.World, func(entry *donburi.Entry) {
		buildInfo := components.BuildInfoRes.Get(entry)

		col := i % 2
		row := i / 2

		// Calculate position for the icon
		iconX := menuX + col*colWidth
		iconY := menuY + row*rowHeight

		// Draw icon
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(iconX), float64(iconY))
		screen.DrawImage(buildInfo.Icon, opts)

		i++
	})
}

package systems

import (
	"fmt"
	"image/color"

	"github.com/gfeyer/ebit/internal/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/font/basicfont"
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
	iconWidth := 32
	iconHeight := 32
	padding := 5
	colWidth := iconWidth + padding
	rowHeight := iconHeight + padding + 15 // Extra space for text

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

		// Draw cost
		costText := fmt.Sprintf("$%d", buildInfo.Cost)
		text.Draw(screen, costText, basicfont.Face7x13, iconX, iconY+iconHeight+12, color.White)

		i++
	})
}

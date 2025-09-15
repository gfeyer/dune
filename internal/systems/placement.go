package systems

import (
	"image/color"

	"github.com/gfeyer/ebit/internal/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi/ecs"
)

func DrawPlacement(ecs *ecs.ECS, screen *ebiten.Image) {
	placementEntry, ok := PlacementQuery.First(ecs.World)
	if !ok {
		return
	}
	placement := components.PlacementRes.Get(placementEntry)

	if placement.IsPlacing {
		mx, my := ebiten.CursorPosition()

		// Get icon dimensions
		bounds := placement.Icon.Bounds()
		width := float32(bounds.Dx())
		height := float32(bounds.Dy())

		// Draw a white rectangle as the footprint
		vector.DrawFilledRect(screen, float32(mx), float32(my), width, height, color.White, false)
	}
}

package systems

import (
	"github.com/gfeyer/ebit/internal/components"
	"github.com/hajimehoshi/ebiten/v2"
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

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(mx), float64(my))
		screen.DrawImage(placement.Icon, opts)
	}
}

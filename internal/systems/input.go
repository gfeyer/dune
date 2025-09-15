package systems

import (
	"github.com/gfeyer/ebit/internal/camera"
	"github.com/gfeyer/ebit/internal/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var (
	qSelectable = donburi.NewQuery(filter.Contains(components.Position, components.Sprite, components.SelectableRes))
)

func UpdateInput(w donburi.World) {
	// Left-click to select
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Get camera and mouse position
		cameraEntry, _ := camera.CameraQuery.First(w)
		cam := camera.CameraRes.Get(cameraEntry)
		mx, my := ebiten.CursorPosition()
		wx, wy := float64(mx)+cam.X, float64(my)+cam.Y

		clickedOnUnit := false
		qSelectable.Each(w, func(entry *donburi.Entry) {
			p := components.Position.Get(entry)
			s := components.Sprite.Get(entry)
			sel := components.SelectableRes.Get(entry)

			// Simple bounding box collision
			bounds := (*s).Bounds()
			if wx >= p.X && wx < p.X+float64(bounds.Dx()) && wy >= p.Y && wy < p.Y+float64(bounds.Dy()) {
				sel.Selected = true
				clickedOnUnit = true
			} else {
				sel.Selected = false
			}
		})

		// If we didn't click on any unit, deselect all
		if !clickedOnUnit {
			qSelectable.Each(w, func(entry *donburi.Entry) {
				components.SelectableRes.Get(entry).Selected = false
			})
		}
	}

	// Right-click to move
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		cameraEntry, _ := camera.CameraQuery.First(w)
		cam := camera.CameraRes.Get(cameraEntry)
		mx, my := ebiten.CursorPosition()
		wx, wy := float64(mx)+cam.X, float64(my)+cam.Y

		qSelectable.Each(w, func(entry *donburi.Entry) {
			if components.SelectableRes.Get(entry).Selected {
				*components.TargetRes.Get(entry) = components.Target{X: wx, Y: wy}
			}
		})
	}
}

package systems

import (
	"github.com/gfeyer/ebit/internal/camera"
	"github.com/gfeyer/ebit/internal/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"image"
)


func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func UpdateInput(ecs *ecs.ECS) {
	dragEntry, _ := QDrag.First(ecs.World)
	drag := components.DragRes.Get(dragEntry)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		drag.IsDragging = true
		drag.StartX, drag.StartY = ebiten.CursorPosition()
		drag.EndX, drag.EndY = drag.StartX, drag.StartY
	}

	if drag.IsDragging {
		drag.EndX, drag.EndY = ebiten.CursorPosition()
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		drag.IsDragging = false

		// Deselect all units first
		QSelectable.Each(ecs.World, func(entry *donburi.Entry) {
			components.SelectableRes.Get(entry).Selected = false
		})

		cameraEntry, _ := camera.CameraQuery.First(ecs.World)
		cam := camera.CameraRes.Get(cameraEntry)

		// If the mouse moved, it's a drag selection
		if abs(drag.StartX-drag.EndX) > 5 || abs(drag.StartY-drag.EndY) > 5 {
			rect := image.Rect(drag.StartX, drag.StartY, drag.EndX, drag.EndY).Canon()
			QSelectable.Each(ecs.World, func(entry *donburi.Entry) {
				p := components.Position.Get(entry)
				screenX, screenY := int(p.X-cam.X), int(p.Y-cam.Y)
				if image.Pt(screenX, screenY).In(rect) {
					components.SelectableRes.Get(entry).Selected = true
				}
			})
		} else { // Otherwise, it's a single click
			mx, my := ebiten.CursorPosition()
			wx, wy := float64(mx)+cam.X, float64(my)+cam.Y

			// Find the top-most unit under the cursor
			var topUnit *donburi.Entry
			QSelectable.Each(ecs.World, func(entry *donburi.Entry) {
				p := components.Position.Get(entry)
				s := components.Sprite.Get(entry)

				bounds := (*s).Bounds()
				if wx >= p.X && wx < p.X+float64(bounds.Dx()) && wy >= p.Y && wy < p.Y+float64(bounds.Dy()) {
					topUnit = entry
				}
			})

			if topUnit != nil {
				components.SelectableRes.Get(topUnit).Selected = true
			}
		}
	}

	// Right-click to move
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		cameraEntry, _ := camera.CameraQuery.First(ecs.World)
		cam := camera.CameraRes.Get(cameraEntry)
		mx, my := ebiten.CursorPosition()
		wx, wy := float64(mx)+cam.X, float64(my)+cam.Y

		QSelectable.Each(ecs.World, func(entry *donburi.Entry) {
			if components.SelectableRes.Get(entry).Selected {
				*components.TargetRes.Get(entry) = components.Target{X: wx, Y: wy}
			}
		})
	}
}

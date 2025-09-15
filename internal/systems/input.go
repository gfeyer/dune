package systems

import (
	"image"

	"github.com/gfeyer/ebit/internal/camera"
	"github.com/gfeyer/ebit/internal/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
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
	}

	if drag.IsDragging {
		drag.EndX, drag.EndY = ebiten.CursorPosition()
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		drag.IsDragging = false

		cameraEntry, _ := camera.CameraQuery.First(ecs.World)
		cam := camera.CameraRes.Get(cameraEntry)

		// Drag selection
		if abs(drag.StartX-drag.EndX) > 5 || abs(drag.StartY-drag.EndY) > 5 {
			// Deselect all units first
			QSelectable.Each(ecs.World, func(entry *donburi.Entry) {
				components.SelectableRes.Get(entry).Selected = false
			})

			rect := image.Rect(drag.StartX, drag.StartY, drag.EndX, drag.EndY).Canon()
			QSelectable.Each(ecs.World, func(entry *donburi.Entry) {
				p := components.Position.Get(entry)
				screenX, screenY := int(p.X-cam.X), int(p.Y-cam.Y)
				if image.Pt(screenX, screenY).In(rect) {
					components.SelectableRes.Get(entry).Selected = true
				}
			})
		} else { // Single-click selection
			mx, my := ebiten.CursorPosition()
			wx, wy := float64(mx)+cam.X, float64(my)+cam.Y

			var clickedUnit *donburi.Entry
			QSelectable.Each(ecs.World, func(entry *donburi.Entry) {
				p := components.Position.Get(entry)
				s := components.Sprite.Get(entry)
				bounds := (*s).Bounds()
				if wx >= p.X && wx < p.X+float64(bounds.Dx()) && wy >= p.Y && wy < p.Y+float64(bounds.Dy()) {
					clickedUnit = entry
				}
			})

			// Deselect all units
			QSelectable.Each(ecs.World, func(entry *donburi.Entry) {
				components.SelectableRes.Get(entry).Selected = false
			})

			// Select only the clicked unit
			if clickedUnit != nil {
				components.SelectableRes.Get(clickedUnit).Selected = true
			}
		}
	}

	// Right-click to move or harvest
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		cameraEntry, _ := camera.CameraQuery.First(ecs.World)
		cam := camera.CameraRes.Get(cameraEntry)
		mx, my := ebiten.CursorPosition()
		wx, wy := float64(mx)+cam.X, float64(my)+cam.Y

		// Check if the click was on a spice field
		var targetSpice *donburi.Entry
		QSpice.Each(ecs.World, func(spiceEntry *donburi.Entry) {
			spicePos := components.Position.Get(spiceEntry)
			spiceSprite := components.Sprite.Get(spiceEntry)
			bounds := (*spiceSprite).Bounds()
			if wx >= spicePos.X && wx < spicePos.X+float64(bounds.Dx()) && wy >= spicePos.Y && wy < spicePos.Y+float64(bounds.Dy()) {
				targetSpice = spiceEntry
			}
		})

		QSelectable.Each(ecs.World, func(entry *donburi.Entry) {
			if components.SelectableRes.Get(entry).Selected {
				unit := components.UnitRes.Get(entry)
				// If it's a harvester and a spice field was clicked, set it as the target
				if unit.Type == components.Harvester && targetSpice != nil {
					harvester := components.HarvesterRes.Get(entry)
					harvester.State = components.StateMovingToSpice
					harvester.TargetSpice = targetSpice.Entity()
					spicePos := components.Position.Get(targetSpice)
					*components.TargetRes.Get(entry) = components.Target{X: spicePos.X, Y: spicePos.Y}
				} else { // Otherwise, it's a normal move command
					*components.TargetRes.Get(entry) = components.Target{X: wx, Y: wy}
					// If it was a harvester, clear its spice target
					if unit.Type == components.Harvester {
						harvester := components.HarvesterRes.Get(entry)
						harvester.State = components.StateIdle
						harvester.TargetSpice = 0
					}
				}
			}
		})
	}
}

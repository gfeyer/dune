package systems

import (
	"math"

	"github.com/gfeyer/ebit/internal/components"
	"github.com/gfeyer/ebit/internal/settings"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var (
	qMovers  = donburi.NewQuery(filter.Contains(components.Position, components.Velocity, components.Sprite))
	qSettings = donburi.NewQuery(filter.Contains(settings.SettingsRes))
)

func UpdateMovement(w donburi.World) {
	const dt = 1.0 / 60.0

	settingsEntry, _ := qSettings.First(w)
	s := settings.SettingsRes.Get(settingsEntry)

	qMovers.Each(w, func(entry *donburi.Entry) {
		p := components.Position.Get(entry)
		v := components.Velocity.Get(entry)
		t := components.TargetRes.Get(entry)

		// If there's a target, move towards it
		if t.X != 0 || t.Y != 0 {
			dx := t.X - p.X
			dy := t.Y - p.Y
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist < 5 { // Arrived
				v.X, v.Y = 0, 0
				*t = components.Target{}
			} else {
				v.X = (dx / dist) * 120
				v.Y = (dy / dist) * 120
			}
		}

		// integrate
		p.X += v.X * dt
		p.Y += v.Y * dt

		// simple map bounds
		if p.X < 0 {
			p.X = 0
		}
		if p.Y < 0 {
			p.Y = 0
		}
		if p.X > float64(s.MapWidth) {
			p.X = float64(s.MapWidth)
		}
		if p.Y > float64(s.MapHeight) {
			p.Y = float64(s.MapHeight)
		}
	})
}

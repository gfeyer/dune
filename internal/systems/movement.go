package systems

import (
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

		// integrate
		p.X += v.X * dt
		p.Y += v.Y * dt

		// bounce off edges (simple bounds)
		sw, sh := (*components.Sprite.Get(entry)).Bounds().Dx(), (*components.Sprite.Get(entry)).Bounds().Dy()
		if p.X < 0 {
			p.X = 0
			v.X = -v.X
		}
		if p.Y < 0 {
			p.Y = 0
			v.Y = -v.Y
		}
		if p.X+float64(sw) > float64(s.ScreenWidth) {
			p.X = float64(s.ScreenWidth - sw)
			v.X = -v.X
		}
		if p.Y+float64(sh) > float64(s.ScreenHeight) {
			p.Y = float64(s.ScreenHeight - sh)
			v.Y = -v.Y
		}
	})
}

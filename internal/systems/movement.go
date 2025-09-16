package systems

import (
	"math"

	"github.com/gfeyer/ebit/internal/components"
	"github.com/gfeyer/ebit/internal/settings"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var (
	// qMovers retrieves all entities that have position, velocity, and a target, making them capable of movement.
	qMovers = donburi.NewQuery(filter.Contains(components.Position, components.Velocity, components.Sprite, components.TargetRes))
	// qSettings retrieves the game settings entity.
	qSettings = donburi.NewQuery(filter.Contains(settings.SettingsRes))
)

// UpdateMovement handles the movement of all units in the game.
// It calculates the velocity needed to reach a target and updates the unit's position accordingly.
func UpdateMovement(ecs *ecs.ECS) {
	const dt = 1.0 / 60.0

	s := settings.GetSettings(ecs.World)
	
	qMovers.Each(ecs.World, func(entry *donburi.Entry) {
		p := components.Position.Get(entry)
		v := components.Velocity.Get(entry)
		t := components.TargetRes.Get(entry)

		// If the entity has a target, calculate the velocity to move towards it.
		if t.X != 0 || t.Y != 0 {
			dx := t.X - p.X
			dy := t.Y - p.Y
			dist := math.Sqrt(dx*dx + dy*dy)

			// If the unit is close enough to the target, stop moving.
			if dist < 5 { // Arrived
				v.X, v.Y = 0, 0
				*t = components.Target{}
			} else {
				// Otherwise, set the velocity to move towards the target at a constant speed.
				v.X = (dx / dist) * 240
				v.Y = (dy / dist) * 240
			}
		}

		// Update the position based on the current velocity and the time delta.
		p.X += v.X * dt
		p.Y += v.Y * dt

		// Ensure the unit stays within the map boundaries.
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

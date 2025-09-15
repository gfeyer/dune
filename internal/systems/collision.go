package systems

import (
	"math"

	"github.com/gfeyer/ebit/internal/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var qCollision = donburi.NewQuery(filter.Contains(components.Position, components.UnitRes))

func ResolveCollisions(ecs *ecs.ECS) {
	qCollision.Each(ecs.World, func(entry *donburi.Entry) {
		p1 := components.Position.Get(entry)

		qCollision.Each(ecs.World, func(other *donburi.Entry) {
			if entry.Entity() == other.Entity() {
				return
			}

			p2 := components.Position.Get(other)

			dx := p1.X - p2.X
			dy := p1.Y - p2.Y
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist < 24 {
				overlap := (24 - dist) / 2
				p1.X += dx / dist * overlap
				p1.Y += dy / dist * overlap
				p2.X -= dx / dist * overlap
				p2.Y -= dy / dist * overlap
			}
		})
	})
}

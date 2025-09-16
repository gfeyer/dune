package systems

import (
	"math"

	"github.com/gfeyer/ebit/internal/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var qCollision = donburi.NewQuery(filter.Contains(components.Position, components.UnitRes, components.Sprite))

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

			s1 := components.Sprite.Get(entry)
			s2 := components.Sprite.Get(other)
			radius1 := float64((*s1).Bounds().Dx()) / 2
			radius2 := float64((*s2).Bounds().Dx()) / 2
			requiredDist := (radius1 + radius2) * 0.5 // Allow 50% overlap

			if dist < requiredDist {
				overlap := (requiredDist - dist) / 2
				p1.X += dx / dist * overlap
				p1.Y += dy / dist * overlap
				p2.X -= dx / dist * overlap
				p2.Y -= dy / dist * overlap
			}
		})
	})
}

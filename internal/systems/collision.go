package systems

import (
	"math"

	"github.com/gfeyer/ebit/internal/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

// qCollision is a query that retrieves all entities with Position, UnitRes, and Sprite components, which are necessary for collision detection.
var qCollision = donburi.NewQuery(filter.Contains(components.Position, components.UnitRes, components.Sprite))

// ResolveCollisions handles the collision detection and resolution between units.
// It iterates through all pairs of units and pushes them apart if they overlap.
func ResolveCollisions(ecs *ecs.ECS) {
	// Iterate over each entity that can collide.
	qCollision.Each(ecs.World, func(entry *donburi.Entry) {
		p1 := components.Position.Get(entry)

		// Compare it with every other entity that can collide.
		qCollision.Each(ecs.World, func(other *donburi.Entry) {
			// Don't check for collision with itself.
			if entry.Entity() == other.Entity() {
				return
			}

			p2 := components.Position.Get(other)

			// Calculate the distance between the two entities.
			dx := p1.X - p2.X
			dy := p1.Y - p2.Y
			dist := math.Sqrt(dx*dx + dy*dy)

			// Calculate the radii of the two entities based on their sprite size.
			s1 := components.Sprite.Get(entry)
			s2 := components.Sprite.Get(other)
			radius1 := float64((*s1).Bounds().Dx()) / 2
			radius2 := float64((*s2).Bounds().Dx()) / 2
			// The required distance is half the sum of the radii, allowing for 50% overlap.
			requiredDist := (radius1 + radius2) * 0.5 // Allow 50% overlap

			// If the distance is less than the required distance, a collision has occurred.
			if dist < requiredDist {
				// Calculate the overlap and push the entities apart.
				overlap := (requiredDist - dist) / 2
				p1.X += dx / dist * overlap
				p1.Y += dy / dist * overlap
				p2.X -= dx / dist * overlap
				p2.Y -= dy / dist * overlap
			}
		})
	})
}

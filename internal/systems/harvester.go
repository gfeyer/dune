package systems

import (
	"math"

	"github.com/gfeyer/ebit/internal/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var (
	qHarvesters = donburi.NewQuery(filter.Contains(components.UnitRes, components.HarvesterRes, components.Position, components.TargetRes))
	qSpice    = donburi.NewQuery(filter.Contains(components.SpiceRes, components.SpiceAmountRes, components.Position))
	qRefinery = donburi.NewQuery(filter.Contains(components.RefineryRes, components.Position))
)

func UpdateHarvester(ecs *ecs.ECS) {
	qHarvesters.Each(ecs.World, func(entry *donburi.Entry) {
		unit := components.UnitRes.Get(entry)
		if unit.Type != components.Harvester {
			return
		}

		harvester := components.HarvesterRes.Get(entry)
		p := components.Position.Get(entry)
		t := components.TargetRes.Get(entry)

		switch harvester.State {
		case components.StateIdle:
			// Find nearest spice
			var nearestSpice *donburi.Entry
			minDist := math.MaxFloat64
			qSpice.Each(ecs.World, func(spiceEntry *donburi.Entry) {
				spicePos := components.Position.Get(spiceEntry)
				dx := spicePos.X - p.X
				dy := spicePos.Y - p.Y
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist < minDist {
					minDist = dist
					nearestSpice = spiceEntry
				}
			})

			if nearestSpice != nil {
				harvester.State = components.StateMovingToSpice
				harvester.TargetSpice = nearestSpice.Entity()
				spicePos := components.Position.Get(nearestSpice)
				t.X, t.Y = spicePos.X, spicePos.Y
			}
		case components.StateMovingToSpice:
			// Check if arrived at spice
			targetSpiceEntry := ecs.World.Entry(harvester.TargetSpice)
			targetSpicePos := components.Position.Get(targetSpiceEntry)
			dx := targetSpicePos.X - p.X
			dy := targetSpicePos.Y - p.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < 32 { // Close enough to harvest
				harvester.State = components.StateHarvesting
				harvester.HarvestTimer = 180 // 3 seconds
			}
		case components.StateHarvesting:
			harvester.HarvestTimer--
			if harvester.HarvestTimer <= 0 {
				targetSpiceEntry := ecs.World.Entry(harvester.TargetSpice)
				spiceAmount := components.SpiceAmountRes.Get(targetSpiceEntry)
				amountToHarvest := 10
				if spiceAmount.Amount < amountToHarvest {
					amountToHarvest = spiceAmount.Amount
				}
				spiceAmount.Amount -= amountToHarvest
				harvester.CarriedAmount += amountToHarvest

				if harvester.CarriedAmount >= harvester.Capacity {
					harvester.State = components.StateMovingToRefinery
					refineryEntry, _ := qRefinery.First(ecs.World)
					refineryPos := components.Position.Get(refineryEntry)
					t.X, t.Y = refineryPos.X, refineryPos.Y
				} else {
					harvester.State = components.StateIdle
				}
			}
		case components.StateMovingToRefinery:
			refineryEntry, _ := qRefinery.First(ecs.World)
			refineryPos := components.Position.Get(refineryEntry)
			dx := refineryPos.X - p.X
			dy := refineryPos.Y - p.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < 64 { // Close enough to unload
				harvester.State = components.StateUnloading
			}
		case components.StateUnloading:
			// For now, just dump the spice and go back to being idle
			harvester.CarriedAmount = 0
			harvester.State = components.StateIdle
		}
	})
}

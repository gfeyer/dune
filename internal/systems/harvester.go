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
			// Harvester is idle, waiting for a command.
			// If it has a target spice, it means it has completed a loop and should go back.
			if harvester.TargetSpice != 0 {
				harvester.State = components.StateMovingToSpice
				targetSpiceEntry := ecs.World.Entry(harvester.TargetSpice)
				spicePos := components.Position.Get(targetSpiceEntry)
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
				// Stop moving
				v := components.Velocity.Get(entry)
				v.X, v.Y = 0, 0
				*t = components.Target{}
			}
		case components.StateHarvesting:
			// If harvester is full, move to refinery
			if harvester.CarriedAmount >= harvester.Capacity {
				harvester.State = components.StateMovingToRefinery
				refineryEntry, _ := qRefinery.First(ecs.World)
				refineryPos := components.Position.Get(refineryEntry)
				t.X, t.Y = refineryPos.X, refineryPos.Y
				return
			}

			// Harvest spice continuously
			targetSpiceEntry := ecs.World.Entry(harvester.TargetSpice)
			// Check if spice field is depleted
			if !targetSpiceEntry.HasComponent(components.SpiceAmountRes) {
				harvester.State = components.StateIdle
				return
			}
			spiceAmount := components.SpiceAmountRes.Get(targetSpiceEntry)
			amountToHarvest := 1
			if spiceAmount.Amount < amountToHarvest {
				amountToHarvest = spiceAmount.Amount
			}

			// If spice field is empty, find a new one
			if amountToHarvest == 0 {
				ecs.World.Remove(targetSpiceEntry.Entity())
				harvester.State = components.StateIdle
				return
			}

			spiceAmount.Amount -= amountToHarvest
			harvester.CarriedAmount += amountToHarvest
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

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
	qSpice      = donburi.NewQuery(filter.Contains(components.SpiceRes, components.SpiceAmountRes, components.Position))
	qRefinery   = donburi.NewQuery(filter.Contains(components.RefineryRes, components.Position))
)

func handleIdle(ecs *ecs.ECS, harvester *components.HarvesterData, p *components.Pos, t *components.Target) {
	// Harvester is idle, waiting for a command.
	// If it has a target spice, it means it has completed a loop and should go back.
	if harvester.TargetSpice != 0 {
		if !ecs.World.Valid(harvester.TargetSpice) {
			if harvester.CarriedAmount > 0 {
				harvester.State = components.StateMovingToRefinery
			} else {
				harvester.State = components.StateIdle
			}
			harvester.TargetSpice = 0
			return
		}
		harvester.State = components.StateMovingToSpice
		targetSpiceEntry := ecs.World.Entry(harvester.TargetSpice)
		spicePos := components.Position.Get(targetSpiceEntry)
		t.X, t.Y = spicePos.X, spicePos.Y
	}

	if harvester.CarriedAmount > 0 {
		harvester.State = components.StateMovingToRefinery
		// Find the closest refinery and set it as the target
		if closestRefinery := findClosestRefinery(ecs, p); closestRefinery != nil {
			harvester.TargetRefinery = closestRefinery.Entity()
			refineryPos := components.Position.Get(closestRefinery)
			t.X, t.Y = refineryPos.X, refineryPos.Y
		}
	}
}

func handleMovingToSpice(ecs *ecs.ECS, entry *donburi.Entry, harvester *components.HarvesterData, p *components.Pos, t *components.Target) {
	// Check if target is still valid
	if !ecs.World.Valid(harvester.TargetSpice) {
		if harvester.CarriedAmount > 0 {
			harvester.State = components.StateMovingToRefinery
		} else {
			harvester.State = components.StateIdle
		}
		harvester.TargetSpice = 0
		return
	}

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
}

func findClosestRefinery(ecs *ecs.ECS, p *components.Pos) *donburi.Entry {
	var closestRefinery *donburi.Entry
	minDist := math.MaxFloat64

	qRefinery.Each(ecs.World, func(refineryEntry *donburi.Entry) {
		refineryPos := components.Position.Get(refineryEntry)
		dx := refineryPos.X - p.X
		dy := refineryPos.Y - p.Y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist < minDist {
			minDist = dist
			closestRefinery = refineryEntry
		}
	})

	return closestRefinery
}

func handleHarvesting(ecs *ecs.ECS, entry *donburi.Entry, harvester *components.HarvesterData, p *components.Pos, t *components.Target) {
	// If harvester is full, move to refinery
	if harvester.CarriedAmount >= harvester.Capacity {
		harvester.State = components.StateMovingToRefinery

		// Find the closest refinery
		if closestRefinery := findClosestRefinery(ecs, p); closestRefinery != nil {
			harvester.TargetRefinery = closestRefinery.Entity()
			refineryPos := components.Position.Get(closestRefinery)
			t.X, t.Y = refineryPos.X, refineryPos.Y
		}
		return
	}

	// Harvest spice continuously
	if !ecs.World.Valid(harvester.TargetSpice) {
		if harvester.CarriedAmount > 0 {
			harvester.State = components.StateMovingToRefinery
		} else {
			harvester.State = components.StateIdle
		}
		harvester.TargetSpice = 0
		return
	}
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

	// If spice field is empty, go idle
	if amountToHarvest == 0 {
		ecs.World.Remove(targetSpiceEntry.Entity())

		// Go idle and wait for a new command.
		harvester.State = components.StateIdle
		harvester.TargetSpice = 0
		*t = components.Target{}
		return
	}

	spiceAmount.Amount -= amountToHarvest
	harvester.CarriedAmount += amountToHarvest
}

func handleMovingToRefinery(ecs *ecs.ECS, harvester *components.HarvesterData, p *components.Pos, t *components.Target) {
	if harvester.TargetRefinery == 0 || !ecs.World.Valid(harvester.TargetRefinery) {
		if closestRefinery := findClosestRefinery(ecs, p); closestRefinery != nil {
			harvester.TargetRefinery = closestRefinery.Entity()
			refineryPos := components.Position.Get(closestRefinery)
			t.X, t.Y = refineryPos.X, refineryPos.Y
		} else {
			harvester.State = components.StateIdle
			return
		}
	}
	targetRefineryEntry := ecs.World.Entry(harvester.TargetRefinery)
	refineryPos := components.Position.Get(targetRefineryEntry)
	dx := refineryPos.X - p.X
	dy := refineryPos.Y - p.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist < 64 { // Close enough to unload
		harvester.State = components.StateUnloading
	}
}

func handleUnloading(ecs *ecs.ECS, harvester *components.HarvesterData) {
	playerEntry, _ := QPlayer.First(ecs.World)
	player := components.PlayerRes.Get(playerEntry)
	player.Money += harvester.CarriedAmount
	harvester.CarriedAmount = 0
	harvester.State = components.StateIdle
}

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
			handleIdle(ecs, harvester, p, t)
		case components.StateMovingToSpice:
			handleMovingToSpice(ecs, entry, harvester, p, t)
		case components.StateHarvesting:
			handleHarvesting(ecs, entry, harvester, p, t)
		case components.StateMovingToRefinery:
			handleMovingToRefinery(ecs, harvester, p, t)
		case components.StateUnloading:
			handleUnloading(ecs, harvester)
		}
	})
}

package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type Pos struct{ X, Y float64 }
type Vel struct{ X, Y float64 }

// Unit components
type UnitType int

const (
	Trike UnitType = iota
	Harvester
)

type Unit struct {
	Type UnitType
}

type Selectable struct {
	Selected bool
}

type HarvesterState int

const (
	StateIdle HarvesterState = iota
	StateMovingToSpice
	StateHarvesting
	StateMovingToRefinery
	StateUnloading
)

type HarvesterData struct {
	State         HarvesterState
	TargetSpice   donburi.Entity
	HarvestTimer  int
	CarriedAmount int
	Capacity      int
}

type Target struct {
	X, Y float64
}

type Minimap struct {
	Width, Height int
	X, Y          int
}

type Drag struct {
	IsDragging     bool
	StartX, StartY int
	EndX, EndY     int
}

type SpiceAmount struct {
	Amount int
}

type Health struct {
	Current int
	Max     int
}

type Refinery struct{}

type Spice struct{}

var (
	Position      = donburi.NewComponentType[Pos]()
	Velocity      = donburi.NewComponentType[Vel]()
	Sprite        = donburi.NewComponentType[*ebiten.Image]()
	UnitRes       = donburi.NewComponentType[Unit]()
	SelectableRes = donburi.NewComponentType[Selectable]()
	TargetRes     = donburi.NewComponentType[Target]()
	MinimapRes    = donburi.NewComponentType[Minimap]()
	DragRes       = donburi.NewComponentType[Drag]()
	SpiceRes      = donburi.NewComponentType[Spice]()
	HarvesterRes  = donburi.NewComponentType[HarvesterData]()
	SpiceAmountRes = donburi.NewComponentType[SpiceAmount]()
	RefineryRes   = donburi.NewComponentType[Refinery]()
	HealthRes     = donburi.NewComponentType[Health]()
)

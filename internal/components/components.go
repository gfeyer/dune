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

type Target struct {
	X, Y float64
}

type Minimap struct {
	Width, Height int
	X, Y          int
}

type Drag struct {
	IsDragging    bool
	StartX, StartY int
	EndX, EndY     int
}

type Spice struct{}

var (
	Position   = donburi.NewComponentType[Pos]()
	Velocity   = donburi.NewComponentType[Vel]()
	Sprite     = donburi.NewComponentType[*ebiten.Image]()
	UnitRes    = donburi.NewComponentType[Unit]()
	SelectableRes = donburi.NewComponentType[Selectable]()
	TargetRes     = donburi.NewComponentType[Target]()
	MinimapRes    = donburi.NewComponentType[Minimap]()
	DragRes       = donburi.NewComponentType[Drag]()
	SpiceRes      = donburi.NewComponentType[Spice]()
)

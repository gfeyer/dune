package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type Pos struct{ X, Y float64 }
type Vel struct{ X, Y float64 }

var (
	Position = donburi.NewComponentType[Pos]()
	Velocity = donburi.NewComponentType[Vel]()
	Sprite   = donburi.NewComponentType[*ebiten.Image]()
)

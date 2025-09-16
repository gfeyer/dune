package fog

import (
	"github.com/gfeyer/ebit/internal/settings"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type VisibilityState int

const (
	Hidden  VisibilityState = iota // Completely black, never seen
	Shroud                         // Seen before, but no units currently there
	Visible                        // Currently visible by a unit
)

// Fog is a resource that holds the state of the fog of war.
type Fog struct {
	Grid     [][]VisibilityState
	TileSize int
	Width    int
	Height   int
}

var FogRes = donburi.NewComponentType[Fog]()

func NewFog(s *settings.Settings, tileSize int) *Fog {
	width := s.MapWidth / tileSize
	height := s.MapHeight / tileSize
	grid := make([][]VisibilityState, height)
	for i := range grid {
		grid[i] = make([]VisibilityState, width)
		for j := range grid[i] {
			grid[i][j] = Hidden
		}
	}
	return &Fog{
		Grid:     grid,
		TileSize: tileSize,
		Width:    width,
		Height:   height,
	}
}

// GetFog gets the fog from the world.
func GetFog(w donburi.World) *Fog {
	entry, _ := donburi.NewQuery(filter.Contains(FogRes)).First(w)
	return FogRes.Get(entry)
}

package systems

import (
	"github.com/yohamta/donburi/ecs"
)

const (
	LayerSprites ecs.LayerID = iota // 0
	LayerUI                         // 1 (in-world UI like health bars)
	LayerMinimap                    // 2
	LayerMenus      // 3 (for future build menus etc.)
	LayerBuildMenuUI // 4
	LayerPlacement   // 5
)


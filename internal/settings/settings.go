package settings

import "github.com/yohamta/donburi"

// Settings is a resource that holds game-wide settings.
type Settings struct {
	ScreenWidth  int
	ScreenHeight int
}

var SettingsRes = donburi.NewComponentType[Settings]()

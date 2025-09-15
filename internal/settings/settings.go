package settings

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

// Settings is a resource that holds game-wide settings.
type Settings struct {
	ScreenWidth  int
	ScreenHeight int
	MapWidth     int
	MapHeight    int
}

var SettingsRes = donburi.NewComponentType[Settings]()

var SettingsQuery = donburi.NewQuery(filter.Contains(SettingsRes))

// GetSettings gets the settings from the world.
func GetSettings(w donburi.World) *Settings {
	entry, _ := SettingsQuery.First(w)
	return SettingsRes.Get(entry)
}

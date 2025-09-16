package game

import (
	"image/color"
	"math/rand"

	"github.com/gfeyer/ebit/internal/camera"
	"github.com/gfeyer/ebit/internal/components"
	"github.com/gfeyer/ebit/internal/factory"
	"github.com/gfeyer/ebit/internal/fog"
	"github.com/gfeyer/ebit/internal/settings"
	"github.com/gfeyer/ebit/internal/systems"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type Game struct {
	ecs *ecs.ECS
}

func NewGame(w, h int) *Game {
	world := donburi.NewWorld()
	ecs := ecs.NewECS(world)

	// Register settings
	e := world.Create(settings.SettingsRes)
	entry := world.Entry(e)
	*settings.SettingsRes.Get(entry) = settings.Settings{
		ScreenWidth:  w,
		ScreenHeight: h,
		MapWidth:     w * 4,
		MapHeight:    h * 4,
	}

	// Register camera
	ce := world.Create(camera.CameraRes)
	centry := world.Entry(ce)
	*camera.CameraRes.Get(centry) = camera.Camera{
		X: float64((w*4)/2 - w/2),
		Y: float64((h*4)/2 - h/2),
	}

	// Create minimap
	mme := world.Create(components.MinimapRes)
	mmentry := world.Entry(mme)
	*components.MinimapRes.Get(mmentry) = components.Minimap{
		Width:  w / 5,
		Height: h / 5,
		X:      w - w/5 - 10,
		Y:      10,
	}

	// Create drag selection
	de := world.Create(components.DragRes)
	dentry := world.Entry(de)
	*components.DragRes.Get(dentry) = components.Drag{}

	// Create placement
	ple := world.Create(components.PlacementRes)
	plentry := world.Entry(ple)
	*components.PlacementRes.Get(plentry) = components.Placement{}

	// Create player
	pe := world.Create(components.PlayerRes)
	pentry := world.Entry(pe)
	*components.PlayerRes.Get(pentry) = components.Player{Money: 1000}

	// Create fog
	fe := world.Create(fog.FogRes)
	fentry := world.Entry(fe)
	*fog.FogRes.Get(fentry) = *fog.NewFog(settings.GetSettings(world), 16)

	// Register systems
	ecs.AddSystem(systems.UpdateMovement)
	ecs.AddSystem(systems.ResolveCollisions)
	ecs.AddSystem(systems.UpdateInput)
	ecs.AddSystem(systems.UpdateBuildInput)
	ecs.AddSystem(camera.Update)
	ecs.AddSystem(systems.UpdateMinimap)
	ecs.AddSystem(systems.UpdateHarvester)
	ecs.AddSystem(systems.UpdateFog)

	// Register renderers
	ecs.AddRenderer(systems.LayerSpice, systems.DrawSpice)
	ecs.AddRenderer(systems.LayerBuildings, systems.DrawBuildings)
	ecs.AddRenderer(systems.LayerUnits, systems.DrawUnits)
	ecs.AddRenderer(systems.LayerUI, systems.DrawUI)
	ecs.AddRenderer(systems.LayerMinimap, systems.DrawMinimap)
	ecs.AddRenderer(systems.LayerBuildMenuUI, systems.DrawBuildMenu)
	ecs.AddRenderer(systems.LayerPlacement, systems.DrawPlacement)
	ecs.AddRenderer(systems.LayerFog, systems.DrawFog)

	// Spawn initial units
	s := settings.GetSettings(world)
	centerX := float64(s.MapWidth) / 2
	centerY := float64(s.MapHeight) / 2
	factory.CreateTrike(world, centerX+50, centerY+50)
	factory.CreateHarvester(world, centerX, centerY-50)
	factory.CreateRefinery(world, centerX-50, centerY-50)

	// Create build options
	minimap := components.MinimapRes.Get(mmentry)
	padding := 5
	iconWidth := (minimap.Width - padding) / 2
	iconHeight := 64
	factory.CreateBuildOption(world, components.BuildingRefinery, "Refinery", 750, iconWidth, iconHeight)
	factory.CreateBuildOption(world, components.BuildingBarracks, "Barracks", 250, iconWidth, iconHeight)

	// Create unit options
	factory.CreateUnitOption(world, components.Harvester, "Harvester", 500, components.BuildingRefinery, iconWidth, iconHeight)
	factory.CreateUnitOption(world, components.Trike, "Trike", 350, components.BuildingBarracks, iconWidth, iconHeight)
	factory.CreateUnitOption(world, components.Quad, "Quad", 800, components.BuildingBarracks, iconWidth, iconHeight)

	// Spawn spice
	for i := 0; i < 50; i++ {
		x := rand.Float64() * float64(s.MapWidth)
		y := rand.Float64() * float64(s.MapHeight)
		factory.CreateSpice(world, x, y)
	}

	return &Game{ecs: ecs}
}

func (g *Game) Update() error {
	g.ecs.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{210, 180, 140, 255}) // sand color
	g.ecs.DrawLayer(systems.LayerSpice, screen)
	g.ecs.DrawLayer(systems.LayerBuildings, screen)
	g.ecs.DrawLayer(systems.LayerUnits, screen)
	g.ecs.DrawLayer(systems.LayerFog, screen)
	g.ecs.DrawLayer(systems.LayerUI, screen)
	g.ecs.DrawLayer(systems.LayerMinimap, screen)
	g.ecs.DrawLayer(systems.LayerBuildMenuUI, screen)
	g.ecs.DrawLayer(systems.LayerPlacement, screen)
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	s := settings.GetSettings(g.ecs.World)
	return s.ScreenWidth, s.ScreenHeight
}

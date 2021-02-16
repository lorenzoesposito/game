package engine

import (
	"fmt"
	. "gfx2"

	. "game.com/lorenzo/game/utils"
)

type Entity struct {
	Index     int // Index in the entity list
	Transform [4][4]float64
	Colors    []Vec3f // Color of this entity
	Mesh      []Tri3f // Triangles of this entity
}

// Global Variables
var SX, SY uint16
var CanvasWidth, CanvasHeight float64 = 2, 2
var MainCamera [4][4]float64 = [4][4]float64{
	[4]float64{1, 0, 0, 0},
	[4]float64{0, 1, 0, 0},
	[4]float64{0, 0, 1, 0},
	[4]float64{0, 0, 0, 1}}

var ProgramState int = 0

var entityList []Entity // List of all existing entities

// Options
var DoDebugging bool = false

func GetProgramState() int {
	return ProgramState
}

func MainMenu() {
	ResizeWindow(800, 800)

	selected := 0

	for {
		UpdateAus()

		Stiftfarbe(255, 255, 0)
		Cls()

		if selected == 0 {
			Stiftfarbe(255, 0, 0)
		} else {
			Stiftfarbe(0, 0, 0)
		}
		SetzeFont("Minecraft.ttf", 30)
		SchreibeFont(SX/2-30, SY/2-100, "Play")

		if selected == 1 {
			Stiftfarbe(255, 0, 0)
		} else {
			Stiftfarbe(0, 0, 0)
		}
		SetzeFont("Minecraft.ttf", 30)
		SchreibeFont(SX/2-56, SY/2-60, "Settings")

		UpdateAn()

		A, B, _ := TastaturLesen1()
		if DoDebugging {
			fmt.Println(A, selected)
		}
		if B == 1 {
			if A == 27 {
				QuitProgram = true
				return
			}
			if A == 13 {
				if selected == 0 {
					ProgramState = 1
					return
				}
				if selected == 1 {
					ProgramState = 2
					return
				}
			}
			if A == 273 {
				selected--
			}
			if A == 274 {
				selected = (selected + 1) % 2
			}
			if selected < 0 {
				selected = 1
			}
		}
	}
}

func SettingsMenu() {
	ResizeWindow(800, 800)

	for {
		UpdateAus()

		Stiftfarbe(255, 255, 0)
		Cls()

		Stiftfarbe(255, 0, 0)

		SetzeFont("Minecraft.ttf", 30)
		SchreibeFont(SX/2-30, SY/2-100, "Settings")

		UpdateAn()

		A, B, _ := TastaturLesen1()
		if DoDebugging {
			fmt.Println(A)
		}
		if B == 1 {
			if A == 27 {
				QuitProgram = true
				return
			}
			if A == 13 {
				ProgramState = 0
				return
			}
		}
	}
}

func GetEntities() []Entity {
	return entityList
}

// Global Functions

func ResizeWindow(sizeX, sizeY uint16) {
	// Clamp Size
	if sizeX > 1920 {
		sizeX = 1920
	}
	if sizeY > 1200 {
		sizeY = 1200
	}
	// Don't Resize if size is same as before
	if SX == sizeX || SY == sizeY {
		return
	}
	// Resize window to new size
	SX, SY = sizeX, sizeY
	Fenster(sizeX, sizeY)
}

func DestroyEntity(i int) {
	entityList[i] = entityList[len(entityList)-1]
	entityList = entityList[:len(entityList)-1]
	for i = i; i < len(entityList); i++ {
		entityList[i].Index = i
	}
}

func InstantiateEntity(e Entity) {
	e.Index = len(entityList)
	entityList = append(entityList, e)
}

func CullCamera(camera Vec3f) bool {
	if camera.Z < 0 {
		return false
	}
	return true
}

func CameraToCanvas(world Vec3f) (canvas Vec3f) {
	canvas = Vec3f{world.X / world.Z, world.Y / world.Z, -world.Z}
	return canvas
}

func CanvasToNdc(canvas Vec3f) (ndc Vec2f) {
	ndc = Vec2f{(canvas.X + CanvasWidth/2) / CanvasWidth, (canvas.Y + CanvasHeight/2) / CanvasHeight}
	return ndc
}

func NdcToRaster(ndc Vec2f) (raster Vec2ui) {
	raster = Vec2ui{uint16(ndc.X * float64(SX)), uint16((1 - ndc.Y) * float64(SY))}
	return raster
}

func DrawTriangle(tRaster Tri2ui) {
	Volldreieck(tRaster.A.X, tRaster.A.Y, tRaster.B.X, tRaster.B.Y, tRaster.C.X, tRaster.C.Y)
}

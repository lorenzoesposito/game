package main

import (
	"fmt"
	. "gfx2"
	"math"
	"net"
	"time"

	. "game.com/lorenzo/game/engine"
	. "game.com/lorenzo/game/utils"
)

func main() {
	for !QuitProgram {
		states[GetProgramState()]()
	}
	FensterAus()
	fmt.Println("about to shutdown...")
	fmt.Scanln()
}

type State func()

var states []State = []State{MainMenu, GameLoop, SettingsMenu}

func GameLoop() {
	QuitGame = false
	addr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	LogFatal(err)
	connection, err := net.DialUDP("udp", nil, addr)
	LogFatal(err)
	defer connection.Close()
	connection.Write([]byte("joined_"))

	// Scene setup
	InstantiateEntity(Entity{0, IdentityMatrix(), Color("test"), Mesh("test")})
	//entityList[0].transform = Translate(entityList[0].transform, Vec3f{0,0,20})

	// Clock zero
	t0 := time.Now()
	t0 = t0
	t1 := time.Now()
	// Seconds elapsed since last frame
	deltaTime := float64(0)

	go KeyboardInput()

	//Camera Rotation values

	// Game loop
	for !QuitGame {
		// Calculate deltaTime
		deltaTime = float64(time.Now().Sub(t1).Nanoseconds()) / 10000000
		if DoDebugging {
			fmt.Println("fps :", 1/deltaTime)
		}
		t1 = time.Now()

		// Move Player
		if DoDebugging {
			fmt.Println(MousePos)
		}
		playerInput := Vec3f{GetAxis(Left, Right), GetAxis(Shift, Space), GetAxis(Down, Up)}
		MainCamera = Translate(MainCamera, Times(playerInput, ValVec(1*deltaTime)))
		GetEntities()[0].Transform = SetPosition(GetEntities()[0].Transform, Vec3f{0, math.Sin(time.Now().Sub(t0).Seconds())*5 - 10, 10})

		// Draw
		UpdateAus()

		Stiftfarbe(0, 0, 0)
		Cls()
		for i := 0; i < len(GetEntities()); i++ {
			// Draw Mesh
			for k := 0; k < len(GetEntities()[i].Mesh); k++ {
				Stiftfarbe(VecToColor(GetEntities()[i].Colors[k]))
				// Local to world
				tri := Tri3fTransform(GetEntities()[i].Mesh[k], GetEntities()[i].Transform)
				// World to Camera
				tri = Tri3fTransform(tri, Inverse(MainCamera))
				// Camera to canvas
				tri = Tri3f{CameraToCanvas(tri.A), CameraToCanvas(tri.B), CameraToCanvas(tri.C)}
				if CullCamera(tri.A) && CullCamera(tri.B) && CullCamera(tri.C) {
					continue
				}
				// Canvas to NDC
				tri2 := Tri2f{CanvasToNdc(tri.A), CanvasToNdc(tri.B), CanvasToNdc(tri.C)}
				// NDC to Raster
				triRaster := Tri2ui{NdcToRaster(tri2.A), NdcToRaster(tri2.B), NdcToRaster(tri2.C)}
				// Convert to uint16
				DrawTriangle(triRaster)
			}
		}
		UpdateAn()
	}
}

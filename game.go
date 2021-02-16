package main

import (
	"fmt"
	. "gfx2"
	"time"

	. "game.com/lorenzo/game/client"
	. "game.com/lorenzo/game/engine"
	. "game.com/lorenzo/game/utils"
)

var client Client
var player Entity

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

func handleMessage(conn string, pos Vec3f) {
	if conn == client.Id {
		player.Transform = player.SetPosition(pos)
	} else {
		if FindEntity(conn) == -1 {
			InstantiateEntity(Entity{conn, IdentityMatrix(), Color("test"), Mesh("test")})
			GetEntities()[FindEntity(conn)].Transform = GetEntities()[FindEntity(conn)].SetPosition(pos)
		} else {
			GetEntities()[FindEntity(conn)].Transform = GetEntities()[FindEntity(conn)].SetPosition(pos)
		}
	}
}

func read() {
	for {
		handleMessage(ParseClient(client.Read()))
	}
}

func GameLoop() {
	QuitGame = false
	client.InitializeClient("192.168.178.48:8080")
	player = Entity{client.Id, IdentityMatrix(), Color("test"), Mesh("test")}

	// Scene setup
	//InstantiateEntity(Entity{1111, IdentityMatrix(), Color("test"), Mesh("test")})
	//entityList[0].transform = Translate(entityList[0].transform, Vec3f{0,0,20})

	// Clock zero
	t0 := time.Now()
	t0 = t0
	t1 := time.Now()
	// Seconds elapsed since last frame
	deltaTime := float64(0)

	go KeyboardInput()
	go read()

	//Camera Rotation values

	// Game loop
	for !QuitGame {
		client.Conn.Write(append([]byte("update_"), BoolsToBytes([]bool{Left, Right, Shift, Space, Down, Up})...))
		// Calculate deltaTime
		deltaTime = float64(time.Now().Sub(t1).Nanoseconds()) / 10000000
		if DoDebugging {
			fmt.Println("fps :", 1/deltaTime)
		}
		t1 = time.Now()

		if DoDebugging {
			fmt.Println(MousePos)
		}

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
		for k := 0; k < len(player.Mesh); k++ {
			Stiftfarbe(VecToColor(player.Colors[k]))
			// Local to world
			tri := Tri3fTransform(player.Mesh[k], player.Transform)
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
		UpdateAn()
	}
}

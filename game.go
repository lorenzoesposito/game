package main

import (
	"fmt"
	"time"

	. "game.com/lorenzo/game/client"
	. "game.com/lorenzo/game/engine"
	. "game.com/lorenzo/game/gfx2"
	. "game.com/lorenzo/game/utils"
)

var server = "localhost:8080"
var client Client
var player Entity
var updates int

//var c = make(chan []string)
var closedConnections []string

func main() {
	for !QuitProgram {
		states[GetProgramState()]()
	}
	FensterAus()
	fmt.Scanln()
}

type State func()

var states []State = []State{MainMenu, GameLoop, SettingsMenu}

func handleMessage(conn string, pos Vec3f) {
	updates++
	if conn == client.Id {
		player.Transform = player.SetPosition(pos)
	} else if conn[:4] == "quit" {
		closedConnections[len(closedConnections)-1] = conn[4:]
	} else {
		if FindEntity(conn) == -1 {
			InstantiateEntity(Entity{conn, IdentityMatrix(), Color("car"), Mesh("car")})
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
	client.InitializeClient(server)
	player = Entity{client.Id, IdentityMatrix(), Color("car"), Mesh("car")}

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
			closedConnections = append(closedConnections, " ")
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
		for _, i := range closedConnections {
			if i != " " {
				DestroyEntity(i)
			}
		}
		closedConnections = []string{}
	}
}

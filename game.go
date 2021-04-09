package main

import (
	"fmt"
	"time"

	. "game.com/lorenzo/game/client"
	. "game.com/lorenzo/game/engine"
	. "game.com/lorenzo/game/gfx"
	. "game.com/lorenzo/game/utils"
)

var (
	server  = "localhost:8080"
	client  Client
	player  Entity
	updates int
	//object  = Dir[:len(Dir)-4] + "game/models/human/Idle/human0"
	object            = Dir[:len(Dir)-4] + "game/models/elegant/Idle/elegantIdle_000001"
	closedConnections []string
)

func main() {
	ParseObjs()
	for !QuitProgram {
		states[GetProgramState()]()
	}
	FensterAus()
}

type State func()

var states []State = []State{MainMenu, GameLoop, SettingsMenu}

func handleMessage(conn string, pos Vec3f, obj string) {
	fmt.Println(conn, pos, obj)
	updates++
	if conn == client.Id {
		player.Transform = player.Update(pos, obj)
	} else if conn[:4] == "quit" {
		closedConnections[len(closedConnections)-1] = conn[4:]
	} else {
		if FindEntity(conn) == -1 {
			InstantiateEntity(Entity{conn, IdentityMatrix(), Color(object), Mesh(object), obj})
			GetEntities()[FindEntity(conn)].Transform = GetEntities()[FindEntity(conn)].Update(pos, obj)
		} else {
			GetEntities()[FindEntity(conn)].Transform = GetEntities()[FindEntity(conn)].Update(pos, obj)
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
	client.InitializeClient(server, object[len(object)-18:])
	player = Entity{client.Id, IdentityMatrix(), Color(object), Mesh(object), object[len(object)-18:]}

	// Clock zero

	//t0 := time.Now()
	t1 := time.Now()
	// Seconds elapsed since last frame
	deltaTime := float64(0)

	go KeyboardInput()
	go read()
	// Game loop
	for !QuitGame {
		player.Animate()
		client.Conn.Write(append([]byte("update_"+player.Object+"_"), BoolsToBytes([]bool{Left, Right, Shift, Space, Down, Up})...))
		fmt.Println("sent this update:", ParseServer(string(append([]byte("update_"+player.Object+"_"), BoolsToBytes([]bool{Left, Right, Shift, Space, Down, Up})...))))
		// Calculate deltaTime
		deltaTime = float64(time.Now().Sub(t1).Nanoseconds()) / 10000000
		if DoDebugging {
			fmt.Println("fps :", 1/deltaTime)
		}
		t1 = time.Now()
		//fmt.Println(int(math.Sin(float64(frames)/5) + 1))
		if DoDebugging {
			fmt.Println(MousePos)
		}
		// Draw
		UpdateAus()
		//	fmt.Println("fps :", 100/deltaTime)
		Stiftfarbe(0, 0, 0)
		Cls()
		for i := 0; i < len(GetEntities()); i++ {
			GetEntities()[i].Animate()
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
		time.Sleep(time.Second / 24)
	}
}

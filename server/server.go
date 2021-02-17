package main

import (
	"fmt"
	"net"

	. "game.com/lorenzo/game/utils"
)

var start = Vec3f{0, 0, 20}
var players = make(map[string]player)

type player struct {
	conn     *net.UDPAddr
	position Vec3f
}

func playerToBytes(pl player) []byte {
	return []byte(fmt.Sprintf("%s_%f_%f_%f_%f_%f_%f", pl.conn, pl.position.X, pl.position.Y, pl.position.Z))
}

func handleMessage(listener *net.UDPConn) {
	buffer := make([]byte, 1024*8)
	n, conn, err := listener.ReadFromUDP(buffer)

	LogFatal(err)
	if err != nil {
		return
	}
	message := string(buffer[:n])

	switch ParseServer(message).MsgType {
	case "joined":
		players[conn.String()] = player{conn, start}
		fmt.Println(conn.String())
	case "update":
		p := players[conn.String()].position
		players[conn.String()] = player{conn, Vec3f{p.X + GetAxis(ParseServer(message).Input[0], ParseServer(message).Input[1]),
			p.Y + GetAxis(ParseServer(message).Input[2], ParseServer(message).Input[3]),
			p.Z + GetAxis(ParseServer(message).Input[4], ParseServer(message).Input[5])}}
	}

	for _, player := range players {
		for _, pl := range players {
			listener.WriteToUDP(playerToBytes(pl), player.conn)
		}
	}
}

func main() {
	s, err := net.ResolveUDPAddr("udp", "localhost:8080")
	LogFatal(err)
	listener, err := net.ListenUDP("udp", s)
	LogFatal(err)
	defer listener.Close()
	for {
		handleMessage(listener)
	}
}

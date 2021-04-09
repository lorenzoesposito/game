package main

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"time"
	"unicode"

	. "game.com/lorenzo/game/utils"
)

var (
	start      = Vec3f{0, 0, 20}
	players    = make(map[string]player)
	closeConns []*string
	listener   *net.UDPConn
)

type player struct {
	conn      *net.UDPAddr
	position  Vec3f
	connected int
}

func playerToBytes(pl player, animation string) []byte {
	return []byte(fmt.Sprintf("%s_%f_%f_%f_", pl.conn, pl.position.X, pl.position.Y, pl.position.Z) + animation)
}
func quit(pl player) []byte {
	return []byte(fmt.Sprintf("quit%s_%f_%f_%f_%f_%f_%f", pl.conn, pl.position.X, pl.position.Y, pl.position.Z))
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
		players[conn.String()] = player{conn, start, 1}
	case "update":
		p := players[conn.String()].position
		players[conn.String()] = player{conn, Vec3f{p.X + GetAxis(ParseServer(message).Input[0], ParseServer(message).Input[1]),
			p.Y + GetAxis(ParseServer(message).Input[2], ParseServer(message).Input[3]),
			p.Z + GetAxis(ParseServer(message).Input[4], ParseServer(message).Input[5])}, players[conn.String()].connected + 1}
	}
	for _, player := range players {
		listener.WriteToUDP(playerToBytes(players[conn.String()], animation(message)), player.conn)
		fmt.Println("message", string(playerToBytes(players[conn.String()], animation(message))))
	}
}

func animation(msg string) string {
	fmt.Println("AAAAAAAA", ParseServer(msg).Object)
	a, b, c := getAnimation(ParseServer(msg).Object)
	fmt.Println(a, b, c)
	var newAnimation string
	switch b {
	case "Idle":
		if GetAxis(ParseServer(msg).Input[4], ParseServer(msg).Input[5]) != 0 {
			newAnimation = a + "Run_000001"
		} else {
			newAnimation = a + b + fmt.Sprintf("_%06d", (c+1)%32+1)
			fmt.Println("anim:", newAnimation)
		}
	case "Run":
		if GetAxis(ParseServer(msg).Input[4], ParseServer(msg).Input[5]) != 0 {
			newAnimation = a + b + fmt.Sprintf("_%06d", (c+1)%16+1)
		} else {
			newAnimation = a + "Idle_000001"
		}
	}

	return newAnimation
}

func f(x float64, n float64) int {
	return int((math.Acos(math.Cos((1/n)*math.Pi*x)) * (n / math.Pi)) + 0.5)
}

func getAnimation(obj string) (string, string, int) {
	str := []rune(obj)
	n := 0
	for i := range str {
		if unicode.IsUpper(str[i]) {
			n = i
		}
	}
	num, _ := strconv.Atoi(obj[len(obj)-6:])
	return obj[:n], obj[n : len(obj)-7], num
}

func setZero() {
	for i := range players {
		players[i] = player{players[i].conn, players[i].position, 0}
	}
}

func c() {
	for i := range players {
		if players[i].connected == 0 {
			for n := range players {

				listener.WriteToUDP(quit(players[i]), players[n].conn)
			}
			delete(players, i)
		}
	}
	setZero()
}

func main() {
	s, err := net.ResolveUDPAddr("udp", "localhost:8080")
	LogFatal(err)
	listener, _ = net.ListenUDP("udp", s)
	defer listener.Close()
	go CallEvery(2*time.Second, c)
	for {
		handleMessage(listener)
	}
}

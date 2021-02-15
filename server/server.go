package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	. "game.com/lorenzo/game/utils"
)

type player struct {
	conn     *net.UDPAddr
	position Vec3f
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleMessage(listener *net.UDPConn) {
	buffer := make([]byte, 1024*8)
	n, conn, err := listener.ReadFromUDP(buffer)

	logFatal(err)
	if err != nil {
		return
	}
	message := string(buffer[:n])
	msgType, color, _ := parse(message)
	_, _, pl := parse(message)

	if msgType == "joined" {

		players[conn.String()] = player{conn, color,
			pl[0],
			pl[1], pl[2], pl[3]}
		fmt.Println(conn.String(), "just joined")

	} else {
		players[conn.String()] = player{conn, color,
			pl[0],
			pl[1], pl[2], pl[3]}
		for _, user := range players {
			for _, pl := range players {
				str := strings.Split(pl.conn.String(), ":")
				str2 := strings.Split(user.conn.String(), ":")
				if str[1] != str2[1] {
					listener.WriteToUDP([]byte(packet(user)), pl.conn)
					fmt.Println("update package sent to ", msgType[:1], ":", pl.conn)
					fmt.Println("package:", packet(user))
				}
			}
		}
	}

}

func packet(player player) string {
	return string(fmt.Sprintf("%s_%s_%s_%s_%s_%s",
		player.conn.String(),
		player.color,
		strconv.Itoa(player.x),
		strconv.Itoa(player.y),
		strconv.Itoa(player.width),
		strconv.Itoa(player.height)))
}

func parse(user string) (string, string, []int) {
	split := strings.Split(user, "_")
	//ID := split[0]
	T := split[0]
	C := split[1]
	X, _ := strconv.Atoi(split[2])
	Y, _ := strconv.Atoi(split[3])
	W, _ := strconv.Atoi(split[4])
	H, _ := strconv.Atoi(split[5])
	return T, C, []int{X, Y, W, H}
}

var (
	players         = make(map[string]player)
	openConnections = make(map[string]*net.UDPAddr)
	newConnection   = make(chan *net.UDPAddr)
	deadConnection  = make(chan *net.UDPAddr)
)

func main() {
	for key := range players {
		delete(players, key)
	}

	s, err := net.ResolveUDPAddr("udp", "localhost:8080")
	logFatal(err)
	listener, err := net.ListenUDP("udp", s)
	logFatal(err)
	defer listener.Close()
	for {
		handleMessage(listener)
	}
}

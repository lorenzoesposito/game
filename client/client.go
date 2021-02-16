package client

import (
	"net"

	. "game.com/lorenzo/game/utils"
)

type Client struct {
	Id   string
	Conn net.Conn
}

func (c *Client) InitializeClient(server string) {
	addr, err := net.ResolveUDPAddr("udp", server)
	LogFatal(err)
	connection, err := net.DialUDP("udp", nil, addr)
	LogFatal(err)
	c.Id = connection.LocalAddr().String()
	c.Conn = connection
	defer connection.Close()

	connection.Write([]byte("joined_"))
}

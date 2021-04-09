package client

import (
	"net"

	. "game.com/lorenzo/game/utils"
)

type Client struct {
	Id   string
	Conn *net.UDPConn
}

func (c *Client) InitializeClient(server string, obj string) {
	addr, err := net.ResolveUDPAddr("udp", server)
	LogFatal(err)
	connection, err := net.DialUDP("udp", nil, addr)
	LogFatal(err)
	c.Id = connection.LocalAddr().String()
	c.Conn = connection

	connection.Write(append([]byte("joined_"+obj+"_"), BoolsToBytes([]bool{false, false, false, false, false, false})...))
}

func (c *Client) Read() string {
	buffer := make([]byte, 1024*8)
	n, _, err := c.Conn.ReadFromUDP(buffer)
	LogFatal(err)
	return string(buffer[:n])
}

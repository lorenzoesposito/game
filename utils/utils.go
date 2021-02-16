package utils

import (
	"log"
	"strconv"
	"strings"
)

type Msg struct {
	MsgType string
	Input   []bool
}

func GetAxis(neg, pos bool) float64 {
	ax := float64(0)
	if neg {
		ax--
	}
	if pos {
		ax++
	}
	return ax
}

func StringToFloat(s string) float64 {
	if s, err := strconv.ParseFloat(s, 64); err == nil {
		return s
	}
	return -1
}

func BoolsToBytes(t []bool) []byte {
	b := make([]byte, (len(t)+7)/8)
	for i, x := range t {
		if x {
			b[i/8] |= 0x80 >> uint(i%8)
		}
	}
	return b
}

func BytesToBools(b []byte) []bool {
	t := make([]bool, 8*len(b))
	for i, x := range b {
		for j := 0; j < 8; j++ {
			if (x<<uint(j))&0x80 == 0x80 {
				t[8*i+j] = true
			}
		}
	}
	return t
}

func ParseServer(msg string) Msg {
	split := strings.Split(msg, "_")
	msgType := split[0]
	return Msg{msgType, BytesToBools([]byte(split[1]))}
}

func ParseClient(msg string) (string, Vec3f) {
	split := strings.Split(msg, "_")
	return split[0], Vec3f{StringToFloat(split[1]), StringToFloat(split[2]), StringToFloat(split[3])}
}

func LogFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

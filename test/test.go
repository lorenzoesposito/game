package main

import (
	"fmt"
	"strconv"
	"unicode"
)

func main() {
	str := "elegantIdle_000001"
	fmt.Println(getAnimation(str))
	fmt.Println(str[len(str)-6:])
}

func getAnimation(obj string) (string, string, int) {
	str := []rune(obj)
	n := 0
	for i := range str {
		if unicode.IsUpper(str[i]) {
			n = i
		}
	}
	num, _ := strconv.Atoi(obj[len(obj)-7:])

	return obj[:n], obj[n : len(obj)-7], num
}

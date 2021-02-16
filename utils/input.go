package utils

import (
	. "game.com/lorenzo/game/gfx2"
)

var MousePos Vec2ui = Vec2ui{400, 400}       // Mouse cursor position
var Shift, Space, Up, Down, Right, Left bool // Keyboard state inputs
var QuitGame, QuitKeyboardInput, QuitMouseInput, QuitProgram bool

func KeyboardInput() {
	QuitKeyboardInput = false
	TastaturpufferAn()
	for !QuitKeyboardInput {
		// Get Keyboard Input from the buffer
		a, b, c := TastaturpufferLesen1()
		// Process it
		ProcessKeyboardInput(a, b, c)
	}
	TastaturpufferAus()
}

func ProcessKeyboardInput(a uint16, b uint8, c uint16) {
	// Turn press input into a bool
	press := true
	if b == 0 {
		press = false
	}

	// Switch appropriate bool on/off
	switch a {
	case 27: // esc
		QuitMouseInput = true
		QuitKeyboardInput = true
		QuitProgram = true
		QuitGame = true
		break
	case 32: // space
		Space = press
		break
	case 304: // space
		Shift = press
		break
	case 273: // Up key
		Up = press
		break
	case 274: // Down key
		Down = press
		break
	case 275: // Right key
		Right = press
		break
	case 276: // Left key
		Left = press
		break
	}
}

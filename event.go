package geui

import (
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Event interface {
	String() string
}

type (

	MouseMove struct {
		X, Y float64
		glfw.MouseButton
	}

	MouseDown struct {
		X, Y float64
		glfw.MouseButton
	}

	MouseUp struct {
		X, Y float64
		glfw.MouseButton
	}

	MouseScroll struct {
		X, Y float64
	}

	KbType struct {
		rune
	}

	// KbDown is an event that happens when a key on the keyboard gets pressed.
	KbDown struct {
		glfw.Key
	}

	// KbUp is an event that happens when a key on the keyboard gets released.
	KbUp struct {
		glfw.Key
	}

	// KbRepeat is an event that happens when a key on the keyboard gets repeated.
	// This happens when its held down for some time.
	KbRepeat struct {
		glfw.Key
	}

	Resize struct {
		X, Y          float64
		Width, Height float64
	}
)

func (mm MouseMove) String() string   { return fmt.Sprintf("mouse/move/%v/%v", mm.X, mm.Y) }
func (md MouseDown) String() string   { return fmt.Sprintf("mouse/down/%v/%v/%v", md.X, md.Y, md.MouseButton) }
func (mu MouseUp) String() string     { return fmt.Sprintf("mouse/up/%v/%v/%v", mu.X, mu.Y, mu.MouseButton) }
func (ms MouseScroll) String() string { return fmt.Sprintf("mouse/scroll/%v/%v", ms.X, ms.Y) }
func (kt KbType) String() string      { return fmt.Sprintf("keyboad/type/%v", kt.rune) }
func (kd KbDown) String() string      { return fmt.Sprintf("keyboad/down/%v", kd.Key) }
func (ku KbUp) String() string        { return fmt.Sprintf("keyboad/up/%v", ku.Key) }
func (kr KbRepeat) String() string    { return fmt.Sprintf("keyboad/repeat/%v", kr.Key) }
func (rs Resize) String() string      { return fmt.Sprintf("viewport/resize/%v/%v", rs.Width, rs.Height) }

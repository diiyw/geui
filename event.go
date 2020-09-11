package geui

import "fmt"

type Event interface {
	String() string
}

type (
	WindowClose struct{}

	MouseMove struct {
		X, Y float64
		Mouse
	}

	MouseDown struct {
		X, Y float64
		Mouse
	}

	MouseUp struct {
		X, Y float64
		Mouse
	}

	MouseScroll struct {
		X, Y float64
	}

	KbType struct {
		rune
	}

	// KbDown is an event that happens when a key on the keyboard gets pressed.
	KbDown struct {
		Key
	}

	// KbUp is an event that happens when a key on the keyboard gets released.
	KbUp struct {
		Key
	}

	// KbRepeat is an event that happens when a key on the keyboard gets repeated.
	// This happens when its held down for some time.
	KbRepeat struct {
		Key
	}

	Resize struct {
		X, Y          float64
		Width, Height float64
	}
)

func (WindowClose) String() string    { return "window/close" }
func (mm MouseMove) String() string   { return fmt.Sprintf("mouse/move/%v/%v", mm.X, mm.Y) }
func (md MouseDown) String() string   { return fmt.Sprintf("mouse/down/%v/%v/%s", md.X, md.Y, md.Mouse) }
func (mu MouseUp) String() string     { return fmt.Sprintf("mouse/up/%v/%v/%s", mu.X, mu.Y, mu.Mouse) }
func (ms MouseScroll) String() string { return fmt.Sprintf("mouse/scroll/%v/%v", ms.X, ms.Y) }
func (kt KbType) String() string      { return fmt.Sprintf("keyboad/type/%v", kt.rune) }
func (kd KbDown) String() string      { return fmt.Sprintf("keyboad/down/%s", kd.Key) }
func (ku KbUp) String() string        { return fmt.Sprintf("keyboad/up/%s", ku.Key) }
func (kr KbRepeat) String() string    { return fmt.Sprintf("keyboad/repeat/%s", kr.Key) }
func (rs Resize) String() string      { return fmt.Sprintf("viewport/resize/%v/%v", rs.Width, rs.Height) }

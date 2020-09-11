package geui

// Button indicates a mouse button in an event.
type Mouse string

// List of all mouse buttons.
const (
	MouseLeft   Mouse = "left"
	MouseRight  Mouse = "right"
	MouseMiddle Mouse = "middle"
)


// Key indicates a keyboard key in an event.
type Key string

// List of all keyboard keys.
const (
	KeyLeft      Key = "left"
	KeyRight     Key = "right"
	KeyUp        Key = "up"
	KeyDown      Key = "down"
	KeyEscape    Key = "escape"
	KeySpace     Key = "space"
	KeyBackspace Key = "backspace"
	KeyDelete    Key = "delete"
	KeyEnter     Key = "enter"
	KeyTab       Key = "tab"
	KeyHome      Key = "home"
	KeyEnd       Key = "end"
	KeyPageUp    Key = "pageup"
	KeyPageDown  Key = "pagedown"
	KeyShift     Key = "shift"
	KeyCtrl      Key = "ctrl"
	KeyAlt       Key = "alt"
)
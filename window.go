package geui

import (
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func init() {
	runtime.LockOSThread()
}

type WindowOption func(*windowOptions)

type windowOptions struct {
	title         string
	width, height float64
	resizable     bool
	borderless    bool
	maximized     bool
}

func Title(title string) WindowOption {
	return func(o *windowOptions) {
		o.title = title
	}
}

// Size option sets the width and height of the window.
func Size(width, height float64) WindowOption {
	return func(o *windowOptions) {
		o.width = width
		o.height = height
	}
}

// Resizable option makes the window resizable by the user.
func Resizable() WindowOption {
	return func(o *windowOptions) {
		o.resizable = true
	}
}

// Borderless option makes the window borderless.
func Borderless() WindowOption {
	return func(o *windowOptions) {
		o.borderless = true
	}
}

// Maximized option makes the window start maximized.
func Maximized() WindowOption {
	return func(o *windowOptions) {
		o.maximized = true
	}
}

func NewWindow(n *Node, options ...WindowOption) (*Window, error) {
	o := windowOptions{
		title:      "",
		width:      640,
		height:     480,
		resizable:  false,
		borderless: false,
		maximized:  false,
	}
	for _, opt := range options {
		opt(&o)
	}

	w := &Window{
		events:      make(chan Event, 16),
		node:        n,
		loadedFonts: make(map[string]*truetype.Font),
	}

	var err error
	w.ctx, err = initGLFW(&o)
	if err != nil {
		panic(err)
	}

	w.canvas = gg.NewContext(int(o.width), int(o.height))

	w.initEvent()
	return w, nil
}

func initGLFW(o *windowOptions) (*glfw.Window, error) {
	err := glfw.Init()
	if err != nil {
		return nil, err
	}
	glfw.WindowHint(glfw.DoubleBuffer, glfw.False)
	if o.resizable {
		glfw.WindowHint(glfw.Resizable, glfw.True)
	} else {
		glfw.WindowHint(glfw.Resizable, glfw.False)
	}
	if o.borderless {
		glfw.WindowHint(glfw.Decorated, glfw.False)
	}
	if o.maximized {
		glfw.WindowHint(glfw.Maximized, glfw.True)
	}
	w, err := glfw.CreateWindow(int(o.width), int(o.height), o.title, nil, nil)
	if err != nil {
		return nil, err
	}
	if o.maximized {
		w, h := w.GetFramebufferSize()
		o.width, o.height = float64(w), float64(h)
	}
	return w, nil
}

type Window struct {
	ctx         *glfw.Window
	events      chan Event
	canvas      *gg.Context
	node        *Node
	newSize     chan image.Rectangle
	loadedFonts map[string]*truetype.Font
}

var buttons = map[glfw.MouseButton]Mouse{
	glfw.MouseButtonLeft:   MouseLeft,
	glfw.MouseButtonRight:  MouseRight,
	glfw.MouseButtonMiddle: MouseMiddle,
}

var keys = map[glfw.Key]Key{
	glfw.KeyLeft:         KeyLeft,
	glfw.KeyRight:        KeyRight,
	glfw.KeyUp:           KeyUp,
	glfw.KeyDown:         KeyDown,
	glfw.KeyEscape:       KeyEscape,
	glfw.KeySpace:        KeySpace,
	glfw.KeyBackspace:    KeyBackspace,
	glfw.KeyDelete:       KeyDelete,
	glfw.KeyEnter:        KeyEnter,
	glfw.KeyTab:          KeyTab,
	glfw.KeyHome:         KeyHome,
	glfw.KeyEnd:          KeyEnd,
	glfw.KeyPageUp:       KeyPageUp,
	glfw.KeyPageDown:     KeyPageDown,
	glfw.KeyLeftShift:    KeyShift,
	glfw.KeyRightShift:   KeyShift,
	glfw.KeyLeftControl:  KeyCtrl,
	glfw.KeyRightControl: KeyCtrl,
	glfw.KeyLeftAlt:      KeyAlt,
	glfw.KeyRightAlt:     KeyAlt,
}

func (w *Window) initEvent() {
	var mx, my float64

	w.ctx.SetCursorPosCallback(func(_ *glfw.Window, x, y float64) {
		mx, my = x, y
		go func() {
			w.events <- MouseMove{
				X: mx,
				Y: my,
			}
		}()
	})

	w.ctx.SetMouseButtonCallback(func(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		btn, ok := buttons[button]
		if !ok {
			return
		}
		switch action {
		case glfw.Press:
			go func() {
				w.events <- MouseDown{
					X:     mx,
					Y:     my,
					Mouse: btn,
				}
			}()
		case glfw.Release:
			go func() {
				w.events <- MouseUp{
					X:     mx,
					Y:     my,
					Mouse: btn,
				}
			}()
		}
	})

	w.ctx.SetScrollCallback(func(_ *glfw.Window, xoff, yoff float64) {
		go func() {
			w.events <- MouseScroll{
				X: xoff,
				Y: yoff,
			}
		}()
	})

	w.ctx.SetCharCallback(func(_ *glfw.Window, r rune) {
		go func() {
			w.events <- KbType{
				r,
			}
		}()
	})

	w.ctx.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
		k, ok := keys[key]
		if !ok {
			return
		}
		switch action {
		case glfw.Press:
			go func() { w.events <- KbDown{k} }()
		case glfw.Release:
			go func() { w.events <- KbUp{k} }()
		case glfw.Repeat:
			go func() { w.events <- KbRepeat{k} }()
		}
	})

	w.ctx.SetFramebufferSizeCallback(func(_ *glfw.Window, width, height int) {
		go func() {
			w.events <- Resize{
				0, 0,
				float64(width), float64(height),
			}
		}()
	})

	w.ctx.SetCloseCallback(func(_ *glfw.Window) {
		go func() { w.events <- WindowClose{} }()
	})
}

func (w *Window) Show() {
	w.ctx.MakeContextCurrent()
	_ = gl.Init()
	w.flush(w.canvas.Image().Bounds())

	var box image.Rectangle
	for !w.ctx.ShouldClose() {
		select {
		case r := <-w.newSize:
			old := w.canvas.Image()
			w.canvas = gg.NewContextForRGBA(image.NewRGBA(r))
			w.canvas.DrawImage(old, 0, 0)
			box = box.Union(r)
		case <-w.events:
			w.flush(box)
		default:
			glfw.PollEvents()
			w.ctx.SwapBuffers()
			time.Sleep(time.Second / 60)
		}
	}
}

func (w *Window) flush(r image.Rectangle) {
	w.render()
	bounds := w.canvas.Image().Bounds()
	r = r.Intersect(bounds)
	if r.Empty() {
		return
	}

	tmp := image.NewRGBA(r)
	draw.Draw(tmp, r, w.canvas.Image(), r.Min, draw.Src)

	gl.DrawBuffer(gl.FRONT)
	gl.Viewport(
		int32(bounds.Min.X),
		int32(bounds.Min.Y),
		int32(bounds.Dx()),
		int32(bounds.Dy()),
	)
	gl.RasterPos2d(
		-1+2*float64(r.Min.X)/float64(bounds.Dx()),
		+1-2*float64(r.Min.Y)/float64(bounds.Dy()),
	)
	gl.PixelZoom(1, -1)
	gl.DrawPixels(
		int32(r.Dx()),
		int32(r.Dy()),
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		unsafe.Pointer(&tmp.Pix[0]),
	)
	gl.Flush()
}

func (w *Window) render() {
	var f func(*Node)
	f = func(n *Node) {
		if n != nil {
			switch n.Type {
			case ElementNode:
				w.canvas.DrawRectangle(n.Model.RelativeX, n.Model.RelativeY, n.Model.Width, n.Model.Height)
				w.canvas.SetHexColor(n.Style.BackgroundColor)
				w.canvas.Fill()
			case CharDataNode:
				f, ok := w.loadedFonts[n.Parent.Style.FontFamily]
				if !ok {
					fontBytes, err := ioutil.ReadFile(n.Parent.Style.FontFamily)
					if err != nil {
						log.Println(err)
						return
					}
					f, err = truetype.Parse(fontBytes)
					if err != nil {
						log.Println(err)
						return
					}
					w.loadedFonts[n.Parent.Style.FontFamily] = f
				}
				face := truetype.NewFace(f, &truetype.Options{
					Size: n.Parent.Style.FontSize,
				})
				w.canvas.SetFontFace(face)
				w.canvas.SetHexColor(n.Parent.Style.FontColor)
				width, height := n.Parent.Model.Width, n.Parent.Model.Height
				w.canvas.DrawStringAnchored(n.Data, n.Parent.Model.RelativeX+width/2, n.Parent.Model.RelativeY+height/2, 0.5, 0.5)
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	f(w.node)
}

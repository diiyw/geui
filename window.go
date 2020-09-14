package geui

import (
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"image"
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
		newSize:     make(chan image.Rectangle),
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
	ctx            *glfw.Window
	events         chan Event
	canvas         *gg.Context
	node           *Node
	newSize        chan image.Rectangle
	loadedFonts    map[string]*truetype.Font
	mouseX, mouseY float64
	active         *Node
}

func (w *Window) initEvent() {
	var mx, my float64

	w.ctx.SetCursorPosCallback(func(_ *glfw.Window, x, y float64) {
		w.mouseX, w.mouseY = x, y
		go func() {
			w.events <- MouseMove{
				X: mx,
				Y: my,
			}
		}()
	})

	w.ctx.SetMouseButtonCallback(func(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		switch action {
		case glfw.Press:
			go func() {
				w.events <- MouseDown{
					X:           mx,
					Y:           my,
					MouseButton: button,
				}
			}()
		case glfw.Release:
			// set active node
			w.active = w.node.GetActiveNode(w.mouseX, w.mouseY)
			go func() {
				w.events <- MouseUp{
					X:           mx,
					Y:           my,
					MouseButton: button,
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
		if w.active != nil {
			w.active.Value = append(w.active.Value, r)
		}
		go func() {
			w.events <- KbType{
				r,
			}
		}()
	})

	w.ctx.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
		switch action {
		case glfw.Press:
			go func() { w.events <- KbDown{key} }()
		case glfw.Release:
			if w.active != nil && key == glfw.KeyBackspace {
				if len(w.active.Value) != 0 {
					w.active.Value = w.active.Value[:len(w.active.Value)-1]
				}
			}
			go func() { w.events <- KbUp{key} }()
		case glfw.Repeat:
			go func() { w.events <- KbRepeat{key} }()
		}
	})

	w.ctx.SetFramebufferSizeCallback(func(_ *glfw.Window, width, height int) {
		go func() {
			w.newSize <- image.Rect(0, 0, width, height)
			w.events <- Resize{
				0, 0,
				float64(width), float64(height),
			}
		}()
	})
}

func (w *Window) Show() {
	w.ctx.MakeContextCurrent()
	_ = gl.Init()
	w.flush()
	var box image.Rectangle
	for !w.ctx.ShouldClose() {
		select {
		case r := <-w.newSize:
			old := w.canvas.Image()
			w.canvas = gg.NewContextForRGBA(image.NewRGBA(r))
			w.canvas.DrawImage(old, 0, 0)
			box = box.Union(r)
		case <-w.events:
			w.flush()
		default:
			glfw.PollEvents()
			w.ctx.SwapBuffers()
			time.Sleep(time.Second / 60)
		}
	}
}

func (w *Window) flush() {
	w.render()

	img := w.canvas.Image().(*image.RGBA)
	bounds := img.Bounds()
	defer func() {
		img = nil
	}()
	gl.DrawBuffer(gl.FRONT)
	gl.Viewport(
		int32(bounds.Min.X),
		int32(bounds.Min.Y),
		int32(bounds.Dx()),
		int32(bounds.Dy()),
	)
	gl.RasterPos2d(
		-1+2*float64(bounds.Min.X)/float64(bounds.Dx()),
		+1-2*float64(bounds.Min.Y)/float64(bounds.Dy()),
	)
	gl.PixelZoom(1, -1)
	gl.DrawPixels(
		int32(bounds.Dx()),
		int32(bounds.Dy()),
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		unsafe.Pointer(&img.Pix[0]),
	)
	gl.Flush()
}

func (w *Window) render() {
	var setFontFace = func(n *Node) {
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
	}
	var f func(*Node)
	f = func(n *Node) {
		if n != nil {
			switch n.Type {
			case ElementNode:
				switch n.Data {
				case "input":
					w.canvas.SetHexColor(n.Style.BorderColor)
					w.canvas.SetLineWidth(n.Style.BorderWidth)
					x, y := n.Model.RelativeX+0.5, n.Model.RelativeY+0.5
					nw, nh := n.Model.Width, n.Model.Height
					w.canvas.DrawLine(x, y, x+nw, y)
					w.canvas.DrawLine(x, y, x, y+nh)
					w.canvas.DrawLine(x+nw, y, x+nw, y+nh)
					w.canvas.DrawLine(x, y+nh, x+nw, y+nh)
					fw, _ := w.canvas.MeasureString(string(n.Value))
					fw /= 2
					afw := fw / float64(len(n.Value))
					v := n.Value
					if fw > nw {
						v = n.Value[len(n.Value)-int(nw/afw)-1:]
						fw = nw - 5
					}
					if n.Focused(w.mouseX, w.mouseY) {
						w.canvas.DrawLine(fw+x+5, y+5, fw+x+5, nh+y-5)
					}
					w.canvas.Stroke()
					setFontFace(n)
					w.canvas.SetHexColor(n.Style.FontColor)
					if len(n.Value) != 0 {
						w.canvas.DrawStringAnchored(string(v), n.Model.RelativeX+6, n.Model.RelativeY+nh/2, 0, 0.4)
					}
				default:
					w.canvas.DrawRectangle(n.Model.RelativeX, n.Model.RelativeY, n.Model.Width, n.Model.Height)
					if n.Focused(w.mouseX, w.mouseY) {
						w.canvas.SetHexColor(n.Style.HoverColor)
					} else {
						w.canvas.SetHexColor(n.Style.BackgroundColor)
					}
					w.canvas.Fill()
				}
			case CharDataNode:
				setFontFace(n)
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

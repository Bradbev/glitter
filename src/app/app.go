package app

import (
	"image"
	"runtime"
	"time"

	"github.com/AllenDang/cimgui-go/backend"
	"github.com/AllenDang/cimgui-go/backend/glfwbackend"
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Bradbev/glitter/src/ren"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type App struct {
	OnPostCreate    func()
	OnBeforeDestroy func()
	OnBeforeRender  func()
	OnPostRender    func()
	OnDrop          func(p []string)
	//OnClose         func(b backend.Backend[sdlbackend.SDLWindowFlags])
	OnClose func(b backend.Backend[glfwbackend.GLFWWindowFlags])

	BgColor     imgui.Vec4
	Icons       *image.RGBA
	WindowTitle string
	Width       int
	Height      int

	//backend backend.Backend[sdlbackend.SDLWindowFlags]
	backend backend.Backend[glfwbackend.GLFWWindowFlags]
}

func (a *App) GetSize() (int32, int32) {
	return a.backend.DisplaySize()
}

func (a *App) RunNoDt(loop func()) {
	a.Run(func(dt float32) {
		loop()
	})
}

func (a *App) SetMousePos(x, y float32) {
	a.backend.SetCursorPos(float64(x), float64(y))
}

func (a *App) Run(loop func(dt float32)) {

	if err := gl.Init(); err != nil {
		panic(err)
	}

	runtime.LockOSThread()

	be := a.backend
	be.SetAfterCreateContextHook(a.OnPostCreate)
	be.SetBeforeDestroyContextHook(a.OnBeforeDestroy)
	be.SetBgColor(a.BgColor)
	be.CreateWindow(a.WindowTitle, a.Width, a.Height)
	be.SetDropCallback(a.OnDrop)
	be.SetCloseCallback(a.OnClose)
	if a.Icons != nil {
		be.SetIcons(a.Icons)
	}

	be.SetBeforeRenderHook(a.OnBeforeRender)
	be.SetAfterRenderHook(a.OnPostRender)

	loopTime := time.Now()
	be.Run(func() {
		dt := float32(time.Now().Sub(loopTime)) / float32(time.Second)
		loopTime = time.Now()
		loop(dt)
	})
}

func nop() {}

func Default() *App {
	//currentBackend, _ := backend.CreateBackend(sdlbackend.NewSDLBackend())
	currentBackend, _ := backend.CreateBackend(glfwbackend.NewGLFWBackend())
	return &App{
		backend:         currentBackend,
		OnPostCreate:    nop,
		OnPostRender:    nop,
		OnBeforeDestroy: nop,
		OnDrop:          func(p []string) {},
		//OnClose:         func(b backend.Backend[sdlbackend.SDLWindowFlags]) {},
		OnClose: func(b backend.Backend[glfwbackend.GLFWWindowFlags]) {},

		WindowTitle: "Title",
		Width:       1200,
		Height:      900,
		BgColor:     imgui.NewVec4(0.45, 0.55, 0.6, 1.0),
	}
}

type ImguiCamera struct {
	Camera       *ren.Camera
	Speed        float32
	MouseSpeed   float32
	LastMousePos mgl32.Vec2
	mousePos     mgl32.Vec2
	mouseDown    bool
}

func (c *ImguiCamera) ProcessInput(app *App, dt float32) {
	cam := c.Camera
	speed := c.Speed * dt
	right := cam.Forward.Cross(cam.Up).Normalize()
	forward2D := cam.Forward.Mul(speed)
	if imgui.IsKeyDown(imgui.KeyA) {
		cam.Position = cam.Position.Sub(right.Mul(speed))
	}
	if imgui.IsKeyDown(imgui.KeyD) {
		cam.Position = cam.Position.Add(right.Mul(speed))
	}
	if imgui.IsKeyDown(imgui.KeyW) {
		cam.Position = cam.Position.Add(forward2D)
	}
	if imgui.IsKeyDown(imgui.KeyS) {
		cam.Position = cam.Position.Sub(forward2D)
	}
	up := cam.Up.Mul(speed)
	if imgui.IsKeyDown(imgui.KeyE) {
		cam.Position = cam.Position.Add(up)
	}
	if imgui.IsKeyDown(imgui.KeyQ) {
		cam.Position = cam.Position.Sub(up)
	}
	anyWindowFocused := imgui.IsWindowFocusedV(imgui.FocusedFlagsAnyWindow)
	if !anyWindowFocused {
		// the first click on a window will enter this path as the window
		// is not yet focused
		if imgui.IsMouseDown(imgui.MouseButtonLeft) && !c.mouseDown {
			c.mousePos = mousePos()
		}
		if c.mouseDown {
			mousePos := mousePos()
			delta := c.mousePos.Sub(mousePos).Mul(c.MouseSpeed)

			wx, wy := app.backend.GetWindowPos()
			//wx, wy = 0, 0 // sdl offsets
			app.SetMousePos(c.mousePos.X()-float32(wx), c.mousePos.Y()-float32(wy))
			//c.mousePos = mousePos
			rot := mgl32.QuatRotate(mgl32.DegToRad(delta.X()*dt), mgl32.Vec3{0, 1, 0})
			cam.Forward = rot.Rotate(cam.Forward).Normalize()
			rot = mgl32.QuatRotate(mgl32.DegToRad(delta.Y()*dt), right)
			cam.Forward = rot.Rotate(cam.Forward).Normalize()
		}
		c.mouseDown = imgui.IsMouseDown(imgui.MouseButtonLeft)
	} else {
		c.mouseDown = false
	}
}

func mousePos() mgl32.Vec2 {
	p := imgui.MousePos()
	mx, my := p.X, p.Y
	//mx, my, _ := sdl.GetGlobalMouseState()
	return mgl32.Vec2{
		float32(mx),
		float32(my),
	}
}

func NewCamera() *ImguiCamera {
	return &ImguiCamera{
		Camera:     ren.NewCamera(),
		Speed:      1,
		MouseSpeed: 10,
	}
}

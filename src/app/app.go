package app

import (
	"image"
	"runtime"

	imgui "github.com/AllenDang/cimgui-go"
	"github.com/AllenDang/cimgui-go/backend"
	"github.com/AllenDang/cimgui-go/backend/sdlbackend"
	"github.com/go-gl/gl/v2.1/gl"
)

type App struct {
	OnPostCreate    func()
	OnBeforeDestroy func()
	OnBeforeRender  func()
	OnPostRender    func()
	OnDrop          func(p []string)
	OnClose         func(b backend.Backend[sdlbackend.SDLWindowFlags])

	BgColor     imgui.Vec4
	Icons       *image.RGBA
	WindowTitle string
	Width       int
	Height      int

	backend backend.Backend[sdlbackend.SDLWindowFlags]
}

func (a *App) GetSize() (int32, int32) {
	return a.backend.DisplaySize()
}

func (a *App) Run(loop func()) {

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

	be.Run(loop)
}

func nop() {}

func Default() *App {
	currentBackend, _ := backend.CreateBackend(sdlbackend.NewSDLBackend())
	return &App{
		backend:         currentBackend,
		OnPostCreate:    nop,
		OnPostRender:    nop,
		OnBeforeDestroy: nop,
		OnDrop:          func(p []string) {},
		OnClose:         func(b backend.Backend[sdlbackend.SDLWindowFlags]) {},

		WindowTitle: "Title",
		Width:       1200,
		Height:      900,
		BgColor:     imgui.NewVec4(0.45, 0.55, 0.6, 1.0),
	}
}

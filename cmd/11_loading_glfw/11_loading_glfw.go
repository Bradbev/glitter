package main

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"runtime"

	"github.com/Bradbev/glitter/src/app"
	"github.com/Bradbev/glitter/src/asset"
	"github.com/Bradbev/glitter/src/ren"
	"github.com/bloeys/assimp-go/asig/asig"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/laher/mergefs"
)

var (
	verticies = []float32{
		0.5, 0.5, 0.0, // top right
		0.5, -0.5, 0.0, // bottom right
		-0.5, -0.5, 0.0, // bottom left
		-0.5, 0.5, 0.0, // top left
	}
	indices = []uint32{
		0, 1, 3, // first triangle
		1, 2, 3, // second triangle
	}
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

//go:embed *.vert *.frag
var embeddedAssets embed.FS
var assets fs.FS

func must[T any](t T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, gl.TRUE)

	window, err := glfw.CreateWindow(1200, 900, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	gl.Init()
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
	gl.Enable(gl.DEPTH_TEST)

	assets = mergefs.Merge(embeddedAssets, //
		os.DirFS("assets"),       // begin run from the root,
		os.DirFS("../../assets")) // run from inside this dir,

	objShader := must(ren.NewProgramFS(assets, "vertex.vert", "obj.frag"))
	//quadShader := must(ren.NewProgramFS(assets, "part001.vert", "part001.frag"))
	vao := ren.NewVertexAttribObject()
	vao.Float32AttribData(vao.NextSlot(), 3, verticies, gl.STATIC_DRAW)
	vao.IndexData(indices, gl.STATIC_DRAW)

	camera := app.NewCamera()
	camera.Camera.Position = mgl32.Vec3{0, 0, 5}
	camera.Camera.LookAt(mgl32.Vec3{0, 0, 0})
	scene := must(asset.ImportFile("/Users/bradbeveridge/dev2/3rdparty/LearnOpenGL/resources/objects/backpack/backpack.obj", asig.PostProcessTriangulate|asig.PostProcessJoinIdenticalVertices))
	scene.Setup()

	sx, sy := int32(1200), int32(900)
	for !window.ShouldClose() {
		// Do OpenGL stuff.
		grey := float32(0.05)
		gl.ClearColor(grey, grey, grey, 0)
		gl.Clear(gl.DEPTH_BUFFER_BIT | gl.COLOR_BUFFER_BIT)

		objShader.UseProgram()
		gl.Viewport(0, 0, sx*2, sy*2)
		camera.Camera.CacheMatricies(float32(sx), float32(sy))
		view, projection := camera.Camera.GetMatrices()

		model := mgl32.Translate3D(0, 0, 0)
		objShader.UniformMatrix4f("model", model)
		objShader.UniformMatrix4f("view", view)
		objShader.UniformMatrix4f("projection", projection)

		scene.Draw(objShader)

		//quadShader.UseProgram()
		//vao.Enable()
		//		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

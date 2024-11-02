// https://learnopengl.com/Lighting/Colors
package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"runtime"

	"github.com/AllenDang/cimgui-go/imgui"

	"github.com/Bradbev/glitter/src/app"
	"github.com/Bradbev/glitter/src/imguix"
	"github.com/Bradbev/glitter/src/ren"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/laher/mergefs"
)

var (
	cube = []float32{
		-0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, 0.5, -0.5,
		0.5, 0.5, -0.5,
		-0.5, 0.5, -0.5,
		-0.5, -0.5, -0.5,

		-0.5, -0.5, 0.5,
		0.5, -0.5, 0.5,
		0.5, 0.5, 0.5,
		0.5, 0.5, 0.5,
		-0.5, 0.5, 0.5,
		-0.5, -0.5, 0.5,

		-0.5, 0.5, 0.5,
		-0.5, 0.5, -0.5,
		-0.5, -0.5, -0.5,
		-0.5, -0.5, -0.5,
		-0.5, -0.5, 0.5,
		-0.5, 0.5, 0.5,

		0.5, 0.5, 0.5,
		0.5, 0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, 0.5,
		0.5, 0.5, 0.5,

		-0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, 0.5,
		0.5, -0.5, 0.5,
		-0.5, -0.5, 0.5,
		-0.5, -0.5, -0.5,

		-0.5, 0.5, -0.5,
		0.5, 0.5, -0.5,
		0.5, 0.5, 0.5,
		0.5, 0.5, 0.5,
		-0.5, 0.5, 0.5,
		-0.5, 0.5, -0.5,
	}
	normals = []float32{
		0.0, 0.0, -1.0,
		0.0, 0.0, -1.0,
		0.0, 0.0, -1.0,
		0.0, 0.0, -1.0,
		0.0, 0.0, -1.0,
		0.0, 0.0, -1.0,

		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,

		-1.0, 0.0, 0.0,
		-1.0, 0.0, 0.0,
		-1.0, 0.0, 0.0,
		-1.0, 0.0, 0.0,
		-1.0, 0.0, 0.0,
		-1.0, 0.0, 0.0,

		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,

		0.0, -1.0, 0.0,
		0.0, -1.0, 0.0,
		0.0, -1.0, 0.0,
		0.0, -1.0, 0.0,
		0.0, -1.0, 0.0,
		0.0, -1.0, 0.0,

		0.0, 1.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 1.0, 0.0,
	}
	colors = []float32{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,
	}
	texCoords = []float32{
		0.0, 0.0,
		1.0, 0.0,
		1.0, 1.0,
		1.0, 1.0,
		0.0, 1.0,
		0.0, 0.0,

		0.0, 0.0,
		1.0, 0.0,
		1.0, 1.0,
		1.0, 1.0,
		0.0, 1.0,
		0.0, 0.0,

		1.0, 0.0,
		1.0, 1.0,
		0.0, 1.0,
		0.0, 1.0,
		0.0, 0.0,
		1.0, 0.0,

		1.0, 0.0,
		1.0, 1.0,
		0.0, 1.0,
		0.0, 1.0,
		0.0, 0.0,
		1.0, 0.0,

		0.0, 1.0,
		1.0, 1.0,
		1.0, 0.0,
		1.0, 0.0,
		0.0, 0.0,
		0.0, 1.0,

		0.0, 1.0,
		1.0, 1.0,
		1.0, 0.0,
		1.0, 0.0,
		0.0, 0.0,
		0.0, 1.0,
	}

	showDemoWindow  bool
	lightColor              = mgl32.Vec3{1, 1, 1}
	objColor                = mgl32.Vec3{1, 0.5, 0.31}
	lightPos                = [3]float32{1.0, 0.0, 0.0}
	blend           float32 = 0.5
	rotation        float32 = 0
	scale           float32 = 0.5
	x, y, z         float32 = 0, 0, 0
	fov             float32 = 45
	showDebugCamera         = true
	normCount       int32   = 1
	normIdx         int32   = 0
)

type TestEdit struct {
	FloatSlider  float32  `ui:"slider:{min:0, max:20}"`
	FloatDrag    float32  `ui:"drag"`
	FloatDragp   *float32 `ui:"drag"`
	FloatInput   float32  `ui:"input"`
	FloatDefault float32
	IntSlider    int32  `ui:"slider:{min:0, max:20}"`
	IntSliderp   *int32 `ui:"slider:{min:0, max:200}"`
	IntDrag      int32  `ui:"drag"`
	IntInput     int32  `ui:"input"`
	IntDefault   int32
	Vec3Slider   mgl32.Vec3 `ui:"slider:{min:0, max:20}"`
	Vec3Drag     mgl32.Vec3 `ui:"drag:{min:0, max:20}"`
	Vec3Input    mgl32.Vec3 `ui:"input"`
	Vec3Default  mgl32.Vec3
	Vec3Color    mgl32.Vec3  `ui:"color"`
	Vec3ColorP   *mgl32.Vec3 `ui:"color"`
	String       string
}

var testEdit = TestEdit{
	IntSliderp: new(int32),
	FloatDragp: new(float32),
	Vec3ColorP: new(mgl32.Vec3),
}

// makeVAO initializes and returns a vertex array from the points provided.
func showWidgetsDemo(camera *ren.Camera) {
	if showDemoWindow {
		imgui.ShowDemoWindowV(&showDemoWindow)
	}

	imgui.SetNextWindowSizeV(imgui.NewVec2(700, 700), imgui.CondOnce)
	imgui.Begin("Window 1")
	imguix.EditStruct("TestEdit", &testEdit)
	imgui.Text(fmt.Sprintf("%v", testEdit))
	imgui.Text(fmt.Sprintf("%v", *testEdit.IntSliderp))
	imgui.Checkbox("Show demo window", &showDemoWindow)

	imguix.TreeNode("Group", func() {
		imgui.ColorEdit3("Obj", ren.Vec3ToRaw(&objColor))
		imgui.ColorEdit3("Light", ren.Vec3ToRaw(&lightColor))
	})

	imgui.SliderFloat3("LightPos", &lightPos, -2, 2)
	imgui.SliderFloat("Blend", &blend, 0, 1)
	imgui.SliderFloat("Scale", &scale, 0, 1)
	imgui.SliderFloat("Rotation", &rotation, 0, 180)
	imgui.SliderFloat("x", &x, -1, 1)
	imgui.SliderFloat("y", &y, -1, 1)
	imgui.SliderFloat("z", &z, -1, 1)
	imgui.SliderFloat("fov", &fov, 10, 180)
	imgui.SliderInt("NormsCount", &normCount, 1, 100)
	imgui.SliderInt("NormsIndex", &normIdx, 0, 100)

	imgui.Checkbox("Show Camera Debug", &showDebugCamera)
	if showDebugCamera {
		imgui.Text(fmt.Sprintf("Pos %v\nForward %v", camera.Position, camera.Forward))
	}

	imgui.End()
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

type Light struct {
	Position mgl32.Vec3
	Ambient  mgl32.Vec3
	Diffuse  mgl32.Vec3
	Specular mgl32.Vec3
}

var light = &Light{
	Ambient:  mgl32.Vec3{0.2, 0.2, 0.2},
	Diffuse:  mgl32.Vec3{0.5, 0.5, 0.5},
	Specular: mgl32.Vec3{1.0, 1.0, 1.0},
}

type Material struct {
	Diffuse   int32
	Specular  int32
	Shininess float32
}

var mat = &Material{
	Diffuse:   0,
	Specular:  1,
	Shininess: 64,
}

func main() {
	assets = mergefs.Merge(embeddedAssets, //
		os.DirFS("assets"),       // begin run from the root,
		os.DirFS("../../assets")) // run from inside this dir,
	a := app.Default()
	var (
		cubeShader    *ren.Program
		lampShader    *ren.Program
		lightVao      *ren.VertexAttribObject
		container     *ren.Texture
		containerSpec *ren.Texture
	)
	camera := app.NewCamera()
	camera.Camera.Position = mgl32.Vec3{3, 2, 0}
	camera.Camera.LookAt(mgl32.Vec3{0, 0, 0})

	a.OnPostCreate = func() {
		version := gl.GoStr(gl.GetString(gl.VERSION))
		log.Println("OpenGL version", version)
		var err error
		container = must(ren.NewTextureFS(assets, "container2.png", gl.REPEAT, gl.REPEAT))
		containerSpec = must(ren.NewTextureFS(assets, "container2_specular.png", gl.REPEAT, gl.REPEAT))
		cubeShader = must(ren.NewProgramFS(assets, "vertex.vert", "fragment.frag"))

		lightVao = ren.NewVertexAttribObject()
		lightVao.Float32AttribData(lightVao.NextSlot(), 3, cube, gl.STATIC_DRAW)
		lightVao.Float32AttribData(lightVao.NextSlot(), 3, normals, gl.STATIC_DRAW)
		lightVao.Float32AttribData(lightVao.NextSlot(), 2, texCoords, gl.STATIC_DRAW)

		lampShader, err = ren.NewProgramFS(assets, "vertex.vert", "cube.frag")
		if err != nil {
			log.Fatal(err)
		}

		runtime.GC()
		gl.Enable(gl.DEPTH_TEST)
	}

	a.Run(func(dt float32) {
		ren.GarbageCollect()
		grey := float32(0.2)
		gl.ClearColor(grey, grey, grey, 0)
		gl.Clear(gl.DEPTH_BUFFER_BIT | gl.COLOR_BUFFER_BIT)

		showWidgetsDemo(camera.Camera)
		sx, sy := a.GetSize()
		// * 2 because retina?  Imgui must be doing retina somewhere else?
		gl.Viewport(0, 0, sx*2, sy*2)

		camera.ProcessInput(a, dt)
		camera.Camera.CacheMatricies(float32(sx), float32(sy))
		view, projection := camera.Camera.GetMatrices()

		// binding the vao also binds the ebo
		lightVao.Enable()

		// show the light position
		model := mgl32.Translate3D(lightPos[0], lightPos[1], lightPos[2])
		model = model.Mul4(mgl32.Scale3D(0.2, 0.2, 0.2))
		lampShader.UseProgram()
		lampShader.UniformMatrix4f("projection", projection)
		lampShader.UniformMatrix4f("view", view)
		lampShader.UniformMatrix4f("model", model)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		// render the normal cube
		cubeShader.UseProgram()
		cubeShader.UniformMatrix4f("projection", projection)
		cubeShader.UniformMatrix4f("view", view)
		cubeShader.UniformVec3("objectColor", objColor)

		light.Position = lightPos
		cubeShader.UniformStruct("light", light)

		container.Bind(gl.TEXTURE0)
		containerSpec.Bind(gl.TEXTURE1)
		cubeShader.UniformStruct("material", mat)

		cubeShader.UniformVec3("viewPos", camera.Camera.Position)
		model = mgl32.HomogRotate3DZ(mgl32.DegToRad(rotation))
		cubeShader.UniformMatrix4f("model", model)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
	})
}

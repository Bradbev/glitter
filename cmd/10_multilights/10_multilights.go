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
	cubePositions = []mgl32.Vec3{
		{0.0, 0.0, 0.0},
		{2.0, 5.0, -15.0},
		{-1.5, -2.2, -2.5},
		{-3.8, -2.0, -12.3},
		{2.4, -0.4, -3.5},
		{-1.7, 3.0, -7.5},
		{1.3, -2.0, -2.5},
		{1.5, 2.0, -2.5},
		{1.5, 0.2, -1.5},
		{-1.3, 1.0, -1.5},
	}
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
	lightPos                = [3]float32{1.0, 0.0, 0.0}
	blend           float32 = 0.5
	rotation        float32 = 0
	fov             float32 = 45
	showDebugCamera         = false
)

var (
	white = mgl32.Vec3{1, 1, 1}
	up    = mgl32.Vec3{0, 0, 1}
)

type DirectionalLight struct {
	Direction mgl32.Vec3 `ui:"slider:{min:-1,max:1}"`
	Ambient   mgl32.Vec3 `ui:"color"`
	Diffuse   mgl32.Vec3 `ui:"color"`
	Specular  mgl32.Vec3 `ui:"color"`
}
type PointLight struct {
	Position mgl32.Vec3 `ui:"slider:{min:-15,max:15}"`

	Constant  float32
	Linear    float32
	Quadratic float32

	Ambient  mgl32.Vec3 `ui:"color"`
	Diffuse  mgl32.Vec3 `ui:"color"`
	Specular mgl32.Vec3 `ui:"color"`
}

var (
	directionalLight = DirectionalLight{
		Direction: mgl32.Vec3{-1, -0.5, 0},
		Ambient:   white.Mul(0.05),
		Diffuse:   white.Mul(0.1),
		Specular:  white,
	}
	defaultPoint = PointLight{
		// see https://learnopengl.com/Lighting/Light-casters for values
		Constant:  1,
		Linear:    0.09,
		Quadratic: 0.032,
		Ambient:   white.Mul(0.0),
		Diffuse:   white.Mul(0.4),
		Specular:  white,
	}
	pointLights = []PointLight{defaultPoint, defaultPoint, defaultPoint, defaultPoint}
)

// makeVAO initializes and returns a vertex array from the points provided.
func showWidgetsDemo(camera *ren.Camera) {
	if showDemoWindow {
		imgui.ShowDemoWindowV(&showDemoWindow)
	}

	imgui.SetNextWindowSizeV(imgui.NewVec2(500, 500), imgui.CondOnce)
	imgui.Begin("Window 1")
	imgui.Checkbox("Show demo window", &showDemoWindow)
	imguix.EditStruct("Directional", &directionalLight)
	for i := range pointLights {
		imguix.EditStruct(fmt.Sprintf("pointLights[%d]", i), &pointLights[i])
	}

	imguix.TreeNode("Group", func() {
		imgui.ColorEdit3("Light", ren.Vec3ToRaw(&lightColor))
	})

	imgui.SliderFloat3("LightPos", &lightPos, -2, 2)
	imgui.SliderFloat("Blend", &blend, 0, 1)
	imgui.SliderFloat("Rotation", &rotation, 0, 180)
	imgui.SliderFloat("fov", &fov, 10, 180)

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
	pointLights[0].Position = mgl32.Vec3{0.7, 0.2, 2.0}
	pointLights[1].Position = mgl32.Vec3{2.3, -3.3, -4.0}
	pointLights[2].Position = mgl32.Vec3{-4.0, 2.0, -12.0}
	pointLights[3].Position = mgl32.Vec3{0.0, 0.0, -3.0}

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

		// show the point lights
		for _, l := range pointLights {
			model := mgl32.Translate3D(l.Position[0], l.Position[1], l.Position[2])
			model = model.Mul4(mgl32.Scale3D(0.2, 0.2, 0.2))
			lampShader.UseProgram()
			lampShader.UniformMatrix4f("projection", projection)
			lampShader.UniformMatrix4f("view", view)
			lampShader.UniformMatrix4f("model", model)
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		// render the normal cube(s)
		for _, pos := range cubePositions {
			cubeShader.UseProgram()
			model := mgl32.Translate3D(pos[0], pos[1], pos[2])
			model = model.Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(rotation)))
			cubeShader.UniformMatrix4f("model", model)
			cubeShader.UniformMatrix4f("view", view)
			cubeShader.UniformMatrix4f("projection", projection)

			cubeShader.UniformVec3("viewPos", camera.Camera.Position) // for specular

			for i, l := range pointLights {
				cubeShader.UniformStruct(fmt.Sprintf("pointLights[%d]", i), &l)
			}
			cubeShader.UniformStruct("directionalLight", &directionalLight)

			container.Bind(gl.TEXTURE0)
			containerSpec.Bind(gl.TEXTURE1)
			cubeShader.UniformStruct("material", mat)

			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}
	})
}

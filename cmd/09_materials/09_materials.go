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
	"github.com/Bradbev/glitter/src/ren"
	"github.com/go-gl/gl/v2.1/gl"
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

	showDemoWindow  bool
	lightColor              = [3]float32{1, 1, 1}
	objColor                = [3]float32{1, 0.5, 0.31}
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

// makeVAO initializes and returns a vertex array from the points provided.
func showWidgetsDemo(camera *ren.Camera) {
	if showDemoWindow {
		imgui.ShowDemoWindowV(&showDemoWindow)
	}

	imgui.SetNextWindowSizeV(imgui.NewVec2(300, 400), imgui.CondOnce)
	imgui.Begin("Window 1")
	imgui.Checkbox("Show demo window", &showDemoWindow)
	imgui.ColorEdit3("Obj", &objColor)
	imgui.ColorEdit3("Light", &lightColor)
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

func main() {
	assets = mergefs.Merge(embeddedAssets, //
		os.DirFS("assets"),       // begin run from the root,
		os.DirFS("../../assets")) // run from inside this dir,
	a := app.Default()
	var (
		cubeShader *ren.Program
		lampShader *ren.Program
		lightVao   *ren.VertexAttribObject
	)
	camera := app.NewCamera()
	camera.Camera.Position = mgl32.Vec3{3, 2, 0}
	camera.Camera.LookAt(mgl32.Vec3{0, 0, 0})

	a.OnPostCreate = func() {
		version := gl.GoStr(gl.GetString(gl.VERSION))
		log.Println("OpenGL version", version)
		var err error
		cubeShader, err = ren.NewProgramFS(assets, "vertex.vert", "fragment.frag")
		if err != nil {
			log.Fatal(err)
		}

		lightVao = ren.NewVertexAttribObject()
		lightVao.Float32AttribData(lightVao.NextSlot(), 3, cube, gl.STATIC_DRAW)
		lightVao.Float32AttribData(lightVao.NextSlot(), 3, normals, gl.STATIC_DRAW)

		lampShader, err = ren.NewProgramFS(assets, "vertex.vert", "cube.frag")
		if err != nil {
			log.Fatal(err)
		}

		runtime.GC()
		gl.Enable(gl.DEPTH_TEST)
	}

	a.Run(func(dt float32) {
		ren.GarbageCollect()
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		showWidgetsDemo(camera.Camera)
		sx, sy := a.GetSize()
		// * 2 because retina?  Imgui must be doing retina somewhere else?
		gl.Viewport(0, 0, sx*2, sy*2)

		camera.ProcessInput(a, dt)
		camera.Camera.CacheMatricies(float32(sx), float32(sy))
		view, projection := camera.Camera.GetMatrices()
		//view := camera.Camera.GetViewMat()
		//projection := camera.Camera.GetProjectionMat(float32(sx), float32(sy))

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

		cubeShader.UseProgram()
		cubeShader.UniformMatrix4f("projection", projection)
		cubeShader.UniformMatrix4f("view", view)
		cubeShader.UniformVec3("objectColor", objColor)

		cubeShader.UniformVec3("light.ambient", mgl32.Vec3{0.2, 0.2, 0.2})
		cubeShader.UniformVec3("light.diffuse", mgl32.Vec3{0.5, 0.5, 0.5}) // darken diffuse light a bit
		cubeShader.UniformVec3("light.specular", mgl32.Vec3{1.0, 1.0, 1.0})
		cubeShader.UniformVec3("light.position", lightPos)

		cubeShader.UniformVec3("material.ambient", mgl32.Vec3{1.0, 0.5, 0.31})
		cubeShader.UniformVec3("material.diffuse", mgl32.Vec3{1.0, 0.5, 0.31})
		cubeShader.UniformVec3("material.specular", mgl32.Vec3{0.5, 0.5, 0.5})
		cubeShader.Uniform1f("material.shininess", 32.0)

		cubeShader.UniformVec3("lightColor", lightColor)
		cubeShader.UniformVec3("viewPos", camera.Camera.Position)
		// render the normal cube
		model = mgl32.HomogRotate3DZ(mgl32.DegToRad(rotation))
		cubeShader.UniformMatrix4f("model", model)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
	})
}

// Add in some matrix math
package main

import (
	"embed"
	"log"
	"runtime"

	"github.com/AllenDang/cimgui-go/imgui"

	"github.com/Bradbev/glitter/src/app"
	"github.com/Bradbev/glitter/src/ren"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
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
	cubeTex = []float32{
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
	colors = []float32{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,
	}

	showDemoWindow bool
	color4         [4]float32 = [4]float32{1, 0, 0, 0}
	blend          float32    = 0.5
	rotation       float32    = 0
	scale          float32    = 0.5
	x, y, z        float32    = 0, 0, 0
	fov            float32    = 45
)

// makeVAO initializes and returns a vertex array from the points provided.
func showWidgetsDemo() {
	if showDemoWindow {
		imgui.ShowDemoWindowV(&showDemoWindow)
	}

	imgui.SetNextWindowSizeV(imgui.NewVec2(300, 300), imgui.CondOnce)
	imgui.Begin("Window 1")
	imgui.Checkbox("Show demo window", &showDemoWindow)
	imgui.ColorEdit4("Color Edit", &color4)
	imgui.SliderFloat("Blend", &blend, 0, 1)
	imgui.SliderFloat("Scale", &scale, 0, 1)
	imgui.SliderFloat("Rotation", &rotation, 0, 180)
	imgui.SliderFloat("x", &x, -1, 1)
	imgui.SliderFloat("y", &y, -1, 1)
	imgui.SliderFloat("z", &z, -1, 1)
	imgui.SliderFloat("fov", &fov, 10, 180)
	imgui.End()
}

//go:embed *.vert *.frag *.jpg *.jpeg
var assets embed.FS

func main() {
	a := app.Default()
	var (
		p       *ren.Program
		wallTex *ren.Texture
		testTex *ren.Texture
		vao     *ren.VertexAttribObject
	)
	camera := app.NewCamera()
	a.OnPostCreate = func() {
		version := gl.GoStr(gl.GetString(gl.VERSION))
		log.Println("OpenGL version", version)
		var err error
		p, err = ren.NewProgramFS(assets, "part007.vert", "part007.frag")
		if err != nil {
			log.Fatal(err)
		}

		wallTex, err = ren.NewTextureFS(assets, "wall.jpg", gl.REPEAT, gl.REPEAT)
		if err != nil {
			log.Fatal(err)
		}
		testTex, err = ren.NewTextureFS(assets, "test.jpeg", gl.REPEAT, gl.REPEAT)
		if err != nil {
			log.Fatal(err)
		}

		vao = ren.NewVertexAttribObject()
		vao.Float32AttribData(vao.NextSlot(), 3, cube, gl.STATIC_DRAW)
		vao.Float32AttribData(vao.NextSlot(), 3, colors, gl.STATIC_DRAW)
		vao.Float32AttribData(vao.NextSlot(), 2, cubeTex, gl.STATIC_DRAW)
		//vao.IndexData(indices, gl.STATIC_DRAW)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		// bind the GL numbered texture unit to our names
		p.UseProgram()
		p.Uniform1i("ourTexture", 0)
		p.Uniform1i("testTexture", 1)

		runtime.GC()
		gl.Enable(gl.DEPTH_TEST)
	}

	a.Run(func(dt float32) {
		ren.GarbageCollect()
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		showWidgetsDemo()
		sx, sy := a.GetSize()
		// * 2 because retina?  Imgui must be doing retina somewhere else?
		gl.Viewport(0, 0, sx*2, sy*2)

		// Camera - move to a struct
		//radius := 10.0
		//t := float64(mgl32.DegToRad(rotation))
		//camX := float32(math.Sin(t) * radius)
		//camZ := float32(math.Cos(t) * radius)
		//view := mgl32.LookAtV(
		//mgl32.Vec3{camX, 0, camZ},
		//mgl32.Vec3{0, 0, 0},
		//mgl32.Vec3{0, 1, 0})
		camera.ProcessInput(dt)
		view := camera.Camera.GetViewMat()
		projection := camera.Camera.GetProjectionMat(float32(sx), float32(sy))
		//projection := mgl32.Perspective(mgl32.DegToRad(fov), float32(sx)/float32(sy), 0.1, 100)

		for i, pos := range cubePositions {
			p.UseProgram()
			//p.Uniform4f("vertexColor", color4[0], color4[1], color4[2], color4[3])
			//view := mgl32.Translate3D(0, 0, -3)

			// the next 4 lines can be reordered to show how order of matrix multiplication
			// impacts the output.
			//trans := mgl32.Ident4() // only needed to make the next 3 lines more uniform
			//trans = trans.Mul4(mgl32.Translate3D(x, y, z))
			//trans = trans.Mul4(mgl32.QuatRotate(mgl32.DegToRad(rotation), mgl32.Vec3{0, 0, 1}).Mat4())
			//trans = trans.Mul4(mgl32.Scale3D(scale, scale, scale))

			p.Uniform1f("blend", blend)
			//p.UniformMatrix4f("transform", trans)
			p.UniformMatrix4f("projection", projection)
			p.UniformMatrix4f("view", view)

			// binding the vao also binds the ebo
			wallTex.Bind(gl.TEXTURE0)
			testTex.Bind(gl.TEXTURE1)
			vao.Enable()

			model := mgl32.Translate3D(pos.X(), pos.Y(), pos.Z())
			angle := (20.0 * float32(i))
			model = model.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(angle), mgl32.Vec3{1, 0.3, 0.5}.Normalize()))
			p.UniformMatrix4f("model", model)

			//gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil)
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}
	})
}

// Add in some matrix math
// https://learnopengl.com/Getting-started/Transformations
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
	verticies = []float32{
		0.5, 0.5, 0.0, // top right
		0.5, -0.5, 0.0, // bottom right
		-0.5, -0.5, 0.0, // bottom left
		-0.5, 0.5, 0.0, // top left
	}
	colors = []float32{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,
	}
	texCoords = []float32{
		1.0, 1.0, // top right
		1.0, 0.0, // bottom right
		0.0, 0.0, // bottom left
		0.0, 1.0, // top left
	}
	indices = []uint32{
		0, 1, 2, // triangle
		2, 3, 0, // triangle
	}
	showDemoWindow bool
	color4         [4]float32 = [4]float32{1, 0, 0, 0}
	blend          float32    = 0.5
	rotation       float32    = 0
	scale          float32    = 0.5
	x, y, z        float32    = 0, 0, 0
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
	a.OnPostCreate = func() {
		version := gl.GoStr(gl.GetString(gl.VERSION))
		log.Println("OpenGL version", version)
		var err error
		p, err = ren.NewProgramFS(assets, "part006.vert", "part006.frag")
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
		vao.Float32AttribData(vao.NextSlot(), 3, verticies, gl.STATIC_DRAW)
		vao.Float32AttribData(vao.NextSlot(), 3, colors, gl.STATIC_DRAW)
		vao.Float32AttribData(vao.NextSlot(), 2, texCoords, gl.STATIC_DRAW)
		vao.IndexData(indices, gl.STATIC_DRAW)

		// bind the GL numbered texture unit to our names
		p.UseProgram()
		p.Uniform1i("ourTexture", 0)
		p.Uniform1i("testTexture", 1)

		runtime.GC()
	}

	a.RunNoDt(func() {
		ren.GarbageCollect()

		showWidgetsDemo()
		sx, sy := a.GetSize()
		// * 2 because retina?  Imgui must be doing retina somewhere else?
		gl.Viewport(0, 0, sx*2, sy*2)
		p.UseProgram()

		//p.Uniform4f("vertexColor", color4[0], color4[1], color4[2], color4[3])
		projection := mgl32.Perspective(mgl32.DegToRad(45), float32(sx)/float32(sy), 0.1, 100)
		model := mgl32.HomogRotate3D(mgl32.DegToRad(rotation), mgl32.Vec3{1, 0, 0})
		view := mgl32.Translate3D(0, 0, -3)
		// the next 4 lines can be reordered to show how order of matrix multiplication
		// impacts the output.
		trans := mgl32.Ident4() // only needed to make the next 3 lines more uniform
		trans = trans.Mul4(mgl32.Translate3D(x, y, z))
		trans = trans.Mul4(mgl32.QuatRotate(mgl32.DegToRad(rotation), mgl32.Vec3{0, 0, 1}).Mat4())
		trans = trans.Mul4(mgl32.Scale3D(scale, scale, scale))
		p.Uniform1f("blend", blend)
		p.UniformMatrix4f("transform", trans)
		p.UniformMatrix4f("projection", projection)
		p.UniformMatrix4f("model", model)
		p.UniformMatrix4f("view", view)

		// binding the vao also binds the ebo
		wallTex.Bind(gl.TEXTURE0)
		testTex.Bind(gl.TEXTURE1)
		vao.Enable()
		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil)
	})
}

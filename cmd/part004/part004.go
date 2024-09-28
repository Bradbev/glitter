// Use the VertexAttribObject struct and draw two triangles with a single texture
// from https://learnopengl.com/Getting-started/Textures
package main

import (
	"embed"
	"log"
	"runtime"

	imgui "github.com/AllenDang/cimgui-go"

	"github.com/Bradbev/glitter/src/app"
	"github.com/Bradbev/glitter/src/ren"
	"github.com/go-gl/gl/v2.1/gl"
)

var (
	verticies = []float32{
		0.50, -0.5, 0.0, // bottom right
		-0.5, -0.5, 0.0, // bottom right
		0.00, 0.50, 0.0, // top
	}
	colors = []float32{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 0.0, 1.0,
	}
	texCoords = []float32{
		0.0, 0.0, // lower-left corner
		1.0, 0.0, // lower-right corner
		0.5, 1.0, // top-center corner
	}
	indices = []uint32{
		0, 1, 2, // triangle
	}
	showDemoWindow bool
	color4         [4]float32 = [4]float32{1, 0, 0, 0}
	blend          float32
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
	imgui.End()
}

//go:embed *.vert *.frag *.jpg
var assets embed.FS

func main() {
	a := app.Default()
	var (
		p   *ren.Program
		t   *ren.Texture
		vao *ren.VertexAttribObject
	)
	a.OnPostCreate = func() {
		version := gl.GoStr(gl.GetString(gl.VERSION))
		log.Println("OpenGL version", version)
		var err error
		p, err = ren.NewProgramFS(assets, "part004.vert", "part004.frag")
		if err != nil {
			log.Fatal(err)
		}

		t, err = ren.NewTextureFS(assets, "wall.jpg", gl.REPEAT, gl.REPEAT)
		if err != nil {
			log.Fatal(err)
		}

		vao = ren.NewVertexAttribObject()
		vao.Float32AttribData(vao.NextSlot(), 3, verticies, gl.STATIC_DRAW)
		vao.Float32AttribData(vao.NextSlot(), 3, colors, gl.STATIC_DRAW)
		vao.Float32AttribData(vao.NextSlot(), 2, texCoords, gl.STATIC_DRAW)
		vao.IndexData(indices, gl.STATIC_DRAW)

		runtime.GC()
	}

	a.Run(func() {
		ren.GarbageCollect()

		showWidgetsDemo()
		x, y := a.GetSize()
		// * 2 because retina?  Imgui must be doing retina somewhere else?
		gl.Viewport(0, 0, x*2, y*2)
		p.UseProgram()

		//p.Uniform4f("vertexColor", color4[0], color4[1], color4[2], color4[3])
		p.Uniform1f("blend", blend)

		// binding the vao also binds the ebo
		t.Bind(0)
		vao.Enable()
		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil)
	})
}

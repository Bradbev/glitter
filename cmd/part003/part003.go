// Use the VertexAttribObject struct and draw two triangles
package main

import (
	"embed"
	"log"
	"runtime"

	"github.com/AllenDang/cimgui-go/imgui"

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
	verticies2 = []float32{
		0.0, -1.0, 0.0, //  bottom right
		-0.5, -0.5, 0.0, // bottom right
		0.5, -0.5, 0.0, // top
	}
	colors = []float32{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 0.0, 1.0,
	}
	indices = []uint32{
		0, 1, 2, // triangle
	}
	showDemoWindow bool
	color4         [4]float32 = [4]float32{1, 0, 0, 0}
)

// makeVAO initializes and returns a vertex array from the points provided.
func showWidgetsDemo() {
	if showDemoWindow {
		imgui.ShowDemoWindowV(&showDemoWindow)
	}

	imgui.SetNextWindowSizeV(imgui.NewVec2(300, 300), imgui.CondOnce)
	imgui.Begin("Window 1")
	imgui.ColorEdit4("Color Edit", &color4)
	imgui.End()
}

//go:embed *.vert *.frag
var shaders embed.FS

func main() {
	a := app.Default()
	var (
		p    *ren.Program
		vao  *ren.VertexAttribObject
		vao2 *ren.VertexAttribObject
	)
	a.OnPostCreate = func() {
		version := gl.GoStr(gl.GetString(gl.VERSION))
		log.Println("OpenGL version", version)
		var err error
		p, err = ren.NewProgramFS(shaders, "part003.vert", "part003.frag")
		if err != nil {
			log.Fatal(err)
		}
		vao = ren.NewVertexAttribObject()
		vao.Float32AttribData(vao.NextSlot(), 3, verticies2, gl.STATIC_DRAW)
		vao.Float32AttribData(vao.NextSlot(), 3, colors, gl.STATIC_DRAW)
		vao.IndexData(indices, gl.STATIC_DRAW)

		vao2 = ren.NewVertexAttribObject()
		vao2.Float32AttribData(vao2.NextSlot(), 3, verticies, gl.STATIC_DRAW)
		vao2.Float32AttribData(vao2.NextSlot(), 3, colors, gl.STATIC_DRAW)
		vao2.IndexData(indices, gl.STATIC_DRAW)
		runtime.GC()
	}

	a.RunNoDt(func() {
		ren.GarbageCollect()

		showWidgetsDemo()
		x, y := a.GetSize()
		// * 2 because retina?  Imgui must be doing retina somewhere else?
		gl.Viewport(0, 0, x*2, y*2)
		p.UseProgram()

		//p.Uniform4f("vertexColor", color4[0], color4[1], color4[2], color4[3])

		// binding the vao also binds the ebo
		vao.Enable()
		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil)

		vao2.Enable()
		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil)
	})
}

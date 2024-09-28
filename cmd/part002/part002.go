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
	vertAndCol = []float32{
		// positions         // colors
		0.50, -0.5, 0.0, 1.0, 0.0, 1.0, // bottom right
		-0.5, -0.5, 0.0, 0.0, 1.0, 0.0, // bottom left
		0.00, 0.50, 0.0, 0.0, 0.0, 1.0, // top
	}
	verticies = []float32{
		0.5, -0.5, 0.0, // top right
		-0.5, -0.5, 0.0, // bottom right
		0.0, 0.5, 0.0, // bottom left
	}
	colors = []float32{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 0.0, 1.0,
	}
	indices = []uint32{
		0, 1, 2, // first triangle
	}
	showDemoWindow bool
	color4         [4]float32 = [4]float32{1, 0, 0, 0}
)

// makeVAO initializes and returns a vertex array from the points provided.
func makeVAO() (uint32, uint32, *ren.VertexAttribObject) {
	const sizeOfFloat = 4
	var vertexAttributeObject uint32
	if false {
		var vbo uint32
		// Make 1 new buffer object
		gl.GenBuffers(1, &vbo)
		// Bind the object to ARRAY_BUFFER, any operations on ARRAY_BUFFER will hit our object
		// ie, copy the data to the GPU
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, sizeOfFloat*len(vertAndCol), gl.Ptr(&vertAndCol[0]), gl.STATIC_DRAW)

		// VAO lets you bind different attribute setups for the VBO (?)
		gl.GenVertexArrays(1, &vertexAttributeObject)
		gl.BindVertexArray(vertexAttributeObject)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		// vertex attrib 0, size of attrib, type of attrib, isNormalized?, stride (0=packed), initial offset
		// Position attrib
		gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 6*sizeOfFloat, 0)
		gl.EnableVertexAttribArray(0)
		// color attrib
		gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 6*sizeOfFloat, 3*sizeOfFloat)
		gl.EnableVertexAttribArray(1)
	} else {
		// VAO lets you bind different attribute setups for the VBO (?)
		gl.GenVertexArrays(1, &vertexAttributeObject)
		gl.BindVertexArray(vertexAttributeObject)

		var points_vbo uint32
		gl.GenBuffers(1, &points_vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, points_vbo)
		gl.BufferData(gl.ARRAY_BUFFER, sizeOfFloat*len(verticies), gl.Ptr(&verticies[0]), gl.STATIC_DRAW)
		gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*sizeOfFloat, 0)
		gl.EnableVertexAttribArray(0)

		var colors_vbo uint32
		gl.GenBuffers(1, &colors_vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, colors_vbo)
		gl.BufferData(gl.ARRAY_BUFFER, sizeOfFloat*len(colors), gl.Ptr(&colors[0]), gl.STATIC_DRAW)
		gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 3*sizeOfFloat, 0)
		gl.EnableVertexAttribArray(1)
	}

	const bytesPerUint32 = 4
	var elementBufferObject uint32
	gl.GenBuffers(1, &elementBufferObject)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, elementBufferObject)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, bytesPerUint32*len(indices), gl.Ptr(&indices[0]), gl.STATIC_DRAW)

	return vertexAttributeObject, elementBufferObject, nil
}

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
		p         *ren.Program
		vaoHandle uint32
		vao       *ren.VertexAttribObject
		ebo       uint32
	)
	a.OnPostCreate = func() {
		version := gl.GoStr(gl.GetString(gl.VERSION))
		log.Println("OpenGL version", version)
		var err error
		p, err = ren.NewProgramFS(shaders, "part002.vert", "part002.frag")
		if err != nil {
			log.Fatal(err)
		}
		vaoHandle, ebo, vao = makeVAO()
		_ = vao
		/*
		 */
		runtime.GC()
	}
	_ = vaoHandle

	a.Run(func() {
		ren.GarbageCollect()

		showWidgetsDemo()
		x, y := a.GetSize()
		// * 2 because retina?  Imgui must be doing retina somewhere else?
		gl.Viewport(0, 0, x*2, y*2)
		p.UseProgram()

		//p.Uniform4f("vertexColor", color4[0], color4[1], color4[2], color4[3])

		// binding the vao also binds the ebo
		gl.BindVertexArray(vaoHandle)
		//vao.Enable()
		_ = ebo
		//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(verticies)/3))
		//gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
		// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE) // wireframe
		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil)

	})
}

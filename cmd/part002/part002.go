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
		0.5, 0.5, 0.0, // top right
		0.5, -0.5, 0.0, // bottom right
		-0.5, -0.5, 0.0, // bottom left
	}
	colors = []float32{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 0.0, 1.0,
	}
	indices = []uint32{
		0, 1, 3, // first triangle
	}
	showDemoWindow bool
	color4         [4]float32 = [4]float32{1, 0, 0, 0}
)

// makeGlObjects initializes and returns a vertex array from the points provided.
func makeGlObjects(points []float32) (uint32, uint32) {
	const bytesPerFloat = 4
	var vertexBufferObject uint32
	// Make 1 new buffer object
	gl.GenBuffers(1, &vertexBufferObject)
	// Bind the object to ARRAY_BUFFER, any operations on ARRAY_BUFFER will hit our object
	// ie, copy the data to the GPU
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.BufferData(gl.ARRAY_BUFFER, bytesPerFloat*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	// VAO lets you bind different attribute setups for the VBO (?)
	var vertexAttributeObject uint32
	gl.GenVertexArrays(1, &vertexAttributeObject)
	gl.BindVertexArray(vertexAttributeObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// vertex attrib 0, size of attrib, type of attrib, isNormalized?, stride (0=packed), initial offset
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	const bytesPerUint32 = 4
	var elementBufferObject uint32
	gl.GenBuffers(1, &elementBufferObject)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, elementBufferObject)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, bytesPerUint32*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	return vertexAttributeObject, elementBufferObject
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
		p   *ren.Program
		vao uint32
		ebo uint32
	)
	a.OnPostCreate = func() {
		version := gl.GoStr(gl.GetString(gl.VERSION))
		log.Println("OpenGL version", version)
		v, err := ren.NewVertexShader(shaders, "part002.vert")
		if err != nil {
			log.Fatal(err)
		}
		f, err := ren.NewFragmentShader(shaders, "part002.frag")
		if err != nil {
			log.Fatal(err)
		}
		p = ren.NewProgram(v, f)
		vao, ebo = makeGlObjects(verticies)
		runtime.GC()
	}

	a.Run(func() {
		ren.GarbageCollect()

		showWidgetsDemo()
		x, y := a.GetSize()
		// * 2 because retina?  Imgui must be doing retina somewhere else?
		gl.Viewport(0, 0, x*2, y*2)
		p.UseProgram()

		p.Uniform4f("vertexColor", color4[0], color4[1], color4[2], color4[3])

		// binding the vao also binds the ebo
		gl.BindVertexArray(vao)
		_ = ebo
		//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))
		//gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
		// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE) // wireframe
		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil)

	})
}

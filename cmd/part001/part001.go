package main

import (
	"embed"
	"log"

	imgui "github.com/AllenDang/cimgui-go"

	"github.com/Bradbev/glitter/src/app"
	"github.com/Bradbev/glitter/src/ren"
	"github.com/go-gl/gl/v2.1/gl"
)

var (
	triangle = []float32{
		-0.5, -0.5, -1.0,
		0.5, -0.5, -1.0,
		0.0, 0.5, -1.0,
	}
)

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	const bytesPerFloat = 4
	var vertexBufferObject uint32
	// Make 1 new buffer object
	gl.GenBuffers(1, &vertexBufferObject)
	// Bind the object to ARRAY_BUFFER, any operations on ARRAY_BUFFER will hit our object
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.BufferData(gl.ARRAY_BUFFER, bytesPerFloat*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vertexAttributeObject uint32
	gl.GenVertexArrays(1, &vertexAttributeObject)
	gl.BindVertexArray(vertexAttributeObject)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vertexAttributeObject
}

//go:embed *.vert *.frag
var shaders embed.FS

func main() {
	a := app.Default()
	var (
		prog uint32
		vao  uint32
	)
	a.OnPostCreate = func() {
		version := gl.GoStr(gl.GetString(gl.VERSION))
		log.Println("OpenGL version", version)
		vertexShader, err := ren.CompileShaderFS(shaders, "part001.vert", gl.VERTEX_SHADER)
		if err != nil {
			log.Fatal(err)
		}
		fragmentShader, err := ren.CompileShaderFS(shaders, "part001.frag", gl.FRAGMENT_SHADER)
		if err != nil {
			panic(err)
		}
		prog = gl.CreateProgram()
		gl.AttachShader(prog, vertexShader)
		gl.AttachShader(prog, fragmentShader)
		gl.LinkProgram(prog)
		vao = makeVao(triangle)
	}

	a.Run(func() {
		d := false
		imgui.ShowDemoWindowV(&d)
		x, y := a.GetSize()
		// * 2 because retina?  Imgui must be doing retina somewhere else?
		gl.Viewport(0, 0, x*2, y*2)
		gl.UseProgram(prog)
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))

	})
}

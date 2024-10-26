package main

import (
	"embed"
	"log"

	"github.com/AllenDang/cimgui-go/imgui"

	"github.com/Bradbev/glitter/src/app"
	"github.com/Bradbev/glitter/src/ren"
	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	verticies = []float32{
		0.5, 0.5, 0.0, // top right
		0.5, -0.5, 0.0, // bottom right
		-0.5, -0.5, 0.0, // bottom left
		-0.5, 0.5, 0.0, // top left
	}
	indices = []uint32{
		0, 1, 3, // first triangle
		1, 2, 3, // second triangle
	}
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

//go:embed *.vert *.frag
var shaders embed.FS

func main() {
	a := app.Default()
	var (
		prog uint32
		vao  uint32
		ebo  uint32
	)
	a.OnPostCreate = func() {
		err := gl.Init()
		if err != nil {
			log.Fatal(err)
		}
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
		gl.DeleteShader(vertexShader)
		gl.DeleteShader(fragmentShader)
		vao, ebo = makeGlObjects(verticies)
	}

	a.RunNoDt(func() {
		d := false
		imgui.ShowDemoWindowV(&d)
		x, y := a.GetSize()
		// * 2 because retina?  Imgui must be doing retina somewhere else?
		gl.Viewport(0, 0, x*2, y*2)
		gl.UseProgram(prog)
		// binding the vao also binds the ebo
		gl.BindVertexArray(vao)
		_ = ebo
		//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))
		//gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
		// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE) // wireframe
		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil)

	})
}

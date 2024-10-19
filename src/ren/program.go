package ren

import (
	"embed"
	"fmt"
	"io/fs"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
)

//go:embed shaders
var embeddedShaders embed.FS

type Program struct {
	Handle   uint32
	uniforms map[string]int32
}

type Shader struct {
	Handle uint32
}

func makeShader(handle uint32, e error) (*Shader, error) {
	if e != nil {
		return nil, e
	}
	s := &Shader{handle}
	runtime.SetFinalizer(s, func(s *Shader) {
		handle := s.Handle
		onMainThread(func() {
			println("Deleting shader", handle)
			gl.DeleteShader(handle)
		})
	})

	return s, nil
}

func NewVertexShader(fsys fs.FS, filename string) (*Shader, error) {
	return makeShader(CompileShaderFS(fsys, filename, gl.VERTEX_SHADER))
}

func NewFragmentShader(fsys fs.FS, filename string) (*Shader, error) {
	return makeShader(CompileShaderFS(fsys, filename, gl.FRAGMENT_SHADER))
}

func (p *Program) UseProgram() {
	gl.UseProgram(p.Handle)
}

func (p *Program) GetUniformLocation(name string) int32 {
	if loc, ok := p.uniforms[name]; ok {
		return loc
	}

	loc := gl.GetUniformLocation(p.Handle, gl.Str(name+"\x00"))
	p.uniforms[name] = loc
	return loc
}

func (p *Program) Uniform1i(name string, v1 int32) {
	gl.Uniform1i(p.GetUniformLocation(name), v1)
}

func (p *Program) Uniform1f(name string, f1 float32) {
	gl.Uniform1f(p.GetUniformLocation(name), f1)
}

func (p *Program) Uniform3f(name string, f1, f2, f3 float32) {
	gl.Uniform3f(p.GetUniformLocation(name), f1, f2, f3)
}

func (p *Program) UniformVec3(name string, v mgl32.Vec3) {
	gl.Uniform3f(p.GetUniformLocation(name), v[0], v[1], v[2])
}

func (p *Program) Uniform4f(name string, f1, f2, f3, f4 float32) {
	gl.Uniform4f(p.GetUniformLocation(name), f1, f2, f3, f4)
}

func (p *Program) UniformMatrix4f(name string, mat mgl32.Mat4) {
	gl.UniformMatrix4fv(p.GetUniformLocation(name), 1, false, &mat[0])
}

func NewProgram(shaders ...*Shader) (*Program, error) {
	handle := gl.CreateProgram()
	p := &Program{
		Handle:   handle,
		uniforms: map[string]int32{},
	}
	runtime.SetFinalizer(p, func(p *Program) {
		handle := p.Handle
		onMainThread(func() {
			gl.DeleteProgram(handle)
		})
	})
	for _, s := range shaders {
		gl.AttachShader(p.Handle, s.Handle)
	}
	gl.LinkProgram(handle)
	var status int32
	gl.GetProgramiv(handle, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(handle, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(handle, logLength, nil, gl.Str(log))

		return nil, fmt.Errorf("failed to link program %v", log)
	}

	return p, nil
}

func NewProgramFS(fsys fs.FS, vertex, fragment string) (*Program, error) {
	v, err := NewVertexShader(fsys, vertex)
	if err != nil {
		return nil, err
	}
	f, err := NewFragmentShader(fsys, fragment)
	if err != nil {
		return nil, err
	}

	return NewProgram(v, f)
}

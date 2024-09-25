package ren

import (
	"io/fs"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
)

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

func (p *Program) Uniform4f(name string, f1, f2, f3, f4 float32) {
	gl.Uniform4f(p.GetUniformLocation(name), f1, f2, f3, f4)
}

func NewProgram(shaders ...*Shader) *Program {
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

	return p
}

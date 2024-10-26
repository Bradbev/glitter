package ren

import (
	"embed"
	"fmt"
	"io/fs"
	"reflect"
	"runtime"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

//go:embed shaders
var embeddedShaders embed.FS

type uniformFunc func(name string, structPtr unsafe.Pointer)

type Program struct {
	// handle is the gl handle returned from CreateProgram
	handle uint32

	// uniformLocations is a cache over gl.GetUniformLocation
	uniformLocations map[string]int32

	// structCache maps from a struct to a function that can be called
	// to load the members of the struct into uniforms
	structCache map[reflect.Type]uniformFunc
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
			//println("Deleting shader", handle)
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
	gl.UseProgram(p.handle)
}

func (p *Program) GetUniformLocation(name string) int32 {
	if loc, ok := p.uniformLocations[name]; ok {
		return loc
	}

	loc := gl.GetUniformLocation(p.handle, gl.Str(name+"\x00"))
	p.uniformLocations[name] = loc
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

// UniformStruct iterates the struct toLoad and loads uniforms from that
// struct into name by joining the field names for gl, ie name.Field1
// Fields in the struct need to be public
// A cache is used so that the struct only needs to be reflected once.
func (p *Program) UniformStruct(name string, toLoad any) {
	f, ok := p.structCache[reflect.TypeOf(toLoad).Elem()]
	if !ok {
		f = p.buildCacheEntry(toLoad)
	}
	f(name, reflect.ValueOf(toLoad).UnsafePointer())
}

// the private load* functions build a uniformFunc for the given type
func loadVec3(p *Program, fieldName string, offset uintptr) uniformFunc {
	return func(name string, structPtr unsafe.Pointer) {
		base := unsafe.Add(structPtr, offset)
		v := (*mgl32.Vec3)(base)
		p.UniformVec3(name+"."+fieldName, *v)
	}
}

func loadInt32(p *Program, fieldName string, offset uintptr) uniformFunc {
	return func(name string, structPtr unsafe.Pointer) {
		base := unsafe.Add(structPtr, offset)
		i := (*int32)(base)
		p.Uniform1i(name+"."+fieldName, *i)
	}
}

func loadFloat32(p *Program, fieldName string, offset uintptr) uniformFunc {
	return func(name string, structPtr unsafe.Pointer) {
		base := unsafe.Add(structPtr, offset)
		i := (*float32)(base)
		p.Uniform1f(name+"."+fieldName, *i)
	}
}

// buildCacheEntry reflects over toLoad and builds a function
// that allows that struct type to be loaded to a shader program.
func (p *Program) buildCacheEntry(toLoad any) uniformFunc {
	ptrTyp := reflect.TypeOf(toLoad)
	if ptrTyp.Kind() != reflect.Pointer {
		panic("must pass in a pointer")
	}
	typ := ptrTyp.Elem()
	if typ.Kind() != reflect.Struct {
		panic("pointer must be to a struct")
	}

	typeOf := reflect.TypeOf
	funcs := map[reflect.Type]func(p *Program, name string, offset uintptr) uniformFunc{
		typeOf(mgl32.Vec3{}): loadVec3,
		typeOf(float32(0)):   loadFloat32,
		typeOf(int32(0)):     loadInt32,
	}

	toCall := []uniformFunc{}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		name := field.Name
		f, ok := funcs[field.Type]
		if !ok {
			panic("unknown type " + field.Type.String())
		}
		toCall = append(toCall, f(p, name, field.Offset))
	}
	ret := func(name string, structPtr unsafe.Pointer) {
		for _, f := range toCall {
			f(name, structPtr)
		}
	}
	p.structCache[typ] = ret
	return ret
}

func NewProgram(shaders ...*Shader) (*Program, error) {
	handle := gl.CreateProgram()
	p := &Program{
		handle:           handle,
		uniformLocations: map[string]int32{},
		structCache:      map[reflect.Type]uniformFunc{},
	}
	runtime.SetFinalizer(p, func(p *Program) {
		handle := p.handle
		onMainThread(func() {
			gl.DeleteProgram(handle)
		})
	})
	for _, s := range shaders {
		gl.AttachShader(p.handle, s.Handle)
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

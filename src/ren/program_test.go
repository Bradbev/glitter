package ren

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/go-gl/mathgl/mgl32"
)

type testStruct struct {
	Float  float32
	Float2 float32
	Vec2   mgl32.Vec2
	Vec3   mgl32.Vec3
	Mat4   mgl32.Mat4
}

type loadFunc func(p *Program, structPtr unsafe.Pointer)
type uniformStructInterp struct {
	cache map[reflect.Type]loadFunc
}

func uFloat32(location string, offset uintptr) func(structBase unsafe.Pointer) {
	return func(structBase unsafe.Pointer) {
		base := unsafe.Add(structBase, offset)
		f := (*float32)(base)
		fmt.Println(base, location, *f)
		//gl.Uniform1f(location, f)
	}
}

func uVec3(location string, offset uintptr) func(structBase unsafe.Pointer) {
	return func(structBase unsafe.Pointer) {
		base := unsafe.Add(structBase, offset)
		f := (*mgl32.Vec3)(base)
		fmt.Println(base, location, f[0], f[1], f[2])
		//gl.Uniform1f(location, f)
	}
}

func listStruct(uniforms any) {
	ptrTyp := reflect.TypeOf(uniforms)
	if ptrTyp.Kind() != reflect.Pointer {
		panic("must pass in a pointer")
	}
	typ := ptrTyp.Elem()
	if typ.Kind() != reflect.Struct {
		panic("pointer must be to a struct")
	}
	ptrToU := reflect.ValueOf(uniforms).UnsafePointer()
	funcs := []func(structBase unsafe.Pointer){}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		name := field.Name
		ftyp := field.Type.String()
		off := field.Offset
		p := unsafe.Add(ptrToU, off)
		fmt.Println(name, ftyp, off, p)
		//location := p.GetUniformLocation(name)
		if field.Type == reflect.TypeOf(float32(0)) {
			funcs = append(funcs, uFloat32(name, off))
		}
		if field.Type == reflect.TypeOf(mgl32.Vec3{}) {
			funcs = append(funcs, uVec3(name, off))
		}

	}
	for _, f := range funcs {
		f(ptrToU)
	}

}

func TestProgramStruct(t *testing.T) {
	ts := testStruct{
		Float:  1,
		Float2: 2,
		Vec3:   mgl32.Vec3{11, 12, 13},
	}
	listStruct(&ts)
	t.Fail()
}

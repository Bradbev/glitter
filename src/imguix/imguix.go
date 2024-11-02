package imguix

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/go-gl/mathgl/mgl32"
)

func TreeNode(name string, body func()) bool {
	if imgui.TreeNodeExStrV(name, imgui.TreeNodeFlagsDefaultOpen) {
		body()
		imgui.TreePop()
		return true
	}
	return false
}

// EditStruct creates an imgui editor for structs.
// The `ui` tag is used is specify tags per field.
// Numeric types can be
// - slider:{min,max}
// - drag:{min,max,speed}
// - input:{step, stepfast} (default)
func EditStruct(name string, toEdit any) {
	TreeNode(name, func() {
		typ := reflect.TypeOf(toEdit).Elem()
		val := reflect.ValueOf(toEdit)

		fn, ok := editFnCache[typ]
		if !ok {
			fn = buildEditFn(typ)
			editFnCache[typ] = fn
		}
		if fn != nil {
			fn(val.UnsafePointer())
		} else {
			imgui.Text(fmt.Sprintf("%v", typ))
		}
	})
}

// private below here
func init() {
	addBuilder[int32](builderInt32, HandlesPointers)
	addBuilder[float32](builderFloat32, HandlesPointers)
	addBuilder[mgl32.Vec3](builderVec3, HandlesPointers)
	addBuilder[string](builderString, HandlesPointers)
}

type editFn func(structPtr unsafe.Pointer)

var editFnCache = map[reflect.Type]editFn{}
var buildFnForType = map[reflect.Type]func(tag string, field reflect.StructField) editFn{}

type builderOpts int

const (
	HandlesPointers builderOpts = iota
)

func addBuilder[T any](builder func(tag string, field reflect.StructField) editFn, opts ...builderOpts) {
	var fake T
	buildFnForType[reflect.TypeOf(fake)] = builder
	for _, opt := range opts {
		switch opt {
		case HandlesPointers:
			buildFnForType[reflect.TypeOf(&fake)] = builder
		}
	}
}

func buildEditFn(typ reflect.Type) editFn {
	fns := []editFn{}
	for i := 0; i < typ.NumField(); i++ {
		fns = append(fns, buildFieldEditFn(typ.Field(i)))
	}
	return func(ptr unsafe.Pointer) {
		for _, fn := range fns {
			if fn != nil {
				fn(ptr)
			}
		}
	}
}

func buildFieldEditFn(field reflect.StructField) editFn {
	if buildFn, ok := buildFnForType[field.Type]; ok {
		return buildFn(field.Tag.Get("ui"), field)
	}
	return func(structPtr unsafe.Pointer) {
		imgui.Text(fmt.Sprintf("Unable to edit %q, unknown type", field.Name))
	}
}

func builderString(tag string, field reflect.StructField) editFn {
	details := numericTagDetailDefault
	kind := MustParseInto(tag, &details)

	target := func(structPtr unsafe.Pointer) *string {
		if field.Type.Kind() == reflect.Pointer {
			return *(**string)(unsafe.Add(structPtr, field.Offset))
		}
		return (*string)(unsafe.Add(structPtr, field.Offset))
	}
	switch kind {
	default:
		return func(structPtr unsafe.Pointer) {
			imgui.InputTextWithHint(field.Name, "", target(structPtr), 0, nil)
		}
	}
}

// vec3TagDetail is populated from struct tags that can be edited.
type vec3TagDetail struct {
	Min      float32
	Max      float32
	Speed    float32
	Step     float32
	StepFast float32
}

var vec3TagDetailDefault = vec3TagDetail{
	Speed: 1,
}

func builderVec3(tag string, field reflect.StructField) editFn {
	details := numericTagDetailDefault
	kind := MustParseInto(tag, &details)

	target := func(structPtr unsafe.Pointer) *[3]float32 {
		if field.Type.Kind() == reflect.Pointer {
			return *(**[3]float32)(unsafe.Add(structPtr, field.Offset))
		}
		return (*[3]float32)(unsafe.Add(structPtr, field.Offset))
	}

	switch kind {
	case "color":
		return func(structPtr unsafe.Pointer) {
			imgui.ColorEdit3(field.Name, target(structPtr))
		}
	case "slider":
		return func(structPtr unsafe.Pointer) {
			imgui.SliderFloat3(field.Name, target(structPtr), details.Min, details.Max)
		}
	case "drag":
		return func(structPtr unsafe.Pointer) {
			imgui.DragFloat3V(field.Name, target(structPtr), details.Speed, details.Min, details.Max, "%.3f", 0)
		}
	case "input":
	default:
		return func(structPtr unsafe.Pointer) {
			imgui.InputFloat3V(field.Name, target(structPtr), "%.3f", 0)
		}
	}
	return nil
}

// numericTagDetail is populated from struct tags that can be edited.
type numericTagDetail struct {
	Min      float32
	Max      float32
	Speed    float32
	Step     float32
	StepFast float32
}

var numericTagDetailDefault = numericTagDetail{
	Speed: 1,
}

func builderFloat32(tag string, field reflect.StructField) editFn {
	details := numericTagDetailDefault
	kind := MustParseInto(tag, &details)

	target := func(structPtr unsafe.Pointer) *float32 {
		if field.Type.Kind() == reflect.Pointer {
			return *(**float32)(unsafe.Add(structPtr, field.Offset))
		}
		return (*float32)(unsafe.Add(structPtr, field.Offset))
	}

	switch kind {
	case "slider":
		return func(structPtr unsafe.Pointer) {
			imgui.SliderFloat(field.Name, target(structPtr), details.Min, details.Max)
		}
	case "drag":
		return func(structPtr unsafe.Pointer) {
			imgui.DragFloatV(field.Name, target(structPtr), details.Speed, details.Min, details.Max, "%.3f", 0)
		}
	case "input":
	default:
		return func(structPtr unsafe.Pointer) {
			imgui.InputFloatV(field.Name, target(structPtr), details.Step, details.StepFast, "%.3f", 0)
		}
	}
	return nil
}

func builderInt32(tag string, field reflect.StructField) editFn {
	details := numericTagDetailDefault
	kind := MustParseInto(tag, &details)

	target := func(structPtr unsafe.Pointer) *int32 {
		if field.Type.Kind() == reflect.Pointer {
			return *(**int32)(unsafe.Add(structPtr, field.Offset))
		}
		return (*int32)(unsafe.Add(structPtr, field.Offset))
	}

	switch kind {
	case "slider":
		return func(structPtr unsafe.Pointer) {
			imgui.SliderInt(field.Name, target(structPtr), int32(details.Min), int32(details.Max))
		}
	case "drag":
		return func(structPtr unsafe.Pointer) {
			imgui.DragIntV(field.Name, target(structPtr), details.Speed, int32(details.Min), int32(details.Max), "%d", 0)
		}
	case "input":
	default:
		return func(structPtr unsafe.Pointer) {
			imgui.InputIntV(field.Name, target(structPtr), int32(details.Step), int32(details.StepFast), 0)
		}
	}
	return nil
}

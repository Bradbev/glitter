package ren

import "github.com/go-gl/mathgl/mgl32"

func Vec3ToRaw(v *mgl32.Vec3) *[3]float32 {
	return (*[3]float32)(v)
}

package asset

import (
	"testing"
	"unsafe"

	"github.com/bloeys/gglm/gglm"
	"github.com/stretchr/testify/assert"
)

func TestSliceCast(t *testing.T) {
	d := [3]float32{1, 2, 3}
	in := []gglm.Vec3{{Data: d}}
	out := castGglmVec3ToFloat32(in)

	assert.Equal(t, unsafe.Pointer(&in[0].Data[0]), unsafe.Pointer(&out[0]))
	assert.Equal(t, 3, len(out))
}

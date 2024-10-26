package ren

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type DebugNormsVisualizer struct {
	program   *Program
	vao       *VertexAttribObject
	lineCount int32
}

func (d *DebugNormsVisualizer) Draw(c *Camera, model mgl32.Mat4) {
	d.DrawRange(c, model, 0, d.lineCount)
}

func (d *DebugNormsVisualizer) DrawRange(c *Camera, model mgl32.Mat4, start, count int32) {
	p := d.program
	view, projection := c.GetMatrices()
	p.UseProgram()
	p.UniformMatrix4f("view", view)
	p.UniformMatrix4f("projection", projection)
	p.UniformMatrix4f("model", model)
	d.vao.Enable()
	toDraw := min(count, d.lineCount)
	gl.DrawArrays(gl.LINES, start, toDraw)
}

func GenerateNormalsVisualizer(verts []float32, norms []float32, scale float32) (*DebugNormsVisualizer, error) {
	if len(verts) != len(norms) {
		return nil, fmt.Errorf("len of verts and norms must be the same %d:%d", len(verts), len(norms))
	}
	normProgram, err := NewProgramFS(embeddedShaders, "shaders/debug.vert", "shaders/debug.frag")
	if err != nil {
		return nil, err
	}

	var lines []float32
	var colors []float32
	surfaceC := mgl32.Vec3{1, 1, 1}
	outC := surfaceC.Mul(0.3)

	for i := 0; i < len(verts)/3; i++ {
		v := mgl32.Vec3{verts[i*3+0], verts[i*3+1], verts[i*3+2]}
		n := v.Add(mgl32.Vec3{norms[i*3+0], norms[i*3+1], norms[i*3+2]})

		lines = append(lines, v[:]...)
		colors = append(colors, surfaceC[:]...)
		lines = append(lines, n[:]...)
		colors = append(colors, outC[:]...)
	}

	vao := NewVertexAttribObject()
	vao.Float32AttribData(vao.NextSlot(), 3, lines, gl.STATIC_DRAW)
	vao.Float32AttribData(vao.NextSlot(), 3, colors, gl.STATIC_DRAW)

	return &DebugNormsVisualizer{
		program:   normProgram,
		vao:       vao,
		lineCount: int32(len(verts)/3) * 2,
	}, nil

}

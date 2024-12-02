package ren

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Mesh struct {
	Points   []float32
	Normals  []float32
	Uvs      [][]float32
	Indicies []uint32
	Textures []TextureAndType

	vao *VertexAttribObject
}

type Scene struct {
	Meshes []*Mesh
}

type TextureAndType struct {
	Texture *Texture
	Type    TextureType
}

type TextureType int

const (
	TexDiffuse TextureType = iota
	TexSpecular
	TexNormal
)

func (t TextureType) String() string {
	switch t {
	case TexDiffuse:
		return "texture_diffuse"
	case TexSpecular:
		return "texture_specular"
	case TexNormal:
		return "texture_normal"
	}
	return ""
}

func (m *Mesh) AddUvs(uvs []float32) {
	m.Uvs = append(m.Uvs, uvs)
}

func (m *Mesh) Setup() {
	vao := NewVertexAttribObject()
	vao.Float32AttribData(vao.NextSlot(), 3, m.Points, gl.STATIC_DRAW)
	vao.Float32AttribData(vao.NextSlot(), 3, m.Normals, gl.STATIC_DRAW)
	vao.Float32AttribData(vao.NextSlot(), 2, m.Uvs[0], gl.STATIC_DRAW)
	vao.IndexData(m.Indicies, gl.STATIC_DRAW)
	m.vao = vao
}

func (m *Mesh) Draw(p *Program) {
	p.UseProgram()
	diffuse := 1
	spec := 1
	normal := 1
	for i, t := range m.Textures {
		switch t.Type {
		case TexDiffuse:
			p.Uniform1i(fmt.Sprintf("texture_diffuse%d", diffuse), int32(i))
			diffuse++
		case TexSpecular:
			p.Uniform1i(fmt.Sprintf("texture_specular%d", spec), int32(i))
			spec++
		case TexNormal:
			p.Uniform1i(fmt.Sprintf("texture_normal%d", spec), int32(i))
			normal++
		}
		t.Texture.Bind(uint32(gl.TEXTURE0 + i))
	}
	m.vao.Enable()
	gl.DrawElements(gl.TRIANGLES, int32(len(m.Indicies)), gl.UNSIGNED_INT, nil)
}

func (s *Scene) Setup() {
	for _, m := range s.Meshes {
		m.Setup()
	}
}

func (s *Scene) Draw(p *Program) {
	for _, m := range s.Meshes {
		m.Draw(p)
	}
}

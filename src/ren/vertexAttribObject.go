package ren

import "github.com/go-gl/gl/v4.1-core/gl"

type VertexAttribObject struct {
	Handle         uint32
	slots          []uint32
	indiciesHandle uint32
}

func NewVertexAttribObject() *VertexAttribObject {
	vao := &VertexAttribObject{}
	gl.GenVertexArrays(1, &vao.Handle)
	gl.BindVertexArray(vao.Handle)
	return vao
}

func (v *VertexAttribObject) NextSlot() int {
	return len(v.slots)
}

func (v *VertexAttribObject) Enable() {
	gl.BindVertexArray(v.Handle)
}

func (v *VertexAttribObject) Float32AttribData(slot int, elementSize int, data []float32, usage uint32) {
	if slot > len(v.slots) {
		panic("Can only add to the next slot")
	}
	gl.BindVertexArray(v.Handle)

	var slotHandle uint32
	if slot == len(v.slots) {
		v.slots = append(v.slots, 0)
		gl.GenBuffers(1, &slotHandle)
		v.slots[slot] = slotHandle
	} else {
		slotHandle = v.slots[slot]
	}

	const sizeOfFloat = 4
	gl.BindBuffer(gl.ARRAY_BUFFER, slotHandle)
	gl.BufferData(gl.ARRAY_BUFFER, sizeOfFloat*len(data), gl.Ptr(data), usage)
	gl.VertexAttribPointerWithOffset(uint32(slot), int32(elementSize), gl.FLOAT, false, int32(elementSize)*sizeOfFloat, 0)
	gl.EnableVertexAttribArray(uint32(slot))
}

func (v *VertexAttribObject) IndexData(indicies []uint32, usage uint32) {
	const sizeOfUint32 = 4
	gl.GenBuffers(1, &v.indiciesHandle)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, v.indiciesHandle)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, sizeOfUint32*len(indicies), gl.Ptr(indicies), usage)
}

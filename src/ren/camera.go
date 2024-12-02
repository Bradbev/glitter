package ren

import "github.com/go-gl/mathgl/mgl32"

type Camera struct {
	Position mgl32.Vec3
	Forward  mgl32.Vec3
	Up       mgl32.Vec3
	Fov      float32
	NearClip float32
	FarClip  float32

	view       mgl32.Mat4
	projection mgl32.Mat4
}

func NewCamera() *Camera {
	return &Camera{
		Forward:  mgl32.Vec3{0, -1, 0},
		Up:       mgl32.Vec3{0, 1, 0},
		Fov:      45,
		NearClip: 0.1,
		FarClip:  100,
	}
}

func (c *Camera) LookAt(location mgl32.Vec3) {
	c.Forward = location.Sub(c.Position).Normalize()
}

func (c *Camera) CacheMatricies(screenX, screenY float32) {
	c.view = c.GetViewMat()
	c.projection = c.GetProjectionMat(screenX, screenY)
}

func (c *Camera) GetMatrices() (view mgl32.Mat4, projection mgl32.Mat4) {
	return c.view, c.projection
}

func (c *Camera) GetViewMat() mgl32.Mat4 {
	return mgl32.LookAtV(
		c.Position,
		c.Position.Add(c.Forward),
		c.Up)
}

func (c *Camera) GetProjectionMat(screenX, screenY float32) mgl32.Mat4 {
	return mgl32.Perspective(mgl32.DegToRad(c.Fov), screenX/screenY, c.NearClip, c.FarClip)
}

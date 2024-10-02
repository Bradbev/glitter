package ren

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"

	"github.com/anthonynsimon/bild/transform"
	"github.com/go-gl/gl/v2.1/gl"
)

type Texture struct {
	handle uint32
}

func NewTextureFS(fsys fs.FS, filename string, wrapR, wrapS int32) (*Texture, error) {
	imgFile, err := fsys.Open(filename)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}
	return NewTexture(img, wrapR, wrapS)
}

func NewTexture(img image.Image, wrapR, wrapS int32) (*Texture, error) {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("Stride does not match")
	}
	rgba = transform.FlipV(rgba)

	tex := &Texture{}
	gl.GenTextures(1, &tex.handle)
	tex.Bind(gl.TEXTURE0)
	defer tex.UnBind()

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, wrapR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, wrapS)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR) // minification filter
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR) // magnification filter
	sz := rgba.Rect.Size()
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.SRGB_ALPHA, int32(sz.X), int32(sz.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

	gl.GenerateMipmap(tex.handle)

	return tex, nil
}

func (t *Texture) Bind(textureUnit uint32) {
	gl.ActiveTexture(textureUnit)
	gl.BindTexture(gl.TEXTURE_2D, t.handle)
}

func (t *Texture) UnBind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

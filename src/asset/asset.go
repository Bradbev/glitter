package asset

import (
	"bytes"
	"fmt"
	"image/png"
	"os"
	"path/filepath"

	"github.com/Bradbev/glitter/src/ren"
	"github.com/bloeys/assimp-go/asig/asig"
	"github.com/bloeys/gglm/gglm"
	"github.com/go-gl/gl/v4.6-core/gl"
)

func castGglmVec3ToFloat32(in []gglm.Vec3) []float32 {
	result := make([]float32, len(in)*3)
	for i, v := range in {
		result[i*3+0] = v.X()
		result[i*3+1] = v.Z()
		result[i*3+2] = v.Y()
	}
	return result

	//data := (*float32)(unsafe.Pointer(&in[0].Data[0]))
	//return unsafe.Slice(data, 3*len(in))
}

func copyGglmVec3ToUv2D(in []gglm.Vec3) []float32 {
	result := make([]float32, len(in)*2)
	for i, v := range in {
		result[i*2] = v.X()
		result[i*2+1] = v.Y()
	}
	return result
}

func ImportFile(file string, postProcessFlags asig.PostProcess) (*ren.Scene, error) {
	scene, release, err := asig.ImportFile(file, postProcessFlags)
	defer release()
	if err != nil {
		return nil, err
	}

	fs := os.DirFS(filepath.Dir(file))
	result := &ren.Scene{}

	for _, m := range scene.Meshes {

		// triangles
		indicies := []uint32{}
		for _, face := range m.Faces {
			for _, index := range face.Indices {
				indicies = append(indicies, uint32(index))
			}
		}

		mesh := &ren.Mesh{
			Points:   castGglmVec3ToFloat32(m.Vertices),
			Normals:  castGglmVec3ToFloat32(m.Normals),
			Indicies: indicies,
		}
		for _, texCoords := range m.TexCoords {
			mesh.AddUvs(copyGglmVec3ToUv2D(texCoords))
		}

		for _, mat := range scene.Materials {
			for i := 0; i < asig.GetMaterialTextureCount(mat, asig.TextureTypeDiffuse); i++ {
				info, err := asig.GetMaterialTexture(mat, asig.TextureTypeDiffuse, uint(i))
				if err != nil {
					continue
				}
				tex, err := ren.NewTextureFS(fs, info.Path, gl.REPEAT, gl.REPEAT)
				if err != nil {
					continue
				}
				mesh.Textures = append(mesh.Textures, ren.TextureAndType{Texture: tex, Type: ren.Diffuse})
			}
		}
		result.Meshes = append(result.Meshes, mesh)
	}

	return result, nil
}

func Test() {
	scene, release, err := asig.ImportFile("obj.obj", asig.PostProcessTriangulate|asig.PostProcessJoinIdenticalVertices)
	if err != nil {
		panic(err)
	}
	defer release()

	fmt.Printf("RootNode: %+v\n\n", scene.RootNode)

	for i := 0; i < len(scene.Meshes); i++ {

		println("Mesh:", i, "; Verts:", len(scene.Meshes[i].Vertices), "; Normals:", len(scene.Meshes[i].Normals), "; MatIndex:", scene.Meshes[i].MaterialIndex)
		for j := 0; j < len(scene.Meshes[i].Vertices); j++ {
			fmt.Printf("V(%v): (%v, %v, %v)\n", j, scene.Meshes[i].Vertices[j].X(), scene.Meshes[i].Vertices[j].Y(), scene.Meshes[i].Vertices[j].Z())
		}
	}

	for i := 0; i < len(scene.Materials); i++ {

		m := scene.Materials[i]
		println("Material:", i, "; Props:", len(scene.Materials[i].Properties))
		texCount := asig.GetMaterialTextureCount(m, asig.TextureTypeDiffuse)
		fmt.Println("Texture count:", texCount)

		if texCount > 0 {

			texInfo, err := asig.GetMaterialTexture(m, asig.TextureTypeDiffuse, 0)
			if err != nil {
				panic(err)
			}

			fmt.Printf("%v", texInfo)
		}
	}

	ts := scene.Textures
	for i := 0; i < len(ts); i++ {
		t := ts[i]

		fmt.Printf("T(%v): Name=%v, Hint=%v, Width=%v, Height=%v, NumTexels=%v\n", i, t.Filename, t.FormatHint, t.Width, t.Height, len(t.Data))

		if t.FormatHint == "png" {
			decodePNG(t.Data)
		}
	}
}

func decodePNG(texels []byte) {

	img, err := png.Decode(bytes.NewReader(texels))
	if err != nil {
		panic("wow2: " + err.Error())
	}

	println("C:", img.At(100, 100))
}

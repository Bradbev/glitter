package asset

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/Bradbev/glitter/src/ren"
	"github.com/bloeys/assimp-go/asig/asig"
	"github.com/bloeys/gglm/gglm"
	"github.com/go-gl/gl/v4.6-core/gl"
)

func castGglmVec3ToFloat32(in []gglm.Vec3) []float32 {
	data := (*float32)(unsafe.Pointer(&in[0].Data[0]))
	return unsafe.Slice(data, 3*len(in))
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
			mesh.Textures = append(mesh.Textures, loadMaterialTextures(fs, mat, asig.TextureTypeDiffuse)...)
			mesh.Textures = append(mesh.Textures, loadMaterialTextures(fs, mat, asig.TextureTypeNormal)...)
			mesh.Textures = append(mesh.Textures, loadMaterialTextures(fs, mat, asig.TextureTypeSpecular)...)
		}
		result.Meshes = append(result.Meshes, mesh)
	}

	return result, nil
}

func loadMaterialTextures(fsys fs.FS, material *asig.Material, textureType asig.TextureType) []ren.TextureAndType {
	result := []ren.TextureAndType{}
	for i := 0; i < asig.GetMaterialTextureCount(material, textureType); i++ {
		info, err := asig.GetMaterialTexture(material, textureType, uint(i))
		if err != nil {
			continue
		}
		tex, err := ren.NewTextureFS(fsys, info.Path, gl.REPEAT, gl.REPEAT)
		if err != nil {
			continue
		}
		result = append(result, ren.TextureAndType{
			Texture: tex,
			Type:    asigTexTypeToRenTexType(textureType),
		})
	}
	return result
}

func asigTexTypeToRenTexType(textureType asig.TextureType) ren.TextureType {
	switch textureType {
	case asig.TextureTypeDiffuse:
		return ren.TexDiffuse
	case asig.TextureTypeSpecular:
		return ren.TexSpecular
	case asig.TextureTypeNormal:
		return ren.TexNormal
	default:
		panic(fmt.Sprintf("Unknown texture %d", textureType))
	}
}

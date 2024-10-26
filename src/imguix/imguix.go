package imguix

import "github.com/AllenDang/cimgui-go/imgui"

func TreeNode(name string, body func()) bool {
	if imgui.TreeNodeExStrV(name, imgui.TreeNodeFlagsDefaultOpen) {
		body()
		imgui.TreePop()
		return true
	}
	return false
}

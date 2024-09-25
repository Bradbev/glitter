#version 410
layout(location = 0) in vec3 vp;

//attribute vec3 v_color;
out vec4 vertexColor;

void main() {
	gl_Position = vec4(vp, 1.0);
	vertexColor = vec4(0.5, 0.0, 0.0, 1.0);
}

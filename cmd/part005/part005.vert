#version 330 core
layout(location = 0) in vec3 aPos;
layout(location = 1) in vec3 aColor;
layout(location = 2) in vec2 aTexCoord;

//attribute vec3 v_color;
out vec4 ourColor;
out vec2 TexCoord;

void main() {
	gl_Position = vec4(aPos, 1.0);
	//ourColor = vec4(0.5, 0.0, 0.0, 1.0);
	ourColor = vec4(aColor, 1);
	TexCoord = aTexCoord;
}

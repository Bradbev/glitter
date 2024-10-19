#version 330 core
layout(location = 0) in vec3 aPos;
layout(location = 1) in vec3 aNormal;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

out vec3 Normal;
out vec3 FragPos;

void main() {
	gl_Position = projection * view * model * vec4(aPos, 1.0);

 	// put the FragPos into world space
	FragPos = vec3(model * vec4(aPos, 1.0));

	// Normals are transformed differently 
	// https://paroj.github.io/gltut/Illumination/Tut09%20Normal%20Transformation.html
	Normal = mat3(transpose(inverse(model))) * aNormal;
}

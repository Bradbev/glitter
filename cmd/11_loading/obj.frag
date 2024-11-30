#version 330 core
out vec4 FragColor;

in vec2 TexCoords;

uniform sampler2D texture_diffuse1;

void main() {
	float gamma = 1.6;
	vec3 result = vec3(texture(texture_diffuse1, TexCoords));
	FragColor = vec4(pow(result, vec3(1.0 / gamma)), 1.0);
}
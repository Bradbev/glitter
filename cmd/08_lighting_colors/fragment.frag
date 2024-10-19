#version 330 core
in vec3 Normal;
in vec3 FragPos;

out vec4 FragColor;

uniform vec3 objectColor;
uniform vec3 lightColor;
uniform vec3 lightPos;

void main() {
	// ambient with no falloff 
	float ambientStrength = 0.1;
	vec3 ambient = ambientStrength * lightColor;

 	// diffuse is the dot of the surface normal and the light direction 0..1
	vec3 norm = normalize(Normal);
	vec3 lightDir = normalize(lightPos - FragPos);
	float diff = max(dot(lightDir, norm), 0.0); // if the dot is -ve it faces away
	vec3 diffuse = diff * lightColor;

	vec3 result = (ambient + diffuse) * objectColor;
	FragColor = vec4(result, 1.0);
}
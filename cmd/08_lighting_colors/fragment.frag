#version 330 core
in vec3 Normal;
in vec3 FragPos;

out vec4 FragColor;

uniform vec3 objectColor;
uniform vec3 lightColor;
uniform vec3 lightPos;
uniform vec3 viewPos;

void main() {
	// ambient with no falloff 
	float ambientStrength = 0.1;
	vec3 ambient = ambientStrength * lightColor;

 	// diffuse is the dot of the surface normal and the light direction 0..1
	vec3 norm = normalize(Normal);
	vec3 lightDir = normalize(lightPos - FragPos);
	float diff = max(dot(lightDir, norm), 0.0); // if the dot is -ve it faces away
	vec3 diffuse = diff * lightColor;

	// specular.
	// Based on the reflected light ray bouncing into the eye.
	// future optimization is to do lighting calcs in view space
	float specularStrength = 0.5;
	vec3 viewDir = normalize(viewPos - FragPos);
	vec3 reflectDir = reflect(-lightDir, norm);
	float spec = pow(max(dot(viewDir, reflectDir), 0.0), 32);
	vec3 specular = specularStrength * spec * lightColor;

	// combine
	vec3 result = (ambient + diffuse + specular) * objectColor;
	FragColor = vec4(result, 1.0);
}
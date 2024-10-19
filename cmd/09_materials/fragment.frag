#version 330 core

struct Material {
	vec3 ambient;
	vec3 diffuse;
	vec3 specular;
	float shininess;
};

struct Light {
	vec3 position;
	vec3 ambient;
	vec3 diffuse;
	vec3 specular;
};

uniform Material material;
uniform Light light;

in vec3 Normal;
in vec3 FragPos;

out vec4 FragColor;

uniform vec3 lightColor;
uniform vec3 viewPos;

void main() {
	// ambient with no falloff 
	vec3 ambient = light.ambient * material.ambient;

 	// diffuse is the dot of the surface normal and the light direction 0..1
	vec3 norm = normalize(Normal);
	vec3 lightDir = normalize(light.position - FragPos);
	float diff = max(dot(lightDir, norm), 0.0); // if the dot is -ve it faces away
	vec3 diffuse = light.diffuse * (diff * material.diffuse);

	// specular.
	// Based on the reflected light ray bouncing into the eye.
	// future optimization is to do lighting calcs in view space
	vec3 viewDir = normalize(viewPos - FragPos);
	vec3 reflectDir = reflect(-lightDir, norm);
	float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shininess);
	vec3 specular = light.specular * (spec * material.specular); 

	// combine
	vec3 result = ambient + diffuse + specular;
	FragColor = vec4(result, 1.0);
}
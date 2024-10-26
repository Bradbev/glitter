#version 330 core

struct Material {
	sampler2D Diffuse;
	sampler2D Specular;
	float Shininess;
};

struct Light {
	vec3 Position;
	vec3 Ambient;
	vec3 Diffuse;
	vec3 Specular;
};

uniform Material material;
uniform Light light;
uniform vec3 viewPos;

in vec3 Normal;
in vec3 FragPos;
in vec2 TexCoords;

out vec4 FragColor;

void main() {
	// ambient with no falloff 
//	vec3 ambient = light.ambient * material.ambient;
	vec3 ambient = light.Ambient * texture(material.Diffuse, TexCoords).rgb;

 	// diffuse is the dot of the surface normal and the light direction 0..1
	vec3 norm = normalize(Normal);
	vec3 lightDir = normalize(light.Position - FragPos);
	float diff = max(dot(lightDir, norm), 0.0); // if the dot is -ve it faces away
//	vec3 diffuse = light.diffuse * (diff * material.diffuse);
	vec3 diffuse = light.Diffuse * diff * texture(material.Diffuse, TexCoords).rgb;

	// specular.
	// Based on the reflected light ray bouncing into the eye.
	// future optimization is to do lighting calcs in view space
	vec3 viewDir = normalize(viewPos - FragPos);
	vec3 reflectDir = reflect(-lightDir, norm);
	float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.Shininess);
	vec3 specular = light.Specular * spec * texture(material.Specular, TexCoords).rgb; 

	// combine
	vec3 result = ambient + diffuse + specular;

 // gamma correct our scene to roughly match what I see online
	float gamma = 1.6;
	FragColor = vec4(pow(result, vec3(1.0 / gamma)), 1.0);
}
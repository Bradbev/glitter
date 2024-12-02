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

struct DirectionalLight {
	vec3 Direction;
	vec3 Ambient;
	vec3 Diffuse;
	vec3 Specular;
};
uniform DirectionalLight directionalLight;

struct PointLight {
	vec3 Position;

	float Constant;
	float Linear;
	float Quadratic;

	vec3 Ambient;
	vec3 Diffuse;
	vec3 Specular;
};  
#define NR_POINT_LIGHTS 4  
uniform PointLight pointLights[NR_POINT_LIGHTS];

uniform Material material;
uniform Light light;
uniform vec3 viewPos;

in vec3 Normal;
in vec3 FragPos;
in vec2 TexCoords;

out vec4 FragColor;

vec3 CalcDirLight(DirectionalLight light, vec3 normal, vec3 viewDir) {
	vec3 lightDir = normalize(-light.Direction);
    // diffuse shading
	float diff = max(dot(normal, lightDir), 0.0);
    // specular shading
	vec3 reflectDir = reflect(-lightDir, normal);
	float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.Shininess);
    // combine results
	vec3 ambient = light.Ambient * vec3(texture(material.Diffuse, TexCoords));
	vec3 diffuse = light.Diffuse * diff * vec3(texture(material.Diffuse, TexCoords));
	vec3 specular = light.Specular * spec * vec3(texture(material.Specular, TexCoords));
	return (ambient + diffuse + specular);
}

vec3 CalcPointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir) {
	vec3 lightDir = normalize(light.Position - fragPos);
    // diffuse shading
	float diff = max(dot(normal, lightDir), 0.0);
    // specular shading
	vec3 reflectDir = reflect(-lightDir, normal);
	float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.Shininess);
    // attenuation
	float distance = length(light.Position - fragPos);
	float attenuation = 1.0 / (light.Constant + light.Linear * distance +
		light.Quadratic * (distance * distance));    
    // combine results
	vec3 ambient = light.Ambient * vec3(texture(material.Diffuse, TexCoords));
	vec3 diffuse = light.Diffuse * diff * vec3(texture(material.Diffuse, TexCoords));
	vec3 specular = light.Specular * spec * vec3(texture(material.Specular, TexCoords));
	ambient *= attenuation;
	diffuse *= attenuation;
	specular *= attenuation;
	return (ambient + diffuse + specular);
}

void main() {
    // properties
	vec3 norm = normalize(Normal);
	vec3 viewDir = normalize(viewPos - FragPos);

    // phase 1: Directional lighting
	vec3 result = CalcDirLight(directionalLight, norm, viewDir);

    // phase 2: Point lights
	for(int i = 0; i < NR_POINT_LIGHTS; i++) {
		result += CalcPointLight(pointLights[i], norm, FragPos, viewDir);
	}

    // phase 3: Spot light
    //result += CalcSpotLight(spotLight, norm, FragPos, viewDir);    

	FragColor = vec4(result, 1.0);
}
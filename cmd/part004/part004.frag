#version 330 core
out vec4 FragColor;
in vec4 ourColor;
in vec2 TexCoord;

uniform sampler2D ourTexture;
//uniform vec4 vertexColor;

uniform float blend;

void main() {
	//FragColor = vertexColor;
//	FragColor = ourColor;
	FragColor = texture(ourTexture, TexCoord);
	// mix the texture with the color attribute
	FragColor = mix(texture(ourTexture, TexCoord), ourColor, blend);
}

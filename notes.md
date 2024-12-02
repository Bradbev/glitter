# Problems

# Solved
- gamma needs to be raised by 1.6
 - not the texture loading
 - not the SDL path
 - GL load format.  I used SRGB_ALPHA, not RGBA

- random crashes at startup
 - appears to be SDL/imgui thing

- I need to swizzle X/Y when loading obj files (because Z up)
- object is mirrored about the Z axis (because Z up)
- Use Y up!

SO far, no GLFW rando crashes

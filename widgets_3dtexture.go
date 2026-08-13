package lvgl

/*
#include <lvgl.h>
*/
import "C"

// Texture3D wraps an lv_3dtexture widget: displays an OpenGL texture the
// caller has already rendered elsewhere (this package doesn't manage an
// OpenGL context or create textures itself — LVGL just needs the numeric
// texture handle, e.g. from a `glGenTextures`/`glTexImage2D` call made
// through whatever Go OpenGL binding the host app is using).
type Texture3D struct{ *Obj }

// NewTexture3D creates a 3D texture display as a child of parent. Set the
// object's size to match the texture with SetSize.
func NewTexture3D(parent *Obj) *Texture3D {
	return &Texture3D{wrapObj(C.lv_3dtexture_create(parent.c))}
}

// SetSrc sets the OpenGL texture handle to display.
func (t *Texture3D) SetSrc(textureID uint32) {
	C.lv_3dtexture_set_src(t.c, C.lv_3dtexture_id_t(textureID))
}

// SetFlip sets whether the texture is flipped horizontally/vertically
// when drawn (OpenGL textures are often upside-down relative to LVGL's
// coordinate system).
func (t *Texture3D) SetFlip(hFlip, vFlip bool) {
	C.lv_3dtexture_set_flip(t.c, C.bool(hFlip), C.bool(vFlip))
}

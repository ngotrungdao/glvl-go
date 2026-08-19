// Package lvgl is a hand-written Go wrapper around LVGL 9.6, targeting the
// SDL2 display backend. It covers the core object model, the event
// system, styles, flex/grid layout, animation, and a broad (not exhaustive)
// set of common widgets.
//
// # Memory ownership
//
// LVGL objects (Obj, and anything embedding it) are owned by LVGL itself:
// call Delete on the root of a subtree to free it, same as lv_obj_delete.
//
// Style and GridDsc are different: LVGL stores a raw, non-refcounted
// pointer to the C memory backing them, so they are allocated on the C
// heap and must be freed explicitly with Delete once they are no longer
// attached to any object. See core_style.go and layouts_grid.go.
//
// # File layout
//
// This is a single Go package (cgo's C pseudo-type isn't shareable across
// package boundaries, so splitting into core/widgets/... sub-packages
// would need unsafe.Pointer round-trips at every cross-package call) —
// files are grouped by filename prefix instead, mirroring
// install/include/lvgl's directory layout: core_*.go (object model,
// events, animation, styles), draw_color.go, layouts_*.go (flex/grid),
// display_sdl.go, font_symbols.go, and widgets_*.go (one per LVGL
// widget). cgo.go/app.go/handle.go are package-level glue with no single
// LVGL header counterpart.
package lvgl

/*
#cgo CFLAGS: -I${SRCDIR}/lvgl-c/install/include/lvgl -DLV_CONF_INCLUDE_SIMPLE
#cgo LDFLAGS: -L${SRCDIR}/lvgl-c/install/lib64 -llvgl -lSDL2 -lwayland-client -lwayland-cursor -lwayland-egl -lxkbcommon -ldrm -lgbm -lfreetype -ljpeg -lpng -lwebp -lm -lgstvideo-1.0 -lgstbase-1.0 -lgstreamer-1.0 -lgobject-2.0 -lglib-2.0 -lavformat -lavcodec -lavutil -lswscale
#include <lvgl.h>
*/
import "C"

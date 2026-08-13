package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// GridAlign mirrors lv_grid_align_t.
type GridAlign uint32

var (
	GridAlignStart        = GridAlign(C.LV_GRID_ALIGN_START)
	GridAlignCenter       = GridAlign(C.LV_GRID_ALIGN_CENTER)
	GridAlignEnd          = GridAlign(C.LV_GRID_ALIGN_END)
	GridAlignStretch      = GridAlign(C.LV_GRID_ALIGN_STRETCH)
	GridAlignSpaceEvenly  = GridAlign(C.LV_GRID_ALIGN_SPACE_EVENLY)
	GridAlignSpaceAround  = GridAlign(C.LV_GRID_ALIGN_SPACE_AROUND)
	GridAlignSpaceBetween = GridAlign(C.LV_GRID_ALIGN_SPACE_BETWEEN)
)

// GridTemplateLast terminates a grid track descriptor array
// (LV_GRID_TEMPLATE_LAST). GridContent sizes a track to fit its content
// (LV_GRID_CONTENT). Both are reimplemented from core/lv_area.h's
// LV_COORD_MAX rather than read via cgo, since they are simple arithmetic
// on it (see CoordMax in core_obj.go).
const (
	GridTemplateLast = CoordMax
	GridContent      = CoordMax - 101
)

// Fr returns a "fractional unit" grid track size, equivalent to the
// LV_GRID_FR(x) macro: the track gets x parts of the remaining free space.
func Fr(x uint8) int32 {
	return int32(C.lv_grid_fr(C.uint8_t(x)))
}

// GridDsc is a grid column/row track descriptor array. LVGL stores a raw
// pointer to it (not a copy), so it is allocated on the C heap; call
// Delete once it is no longer attached to any object via SetGridDscArray.
type GridDsc struct {
	c *C.int32_t
}

// NewGridDsc builds a track descriptor from track sizes (in pixels, or via
// Pct/Fr/GridContent), automatically appending the GridTemplateLast
// terminator LVGL requires.
func NewGridDsc(tracks ...int32) *GridDsc {
	n := len(tracks) + 1
	c := (*C.int32_t)(C.malloc(C.size_t(n) * C.size_t(unsafe.Sizeof(C.int32_t(0)))))
	buf := unsafe.Slice((*int32)(unsafe.Pointer(c)), n)
	copy(buf, tracks)
	buf[n-1] = GridTemplateLast
	return &GridDsc{c: c}
}

// Delete frees the C memory backing the descriptor. Only call this once
// the descriptor is detached from every object it was set on.
func (g *GridDsc) Delete() {
	if g.c == nil {
		return
	}
	C.free(unsafe.Pointer(g.c))
	g.c = nil
}

// SetGridDscArray configures the object to lay out its children on a grid
// with the given column and row tracks. The Obj keeps a reference to cols
// and rows so they aren't garbage-collected while still needed, but does
// not take ownership: the caller is still responsible for calling Delete
// on them once the object no longer uses this layout.
func (o *Obj) SetGridDscArray(cols, rows *GridDsc) {
	C.lv_obj_set_grid_dsc_array(o.c, cols.c, rows.c)
	o.gridCols, o.gridRows = cols, rows
}

// SetGridAlign sets how the whole grid is aligned within the object when
// it doesn't fill it exactly.
func (o *Obj) SetGridAlign(columnAlign, rowAlign GridAlign) {
	C.lv_obj_set_grid_align(o.c, C.lv_grid_align_t(columnAlign), C.lv_grid_align_t(rowAlign))
}

// SetGridCell places the object within its parent's grid at the given
// column/row position and span, with per-axis alignment.
func (o *Obj) SetGridCell(colAlign GridAlign, colPos, colSpan int32, rowAlign GridAlign, rowPos, rowSpan int32) {
	C.lv_obj_set_grid_cell(o.c, C.lv_grid_align_t(colAlign), C.int32_t(colPos), C.int32_t(colSpan),
		C.lv_grid_align_t(rowAlign), C.int32_t(rowPos), C.int32_t(rowSpan))
}

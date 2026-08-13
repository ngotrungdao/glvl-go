package lvgl

/*
#include <lvgl.h>
*/
import "C"

// TileView wraps an lv_tileview widget: a grid of full-size "tiles" the
// user swipes between.
type TileView struct{ *Obj }

// NewTileView creates a tile view as a child of parent.
func NewTileView(parent *Obj) *TileView {
	return &TileView{wrapObj(C.lv_tileview_create(parent.c))}
}

// AddTile adds a tile at the given column/row, allowing swiping in the
// given direction(s) (see Dir* constants in core_dir.go, e.g. DirHor).
func (t *TileView) AddTile(col, row uint8, dir Dir) *Obj {
	return wrapObj(C.lv_tileview_add_tile(t.c, C.uint8_t(col), C.uint8_t(row), C.lv_dir_t(dir)))
}

// SetTile scrolls to the given tile, optionally animating.
func (t *TileView) SetTile(tile *Obj, animate bool) {
	C.lv_tileview_set_tile(t.c, tile.c, C.lv_anim_enable_t(animate))
}

// TileActive returns the currently visible tile.
func (t *TileView) TileActive() *Obj { return wrapObj(C.lv_tileview_get_tile_active(t.c)) }

// SetTileByIndex scrolls to the tile at the given column/row, optionally
// animating.
func (t *TileView) SetTileByIndex(col, row uint32, animate bool) {
	C.lv_tileview_set_tile_by_index(t.c, C.uint32_t(col), C.uint32_t(row), C.lv_anim_enable_t(animate))
}

package lvgl

/*
#include <lvgl.h>

extern lv_obj_tree_walk_res_t goTreeWalkTrampoline(lv_obj_t *obj, void *user_data);
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

// TreeWalkResult mirrors lv_obj_tree_walk_res_t, returned from a TreeWalk
// callback to control traversal.
type TreeWalkResult uint32

const (
	// TreeWalkNext continues the walk normally (into children, then siblings).
	TreeWalkNext TreeWalkResult = iota
	// TreeWalkSkipChildren continues the walk but skips this node's children.
	TreeWalkSkipChildren
	// TreeWalkEnd stops the walk immediately.
	TreeWalkEnd
)

type treeWalkCallback struct {
	fn func(*Obj) TreeWalkResult
}

// TreeWalk calls fn for the object and every descendant, depth-first,
// stopping early if fn returns TreeWalkEnd or skipping a subtree if it
// returns TreeWalkSkipChildren.
func (o *Obj) TreeWalk(fn func(*Obj) TreeWalkResult) {
	h := cgo.NewHandle(&treeWalkCallback{fn: fn})
	defer h.Delete()
	C.lv_obj_tree_walk(o.c, C.lv_obj_tree_walk_cb_t(C.goTreeWalkTrampoline), unsafeHandlePointer(h))
}

//export goTreeWalkTrampoline
func goTreeWalkTrampoline(obj *C.lv_obj_t, userData unsafe.Pointer) C.lv_obj_tree_walk_res_t {
	h := cgo.Handle(uintptr(userData))
	cb := h.Value().(*treeWalkCallback)
	return C.lv_obj_tree_walk_res_t(cb.fn(wrapObj(obj)))
}

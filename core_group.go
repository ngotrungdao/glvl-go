package lvgl

/*
#include <lvgl.h>
*/
import "C"

// GroupRefocusPolicy mirrors lv_group_refocus_policy_t.
type GroupRefocusPolicy uint32

var (
	GroupRefocusPolicyNext = GroupRefocusPolicy(C.LV_GROUP_REFOCUS_POLICY_NEXT)
	GroupRefocusPolicyPrev = GroupRefocusPolicy(C.LV_GROUP_REFOCUS_POLICY_PREV)
)

// Group wraps an lv_group_t: an ordered set of objects that keyboard/
// encoder input devices can move focus between (via Indev), for
// interfaces without a pointer/touch device. Attach an indev to a group
// with Indev.SetGroup.
type Group struct {
	c *C.lv_group_t
}

// NewGroup creates a new, empty input group.
func NewGroup() *Group {
	c := C.lv_group_create()
	if c == nil {
		return nil
	}
	return &Group{c: c}
}

// Delete frees the group. Objects in it are not deleted, just removed
// from the group.
func (g *Group) Delete() {
	if g.c == nil {
		return
	}
	C.lv_group_delete(g.c)
	g.c = nil
}

// SetDefaultGroup sets g as the default group new indevs attach to
// automatically. Pass nil to clear the default.
func SetDefaultGroup(g *Group) {
	var c *C.lv_group_t
	if g != nil {
		c = g.c
	}
	C.lv_group_set_default(c)
}

// DefaultGroup returns the current default group, if any.
func DefaultGroup() *Group {
	c := C.lv_group_get_default()
	if c == nil {
		return nil
	}
	return &Group{c: c}
}

// AddObj adds an object to the end of the group.
func (g *Group) AddObj(o *Obj) { C.lv_group_add_obj(g.c, o.c) }

// RemoveFromGroup removes an object from whichever group it's in, if any.
func RemoveFromGroup(o *Obj) { C.lv_group_remove_obj(o.c) }

// RemoveAllObjs removes every object from the group.
func (g *Group) RemoveAllObjs() { C.lv_group_remove_all_objs(g.c) }

// FocusObj focuses a specific object (it must already be in a group).
func FocusObj(o *Obj) { C.lv_group_focus_obj(o.c) }

// FocusNext moves focus to the next object in the group.
func (g *Group) FocusNext() { C.lv_group_focus_next(g.c) }

// FocusPrev moves focus to the previous object in the group.
func (g *Group) FocusPrev() { C.lv_group_focus_prev(g.c) }

// FocusFreeze enables/disables temporarily locking focus on the current
// object, ignoring FocusNext/FocusPrev.
func (g *Group) FocusFreeze(enable bool) { C.lv_group_focus_freeze(g.c, C.bool(enable)) }

// SetRefocusPolicy sets which neighbor gets focus when the currently
// focused object is deleted or removed from the group.
func (g *Group) SetRefocusPolicy(policy GroupRefocusPolicy) {
	C.lv_group_set_refocus_policy(g.c, C.lv_group_refocus_policy_t(policy))
}

// SetEditing switches the group between navigate mode and edit mode
// (relevant for widgets like sliders/rollers that consume encoder input
// differently once "entered").
func (g *Group) SetEditing(edit bool) { C.lv_group_set_editing(g.c, C.bool(edit)) }

// SetWrap enables/disables wrapping focus from the last object back to
// the first (and vice versa).
func (g *Group) SetWrap(enable bool) { C.lv_group_set_wrap(g.c, C.bool(enable)) }

// Focused returns the currently focused object, if any.
func (g *Group) Focused() *Obj { return wrapObj(C.lv_group_get_focused(g.c)) }

// Editing reports whether the group is in edit mode.
func (g *Group) Editing() bool { return bool(C.lv_group_get_editing(g.c)) }

// Wrap reports whether focus wraps from last to first object.
func (g *Group) Wrap() bool { return bool(C.lv_group_get_wrap(g.c)) }

// ObjCount returns how many objects are in the group.
func (g *Group) ObjCount() uint32 { return uint32(C.lv_group_get_obj_count(g.c)) }

// ObjByIndex returns the object at the given index in the group.
func (g *Group) ObjByIndex(index uint32) *Obj {
	return wrapObj(C.lv_group_get_obj_by_index(g.c, C.uint32_t(index)))
}

// GroupCount returns the total number of groups that currently exist.
func GroupCount() uint32 { return uint32(C.lv_group_get_count()) }

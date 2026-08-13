package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// TabView wraps an lv_tabview widget.
type TabView struct{ *Obj }

// NewTabView creates a tab view as a child of parent.
func NewTabView(parent *Obj) *TabView {
	return &TabView{wrapObj(C.lv_tabview_create(parent.c))}
}

// AddTab adds a new tab with the given name and returns its content
// container, to which children should be added.
func (t *TabView) AddTab(name string) *Obj {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return wrapObj(C.lv_tabview_add_tab(t.c, cName))
}

// SetActive switches to the tab at the given index, optionally animating.
func (t *TabView) SetActive(idx uint32, animate bool) {
	C.lv_tabview_set_active(t.c, C.uint32_t(idx), C.lv_anim_enable_t(animate))
}

// TabActive returns the index of the currently active tab.
func (t *TabView) TabActive() uint32 { return uint32(C.lv_tabview_get_tab_active(t.c)) }

// SetTabText renames a tab.
func (t *TabView) SetTabText(idx uint32, name string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.lv_tabview_set_tab_text(t.c, C.uint32_t(idx), cName)
}

// SetTabBarPosition sets which side of the tab view the tab bar sits on.
func (t *TabView) SetTabBarPosition(dir Dir) {
	C.lv_tabview_set_tab_bar_position(t.c, C.lv_dir_t(dir))
}

// SetTabBarSize sets the tab bar's thickness in pixels.
func (t *TabView) SetTabBarSize(size int32) {
	C.lv_tabview_set_tab_bar_size(t.c, C.int32_t(size))
}

// TabCount returns the number of tabs.
func (t *TabView) TabCount() uint32 { return uint32(C.lv_tabview_get_tab_count(t.c)) }

// TabButton returns the tab bar button object at the given index.
func (t *TabView) TabButton(idx int32) *Obj {
	return wrapObj(C.lv_tabview_get_tab_button(t.c, C.int32_t(idx)))
}

// Content returns the container that holds all tab content panels.
func (t *TabView) Content() *Obj { return wrapObj(C.lv_tabview_get_content(t.c)) }

// TabBar returns the tab bar object (holds the tab buttons).
func (t *TabView) TabBar() *Obj { return wrapObj(C.lv_tabview_get_tab_bar(t.c)) }

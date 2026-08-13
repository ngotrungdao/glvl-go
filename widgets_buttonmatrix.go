package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// ButtonMatrixCtrl mirrors lv_buttonmatrix_ctrl_t control bits (width
// values 1-15 are stored in the low bits; OR in the flags below).
type ButtonMatrixCtrl uint16

var (
	ButtonMatrixCtrlNone      = ButtonMatrixCtrl(C.LV_BUTTONMATRIX_CTRL_NONE)
	ButtonMatrixCtrlHidden    = ButtonMatrixCtrl(C.LV_BUTTONMATRIX_CTRL_HIDDEN)
	ButtonMatrixCtrlNoRepeat  = ButtonMatrixCtrl(C.LV_BUTTONMATRIX_CTRL_NO_REPEAT)
	ButtonMatrixCtrlDisabled  = ButtonMatrixCtrl(C.LV_BUTTONMATRIX_CTRL_DISABLED)
	ButtonMatrixCtrlCheckable = ButtonMatrixCtrl(C.LV_BUTTONMATRIX_CTRL_CHECKABLE)
	ButtonMatrixCtrlChecked   = ButtonMatrixCtrl(C.LV_BUTTONMATRIX_CTRL_CHECKED)
	ButtonMatrixCtrlClickTrig = ButtonMatrixCtrl(C.LV_BUTTONMATRIX_CTRL_CLICK_TRIG)
	ButtonMatrixCtrlPopover   = ButtonMatrixCtrl(C.LV_BUTTONMATRIX_CTRL_POPOVER)
	ButtonMatrixCtrlRecolor   = ButtonMatrixCtrl(C.LV_BUTTONMATRIX_CTRL_RECOLOR)
)

// ButtonMatrixWidth returns the control bits setting a button's relative
// width (1-15) within its row.
func ButtonMatrixWidth(w uint8) ButtonMatrixCtrl { return ButtonMatrixCtrl(w & 0x0F) }

// ButtonMatrix wraps an lv_buttonmatrix widget: a grid of buttons defined
// by a text map.
//
// LVGL keeps a reference to the map array passed to SetMap rather than
// copying it (a well known LVGL requirement: the map must outlive the
// widget), so SetMap builds and pins a NUL-terminated C string array on
// the C heap, freeing the previous one on replacement.
type ButtonMatrix struct {
	*Obj
	mapPtrs []*C.char
	mapArr  *C.char
}

// NewButtonMatrix creates a button matrix as a child of parent.
func NewButtonMatrix(parent *Obj) *ButtonMatrix {
	return &ButtonMatrix{Obj: wrapObj(C.lv_buttonmatrix_create(parent.c))}
}

// SetMap sets the button matrix's layout: each row is a slice of button
// labels, rendered with a "\n" appended after the last button in a row
// (except the final row), matching lv_buttonmatrix's map format.
func (bm *ButtonMatrix) SetMap(rows [][]string) {
	var flatEntries []string
	for ri, row := range rows {
		for bi, label := range row {
			if bi == len(row)-1 && ri != len(rows)-1 {
				flatEntries = append(flatEntries, label+"\n")
			} else {
				flatEntries = append(flatEntries, label)
			}
		}
	}

	cStrs := make([]*C.char, len(flatEntries)+1) // +1 for the NULL terminator
	for i, s := range flatEntries {
		cStrs[i] = C.CString(s)
	}
	cStrs[len(flatEntries)] = nil

	arr := (**C.char)(C.malloc(C.size_t(len(cStrs)) * C.size_t(unsafe.Sizeof((*C.char)(nil)))))
	slice := unsafe.Slice(arr, len(cStrs))
	copy(slice, cStrs)

	C.lv_buttonmatrix_set_map(bm.c, arr)

	bm.freeMap()
	bm.mapPtrs = cStrs
	bm.mapArr = (*C.char)(unsafe.Pointer(arr))
}

func (bm *ButtonMatrix) freeMap() {
	for _, p := range bm.mapPtrs {
		if p != nil {
			C.free(unsafe.Pointer(p))
		}
	}
	if bm.mapArr != nil {
		C.free(unsafe.Pointer(bm.mapArr))
	}
	bm.mapPtrs = nil
	bm.mapArr = nil
}

// SetButtonCtrl sets control flags/width on a single button.
func (bm *ButtonMatrix) SetButtonCtrl(btnID uint32, ctrl ButtonMatrixCtrl) {
	C.lv_buttonmatrix_set_button_ctrl(bm.c, C.uint32_t(btnID), C.lv_buttonmatrix_ctrl_t(ctrl))
}

// ClearButtonCtrl clears control flags on a single button.
func (bm *ButtonMatrix) ClearButtonCtrl(btnID uint32, ctrl ButtonMatrixCtrl) {
	C.lv_buttonmatrix_clear_button_ctrl(bm.c, C.uint32_t(btnID), C.lv_buttonmatrix_ctrl_t(ctrl))
}

// SetButtonCtrlAll sets control flags on every button.
func (bm *ButtonMatrix) SetButtonCtrlAll(ctrl ButtonMatrixCtrl) {
	C.lv_buttonmatrix_set_button_ctrl_all(bm.c, C.lv_buttonmatrix_ctrl_t(ctrl))
}

// ClearButtonCtrlAll clears control flags on every button.
func (bm *ButtonMatrix) ClearButtonCtrlAll(ctrl ButtonMatrixCtrl) {
	C.lv_buttonmatrix_clear_button_ctrl_all(bm.c, C.lv_buttonmatrix_ctrl_t(ctrl))
}

// SetButtonWidth sets a button's relative width within its row (1-15;
// see also ButtonMatrixWidth).
func (bm *ButtonMatrix) SetButtonWidth(btnID, width uint32) {
	C.lv_buttonmatrix_set_button_width(bm.c, C.uint32_t(btnID), C.uint32_t(width))
}

// HasButtonCtrl reports whether all the given control flags are set on a button.
func (bm *ButtonMatrix) HasButtonCtrl(btnID uint32, ctrl ButtonMatrixCtrl) bool {
	return bool(C.lv_buttonmatrix_has_button_ctrl(bm.c, C.uint32_t(btnID), C.lv_buttonmatrix_ctrl_t(ctrl)))
}

// OneChecked reports whether the matrix is restricted to one checked
// button at a time.
func (bm *ButtonMatrix) OneChecked() bool { return bool(C.lv_buttonmatrix_get_one_checked(bm.c)) }

// SetOneChecked restricts the matrix to at most one checked button at a
// time (radio-button behavior), when buttons are made checkable.
func (bm *ButtonMatrix) SetOneChecked(enable bool) {
	C.lv_buttonmatrix_set_one_checked(bm.c, C.bool(enable))
}

// SelectedButton returns the index of the last pressed/selected button.
func (bm *ButtonMatrix) SelectedButton() uint32 {
	return uint32(C.lv_buttonmatrix_get_selected_button(bm.c))
}

// ButtonText returns the text of the button at the given index.
func (bm *ButtonMatrix) ButtonText(btnID uint32) string {
	return C.GoString(C.lv_buttonmatrix_get_button_text(bm.c, C.uint32_t(btnID)))
}

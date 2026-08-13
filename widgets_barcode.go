package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

var errBarcodeUpdate = errors.New("lvgl: barcode encoding failed (data too long or invalid for the selected encoding)")

// BarcodeEncoding mirrors lv_barcode_encoding_t.
type BarcodeEncoding uint32

var (
	BarcodeEncodingCode128GS1 = BarcodeEncoding(C.LV_BARCODE_ENCODING_CODE128_GS1)
	BarcodeEncodingCode128Raw = BarcodeEncoding(C.LV_BARCODE_ENCODING_CODE128_RAW)
)

// Barcode wraps an lv_barcode widget.
type Barcode struct{ *Obj }

// NewBarcode creates a barcode as a child of parent. It has no default
// size: call SetSize (inherited from Obj) to give it a non-zero
// width/height before calling Update, or Update fails with a
// "buf_size"-related error.
func NewBarcode(parent *Obj) *Barcode {
	return &Barcode{wrapObj(C.lv_barcode_create(parent.c))}
}

// SetDarkColor sets the color of the barcode's dark bars.
func (b *Barcode) SetDarkColor(c Color) { C.lv_barcode_set_dark_color(b.c, c.toC()) }

// SetLightColor sets the color of the barcode's light bars.
func (b *Barcode) SetLightColor(c Color) { C.lv_barcode_set_light_color(b.c, c.toC()) }

// SetScale sets the width multiplier of each bar.
func (b *Barcode) SetScale(scale uint16) { C.lv_barcode_set_scale(b.c, C.uint16_t(scale)) }

// SetEncoding sets which Code 128 variant to encode with.
func (b *Barcode) SetEncoding(enc BarcodeEncoding) {
	C.lv_barcode_set_encoding(b.c, C.lv_barcode_encoding_t(enc))
}

// SetDirection sets whether bars run horizontally or vertically.
func (b *Barcode) SetDirection(dir Dir) { C.lv_barcode_set_direction(b.c, C.lv_dir_t(dir)) }

// SetTiled enables/disables tiling the barcode to fill the widget's area.
func (b *Barcode) SetTiled(tiled bool) { C.lv_barcode_set_tiled(b.c, C.bool(tiled)) }

// Scale returns the current bar width multiplier.
func (b *Barcode) Scale() uint16 { return uint16(C.lv_barcode_get_scale(b.c)) }

// Update encodes the given text and (re)renders the barcode.
func (b *Barcode) Update(data string) error {
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))
	if C.lv_barcode_update(b.c, cData) != C.LV_RESULT_OK {
		return errBarcodeUpdate
	}
	return nil
}

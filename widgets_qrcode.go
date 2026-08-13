package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// QRCode wraps an lv_qrcode widget.
type QRCode struct{ *Obj }

// NewQRCode creates a QR code as a child of parent. Configure its size and
// colors with SetSize/SetDarkColor/SetLightColor before calling SetData.
func NewQRCode(parent *Obj) *QRCode {
	return &QRCode{wrapObj(C.lv_qrcode_create(parent.c))}
}

// SetSize sets the QR code's width/height in pixels (it's always square).
func (q *QRCode) SetSize(size int32) { C.lv_qrcode_set_size(q.c, C.int32_t(size)) }

// SetDarkColor sets the color of the QR code's dark modules.
func (q *QRCode) SetDarkColor(c Color) { C.lv_qrcode_set_dark_color(q.c, c.toC()) }

// SetLightColor sets the color of the QR code's light modules.
func (q *QRCode) SetLightColor(c Color) { C.lv_qrcode_set_light_color(q.c, c.toC()) }

// SetData encodes the given text and (re)renders the QR code.
func (q *QRCode) SetData(data string) {
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))
	C.lv_qrcode_set_data(q.c, cData)
}

// SetQuietZone enables/disables the light-colored border LVGL adds
// around the QR code by default.
func (q *QRCode) SetQuietZone(enable bool) { C.lv_qrcode_set_quiet_zone(q.c, C.bool(enable)) }

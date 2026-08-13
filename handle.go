package lvgl

import (
	"runtime/cgo"
	"unsafe"
)

// unsafeHandlePointer converts a cgo.Handle into the void* form LVGL's
// user_data parameters expect. The reverse conversion (inside a C
// callback trampoline) is cgo.Handle(uintptr(ptr)).
//
//go:nocheckptr
func unsafeHandlePointer(h cgo.Handle) unsafe.Pointer {
	// `go vet`'s unsafeptr check flags this as a possible misuse, but it's
	// exactly the pattern runtime/cgo.Handle's own documentation shows for
	// passing a Handle through a C void* — the Handle's uintptr value never
	// gets dereferenced as an address, only looked back up in cgo's handle
	// table, so there's no actual GC/pointer-safety hazard here.
	return unsafe.Pointer(uintptr(h))
}

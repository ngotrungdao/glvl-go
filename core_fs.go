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

// FSMode mirrors lv_fs_mode_t; values can be OR'd (e.g. read+write).
type FSMode uint32

var (
	FSModeWrite = FSMode(C.LV_FS_MODE_WR)
	FSModeRead  = FSMode(C.LV_FS_MODE_RD)
)

// FSWhence mirrors lv_fs_whence_t.
type FSWhence uint32

var (
	FSSeekSet = FSWhence(C.LV_FS_SEEK_SET)
	FSSeekCur = FSWhence(C.LV_FS_SEEK_CUR)
	FSSeekEnd = FSWhence(C.LV_FS_SEEK_END)
)

var errFS = errors.New("lvgl: filesystem operation failed")

// IsFSReady reports whether the drive letter (e.g. 'S') is registered
// and ready.
func IsFSReady(letter byte) bool { return bool(C.lv_fs_is_ready(C.char(letter))) }

// File wraps an lv_fs_file_t, opened through LVGL's virtual filesystem
// layer (a path like "S:/dir/file.bin", where the drive letter selects
// among LVGL's registered FS drivers — LV_USE_FS_STDIO is enabled in
// this build). Not related to Go's own os.File.
type File struct {
	c *C.lv_fs_file_t
}

// FSOpen opens a file through LVGL's virtual filesystem.
func FSOpen(path string, mode FSMode) (*File, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	c := (*C.lv_fs_file_t)(C.malloc(C.sizeof_lv_fs_file_t))
	if C.lv_fs_open(c, cPath, C.lv_fs_mode_t(mode)) != C.LV_FS_RES_OK {
		C.free(unsafe.Pointer(c))
		return nil, errFS
	}
	return &File{c: c}, nil
}

// Close closes the file and frees its handle.
func (f *File) Close() error {
	if f.c == nil {
		return nil
	}
	res := C.lv_fs_close(f.c)
	C.free(unsafe.Pointer(f.c))
	f.c = nil
	if res != C.LV_FS_RES_OK {
		return errFS
	}
	return nil
}

// Read reads up to len(buf) bytes into buf, returning how many were read.
func (f *File) Read(buf []byte) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	var br C.uint32_t
	res := C.lv_fs_read(f.c, unsafe.Pointer(&buf[0]), C.uint32_t(len(buf)), &br)
	if res != C.LV_FS_RES_OK {
		return int(br), errFS
	}
	return int(br), nil
}

// Write writes buf to the file, returning how many bytes were written.
func (f *File) Write(buf []byte) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	var bw C.uint32_t
	res := C.lv_fs_write(f.c, unsafe.Pointer(&buf[0]), C.uint32_t(len(buf)), &bw)
	if res != C.LV_FS_RES_OK {
		return int(bw), errFS
	}
	return int(bw), nil
}

// Seek moves the file position.
func (f *File) Seek(pos uint32, whence FSWhence) error {
	if C.lv_fs_seek(f.c, C.uint32_t(pos), C.lv_fs_whence_t(whence)) != C.LV_FS_RES_OK {
		return errFS
	}
	return nil
}

// Tell returns the current file position.
func (f *File) Tell() (uint32, error) {
	var pos C.uint32_t
	if C.lv_fs_tell(f.c, &pos) != C.LV_FS_RES_OK {
		return 0, errFS
	}
	return uint32(pos), nil
}

// Size returns the file's total size in bytes.
func (f *File) Size() (uint32, error) {
	var size C.uint32_t
	if C.lv_fs_get_size(f.c, &size) != C.LV_FS_RES_OK {
		return 0, errFS
	}
	return uint32(size), nil
}

// FSPathSize returns the size in bytes of the file at path, without
// opening it.
func FSPathSize(path string) (uint32, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	var size C.uint32_t
	if C.lv_fs_path_get_size(cPath, &size) != C.LV_FS_RES_OK {
		return 0, errFS
	}
	return uint32(size), nil
}

// FSLoadToBuf reads the entire file at path into a new byte slice.
func FSLoadToBuf(path string) ([]byte, error) {
	size, err := FSPathSize(path)
	if err != nil {
		return nil, err
	}
	if size == 0 {
		return []byte{}, nil
	}
	buf := make([]byte, size)
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	if C.lv_fs_load_to_buf(unsafe.Pointer(&buf[0]), C.uint32_t(size), cPath) != C.LV_FS_RES_OK {
		return nil, errFS
	}
	return buf, nil
}

// FSDir wraps an lv_fs_dir_t for listing a directory's entries.
type FSDir struct {
	c *C.lv_fs_dir_t
}

// FSDirOpen opens a directory for listing through LVGL's virtual
// filesystem.
func FSDirOpen(path string) (*FSDir, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	c := (*C.lv_fs_dir_t)(C.malloc(C.sizeof_lv_fs_dir_t))
	if C.lv_fs_dir_open(c, cPath) != C.LV_FS_RES_OK {
		C.free(unsafe.Pointer(c))
		return nil, errFS
	}
	return &FSDir{c: c}, nil
}

// Read returns the next entry's name ("" once the listing is exhausted).
// Directory names are prefixed with "/" by LVGL's convention.
func (d *FSDir) Read() (string, error) {
	const bufSize = 256
	buf := (*C.char)(C.malloc(bufSize))
	defer C.free(unsafe.Pointer(buf))
	if C.lv_fs_dir_read(d.c, buf, bufSize) != C.LV_FS_RES_OK {
		return "", errFS
	}
	return C.GoString(buf), nil
}

// Close closes the directory listing and frees its handle.
func (d *FSDir) Close() error {
	if d.c == nil {
		return nil
	}
	res := C.lv_fs_dir_close(d.c)
	C.free(unsafe.Pointer(d.c))
	d.c = nil
	if res != C.LV_FS_RES_OK {
		return errFS
	}
	return nil
}

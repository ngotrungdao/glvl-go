package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// TableCellCtrl mirrors lv_table_cell_ctrl_t control bits.
type TableCellCtrl uint32

var (
	TableCellCtrlNone       = TableCellCtrl(C.LV_TABLE_CELL_CTRL_NONE)
	TableCellCtrlMergeRight = TableCellCtrl(C.LV_TABLE_CELL_CTRL_MERGE_RIGHT)
	TableCellCtrlTextCrop   = TableCellCtrl(C.LV_TABLE_CELL_CTRL_TEXT_CROP)
)

// Table wraps an lv_table widget.
type Table struct{ *Obj }

// NewTable creates a table as a child of parent.
func NewTable(parent *Obj) *Table {
	return &Table{wrapObj(C.lv_table_create(parent.c))}
}

// SetRowCount sets the number of rows.
func (t *Table) SetRowCount(n uint32) { C.lv_table_set_row_count(t.c, C.uint32_t(n)) }

// SetColumnCount sets the number of columns.
func (t *Table) SetColumnCount(n uint32) { C.lv_table_set_column_count(t.c, C.uint32_t(n)) }

// SetColumnWidth sets the pixel width of a column.
func (t *Table) SetColumnWidth(col uint32, w int32) {
	C.lv_table_set_column_width(t.c, C.uint32_t(col), C.int32_t(w))
}

// SetCellValue sets the text of a cell.
func (t *Table) SetCellValue(row, col uint32, text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.lv_table_set_cell_value(t.c, C.uint32_t(row), C.uint32_t(col), cText)
}

// CellValue returns the text of a cell.
func (t *Table) CellValue(row, col uint32) string {
	return C.GoString(C.lv_table_get_cell_value(t.c, C.uint32_t(row), C.uint32_t(col)))
}

// SetCellCtrl sets control flags on a cell (e.g. merge with the next
// column, crop overflowing text).
func (t *Table) SetCellCtrl(row, col uint32, ctrl TableCellCtrl) {
	C.lv_table_set_cell_ctrl(t.c, C.uint32_t(row), C.uint32_t(col), C.lv_table_cell_ctrl_t(ctrl))
}

// ClearCellCtrl clears control flags on a cell.
func (t *Table) ClearCellCtrl(row, col uint32, ctrl TableCellCtrl) {
	C.lv_table_clear_cell_ctrl(t.c, C.uint32_t(row), C.uint32_t(col), C.lv_table_cell_ctrl_t(ctrl))
}

// HasCellCtrl reports whether all the given control flags are set on a cell.
func (t *Table) HasCellCtrl(row, col uint32, ctrl TableCellCtrl) bool {
	return bool(C.lv_table_has_cell_ctrl(t.c, C.uint32_t(row), C.uint32_t(col), C.lv_table_cell_ctrl_t(ctrl)))
}

// SetSelectedCell selects a cell (for keyboard/encoder navigation and
// EventValueChanged reporting).
func (t *Table) SetSelectedCell(row, col uint16) {
	C.lv_table_set_selected_cell(t.c, C.uint16_t(row), C.uint16_t(col))
}

// RowCount returns the number of rows.
func (t *Table) RowCount() uint32 { return uint32(C.lv_table_get_row_count(t.c)) }

// ColumnCount returns the number of columns.
func (t *Table) ColumnCount() uint32 { return uint32(C.lv_table_get_column_count(t.c)) }

// ColumnWidth returns a column's pixel width.
func (t *Table) ColumnWidth(col uint32) int32 {
	return int32(C.lv_table_get_column_width(t.c, C.uint32_t(col)))
}

// SelectedCell returns the currently selected cell's row/column.
func (t *Table) SelectedCell() (row, col uint32) {
	var r, c C.uint32_t
	C.lv_table_get_selected_cell(t.c, &r, &c)
	return uint32(r), uint32(c)
}

package lvgl

/*
#include <lvgl.h>
*/
import "C"

// ChartType mirrors lv_chart_type_t.
type ChartType uint32

var (
	ChartTypeNone    = ChartType(C.LV_CHART_TYPE_NONE)
	ChartTypeLine    = ChartType(C.LV_CHART_TYPE_LINE)
	ChartTypeCurve   = ChartType(C.LV_CHART_TYPE_CURVE)
	ChartTypeBar     = ChartType(C.LV_CHART_TYPE_BAR)
	ChartTypeStacked = ChartType(C.LV_CHART_TYPE_STACKED)
	ChartTypeScatter = ChartType(C.LV_CHART_TYPE_SCATTER)
)

// ChartAxis mirrors lv_chart_axis_t.
type ChartAxis uint32

var (
	ChartAxisPrimaryY   = ChartAxis(C.LV_CHART_AXIS_PRIMARY_Y)
	ChartAxisSecondaryY = ChartAxis(C.LV_CHART_AXIS_SECONDARY_Y)
	ChartAxisPrimaryX   = ChartAxis(C.LV_CHART_AXIS_PRIMARY_X)
	ChartAxisSecondaryX = ChartAxis(C.LV_CHART_AXIS_SECONDARY_X)
)

// ChartUpdateMode mirrors lv_chart_update_mode_t.
type ChartUpdateMode uint32

var (
	ChartUpdateModeShift    = ChartUpdateMode(C.LV_CHART_UPDATE_MODE_SHIFT)
	ChartUpdateModeCircular = ChartUpdateMode(C.LV_CHART_UPDATE_MODE_CIRCULAR)
)

// Chart wraps an lv_chart widget.
type Chart struct{ *Obj }

// NewChart creates a chart as a child of parent.
func NewChart(parent *Obj) *Chart {
	return &Chart{wrapObj(C.lv_chart_create(parent.c))}
}

// SetType sets how series data is rendered (line, bar, scatter, ...).
func (c *Chart) SetType(t ChartType) { C.lv_chart_set_type(c.c, C.lv_chart_type_t(t)) }

// SetPointCount sets how many data points each series holds.
func (c *Chart) SetPointCount(n uint32) { C.lv_chart_set_point_count(c.c, C.uint32_t(n)) }

// PointCount returns how many data points each series holds.
func (c *Chart) PointCount() uint32 { return uint32(C.lv_chart_get_point_count(c.c)) }

// SetAxisRange sets the min/max value range for the given axis.
func (c *Chart) SetAxisRange(axis ChartAxis, min, max int32) {
	C.lv_chart_set_axis_range(c.c, C.lv_chart_axis_t(axis), C.int32_t(min), C.int32_t(max))
}

// SetUpdateMode sets whether new points shift old ones out or overwrite
// circularly.
func (c *Chart) SetUpdateMode(mode ChartUpdateMode) {
	C.lv_chart_set_update_mode(c.c, C.lv_chart_update_mode_t(mode))
}

// SetDivLineCount sets the number of horizontal/vertical division lines.
func (c *Chart) SetDivLineCount(hor, ver uint32) {
	C.lv_chart_set_div_line_count(c.c, C.uint32_t(hor), C.uint32_t(ver))
}

// Refresh forces the chart to redraw from current data.
func (c *Chart) Refresh() { C.lv_chart_refresh(c.c) }

// PressedPoint returns the index of the point currently under the
// press/cursor, or a sentinel LVGL-internal "not found" value.
func (c *Chart) PressedPoint() uint32 { return uint32(C.lv_chart_get_pressed_point(c.c)) }

// ChartSeries wraps an lv_chart_series_t. It is owned by the Chart it was
// added to; there is no separate Delete (use Chart.RemoveSeries).
type ChartSeries struct {
	c     *C.lv_chart_series_t
	chart *Chart
}

// AddSeries adds a new data series with the given line/bar color on the
// given axis.
func (c *Chart) AddSeries(color Color, axis ChartAxis) *ChartSeries {
	s := C.lv_chart_add_series(c.c, color.toC(), C.lv_chart_axis_t(axis))
	if s == nil {
		return nil
	}
	return &ChartSeries{c: s, chart: c}
}

// RemoveSeries removes a data series from the chart.
func (c *Chart) RemoveSeries(s *ChartSeries) {
	C.lv_chart_remove_series(c.c, s.c)
}

// HideSeries shows/hides a series without removing it.
func (c *Chart) HideSeries(s *ChartSeries, hide bool) {
	C.lv_chart_hide_series(c.c, s.c, C.bool(hide))
}

// SetSeriesColor changes a series' color.
func (c *Chart) SetSeriesColor(s *ChartSeries, color Color) {
	C.lv_chart_set_series_color(c.c, s.c, color.toC())
}

// SeriesColor returns a series' current color.
func (c *Chart) SeriesColor(s *ChartSeries) Color {
	return colorFromC(C.lv_chart_get_series_color(c.c, s.c))
}

// SetAllValues sets every point in a series to the same value.
func (c *Chart) SetAllValues(s *ChartSeries, value int32) {
	C.lv_chart_set_all_values(c.c, s.c, C.int32_t(value))
}

// SetNextValue appends a value to the series, shifting out the oldest
// point once PointCount is reached (or overwriting circularly in
// ChartUpdateModeCircular).
func (c *Chart) SetNextValue(s *ChartSeries, value int32) {
	C.lv_chart_set_next_value(c.c, s.c, C.int32_t(value))
}

// SetNextValue2 appends an (x, y) point pair, for ChartTypeScatter series.
func (c *Chart) SetNextValue2(s *ChartSeries, x, y int32) {
	C.lv_chart_set_next_value2(c.c, s.c, C.int32_t(x), C.int32_t(y))
}

// SetSeriesValueByID sets a single point's value by index.
func (c *Chart) SetSeriesValueByID(s *ChartSeries, id uint32, value int32) {
	C.lv_chart_set_series_value_by_id(c.c, s.c, C.uint32_t(id), C.int32_t(value))
}

// SetSeriesValues replaces all of a series' point values at once. len(values)
// should match PointCount.
func (c *Chart) SetSeriesValues(s *ChartSeries, values []int32) {
	if len(values) == 0 {
		return
	}
	C.lv_chart_set_series_values(c.c, s.c, (*C.int32_t)(&values[0]), C.size_t(len(values)))
}

// ChartCursor wraps an lv_chart_cursor_t. It is owned by the Chart it was
// added to; there is no separate Delete (use Chart.RemoveCursor).
type ChartCursor struct {
	c *C.lv_chart_cursor_t
}

// AddCursor adds a crosshair cursor of the given color, free to move
// along the given direction(s) (DirHor/DirVer/DirAll).
func (c *Chart) AddCursor(color Color, dir Dir) *ChartCursor {
	cur := C.lv_chart_add_cursor(c.c, color.toC(), C.lv_dir_t(dir))
	if cur == nil {
		return nil
	}
	return &ChartCursor{c: cur}
}

// RemoveCursor removes a cursor from the chart.
func (c *Chart) RemoveCursor(cur *ChartCursor) { C.lv_chart_remove_cursor(c.c, cur.c) }

// SetCursorPos moves a cursor to the given pixel position within the chart.
func (c *Chart) SetCursorPos(cur *ChartCursor, pos Point) {
	p := C.lv_point_t{x: C.int32_t(pos.X), y: C.int32_t(pos.Y)}
	C.lv_chart_set_cursor_pos(c.c, cur.c, &p)
}

// SetCursorPoint snaps a cursor to a specific data point of a series.
func (c *Chart) SetCursorPoint(cur *ChartCursor, s *ChartSeries, id uint32) {
	C.lv_chart_set_cursor_point(c.c, cur.c, s.c, C.uint32_t(id))
}

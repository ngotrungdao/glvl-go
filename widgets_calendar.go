package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// CalendarDate mirrors lv_calendar_date_t.
type CalendarDate struct {
	Year       uint16
	Month, Day uint8 // Month is 1-12, Day is 1-31
}

func (d CalendarDate) toC() C.lv_calendar_date_t {
	return C.lv_calendar_date_t{year: C.uint16_t(d.Year), month: C.uint8_t(d.Month), day: C.uint8_t(d.Day)}
}

func calendarDateFromC(d C.lv_calendar_date_t) CalendarDate {
	return CalendarDate{Year: uint16(d.year), Month: uint8(d.month), Day: uint8(d.day)}
}

// Calendar wraps an lv_calendar widget.
//
// LVGL keeps a pointer to the array passed to SetHighlightedDates rather
// than copying it, so it's pinned on the C heap like Style/GridDsc,
// freed on Delete via the highlighted field below (there's no separate
// Delete for just the highlighted dates — call SetHighlightedDates(nil)
// to clear and free them).
type Calendar struct {
	*Obj
	highlighted *C.lv_calendar_date_t
}

// NewCalendar creates a calendar as a child of parent.
func NewCalendar(parent *Obj) *Calendar {
	return &Calendar{Obj: wrapObj(C.lv_calendar_create(parent.c))}
}

// SetTodayDate marks the given date as "today".
func (c *Calendar) SetTodayDate(year, month, day uint32) {
	C.lv_calendar_set_today_date(c.c, C.uint32_t(year), C.uint32_t(month), C.uint32_t(day))
}

// SetTodayYear sets just the year of "today".
func (c *Calendar) SetTodayYear(year uint32) { C.lv_calendar_set_today_year(c.c, C.uint32_t(year)) }

// SetTodayMonth sets just the month of "today".
func (c *Calendar) SetTodayMonth(month uint32) {
	C.lv_calendar_set_today_month(c.c, C.uint32_t(month))
}

// SetTodayDay sets just the day of "today".
func (c *Calendar) SetTodayDay(day uint32) { C.lv_calendar_set_today_day(c.c, C.uint32_t(day)) }

// SetMonthShown sets which year/month is currently displayed.
func (c *Calendar) SetMonthShown(year, month uint32) {
	C.lv_calendar_set_month_shown(c.c, C.uint32_t(year), C.uint32_t(month))
}

// SetShownYear sets just the year currently displayed.
func (c *Calendar) SetShownYear(year uint32) { C.lv_calendar_set_shown_year(c.c, C.uint32_t(year)) }

// SetShownMonth sets just the month currently displayed.
func (c *Calendar) SetShownMonth(month uint32) {
	C.lv_calendar_set_shown_month(c.c, C.uint32_t(month))
}

// SetHighlightedDates sets the list of dates to highlight (e.g. dates
// with events). Pass nil to clear.
func (c *Calendar) SetHighlightedDates(dates []CalendarDate) {
	var arr *C.lv_calendar_date_t
	if len(dates) > 0 {
		arr = (*C.lv_calendar_date_t)(C.malloc(C.size_t(len(dates)) * C.size_t(unsafe.Sizeof(C.lv_calendar_date_t{}))))
		buf := unsafe.Slice(arr, len(dates))
		for i, d := range dates {
			buf[i] = d.toC()
		}
	}
	C.lv_calendar_set_highlighted_dates(c.c, arr, C.size_t(len(dates)))

	if c.highlighted != nil {
		C.free(unsafe.Pointer(c.highlighted))
	}
	c.highlighted = arr
}

// AddHeaderArrow adds a header with previous/next month arrow buttons.
func (c *Calendar) AddHeaderArrow() *Obj {
	return wrapObj(C.lv_calendar_add_header_arrow(c.c))
}

// AddHeaderDropdown adds a header with dropdowns to pick the year/month
// directly.
func (c *Calendar) AddHeaderDropdown() *Obj {
	return wrapObj(C.lv_calendar_add_header_dropdown(c.c))
}

// BtnMatrix returns the internal ButtonMatrix used to draw the day grid.
func (c *Calendar) BtnMatrix() *Obj { return wrapObj(C.lv_calendar_get_btnmatrix(c.c)) }

// TodayDate returns the date currently marked as "today".
func (c *Calendar) TodayDate() CalendarDate {
	return calendarDateFromC(*C.lv_calendar_get_today_date(c.c))
}

// ShowedDate returns the year/month currently displayed.
func (c *Calendar) ShowedDate() CalendarDate {
	return calendarDateFromC(*C.lv_calendar_get_showed_date(c.c))
}

// PressedDate returns the date the user last pressed, if any.
func (c *Calendar) PressedDate() (date CalendarDate, ok bool) {
	var d C.lv_calendar_date_t
	if C.lv_calendar_get_pressed_date(c.c, &d) != C.LV_RESULT_OK {
		return CalendarDate{}, false
	}
	return calendarDateFromC(d), true
}

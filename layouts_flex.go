package lvgl

/*
#include <lvgl.h>
*/
import "C"

// FlexFlow mirrors lv_flex_flow_t.
type FlexFlow uint32

var (
	FlexFlowRow               = FlexFlow(C.LV_FLEX_FLOW_ROW)
	FlexFlowColumn            = FlexFlow(C.LV_FLEX_FLOW_COLUMN)
	FlexFlowRowWrap           = FlexFlow(C.LV_FLEX_FLOW_ROW_WRAP)
	FlexFlowRowReverse        = FlexFlow(C.LV_FLEX_FLOW_ROW_REVERSE)
	FlexFlowRowWrapReverse    = FlexFlow(C.LV_FLEX_FLOW_ROW_WRAP_REVERSE)
	FlexFlowColumnWrap        = FlexFlow(C.LV_FLEX_FLOW_COLUMN_WRAP)
	FlexFlowColumnReverse     = FlexFlow(C.LV_FLEX_FLOW_COLUMN_REVERSE)
	FlexFlowColumnWrapReverse = FlexFlow(C.LV_FLEX_FLOW_COLUMN_WRAP_REVERSE)
)

// FlexAlign mirrors lv_flex_align_t.
type FlexAlign uint32

var (
	FlexAlignStart        = FlexAlign(C.LV_FLEX_ALIGN_START)
	FlexAlignEnd          = FlexAlign(C.LV_FLEX_ALIGN_END)
	FlexAlignCenter       = FlexAlign(C.LV_FLEX_ALIGN_CENTER)
	FlexAlignSpaceEvenly  = FlexAlign(C.LV_FLEX_ALIGN_SPACE_EVENLY)
	FlexAlignSpaceAround  = FlexAlign(C.LV_FLEX_ALIGN_SPACE_AROUND)
	FlexAlignSpaceBetween = FlexAlign(C.LV_FLEX_ALIGN_SPACE_BETWEEN)
)

// SetFlexFlow makes the object lay out its children using flexbox, in the
// given flow direction.
func (o *Obj) SetFlexFlow(flow FlexFlow) {
	C.lv_obj_set_flex_flow(o.c, C.lv_flex_flow_t(flow))
}

// SetFlexAlign sets how children are aligned along the main axis, cross
// axis, and (when wrapping) across tracks.
func (o *Obj) SetFlexAlign(mainPlace, crossPlace, trackCrossPlace FlexAlign) {
	C.lv_obj_set_flex_align(o.c, C.lv_flex_align_t(mainPlace), C.lv_flex_align_t(crossPlace),
		C.lv_flex_align_t(trackCrossPlace))
}

// SetFlexGrow sets how much this object grows relative to its flex
// siblings when there is extra space along the main axis.
func (o *Obj) SetFlexGrow(grow uint8) {
	C.lv_obj_set_flex_grow(o.c, C.uint8_t(grow))
}

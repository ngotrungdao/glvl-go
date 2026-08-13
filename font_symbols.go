package lvgl

// Symbol* are a handful of LVGL's built-in icon-font glyphs
// (font/lv_symbol_def.h's LV_SYMBOL_* macros), usable anywhere a widget
// accepts a symbol string (e.g. Win.AddButton, List.AddButton,
// ImageButton.SetSrc). These are reimplemented as plain Go string
// constants from LVGL's UTF-8 byte sequences, since preprocessor string
// macros aren't accessible via cgo.
const (
	SymbolOK        = "\xEF\x80\x8C"
	SymbolClose     = "\xEF\x80\x8D"
	SymbolSettings  = "\xEF\x80\x93"
	SymbolHome      = "\xEF\x80\x95"
	SymbolRefresh   = "\xEF\x80\xA1"
	SymbolLeft      = "\xEF\x81\x93"
	SymbolRight     = "\xEF\x81\x94"
	SymbolPlus      = "\xEF\x81\xA7"
	SymbolMinus     = "\xEF\x81\xA8"
	SymbolWarning   = "\xEF\x81\xB1"
	SymbolUp        = "\xEF\x81\xB7"
	SymbolDown      = "\xEF\x81\xB8"
	SymbolSave      = "\xEF\x83\x87"
	SymbolBell      = "\xEF\x83\xB3"
	SymbolTrash     = "\xEF\x8B\xAD"
	SymbolEdit      = "\xEF\x8C\x84"
	SymbolBackspace = "\xEF\x95\x9A"
)

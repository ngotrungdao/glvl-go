#----------------------------------------------------------------
# Generated CMake target import file.
#----------------------------------------------------------------

# Commands may need to know the format version.
set(CMAKE_IMPORT_FILE_VERSION 1)

# Import target "lvgl::lvgl" for configuration ""
set_property(TARGET lvgl::lvgl APPEND PROPERTY IMPORTED_CONFIGURATIONS NOCONFIG)
set_target_properties(lvgl::lvgl PROPERTIES
  IMPORTED_LINK_INTERFACE_LANGUAGES_NOCONFIG "ASM;C;CXX"
  IMPORTED_LOCATION_NOCONFIG "${_IMPORT_PREFIX}/lib64/liblvgl.a"
  )

list(APPEND _cmake_import_check_targets lvgl::lvgl )
list(APPEND _cmake_import_check_files_for_lvgl::lvgl "${_IMPORT_PREFIX}/lib64/liblvgl.a" )

# Commands beyond this point should not need to know the version.
set(CMAKE_IMPORT_FILE_VERSION)

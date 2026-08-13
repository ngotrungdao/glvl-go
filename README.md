# lvgl (Go)

A hand-written Go/cgo wrapper around [LVGL](https://lvgl.io/) 9.6, targeting the **Wayland**
display backend. It covers the core object model, the event system, styles, flex/grid
layout, animation, and about twenty of the most common widgets. It is a broad, practical
wrapper, not a complete 1:1 binding of LVGL's ~216 headers.

## Requirements

- A prebuilt LVGL 9.6 static library and headers. This package expects them under
  `lvgl-c/install/` (a plain directory living inside this repo, built via CMake into
  `lvgl-c/build/` — not an external symlink). Adjust the paths in `cgo.go` if yours live
  elsewhere. This build has `LV_USE_THORVG=0`, so no `libstdc++`/ThorVG dependency is needed.
- Wayland client dev libraries: `wayland-client`, `wayland-cursor`, `wayland-egl`,
  `wayland-protocols`, `libxkbcommon` (Fedora: `wayland-devel`, `wayland-protocols-devel`,
  `libxkbcommon-devel` — `wayland-egl` ships as part of `wayland-devel`). `wayland-egl` is a
  genuine runtime dependency for this `lvgl-c` build, not just a link-time one: it has
  `LV_USE_DRAW_NANOVG=1`/`LV_USE_DRAW_SW=0` and `LV_WAYLAND_USE_SHM=0`, so LVGL rasterizes with
  NanoVG on OpenGL ES and presents frames to the compositor as EGL-backed `wl_buffer`s rather
  than plain SHM — confirmed with `WAYLAND_DEBUG=1`, which shows the window's `wl_surface`
  attaching buffers from a "mesa egl surface queue". (An earlier LVGL build used underneath
  this wrapper rendered via plain SHM software rasterization instead; that is no longer the
  case after switching to the `lvgl-c` build.)
- FreeType, libjpeg, libpng, libwebp dev packages, even though this wrapper only uses the
  Wayland backend: the prebuilt `liblvgl.a` has unconditional references to these libraries'
  symbols from its image-decoder modules (Fedora: `freetype-devel`, `libjpeg-turbo-devel`,
  `libpng-devel`, `libwebp-devel`).
- Go 1.21+ (for `runtime/cgo.Handle`).

## File layout

This is a single Go package (`package lvgl`) — Go can't split one package across
subdirectories, and cgo's `C` pseudo-type isn't shareable across real Go sub-packages without
`unsafe.Pointer` round-trips at every cross-package call, so files are grouped by filename
prefix instead, mirroring `install/include/lvgl`'s directory layout:

| Prefix | Mirrors | Contents |
|---|---|---|
| `core_*.go` | mostly `core/` | object model, events, animation, styles, tick/timer, init, scroll, tree walk, timers, group, observer/subject (`core_fs.go`/`core_indev.go` are the two exceptions, mirroring `fs/`/`indev/` instead — grouped under `core_` by theme, not exact header directory) |
| `draw_color.go` | `draw/` | `Color` |
| `layouts_*.go` | `layouts/` | flex, grid |
| `display_wayland.go` | `display/` + `drivers/wayland/` | `Display`, window creation |
| `font_symbols.go` | `font/` | `Symbol*` icon-font glyph constants |
| `font_font.go` | `font/` + `libs/freetype/` + `libs/tiny_ttf/` | `Font`, built-in/FreeType/TinyTTF loading |
| `widgets_*.go` | `widgets/` | one file per LVGL widget |

`cgo.go`, `app.go`, `handle.go` are package-level glue (build flags, the main loop, the
`cgo.Handle` helper) with no single LVGL header counterpart.

## Quick start

```go
package main

import "lvgl"

func main() {
	lvgl.Init()

	disp := lvgl.WaylandWindowCreate(480, 320, "my app")
	label := lvgl.NewLabel(disp.ScreenActive())
	label.SetText("Hello!")
	label.Center()

	lvgl.Run(disp) // blocks, driving LVGL's tick/timer loop until the window closes
}
```

Run the bundled examples:

```
go run ./example/basic    # a few interactive widgets
go run ./example/gallery  # every wrapped widget on one screen
```

## Memory ownership

- **`Obj`** (and anything embedding it, i.e. every widget type) is owned by LVGL. Call
  `Delete()` to free it and its children, same as `lv_obj_delete`.
- **`Style`** and **`GridDsc`** are different: LVGL keeps a raw, non-refcounted pointer to
  their C memory for as long as they're attached to an object, so they're allocated on the C
  heap rather than as ordinary Go memory. **Call `Delete()` on them explicitly** once they're
  no longer attached to anything (a style or grid descriptor reused across many objects only
  needs to be deleted once, after the last object using it is gone). A missed `Delete()` logs
  a one-time "finalized without an explicit Delete()" warning as a leak diagnostic — it does
  not crash, but the C memory is leaked.
- Event and animation callbacks are tracked per-`Obj` and freed automatically when the object
  is deleted (via an internal `LV_EVENT_DELETE` hook) or when an animation completes/is
  deleted — no manual cleanup needed for those.

## Event system

```go
btn.AddEventCB(lvgl.EventClicked, func(e *lvgl.Event) {
	// e.Target(), e.Code(), ...
})
```

Internally this uses a single `//export`ed C trampoline per callback kind (event, animation)
combined with `runtime/cgo.Handle`, since C function pointers can't be Go closures. See
`core_event.go` and `core_anim.go`.

## Styles

```go
s := lvgl.NewStyle()
s.SetBgColor(lvgl.ColorHex(0x2b2f3a))
s.SetRadius(8)
obj.AddStyle(s, lvgl.Selector(lvgl.PartMain))
// ... later, once detached from every object:
s.Delete()
```

`core_style.go`/`core_style_obj.go` wrap 123 (`Style`) / 122 (`Obj.SetStyle*`) of LVGL's 129
`lv_style_set_*` properties — essentially everything except 5 that take pointer types needing
their own dedicated supporting API (`lv_anim_t*`, `lv_style_transition_dsc_t*`,
`lv_grad_dsc_t*`, `lv_color_filter_dsc_t*`, `lv_image_colorkey_t*`), a low-level raw layout-ID
setter normally driven automatically by `SetFlexFlow`/`SetGridDscArray`, and (on `Obj` only)
the style-based grid track descriptor setters, since `Obj.SetGridDscArray` already covers the
same need via the more common direct API. Image-source/bitmap-mask string properties (e.g.
`SetBgImageSrc`) follow the same C-heap-pinning rule as `Image.SetSrcPath`. To add the
remaining ones, follow the same one-line pattern; each existing setter is a direct, mechanical
translation of the matching `core/lv_style_gen.h` / `core/lv_obj_style_gen.h` signature.

## Fonts

Text color is a style property (`SetTextColor`/`SetStyleTextColor`, above). Font family *and*
size together are a separate `Font` value passed to `SetTextFont`/`SetStyleTextFont` — LVGL has
no separate "set font size" call, size is baked into which font object you use:

```go
label.SetStyleTextFont(lvgl.FontMontserrat24, lvgl.Selector(lvgl.PartMain))
label.SetStyleTextColor(lvgl.ColorHex(0xff8800), lvgl.Selector(lvgl.PartMain))
```

`font_font.go` wraps three ways to get a `Font`:

- **Built-in bitmap fonts**: `FontMontserrat8` through `FontMontserrat48` (every even size in
  that range enabled in this LVGL build) — no `Delete()` needed, these are static library data.
- **`FreeTypeFontCreate(path, renderMode, size, style)`**: loads a `.ttf`/`.otf` file at any
  pixel size via FreeType. Call `FreeTypeInit(maxGlyphCnt)` once first (LVGL may report it as
  already initialized internally — treat that as non-fatal and try creating a font anyway, see
  the doc comment on `FreeTypeInit`).
- **`TinyTTFCreateFile`/`TinyTTFCreateData`**: loads a `.ttf` at any pixel size via LVGL's own
  lightweight TrueType renderer, no FreeType/init step needed.

Fonts loaded via FreeType/TinyTTF own C-heap resources and must be `Delete()`d once no longer
set on any style/object, same rule as `Style`/`GridDsc`.

## Layout

Flex (`layouts_flex.go`) and grid (`layouts_grid.go`) are both wrapped. Grid track descriptor
arrays follow the same C-heap-pinning rule as `Style` — see `GridDsc` and its `Delete()`.

## Core subsystems

Beyond the object model/events/styles/layout/animation above, these previously-unwrapped LVGL
subsystems are now covered:

- **Scroll** (`core_scroll.go`) — `Obj.ScrollTo`/`ScrollBy`/`SetScrollDir`/`SetScrollSnapX`/`Y`,
  scroll position/extent getters.
- **Tree walk** (`core_tree.go`) — `Obj.TreeWalk(fn)` visits a subtree depth-first via a
  one-shot `cgo.Handle` callback (freed right after the call, unlike the persistent
  per-object handles events use).
- **Custom timers** (`core_timer.go`) — `NewTimer(period, fn)` schedules a callback through
  LVGL's own timer loop. Unlike `Anim`, `lv_timer_t` has no "deleted" callback hook, so a
  timer with a finite `SetRepeatCount` and `SetAutoDelete` left on its default leaks its
  `cgo.Handle` when LVGL auto-deletes it — prefer `SetAutoDelete(false)` + explicit `Delete()`
  for anything long-running. **Never call `Delete()` yourself while `SetAutoDelete` is also on
  with a finite repeat count** — the two race and it was observed to hang the process outright
  rather than crash cleanly.
- **Animation easing** (`Anim.SetPath` in `core_anim.go`) — selects one of LVGL's built-in
  easing curves (`AnimPathLinear`, `EaseIn`, `EaseOut`, `EaseInOut`, `Overshoot`, `Bounce`,
  `Step`) via a small C shim, since cgo can't take the address of a named C function directly.
- **Group + Indev** (`core_group.go`, `core_indev.go`) — `Group` for keypad/encoder focus
  navigation, `Indev` for direct input-device access (cursor, key/point state, scroll/gesture
  direction). `Display.Pointer()`/`Keyboard()`/`Touchscreen()` expose the Wayland driver's
  input devices to attach a `Group` to.
- **Observer/Subject data binding** (`core_observer.go`) — LVGL 9's reactive binding system:
  `Subject` (int/float/color/string values) plus `AddObserver` for custom callbacks, and
  `BindValue`/`BindText`/`BindChecked` on `Slider`/`Bar`/`Arc`/`Roller`/`SpinBox`/`Dropdown`/
  `Label`/`Obj` (checkbox/switch) to wire a widget to a `Subject` with no manual event
  plumbing. `Subject` follows the same C-heap-pinning/`Delete()` rule as `Style`.
- **Translation** (`core_translation.go`) — `SetLanguage`/`Language`/`Translate` (the
  read/use side only; registering translation packs themselves needs `lv_translation_pack_t`
  C data structures this pass didn't build supporting API for).
- **Filesystem** (`core_fs.go`) — `FSOpen`/`File.Read`/`Write`/`Seek` (satisfies `io.Reader`/
  `io.Writer`) and `FSDirOpen`/`FSDir.Read` through LVGL's virtual FS layer (works with any
  path using a registered drive letter, e.g. `"A:/tmp/file.txt"` with this build's stdio
  driver). Registering a *custom* FS driver isn't wrapped (a bigger, separate API).

Not covered: **vector graphics/SVG drawing** (`draw/lv_draw_vector.h`, 572 lines) — a
genuinely separate, large 2D drawing API (matrices, paths, gradients, layers), comparable in
scope to wrapping a mini Cairo/Skia; deliberately out of scope for this pass rather than
wrapped shallowly.

## Widgets

39 of LVGL's 45 `widgets/lv_*.h` headers are covered: 37 get their own `widgets_*.go` file —
Label, Button, Slider, Checkbox, Switch, Dropdown, Roller, TextArea, Arc, ArcLabel, Bar, Image,
ImageButton, AnimImage, Canvas, Gif, List, Table, TabView, MsgBox, Spinner, SpinBox, Chart,
Calendar, Win, Keyboard, ButtonMatrix, Led, Line, TileView, Scale, SpanGroup, QRCode, Barcode,
GStreamer, Texture3D, FFmpeg — and 2 more (`lv_calendar_header_arrow.h`, `lv_calendar_header_dropdown.h`)
are covered as `Calendar.AddHeaderArrow()`/`AddHeaderDropdown()` methods rather than separate
types, since they're addons to an existing calendar, not standalone widgets.

`Texture3D` (`lv_3dtexture`) displays an OpenGL texture the host app renders itself elsewhere
(via whatever Go OpenGL binding it uses) — LVGL doesn't manage GL context/texture creation for
it, just draws the finished texture through its own renderer, so it needed no new
dependencies and links cleanly (verified with `nm` before writing any code, same discipline as
GStreamer/Lottie).

Coverage per widget was audited against its full header and extended well beyond "most
common" — most wrapped widgets now expose close to their entire practical setter/getter
surface (~340 of ~400 detected header functions, up from ~100). Deliberately skipped:
deprecated `_fmt`/`_vfmt` varargs text setters (use Go string formatting before calling the
plain setter instead), translation-tag APIs (`LV_USE_TRANSLATION` integration), and a few
advanced/rare entry points (e.g. `ButtonMatrix`'s `ctrl_map` array setter, `Canvas`'s raw
draw-layer functions, `Scale`'s `text_src` array, `Label`'s letter-position hit-testing).

`GStreamer` needed extra system libraries beyond what the rest of the wrapper requires:
`gstreamer1-devel`, `gstreamer1-plugins-base-devel`, and `glib2-devel` (Fedora names) —
`-lgstvideo-1.0 -lgstbase-1.0 -lgstreamer-1.0 -lgobject-2.0 -lglib-2.0` were added to
`cgo.go`'s LDFLAGS accordingly (get the flags for your own system via
`pkg-config --libs gstreamer-1.0 gstreamer-video-1.0`). See Known limitations below before
relying on `GStreamer.Play()`.

`FFmpeg` (`widgets_ffmpeg.go`) similarly needs `ffmpeg-free-devel` (Fedora; pulls in
`libavformat-free-devel`/`libavcodec-free-devel`/`libavutil-free-devel`/`libswscale-free-devel`,
no RPM Fusion needed) — `-lavformat -lavcodec -lavutil -lswscale` were added to `cgo.go`'s
LDFLAGS accordingly. Call `FFmpegInit()` once before `FFmpegGetFrameNum`/`NewFFmpeg`, same
pattern as `FreeTypeInit`. See Known limitations below before playing a non-square video.

`lv_list` and `lv_win` are deprecated upstream in favor of building the same layout from a
flex column; they're still wrapped (as `List`/`Win`) since they work, but prefer
`SetFlexFlow(FlexFlowColumn)` for new code. `lv_menu` is deprecated more strongly (LVGL
recommends building menu navigation entirely from base widgets) and isn't wrapped at all.

The remaining 6 headers, and why:

| Header | Status | Why not wrapped |
|---|---|---|
| `lv_gltf.h` | enabled in config, but **not actually linkable** | its object files (`lv_gltf_view.cpp.o` and related) reference ~50 undefined `glad_gl*` OpenGL ES loader symbols, confirmed absent from every `.a` in this build tree (`liblvgl.a`, `liblvgl_thorvg.a`, checked via `nm`/`ar p`). `glad` is normally vendored generated source, not an installable system package, so — unlike Wayland/FreeType/GStreamer — there's no `-devel` package that fixes this; it needs LVGL rebuilt with the glad loader source included. |
| `lv_menu.h` | enabled, linkable | deprecated upstream, see above |
| `lv_calendar_chinese.h` | **disabled** (`LV_USE_CALENDAR_CHINESE`) | no symbol in `liblvgl.a`; needs LVGL rebuilt with it on |
| `lv_ime_pinyin.h` | enabled, linkable (`lv_ime_pinyin_create` etc. confirmed via `nm`) | not wrapped by choice — `lv_ime_pinyin_set_dict` takes a raw `lv_pinyin_dict_t*` array of `{py, py_mb}` C string pairs that LVGL keeps a live pointer to, needing the same kind of C-heap-pinned array type as `GridDsc`/`Style`; deferred rather than wrapped shallowly |
| `lv_lottie.h` | **disabled** (architecturally, not just a config flag) | `LV_USE_LOTTIE` depends on `LV_DRAW_HAS_VECTOR_SUPPORT`, which `LV_USE_THORVG` only grants `if LV_USE_DRAW_SW` — this build uses NanoVG with `LV_USE_DRAW_SW` off, so the dependency can't be satisfied without switching the whole draw engine back to SW |
| `lv_rlottie.h` | **disabled** | no symbol in `liblvgl.a`; unlike FFmpeg, this LVGL tree has no CMake logic anywhere that links `-lrlottie` even when the Kconfig flag is set, so enabling it would compile but leave dangling symbols |

"disabled" here was verified with `nm liblvgl.a | grep <symbol>`, not just by reading the
`LV_USE_*` config — an earlier pass wrongly wrapped `Lottie` based on a misread config value;
it didn't actually have a linkable symbol and was removed.

## Known limitations

- **Rendering goes through NanoVG on OpenGL ES, not the software rasterizer** — this `lvgl-c`
  build has `LV_USE_DRAW_SW=0`, so there's no plain-CPU fallback path; every widget draw and
  the window's own Wayland presentation depend on a working EGL/GLES context. Verified working
  end-to-end (`go run ./example/gallery` renders every widget correctly, confirmed visually)
  but this is why the `fbo_create_cb: Failed to create FBO` errors mentioned below aren't a
  one-off — they're this build's normal draw engine hitting its cache limits, not a rarely-used
  side path.
- **`FFmpeg` fails to decode any non-square video** — confirmed as a real bug in LVGL's own
  `lv_ffmpeg.c`, not the Go wrapper: a 320×320 test file (`ffmpeg -f lavfi -i
  testsrc=size=320x320 ...`) plays correctly end-to-end (verified visually, real decoded frames
  in a live window), while the identical test with `size=320x240` fails on every single frame
  with `ffmpeg_output_video_frame: Width, height and pixel format have to be constant in a
  rawvideo file, but ... changed: old: width = 240, height = 320 ... new: width = 320, height =
  240` — note the exact width/height *transpose* between "old" and "new", on 100% of frames,
  not intermittently. This points at LVGL's internal rawvideo pipe initializing with width/height
  swapped for non-square input, not anything file- or codec-specific. Stick to square sources
  (or pre-pad to square) until this is fixed upstream.
- **`GStreamer.Play()` hangs indefinitely specifically for a `GStreamerFactoryTestVideo`
  (`videotestsrc`) source** — confirmed NOT a general GStreamer/widget problem: playing a real
  video file (`GStreamerFactoryFile`) works correctly end-to-end, verified visually (actual
  decoded video frames render into the widget in a live window) and via state/position/duration
  polling (state reaches `GStreamerStatePlaying`, duration/position report real, advancing
  values). The `videotestsrc`-specific hang wasn't narrowed down further (suspected caps
  negotiation quirk between `videotestsrc`'s raw output and the widget's internal
  `videoconvert → appsink` chain, but no `.c` source was available to confirm) — avoid
  `videotestsrc` as a source until that's understood; real file/URI sources are fine.
- **`TimerHandler()` hangs (or badly stalls) once ~60+ widgets have accumulated on one screen**
  — confirmed *not* a bug in `Timer`/`Anim`/the wrapper: an isolated minimal repro (one timer,
  one fresh screen) passed 100% across 8+ stress runs, while the same checks embedded in
  `example/headless`'s heavy, ever-growing scene (80+ widgets by the time it would reach them)
  hung on the majority of runs. A logged `lv_image_src_get_type: ... invalid magic` error at
  the point things went bad suggests an interaction with this environment's already-flaky
  NanoVG image/FBO cache (see the QRCode/Barcode limitation below), not anything specific to
  timers — circumstantial, not confirmed. **Do not "fix" a hang like this with a goroutine +
  timeout that abandons the stuck call**: an earlier pass here tried exactly that, and it
  introduced a worse, silent bug — LVGL is not thread-safe, and the abandoned goroutine kept
  calling `TimerHandler()` in the background, racing the main goroutine's subsequent calls and
  tripping `lv_inv_area: ... rendering_in_progress` assertions. The actual fix was isolating
  `Timer`/`Anim`-easing/`Group`/`Indev`/`Observer`-`Subject` checks into their own lightweight
  program (`example/coresubsystems`) against a fresh, small scene, where they're fully
  reliable — real apps should keep this in mind if a single screen accumulates a very large
  number of live widgets over time in this environment.
- Wayland-only: no SDL2/DRM/X11 display backend wrapper (LVGL supports them, this package
  doesn't wire them up).
- **QRCode and Barcode fail to render** with this prebuilt `liblvgl.a`: LVGL logs
  `data size (...) is larger than max size (4096)` / `Failed to open image` at draw time. This
  is `LV_CACHE_DEF_SIZE` being configured too small in the library build to hold a rendered
  QR/barcode bitmap — a build-config issue in the prebuilt static library, not something this
  Go wrapper can work around. Fixing it means rebuilding LVGL with a larger
  `LV_CACHE_DEF_SIZE`. Everything else (including Line, which is vector-drawn rather than
  image-cached) renders fine.
- **Not wrapped, genuinely disabled in this build** (`LV_USE_* = 0`, confirmed by checking for
  their symbols with `nm` on `liblvgl.a` — the header prototypes are preprocessed out
  entirely, so wrapping them wouldn't even compile): FFmpeg, IME Pinyin, RLottie, and the
  built-in (non-rlottie) Lottie widget. These would need LVGL itself rebuilt with those
  options on, not a Go-side fix.
- **Not wrapped, enabled in config but not actually linkable**: `lv_gltf_*` is compiled into
  `liblvgl.a`'s object files, but they reference ~50 undefined `glad_gl*` OpenGL ES loader
  symbols not present anywhere in this build tree — see the widgets table above. This is a
  genuine build-config gap, not a scope decision like `lv_menu`. GStreamer
  (`lv_gstreamer_create`) is also compiled in and *is* wrapped (see above) — it needed
  `gstreamer-1.0-devel` and friends installed to actually link once called, same class of
  dependency gap as the earlier Wayland/FreeType issues, but resolvable via package install
  rather than an LVGL rebuild.
- `go test` doesn't work for this package in some Go toolchain builds (e.g. ones built with
  `GOEXPERIMENT=nodwarf5`) — cgo test support is disabled entirely
  (`use of cgo in test ... not supported`), not something this package can work around.
  `example/headless/main.go` exists as a manual regression check via `go run` instead, and
  should be extended alongside new functionality rather than added as `_test.go` files if
  your toolchain has the same limitation.

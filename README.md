# ux
Go user experience framework for macOS (and at some point in the future,
Windows and Linux).

An example application can be found in the `example` directory.

#### Notes
- Windows and Linux platform support is not currently functional. Windows can
  be created but drawing is largely broken or non-existent.
- All target platforms can be built from macOS. Unfortunately, due to the need
  for cgo on macOS, the macOS target platform cannot be built on other
  platforms.

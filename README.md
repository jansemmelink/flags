# flags
This library parses flag sets similar to the built-in Go library, but a bit different and more specific to what I need:
* Options can be defined with both short (e.g. -d) or long (--word) formats
* Short form value must be written after a space (e.g. -n 10)
* Long form value must be written with an equal sign (e.g. --limit=10)

The library currently supports only a few basic types (bool, int, string) but can be extended easily.

Some goals to be supported soon:
* Partial parsing (parse known common options, then parse the remaining flags to another flag set). This is useful when you write several programs that should support a common set of flags as well as their own set, which may not be defined initially.
* Support for validation functions
* Support for lose standing items, typically at the end, but not necessarily
* Support for sequential items

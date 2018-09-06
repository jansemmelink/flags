# Purpose
Command line style flag parsing, similar to built-in Go "flags", but different:
* Enforces short form -n <value> and/or long form --limit=<value>
* Slightly different usage format
* Groups of flags, e.g. in different libraries
* Support for flags defined after initial parsing

# Overview
This library parses flag sets similar to the built-in Go library, but a bit different and more specific to what I need, as listed above. The library is still under construction, but tagged versions will work. I use git flow for development.

# Already Supported:
* short and/or long format options
* bool, int and string flags

# Soon to be supported:
* value validation functions
* grouping of flags
* lose items and lists on the command line
* sequential items

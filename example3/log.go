package main

import "github.com/jansemmelink/flags"

var logFlags *flags.Set

func init() {
	logFlags = flags.NewSet("logging", "Control log output")
	logFlags.Bool("-d", "--debug", false, "Run in DEBUG mode")
	logFlags.String("", "--logfile", "", "Output log to this file instead of stderr")

	flags.AddSet(*logFlags)
}

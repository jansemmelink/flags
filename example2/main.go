//Package main shows how one uses "github.com/jansemmelink/flags" to do
//flag parsing in stages, e.g. when all your images should support a subset
//of flags, those can be parsed up front and other flags can be parsed
//after that when you figured out which other flags are required.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jansemmelink/flags"
)

//say we have a common set of flags to control log level and output
var logFlags = flags.NewSet("logging", "Control log output")

func init() {
	logFlags.Bool("-d", "--debug", false, "Run in DEBUG mode")
	logFlags.String("", "--logfile", "", "Output log to this file instead of stderr")
}

func main() {
	args := os.Args[1:]
	var err error

	//parse only known arguments then replace args with what remains
	if args, err = logFlags.ParseKnown(args); err != nil {
		logFlags.PrintUsage(os.Stderr)
		panic(fmt.Sprintf("Failed to parse log flags: %v", err))
	}

	//show what we got so far
	if logFlags.Flag("-d").Value().(bool) {
		log.Printf("DEBUG Mode\n")
	} else {
		log.Printf("NOT DEBUG Mode\n")
	}

	logFilename := logFlags.Flag("--logfile").Value().(string)
	if logFilename != "" {
		log.Printf("Logfile = %s\n", logFilename)
	} else {
		log.Printf("Log to terminal stderr\n")
	}

	//show what args remained
	if len(args) > 0 {
		log.Printf("Remaining arguments: %v\n", args)
	} else {
		log.Printf("No remaining arguments\n")
	}
} //main()

package flags

import (
	"fmt"
	"os"
	"path"
)

var (
	defaultSet = NewSet("", "")
)

//Bool in the default set
func Bool(short, long string, init bool, doc string) *FlagDescription {
	newFlagPtr, err := defaultSet.Bool(short, long, init, doc)
	if err != nil {
		panic(fmt.Sprintf("Failed to define flag: %v", err))
	}
	return newFlagPtr
} //Bool()

//Int in the default set
func Int(short, long string, init int, doc string) *FlagDescription {
	newFlagPtr, err := defaultSet.Int(short, long, init, doc)
	if err != nil {
		panic(fmt.Sprintf("Failed to define flag: %v", err))
	}
	return newFlagPtr
} //Int()

//String in the default set
func String(short, long string, init string, doc string) *FlagDescription {
	newFlagPtr, err := defaultSet.String(short, long, init, doc)
	if err != nil {
		panic(fmt.Sprintf("Failed to define flag: %v", err))
	}
	return newFlagPtr
} //String()

//DefaultSet to get read access to the default set
func DefaultSet() Set {
	return *defaultSet
}

//Usage writes program usage for the default set to stderr
func Usage(errorMsg string) {
	if errorMsg != "" {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", errorMsg)
	}
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", path.Base(os.Args[0]))
	defaultSet.PrintUsage(os.Stderr)
	os.Exit(-1)
} //Usage()

//Parse the default set of command line options
func Parse() {
	//Args[0] is the program executable, start from 1
	//progName := path.Base(os.Args[0])
	if os.Args == nil || len(os.Args) < 1 {
		panic("Cannot access program arguments")
	}
	//if "?" is specified or --help, display usage info without an error
	for _, opt := range os.Args[1:] {
		if opt == "?" || opt == "--help" {
			Usage("")
		}
	}

	//do normal flag set parsing and fail with usage screen and exit code 1 on error
	if err := defaultSet.Parse(os.Args[1:]); err != nil {
		Usage(err.Error())
	}
} //Parse()

//Flag to get a named flag by short/long option
//Use it e.g. like this:   if flags.GetFlag("-d").GetValue().(bool) { ... defbug is on ... }
func Flag(n string) FlagDescription {
	return defaultSet.Flag(n)
}

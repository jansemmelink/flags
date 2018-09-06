//Package main shows how one uses "github.com/jansemmelink/flags" to do
//flag parsing in stages, e.g. when all your images should support a subset
//of flags, those can be parsed up front and other flags can be parsed
//after that when you figured out which other flags are required.
//
// Example:
//----------------------------------------------------------------------------
//	./example2 -g -d true -h --logfile=/tmp/s -o upd -n Joe -v Sam
//	2018/09/06 12:04:51 DEBUG Mode
//	2018/09/06 12:04:51 Logfile = /tmp/s
//	2018/09/06 12:04:51 Remaining arguments: [-g -h -o upd -n Joe -v Sam]
//	2018/09/06 12:04:51 Operation="upd"
//	2018/09/06 12:04:51 Remaining arguments: [-g -h -n Joe -v Sam]
//	2018/09/06 12:04:51 DOING oper=upd with flag{-n:Joe}flag{-v:Sam}
//	2018/09/06 12:04:51 Remaining arguments: [-g -h]
//----------------------------------------------------------------------------
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

	//now say wait to determine the operation from -o <oper> or --oper=<oper>
	//then continue parsing operation specified options, we can do this:

	operSpecificFlags := make(map[string]*flags.Set)

	operSpecificFlags["add"] = flags.NewSet("add", "Add")
	operSpecificFlags["add"].String("-n", "--name", "", "Name to add")

	operSpecificFlags["get"] = flags.NewSet("get", "Get")

	operSpecificFlags["upd"] = flags.NewSet("upd", "Update")
	operSpecificFlags["upd"].String("-n", "--name", "", "Name to update")
	operSpecificFlags["upd"].String("-v", "--value", "", "New Value to use")

	operSpecificFlags["del"] = flags.NewSet("del", "Delete")
	operSpecificFlags["del"].String("-n", "--name", "", "Name to delete")

	operNames := []string{}
	for n := range operSpecificFlags {
		operNames = append(operNames, n)
	}

	operSelectorSet := flags.NewSet("operation", "Operation")
	operSelectorSet.String("-o", "--oper", "", fmt.Sprintf("Select operation to do %v", operNames))

	args, err = operSelectorSet.ParseKnown(args)
	if err != nil {
		operSelectorSet.PrintUsage(os.Stderr)
		panic(fmt.Sprintf("Failed to parse operation flags: %v", err))
	}

	operName := operSelectorSet.Flag("-o").Value().(string)
	log.Printf("Operation=\"%s\"", operName)

	//show what args remained
	if len(args) > 0 {
		log.Printf("Remaining arguments: %v\n", args)
	} else {
		log.Printf("No remaining arguments\n")
	}

	operSpecificSet, ok := operSpecificFlags[operName]
	if !ok {
		log.Printf("ERROR: Unknown operation \"%s\", expected %v", operName, operNames)
	}

	args, err = operSpecificSet.ParseKnown(args)
	if err != nil {
		operSpecificSet.PrintUsage(os.Stderr)
		panic(fmt.Sprintf("Failed to parse operation specific flags: %v", err))
	}

	log.Printf("DOING oper=%s with %+v", operName, operSpecificSet)

	//show what args remained
	if len(args) > 0 {
		log.Printf("Remaining arguments: %v\n", args)
	} else {
		log.Printf("No remaining arguments\n")
	}

	//NOTE: The usage in each case shows only that set, not the complete set... solved in next examples.
} //main()

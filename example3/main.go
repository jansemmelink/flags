//Package main here shows the same set of log and operation options
//as in example2, but much shorter to write and with complete usage info
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jansemmelink/flags"
)

func main() {
	//operations were registered in init() functions of each module
	//and options were added to the default flag set, so we can just parse
	//and see usage with '?' option
	log.Printf("Registered operations: %v\n", opers)

	//parse command line: we make a copy of default flag set, which will have log added to it
	//then add all the registered operations to it:
	flagSet := addOpersToFlagSet(flags.DefaultSet())
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		panic(fmt.Sprintf("Failed to parse: %v", err))
	}

	log.Printf("Success\n")
} //main()

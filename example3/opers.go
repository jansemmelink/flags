package main

import (
	"fmt"

	"github.com/jansemmelink/flags"
)

type oper struct {
	flagset *flags.Set
}

var (
	opers = make(map[string]oper, 0)
)

func addOper(name string, flagset *flags.Set) {
	if name == "" {
		panic(fmt.Sprintf("Cannot add operation without a name"))
	}
	if _, ok := opers[name]; ok {
		panic(fmt.Sprintf("Duplicate operation \"%s\" already registered", name))
	}
	opers[name] = oper{flagset}

	//if has flags, add to the default flags for the command line
	/*	if flagset != nil {
		flags.AddSet(*flagset)
	}*/
}

func addOpersToFlagSet(set flags.Set) *flags.Set {
	operFlagPtr, err := set.Group("-o", "--oper", "Select operation")
	if err != nil {
		panic(fmt.Sprintf("Failed to add oper selector to flag set: %v", err))
	}
	for _, oper := range opers {
		if err := operFlagPtr.Add(oper.flagset); err != nil {
			panic(fmt.Sprintf("Failed to add oper %v to flag set: %v", oper.flagset, err))
		}
	}
	return &set
}

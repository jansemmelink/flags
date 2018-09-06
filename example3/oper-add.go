package main

import (
	"github.com/jansemmelink/flags"
)

func init() {
	flagset := flags.NewSet("add", "Add user")
	flagset.String("-n", "--name", "", "Name to add")
	addOper("add", flagset)
}

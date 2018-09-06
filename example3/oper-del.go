package main

import (
	"github.com/jansemmelink/flags"
)

func init() {
	flagset := flags.NewSet("del", "Delete a user")
	flagset.String("-n", "--name", "", "Name to delete")
	addOper("del", flagset)
}

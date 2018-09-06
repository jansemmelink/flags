//Package main shows how one uses "github.com/jansemmelink/flags" for very
//basic command line option parsing
//Run this with option ? to see usage, or make a mistake to see usage
//    $ ./example1 ?
//    Usage: example1 [options]
//    -d	--debug 	false	Run in debug mode
//    -e	--error 	false	Error stack dump
//      	--input 	     	Input filename
//    -o	--output	     	Output filename
//    -l	--limit 	    2	Limit nr of records
//
//Run this with some option values:
//    ./example1 --input=/tmp/in --output=/tmp/out -e -d false
//    2018/09/06 07:30:59 DEBUG OFF
//    2018/09/06 07:30:59 Error dump is true
//    2018/09/06 07:30:59 Limit=2
//    2018/09/06 07:30:59 Input=/tmp/in
//    2018/09/06 07:30:59 Output=/tmp/out
//    2018/09/06 07:30:59 Flags: flag{-d:false}flag{-e:true}flag{--input:/tmp/in}flag{-o:/tmp/out}flag{-l:2}
//
//Note that bool flags can be defined with/without value
//When the value is ommitted, true is assumed, e.g. -d can enable debugging, or -d true or --debug or --debug=true
package main

import (
	"log"

	"github.com/jansemmelink/flags"
)

func main() {
	//define your command line options like this, noting that:
	// * you have the option to keep a local variable, or not
	// * you can later retrieve this flag with either Flag("-d") or Flag("--debug")
	// * you only need short or long or both (at least one)
	// * short option must be a '-' followed by a letter or digit
	// * long option must be "--" followed by a word, starting and ending with
	//       a letter or digit, and the rest of the word may also include
	//       dashes, underscores or dots (no other symbols)
	// * long and short options must be unique
	debugFlag := flags.Bool("-d", "--debug", false, "Run in debug mode")

	//define a few more options, currently supporting bool, int and string:
	errorFlag := flags.Bool("-e", "--error", false, "Error stack dump")
	flags.String("", "--input", "", "Input filename")
	flags.String("-o", "--output", "", "Output filename")
	flags.Int("-l", "--limit", 2, "Limit nr of records")

	//now parse the command line options into those flags
	//this is default parsing that will fail with panic()
	//on any unknown or invalid options
	flags.Parse()

	//retrieve the values using either the reference stored when the flag
	//was defined like this:
	if debugFlag.Value().(bool) {
		log.Printf("DEBUG ON\n")
	} else {
		log.Printf("DEBUG OFF\n")
	}

	//another reference:
	log.Printf("Error dump is %v", errorFlag.Value().(bool))

	//or retrieve the value with the short or long option name and type assertion:
	log.Printf("Limit=%d", flags.Flag("-l").Value().(int))
	log.Printf("Input=%s", flags.Flag("--input").Value().(string))
	log.Printf("Output=%s", flags.Flag("--output").Value().(string))

	//print all flags with their values:
	log.Printf("Flags: %+v\n", flags.DefaultSet())
} //main()

package flags

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	shortValidationPattern = regexp.MustCompile("^-[a-zA-Z0-9]$")
	longValidationPattern  = regexp.MustCompile("^--[a-zA-Z0-9][a-zA-Z0-9_.-]*[a-zA-Z0-9]$")
)

//FlagDescription ...
type FlagDescription struct {
	index int
	short string
	long  string
	value interface{}
	doc   string
}

//Set of flags
type Set struct {
	name  string
	doc   string
	flags []FlagDescription
	short map[string]*FlagDescription
	long  map[string]*FlagDescription
}

//NewSet to create a new set
func NewSet(name, doc string) *Set {
	return &Set{
		name:  name,
		doc:   doc,
		flags: make([]FlagDescription, 0),
		short: make(map[string]*FlagDescription),
		long:  make(map[string]*FlagDescription),
	}
}

//newFlag validates and creates a new flag description
//short must be "-C" where C is a letter or digit
//long must be "--ABC" when ABC is a word consisting of 2 or more characters,
//   starting with a letter or digit, followed by more letters, digits, dashes, dots or underscores
//   and ending again with a letter or digit.
func newFlag(short, long string, value interface{}, doc string) (FlagDescription, error) {
	if short != "" && !shortValidationPattern.MatchString(short) {
		return FlagDescription{}, fmt.Errorf("Short option %s must be \"-<letter|digit>\"", short)
	}
	if long != "" && !longValidationPattern.MatchString(long) {
		return FlagDescription{}, fmt.Errorf("Long option %s must be \"--<word>\" that starts and ends with a letters or digits and allows '_', '-' and '.' in the middle", long)
	}
	f := FlagDescription{
		short: short,
		long:  long,
		value: value,
		doc:   doc,
	}
	return f, nil
} //newFlag()

//Bool adds a bool flag to the set
func (set *Set) Bool(short, long string, init bool, doc string) (*FlagDescription, error) {
	if set == nil {
		return nil, fmt.Errorf("Set.Bool() called on set==nil")
	}
	//create the new flag
	value := init
	newFlag, err := newFlag(short, long, value, doc)
	if err != nil {
		return nil, fmt.Errorf("Set.Bool() cannot add %s %s: %v", short, long, err)
	}
	//add
	newFlagPtr, err := set.Add(newFlag)
	if err != nil {
		return nil, fmt.Errorf("Set.Bool() cannot add %s %s: %v", short, long, err)
	}
	return newFlagPtr, nil
} //Set.Bool()

//Int adds an integer flag to the set
func (set *Set) Int(short, long string, init int, doc string) (*FlagDescription, error) {
	if set == nil {
		return nil, fmt.Errorf("Set.Int() called on set==nil")
	}
	//create the new flag
	value := init
	newFlag, err := newFlag(short, long, value, doc)
	if err != nil {
		return nil, fmt.Errorf("Set.Int() cannot add %s %s: %v", short, long, err)
	}
	//add
	newFlagPtr, err := set.Add(newFlag)
	if err != nil {
		return nil, fmt.Errorf("Set.Int() cannot add %s %s: %v", short, long, err)
	}
	return newFlagPtr, nil
} //Set.Int()

//String adds a string flag to the set
func (set *Set) String(short, long string, init string, doc string) (*FlagDescription, error) {
	if set == nil {
		return nil, fmt.Errorf("Set.String() called on set==nil")
	}
	//create the new flag
	value := init
	newFlag, err := newFlag(short, long, value, doc)
	if err != nil {
		return nil, fmt.Errorf("Set.String() cannot add %s %s: %v", short, long, err)
	}
	//add
	newFlagPtr, err := set.Add(newFlag)
	if err != nil {
		return nil, fmt.Errorf("Set.String() cannot add %s %s: %v", short, long, err)
	}
	return newFlagPtr, nil
} //Set.String()

//Add another flag description
func (set *Set) Add(flag FlagDescription) (*FlagDescription, error) {
	if flag.short == "" && flag.long == "" {
		return nil, fmt.Errorf("Flag without short/long option")
	}
	if flag.value == nil {
		return nil, fmt.Errorf("Flag without valuePtr")
	}
	if flag.doc == "" {
		return nil, fmt.Errorf("Flag without documentation")
	}
	if flag.short != "" {
		if _, ok := set.short[flag.short]; ok {
			return nil, fmt.Errorf("Duplicate short option %s", flag.short)
		}
	}
	if flag.long != "" {
		if _, ok := set.long[flag.long]; ok {
			return nil, fmt.Errorf("Duplicate long option %s", flag.long)
		}
	}
	flag.index = len(set.flags)
	set.flags = append(set.flags, flag)
	newFlagPtr := &set.flags[len(set.flags)-1]
	if flag.short != "" {
		set.short[flag.short] = newFlagPtr
	}
	if flag.long != "" {
		set.long[flag.long] = newFlagPtr
	}
	return newFlagPtr, nil
} //Set.Add()

//Flag to return a flag description
func (set Set) Flag(n string) FlagDescription {
	flag, ok := set.short[n]
	if !ok {
		flag, ok = set.long[n]
		if !ok {
			return FlagDescription{}
		}
	}
	return *flag
} //Set.Flag()

//Parse and return error if found any unknown options
func (set *Set) Parse(options []string) error {
	remainingArgs, err := set.ParseKnown(options)
	if err != nil {
		return err
	}
	if len(remainingArgs) > 0 {
		return fmt.Errorf("Unknown options: %v", remainingArgs)
	}
	return nil
} //Set.Parse()

//ParseKnown process all known arguments and return the remaining/unknown args
//but return error on invalid arguments
func (set *Set) ParseKnown(options []string) ([]string, error) {
	remainingArgs := make([]string, 0)
	skip := 0
	for i, opt := range options {
		if skip > 0 {
			skip--
			continue
		} //if option already parsed as value in previous loop

		flag, ok := set.short[opt]
		valueString := ""
		if ok {
			//found short option match, value in next opt element
			if i < len(options)-1 {
				valueString = options[i+1]
				skip = 1
			}
		} else {
			//not a short option, may be a long options "--word=value"
			//so we need to match "--word"
			ss := strings.SplitN(opt, "=", 2)
			dashDashWord := ss[0]
			if len(ss) > 1 {
				valueString = ss[1]
			}
			flag, ok = set.long[dashDashWord]
			if !ok {
				//unknown option: add to remain and move on
				remainingArgs = append(remainingArgs, opt)
				skip = 0
				continue
			}
		} //if not short

		switch v := flag.value.(type) {
		case bool:
			//if next option is "true" or "false", parse the value
			if valueString == "true" {
				flag.value = true
			} else if valueString == "false" {
				flag.value = false
			} else {
				//not using next option as valueString
				flag.value = true
				skip = 0
			}
		case int:
			intValue, err := strconv.Atoi(valueString)
			if err != nil {
				return remainingArgs, fmt.Errorf("Expecting %s <integer> or %s=<integer>", flag.short, flag.long)
			}
			flag.value = intValue
		case string:
			flag.value = valueString
		default:
			return remainingArgs, fmt.Errorf("Sorry, flags of type %T is not yet fully supported", v)
		}

		//variable flag is local, we must update the flag in the set as well
		set.flags[flag.index].value = flag.value
	} //for each option specified
	return remainingArgs, nil
} //Set.Parse()

//Format to write the set into text
func (set Set) Format(state fmt.State, c rune) {
	s := ""
	for _, flag := range set.flags {
		s += fmt.Sprintf("flag{%n:%v}", flag, flag.value)
	}
	state.Write([]byte(s))
} //Set.Format()

//PrintUsage prints all flags as one would normally print them in command line usage output
func (set Set) PrintUsage(f *os.File) {
	longLen := 0
	valueLen := 0
	for _, flag := range set.flags {
		l := len(flag.long)
		if l > longLen {
			longLen = l
		}
		if flag.value != nil {
			v := fmt.Sprintf("%v", flag.value)
			vl := len(v)
			if vl > valueLen {
				valueLen = vl
			}
		}
	} //for each flag

	for _, flag := range set.flags {
		fmt.Fprintf(f, "\t%s\t%-*.*s\t%*v\t%s\n",
			flag.short,
			longLen,
			longLen,
			flag.long,
			valueLen,
			flag.value,
			flag.doc)
	} //for each flag
	return
} //Set.PrintUsage()

//Format to write the flag into text
func (f FlagDescription) Format(state fmt.State, c rune) {
	s := ""
	S := ""
	if f.short != "" {
		s = f.short
		S := f.short
		if f.long != "" {
			S += " (" + f.long + ")"
		}
	} else {
		s = f.long
		S = f.long
	}
	switch c {
	case 'n':
		state.Write([]byte(s))
	case 'N':
		state.Write([]byte(S))
	default:
		state.Write([]byte(fmt.Sprintf("flag(%s)", s)))
	}
} //FlagDescription.Format()

//Value to get the parsed value of the flag
func (f FlagDescription) Value() interface{} {
	return f.value
} //FlagDescription.Value()

//return true if string consists only of alpha-numeric characters: 0-9,a-z,A-Z
func onlyAlnum(s string) bool {
	for i, c := range s {
		if !unicode.IsDigit(c) && !unicode.IsLetter(c) {
			log.Printf("option \"%s\"[%d]=%v is invalid.", s, i, c)
			return false
		}
	}
	return true
} //onlyAlnum()

/*
...todo: example to combine flag sets
parse flag set to remove things used and later parse the rest
make groups of commands in the usage section
lose standing items, and list of items, sequential items, ...defaultSet
validation functions
*/

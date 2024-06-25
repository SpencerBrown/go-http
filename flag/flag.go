package flag

import (
	"fmt"
	"strings"
	"text/tabwriter"
)

type Flag struct {
	ShortName   string // short flag name
	Description string // description of flag
	Value       any    // default value and type of flag
}

type Flags map[string]Flag

type FlagTypes interface {
	int | int64 | string | bool
}

// NewFlags creates a new set of flags. 
func NewFlags() Flags {
	return make(Flags)
}

// NewFlag creates a new flag. It is a generic function that sets the default value,
// whose type is carried because it is an interface{}.
// It returns a new *Flags because 
func NewFlag[V FlagTypes](f Flags, name string, shortName string, description string, value V) {
	if f == nil {
		panic("flag.NewFlag called with nil Flags")
	}
	f[name] = Flag{
		ShortName:   shortName,
		Description: description,
		Value:       value,
	}
}

// GetFlag is a generic type to get the value of a flag.
// ok is false if the type of the value is not what was expected.
func GetFlag[V FlagTypes](f Flags, name string) (V, bool) {
	v := f[name]
	vv, ok := v.Value.(V)
	return vv, ok
}

// GetFlagMust is a generic type to get the value of a flag.
// It panics if the type of the flag value is not what was expected.
func GetFlagMust[V FlagTypes](f Flags, name string) V {
	v := f[name].Value
	vv, ok := v.(V)
	if !ok {
		var wantV V
		panic(fmt.Sprintf("flag.GetFlagMust internal error: flag %s is type %T, tried to get as type %T", name, v, wantV))
	}
	return vv
}

// GetFlags parses the command line args and sets flags accordingly
// Flag parsing stops just before the first non-flag argument ("-" is a non-flag argument) or after the terminator "--",
// and the Args slice is set to the remaining command line arguments.
// The flag can be --name or -shortname, the value can have an = or not.
// You must use the --flag=false form to turn off a boolean flag.
// Integer flags accept 1234, 0664, 0x1234 and may be negative.
// Boolean flags may be 1, 0, t, f, T, F, true, false, TRUE, FALSE, True, False.
// Duration flags accept any input valid for time.ParseDuration.
// []string flags accept a list of comma-separated strings.
// --help automatically prints out the flags.
func GetFlags(fs Flags) error {
	fmt.Println(fs)
	return nil
}

func (fs Flags) String() string {
	s := strings.Builder{}
	s.WriteString("Flags:\n")
	w := tabwriter.NewWriter(&s, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Name\tShort\tDefault\tType\tDescription")
	for n, f := range fs {
		fmt.Fprintf(w, "%s\t%s\t%v\t%T\t%s\n", n, f.ShortName, f.Value, f.Value, f.Description)
	}
	w.Flush()
	return s.String()
}

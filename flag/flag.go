package flag

import (
	"fmt"
	"strings"
	"text/tabwriter"
)

// Flag is a single flag.
type Flag struct {
	name        string   // name of flag
	alias       []string // alias names
	shortName   string   // short flag name (must be single character)
	description string   // description of flag
	value       any      // default value and type of flag
}

// Flags is a set of Flag.
type Flags map[string]*Flag

// FlagTypes is a constraint on the types of flag values.
type FlagTypes interface {
	int | int64 | string | bool
}

// NewFlags creates a new set of flags.
func NewFlags() Flags {
	return make(Flags)
}

// Name returns the name of a Flag.
func (f *Flag) Name() string {
	return f.name
}

// Alias returns the aliases of a Flag.
func (f *Flag) Alias() []string {
	return f.alias
}

// ShortName returns the one-character short name of a flag.
func (f *Flag) ShortName() string {
	return f.shortName
}

// Description returns the description of a flag.
func (f *Flag) Description() string {
	return f.description
}

// NewFlag creates a new flag and adds it to a set of flags.
// It is a generic function that sets the default value
// whose type is carried because it is an interface{}.
// The short name can be null "" meaning no short name, or a single character.
// The name, short name, and all aliases must not be blank.
// The caller must ensure that the flag name, aliases, and short name
// don't conflict with any existing flags in the set, or with each other;
// otherwise panic ensues.
func NewFlag[V FlagTypes](flgs Flags, nm string, al []string, sn string, desc string, value V) *Flag {
	// do basic checks of the parameters
	if flgs == nil {
		panic("flag.NewFlag called with nil Flags")
	}
	if len(sn) > 1 {
		panic("flag.NewFlag called with shortName of 2 characters or more")
	}
	if len(sn) == 1 && len(strings.TrimSpace(sn)) == 0 {
		panic("flag.NewFlag called with blank short name")
	}
	if len(strings.TrimSpace(nm)) == 0 {
		panic("flag.NewFlag called with blank flag name")
	}
	for _, alias := range al {
		if len(strings.TrimSpace(alias)) == 0 {
			panic("flag.NewFlag called with a blank alias")
		}
	}
	// ensure no duplicates among name, aliases, and shortname
	checker := make([]string, 0)
	checker = append(checker, nm)
	checker = append(checker, sn)
	checker = append(checker, al...)
	chk := make(map[string]struct{})
	for _, str := range checker {
		_, ok := chk[str]
		if ok {
			panic(fmt.Sprintf("flag.NewFlag: duplicate name %s", str))
		}
		chk[str] = struct{}{}
	}
	// create and return the flag
	flg := &Flag{
		name:        nm,
		alias:       al,
		shortName:   sn,
		description: desc,
		value:       value,
	}
	for flgName, flgValue := range flgs {
		if flgName == flg.name {
			panic(fmt.Sprintf("flag.NewFlag: attempt to add already existing flag %s", flgName))
		}
		for _, newAlias := range flg.alias {
			for _, oldAlias := range flgValue.alias {
				if oldAlias == newAlias {
					panic(fmt.Sprintf("flag.NewFlag: attempt to add flag %s with alias %s which is also an alias for flag %s", flg.name, newAlias, flgName))
				}
			}
		}
		if flgValue.shortName == flg.shortName {
			panic(fmt.Sprintf("flag.NewFlag: attempt to add flag %s with identical shortname %s as flag %s", flg.name, string(flg.shortName), flgName))
		}
	}
	flgs[nm] = flg
	return flg
}

// GetFlagOK gets a flag, returning ok as false if the flag does not exist.
func GetFlagOK(f Flags, name string) (*Flag, bool) {
	flg, ok := f[name]
	return flg, ok
}

// GetFlag gets a flag, panics if the flag does not exist.
func GetFlag(f Flags, name string) *Flag {
	flg, ok := f[name]
	if ok {
		return flg
	} else {
		panic(fmt.Sprintf("flag.GetFlag internal error: flag %s does not exist", name))
	}
}

// GetValueOK is a generic function to get a flag and the value of a flag.
// ok is false if the flag does not exist, or the type of the value is not what was expected.
func GetValueOK[V FlagTypes](f *Flag) (V, bool) {
	v, ok := f.value.(V)
	return v, ok
}

// GetValue is a generic function to get a Flag object, and the properly typed value of the flag.
// It panics if the flag does not exist, or type of the flag value is not what was expected.
func GetValue[V FlagTypes](f *Flag) V {
	v, ok := f.value.(V)
	if !ok {
		var wantV V
		panic(fmt.Sprintf("flag.GetValue internal error: for flag %s, value is type %T, tried to get as type %T", f.name, f.value, wantV))
	}
	return v
}

// ParseFlags parses the command line args and sets flags accordingly
// Flag parsing stops just before the first non-flag argument ("-" is a non-flag argument) or after the terminator "--",
// and the Args slice is set to the remaining command line arguments.
// The flag can be --name or -shortname, the value can have an = or not.
// You must use the --flag=false form to turn off a boolean flag.
// Integer flags accept 1234, 0664, 0x1234 and may be negative.
// Boolean flags may be 1, 0, t, f, T, F, true, false, TRUE, FALSE, True, False.
// Duration flags accept any input valid for time.ParseDuration.
// []string flags accept a list of comma-separated strings.
// --help automatically prints out the flags.
func ParseFlags(fs Flags) error {
	fmt.Println(fs)
	return nil
}

func (fs Flags) String() string {
	s := strings.Builder{}
	s.WriteString("Flags:\n")
	w := tabwriter.NewWriter(&s, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Name\tShort\tAliases\tDefault\tType\tDescription")
	for n, f := range fs {
		fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%T\t%s\n", n, f.shortName, strings.Join(f.alias, ","), f.value, f.value, f.description)
	}
	w.Flush()
	return s.String()
}

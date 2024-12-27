package flag

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"unicode/utf8"
)

// Flag is a single flag.
// The aliases and the name must be different from each other, and cannot be a single character.
// The differences must be case insensitive; the names and aliases are converted to lower case.
// The short name must be a single character, or null "". It is case sensitive.
type Flag struct {
	name        string   // name of flag
	alias       []string // alias names
	shortName   string   // short flag name (must be single character)
	description string   // description of flag
	value       any      // default value and type of flag
}

// Flags is a set of *Flag.
// The key is the name of the flag.
// The names and aliases and short names must be unique among all flags in the set.
type Flags map[string]*Flag

// FlagTypes is a constraint on the types of flag values.
type FlagTypes interface {
	int | int64 | string | bool
}

// NewFlags creates a new empty set of flags.
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
// whose type is carried because it is saved as an interface{}.
// The name, short name, and all aliases must not be blank. Blanks are trimmed.
// If anthing is not valid, it panics.
func NewFlag[V FlagTypes](flgs Flags, nm string, al []string, sn string, desc string, value V) *Flag {
	// do basic checks of the parameters
	if flgs == nil {
		panic("flag.NewFlag called with nil Flags")
	}
	shortname := strings.TrimSpace(sn)
	shortnameLength := utf8.RuneCountInString(shortname)
	name := strings.ToLower(strings.TrimSpace(nm))
	nameLength := utf8.RuneCountInString(name)
	if shortnameLength > 1 {
		panic(fmt.Sprintf("flag.NewFlag called with shortName of 2 runes or more: %s", shortname))
	}
	if nameLength == 0 {
		panic("flag.NewFlag called with blank flag name")
	}
	if nameLength == 1 {
		panic(fmt.Sprintf("flag.NewFlag called with a single-rune flag name: %s", name))
	}
	aliases := make([]string, 0)
	for _, aliasuntrimmed := range al {
		alias := strings.ToLower(strings.TrimSpace(aliasuntrimmed))
		aliasLength := utf8.RuneCountInString(alias)
		if aliasLength == 0 {
			panic("flag.NewFlag called with a blank alias")
		}
		if aliasLength == 1 {
			panic(fmt.Sprintf("flag.NewFlag called with a single-rune alias: %s", alias))
		}
		aliases = append(aliases, alias)
	}
	// ensure no duplicates among name and aliases for this flag
	checker := make([]string, 0)
	checker = append(checker, name)
	checker = append(checker, aliases...)
	chk := make(map[string]struct{})
	for _, str := range checker {
		_, ok := chk[str]
		if ok {
			panic(fmt.Sprintf("flag.NewFlag: duplicate name/alias %s", str))
		}
		chk[str] = struct{}{}
	}
	// all seems OK for this flag, create the flag
	flg := &Flag{
		name:        name,
		alias:       aliases,
		shortName:   shortname,
		description: desc,
		value:       value,
	}
	// ensure no conflicts with existing flags
	for flgName, flgValue := range flgs {
		if flgName == flg.name {
			panic(fmt.Sprintf("flag.NewFlag: attempt to add already existing flag name %s", flgName))
		}
		for _, newAlias := range flg.alias {
			if newAlias == flgValue.name {
				panic(fmt.Sprintf("flag.NewFlag: attempt to add alias %s of flag %s which is also the name of another flag", newAlias, flg.name))
			}
			for _, oldAlias := range flgValue.alias {
				if oldAlias == newAlias {
					panic(fmt.Sprintf("flag.NewFlag: attempt to add flag %s with alias %s which is also an alias for flag %s", flg.name, newAlias, flgName))
				}
			}
		}
		if len(flg.shortName) > 0 && flgValue.shortName == flg.shortName {
			panic(fmt.Sprintf("flag.NewFlag: attempt to add flag %s with identical shortname %s as flag %s", flg.name, string(flg.shortName), flgName))
		}
	}
	// no conflicts, add the flag
	flgs[nm] = flg
	return flg
}

// GetFlagOK gets a flag by name, returning ok as false if the flag does not exist.
func GetFlagOK(f Flags, name string) (*Flag, bool) {
	flg, ok := f[name]
	return flg, ok
}

// GetFlag gets a flag by name, panics if the flag does not exist.
func GetFlag(f Flags, name string) *Flag {
	flg, ok := f[name]
	if ok {
		return flg
	} else {
		panic(fmt.Sprintf("flag.GetFlag internal error: flag %s does not exist", name))
	}
}

// GetValueOK is a generic function to get the value of a flag.
// ok is false if the type of the value is not what was expected.
func GetValueOK[V FlagTypes](f *Flag) (V, bool) {
	v, ok := f.value.(V)
	return v, ok
}

// GetValue is a generic function to get the properly typed value of the flag.
// It panics if the type of the flag value is not what was expected.
func GetValue[V FlagTypes](f *Flag) V {
	v, ok := f.value.(V)
	if !ok {
		var wantV V
		panic(fmt.Sprintf("flag.GetValue internal error: for flag %s, value is type %T, tried to get as type %T", f.name, f.value, wantV))
	}
	return v
}

// GetValueAny gets the value of a flag as an interface{}.
func (f *Flag) GetValueAny() any {
	return f.value
}

// FindFlag finds a flag within a flag set by name or alias or shortname.
// It returns nil if the flag is not found.
func (flgs Flags) FindFlag(nm string) *Flag {
	name := strings.ToLower(strings.TrimSpace(nm))
	for _, f := range flgs {
		if f.name == name || f.shortName == name {
			return f
		}
		for _, alias := range f.alias {
			if alias == name {
				return f
			}
		}
	}
	return nil
}

// String returns a string representation of a Flag, useful for debugging.
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

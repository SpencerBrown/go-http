package option

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"unicode"
	"unicode/utf8"
)

// Option is a single option.
// The aliases and the name must be different from each other, and cannot be a single character.
// The differences must be case insensitive; the names and aliases are converted to lower case.
// The short name must be a single character, or empty "" meaning no shortname. It is case sensitive.
type Option struct {
	name            string   // name of option
	aliases         []string // alias names
	shortName       rune     // short option name (0 if none)
	shortAliases    []rune   // short option aliases
	description     string   // description of option
	longDescription string   // long description of option
	value           any      // default value and type of option; also holds the current value
	// the value is an interface and its type is the type of the value, constrainted to OptionTypes
	handler OptionHandler // handler to call for this option, or nil if none
}

// OptionHandler is a function that handles an option when it is set.
// It returns an error if there was a problem handling the option.
type OptionHandler func(opt *Option) error

// Options is a set of *Option.
// The key is the name of the option.
// The names and aliases and short names and aliases must be unique among all options in the set.
type Options map[string]*Option

// OptionTypes is a constraint on the types of option values.
// TODO support float, float64, []string, duration
type OptionTypes interface {
	int | int64 | string | bool
}

// NewOptions creates a new empty set of options.
func NewOptions() Options {
	return make(Options)
}

// Name returns the name of a Option.
func (opt *Option) Name() string {
	return opt.name
}

// Alias returns the aliases of a Option.
func (opt *Option) Alias() []string {
	return opt.aliases
}

// ShortName returns the one-character (rune) short name of a option. Zero rune if none.
func (opt *Option) ShortName() rune {
	return opt.shortName
}

// ShortNameAliases returns the one-character (rune) short name aliases of a option. Empty if none.
func (opt *Option) ShortNameAliases() []rune {
	return opt.shortAliases
}

// Description returns the description of a option.
func (opt *Option) Description() string {
	return opt.description
}

// LongDescription returns the long description of a option.
func (opt *Option) LongDescription() string {
	return opt.longDescription
}

// NewOption creates a new option.
// It is a generic function that sets the default value whose type is carried because it is saved as an interface{}.
// The name and all aliases must not be all whitespace. Whitespace is trimmed.
// The name and aliases are case insensitive, must be at least two characters, and must be unique.
// The short name is case sensitive and must be a single character or 0 if none.
// There must be a short name if there are short name aliases.
// The short name and all short name aliases must be non-whitespace characters and unique.
// The description and long description can be empty strings.
// The value must be one of the types in OptionTypes: int, int64, string, or bool.
// Unicode runes and strings are supported.
// If anthing is not valid, it panics.
func NewOption[V OptionTypes](nm string, al []string, sn rune, sa []rune, desc string, longdesc string, value V, handler OptionHandler) *Option {
	// Trim and lowercase the name and aliases, check for duplicates or single character names/aliases
	name := strings.ToLower(strings.TrimSpace(nm)) // names are case insensitive
	nameLength := utf8.RuneCountInString(name)
	if nameLength == 0 {
		panic("option.NewOption called with blank option name")
	}
	if nameLength == 1 {
		panic(fmt.Sprintf("option.NewOption called with a single-rune option name: %s", name))
	}
	aliases := make([]string, 0)
	for _, aliasuntrimmed := range al {
		alias := strings.ToLower(strings.TrimSpace(aliasuntrimmed))
		aliasLength := utf8.RuneCountInString(alias)
		if aliasLength == 0 {
			panic("option.NewOption called with a blank alias")
		}
		if aliasLength == 1 {
			panic(fmt.Sprintf("option.NewOption called with a single-rune alias: %s", alias))
		}
		aliases = append(aliases, alias)
	}
	// ensure no duplicates among name and aliases for this option
	// note: shortnames cannot be a duplicate because names and aliases must be at least two characters
	allnames := make([]string, 0)
	allnames = append(allnames, name)
	allnames = append(allnames, aliases...)
	chk := make(map[string]struct{})
	for _, thisname := range allnames {
		_, ok := chk[thisname]
		if ok {
			panic(fmt.Sprintf("option.NewOption: duplicate name/alias %s", thisname))
		}
		chk[thisname] = struct{}{}
	}
	// check shortname and shortname aliases
	// if there's no shortname, there cannot be any shortname aliases
	if sn == 0 {
		if len(sa) > 0 {
			panic("option.NewOption called with shortname aliases but no shortname")
		}
	} else {
		// we have a shortname, check that it's not whitespace, and not a duplicate of any of the short aliases
		if unicode.IsSpace(sn) {
			panic("option.NewOption called with a whitespace shortname")
		}
		allshortnames := make([]rune, 0)
		allshortnames = append(allshortnames, sn) // include the shortname itself in the list to check for duplicates
		allshortnames = append(allshortnames, sa...)
		achk := make(map[rune]struct{})
		for _, r := range allshortnames {
			if r == 0 {
				panic("option.NewOption called with a zero rune shortname alias")
			}
			if unicode.IsSpace(r) {
				panic("option.NewOption called with a whitespace shortname alias")
			}
			_, ok := achk[r]
			if ok {
				panic(fmt.Sprintf("option.NewOption: duplicate shortname/alias %c", r))
			}
			achk[r] = struct{}{}
		}
	}
	// note that the value is constrained by the compiler to be one of the allowed types
	// if we got here, all seems OK for this option, create the option
	opt := &Option{
		name:            name,
		aliases:         aliases,
		shortName:       sn,
		shortAliases:    sa,
		description:     desc,
		longDescription: longdesc,
		value:           value,
		handler:         handler,
	}
	return opt
}

// GetOptionOK gets a option by name, returning ok as false if the option does not exist.
func GetOptionOK(f Options, name string) (*Option, bool) {
	opt, ok := f[name]
	return opt, ok
}

// GetOption gets a option by name, panics if the option does not exist.
func GetOption(f Options, name string) *Option {
	opt, ok := f[name]
	if ok {
		return opt
	} else {
		panic(fmt.Sprintf("option.GetOption internal error: option %s does not exist", name))
	}
}

// GetValueOK is a generic function to get the value of a option.
// ok is false if the type of the value is not what was expected.
func GetValueOK[V OptionTypes](f *Option) (V, bool) {
	v, ok := f.value.(V)
	return v, ok
}

// GetValue is a generic function to get the properly typed value of the option.
// It panics if the type of the option value is not what was expected.
func GetValue[V OptionTypes](f *Option) V {
	v, ok := f.value.(V)
	if !ok {
		var wantV V
		panic(fmt.Sprintf("option.GetValue internal error: for option %s, value is type %T, tried to get as type %T", f.name, f.value, wantV))
	}
	return v
}

// GetValueAny gets the value of a option as an interface{}.
func (opt *Option) GetValueAny() any {
	return opt.value
}

// AddOption adds a option to a set of options.
func (opts Options) AddOption(opt *Option) Options {
	// ensure no conflicts with existing options: names, aliases, short names
	for optName, optValue := range opts {
		if optName == opt.name {
			panic(fmt.Sprintf("option.AddOption: attempt to add already existing option name %s", optName))
		}
		for _, newAlias := range opt.aliases {
			if newAlias == optValue.name {
				panic(fmt.Sprintf("option.AddOption: attempt to add alias %s of option %s which is also the name of another option", newAlias, opt.name))
			}
			for _, oldAlias := range optValue.aliases {
				if oldAlias == newAlias {
					panic(fmt.Sprintf("option.AddOption: attempt to add option %s with alias %s which is also an alias for option %s", opt.name, newAlias, optName))
				}
			}
		}
		if opt.shortName != 0 && opt.shortName == optValue.shortName {
			panic(fmt.Sprintf("option.AddOption: attempt to add option %s with identical shortname %s as option %s", opt.name, string(opt.shortName), optName))
		}
		for _, newShortAlias := range opt.shortAliases {
			if newShortAlias == optValue.shortName {
				panic(fmt.Sprintf("option.AddOption: attempt to add shortname alias %c of option %s which is also the shortname of another option", newShortAlias, opt.name))
			}
			for _, oldShortAlias := range optValue.shortAliases {
				if oldShortAlias == newShortAlias {
					panic(fmt.Sprintf("option.AddOption: attempt to add option %s with shortname alias %c which is also a shortname alias for option %s", opt.name, newShortAlias, optName))
				}
			}
		}
	}
	// no conflicts, add the option
	opts[opt.name] = opt
	return opts
}

// FindOption finds a option within a option set by name or alias
// It returns nil if the option is not found.
func (opts Options) FindOption(nm string) *Option {
	name := strings.ToLower(strings.TrimSpace(nm))
	for _, f := range opts {
		if f.name == name {
			return f
		}
		for _, alias := range f.aliases {
			if alias == name {
				return f
			}
		}
	}
	return nil
}

// CheckOptions compares two sets of options to ensure there are no duplicates
// among the names, aliases, short names, and short aliases..
// It returns true if there are no duplicates, false if there are.
func CheckOptionsForDuplicates(opts1 Options, opts2 Options) bool {
	// gather all names and aliases from both sets, check for duplicates
	allNames := make([]string, 0)
	for name, opt := range opts1 {
		allNames = append(allNames, name)
		allNames = append(allNames, opt.aliases...)
	}
	for name, opt := range opts2 {
		allNames = append(allNames, name)
		allNames = append(allNames, opt.aliases...)
	}
	checker := make(map[string]struct{})
	for _, nm := range allNames {
		_, ok := checker[nm]
		if ok {
			return false
		}
		checker[nm] = struct{}{}
	}
	// gather all short names and short aliases from both sets, check for duplicates
	allShortNames := make([]rune, 0)
	for _, opt := range opts1 {
		if opt.shortName != 0 {
			allShortNames = append(allShortNames, opt.shortName)
		}
		allShortNames = append(allShortNames, opt.shortAliases...)
	}
	for _, opt := range opts2 {
		if opt.shortName != 0 {
			allShortNames = append(allShortNames, opt.shortName)
		}
		allShortNames = append(allShortNames, opt.shortAliases...)
	}
	shortChecker := make(map[rune]struct{})
	for _, r := range allShortNames {
		_, ok := shortChecker[r]
		if ok {
			return false
		}
		shortChecker[r] = struct{}{}
	}
	// if we got here, no duplicates
	return true
}

// MergeOptions merges one set of options into another.
// The options in the second set are copied into the first set.
// It is assumed there are no duplicates. Use CheckOptionsForDuplicates to check.
// If there are, the second set will overwrite the first.
func MergeOptions(opts1 Options, opts2 Options) {
	for name, opt := range opts2 {
		opts1[name] = opt
	}
}

// ParseValue sets the value of a option from a string.
func (opt *Option) ParseValue(s string) error {
	switch v := opt.value.(type) {
	case int:
		n, err := fmt.Sscanf(s, "%d", &v)
		if err != nil || n != 1 {
			return fmt.Errorf("option.ParseValue: could not parse %s as int", s)
		}
		opt.value = v
	case int64:
		n, err := fmt.Sscanf(s, "%d", &v)
		if err != nil || n != 1 {
			return fmt.Errorf("option.ParseValue: could not parse %s as int64", s)
		}
		opt.value = v
	case string:
		opt.value = s
	case bool:
		switch s {
		case "true", "True", "TRUE", "t", "T", "1":
			opt.value = true
		case "false", "False", "FALSE", "f", "F", "0":
			opt.value = false
		default:
			return fmt.Errorf("option.ParseValue: could not parse %s as bool", s)
		}
	default:
		panic(fmt.Sprintf("option.ParseValue internal error: unknown type %T", v))
	}
	return nil
}

// String returns a string representation of a Option, useful for debugging.
func (fs Options) String() string {
	s := strings.Builder{}
	s.WriteString("Options:\n")
	w := tabwriter.NewWriter(&s, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Name\tAliases\tShort\tShortAliases\tDefault\tType\tDescription\tLongDescription")
	for n, f := range fs {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\t%T\t%s\t%s\n", n, strings.Join(f.aliases, ","), string(f.shortName), string(f.shortAliases), f.value, f.value, f.description, f.longDescription)
	}
	w.Flush()
	return s.String()
}

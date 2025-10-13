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
// You must use the --flag=false form to turn off a boolean flag.
// -- is used to separate the flags from the arguments.
// Integer flags accept 1234, 0664, 0x1234 and may be negative.
// Boolean flags may be 1, 0, t, f, T, F, true, false, TRUE, FALSE, True, False.
// TODO Duration flags accept any input valid for time.ParseDuration.
// TODO []string flags accept a list of comma-separated strings.
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
// The names and aliases and short names and aliases must be unique among all options in the set.
type Options []*Option

// OptionTypes is a constraint on the types of option values.
// TODO support float, float64, []string, duration
type OptionTypes interface {
	int | int64 | string | bool
}

// NewOptions creates a new empty set of options.
func NewOptions() Options {
	return make(Options, 0)
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
// The short name and short name aliases are runes, are case sensitive, or can be 0 if none.
// There must be a short name if there are short name aliases.
// The short name and all short name aliases must be non-whitespace characters and unique.
// The description and long description can be empty strings.
// The value must be one of the types in OptionTypes: int, int64, string, or bool.
// Unicode runes and strings are supported.
// Returns an error if anything is not valid.
func NewOption[V OptionTypes](nm string, al []string, sn rune, sa []rune, desc string, longdesc string, value V, handler OptionHandler) (*Option, error) {
	// Trim and lowercase the name and aliases, check for duplicates or single character names/aliases
	name := strings.ToLower(strings.TrimSpace(nm)) // names are case insensitive
	nameLength := utf8.RuneCountInString(name)
	if nameLength == 0 {
		return nil, fmt.Errorf("option.NewOption called with blank option name")
	}
	if nameLength == 1 {
		return nil, fmt.Errorf("option.NewOption called with a single-rune option name: %s", name)
	}
	aliases := make([]string, 0)
	for _, aliasuntrimmed := range al { // note if al is nil, this loop is skipped
		alias := strings.ToLower(strings.TrimSpace(aliasuntrimmed))
		aliasLength := utf8.RuneCountInString(alias)
		if aliasLength == 0 {
			return nil, fmt.Errorf("option.NewOption called with a blank alias")
		}
		if aliasLength == 1 {
			return nil, fmt.Errorf("option.NewOption called with a single-rune alias: %s", alias)
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
			return nil, fmt.Errorf("option.NewOption: duplicate name/alias %s", thisname)
		}
		chk[thisname] = struct{}{}
	}
	// check shortname and shortname aliases
	// if there's no shortname, there cannot be any shortname aliases
	shortAliases := sa
	if shortAliases == nil {
		shortAliases = make([]rune, 0)
	}
	if sn == 0 {
		if len(shortAliases) > 0 {
			return nil, fmt.Errorf("option.NewOption called with shortname aliases but no shortname")
		}
	} else {
		// we have a shortname, check that it's not whitespace, and not a duplicate of any of the short aliases
		if unicode.IsSpace(sn) {
			return nil, fmt.Errorf("option.NewOption called with a whitespace shortname")
		}
		allshortnames := make([]rune, 0)
		allshortnames = append(allshortnames, sn) // include the shortname itself in the list to check for duplicates
		allshortnames = append(allshortnames, shortAliases...)
		achk := make(map[rune]struct{})
		for _, r := range allshortnames {
			if r == 0 {
				return nil, fmt.Errorf("option.NewOption called with a zero rune shortname alias")
			}
			if unicode.IsSpace(r) {
				return nil, fmt.Errorf("option.NewOption called with a whitespace shortname alias")
			}
			_, ok := achk[r]
			if ok {
				return nil, fmt.Errorf("option.NewOption: duplicate shortname/alias %c", r)
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
		shortAliases:    shortAliases,
		description:     desc,
		longDescription: longdesc,
		value:           value,
		handler:         handler,
	}
	return opt, nil
}

// NewOptionMust is like NewOption but panics if there is an error.
func NewOptionMust[V OptionTypes](nm string, al []string, sn rune, sa []rune, desc string, longdesc string, value V, handler OptionHandler) *Option {
	opt, err := NewOption(nm, al, sn, sa, desc, longdesc, value, handler)
	if err != nil {
		panic(err)
	}
	return opt
}

// AddOption adds a option to a set of options.
// It returns an error if the option name or any alias or short name or short alias
// conflicts with an existing option in the set.
func (opts *Options) AddOption(opt *Option) error {
	// ensure no conflicts with existing options: names, aliases, short names
	for _, oldOpt := range *opts {
		if oldOpt.name == opt.name {
			return fmt.Errorf("option.AddOption: attempt to add already existing option name %s", oldOpt.name)
		}
		for _, newAlias := range opt.aliases {
			if newAlias == oldOpt.name {
				return fmt.Errorf("option.AddOption: attempt to add alias %s of option %s which is also the name of another option", newAlias, opt.name)
			}
			for _, oldAlias := range oldOpt.aliases {
				if oldAlias == newAlias {
					return fmt.Errorf("option.AddOption: attempt to add option %s with alias %s which is also an alias for option %s", opt.name, newAlias, oldOpt.name)
				}
			}
		}
		if opt.shortName != 0 {
			if opt.shortName == oldOpt.shortName {
				return fmt.Errorf("option.AddOption: attempt to add option %s with identical shortname %s as option %s", opt.name, string(opt.shortName), oldOpt.name)
			}
			for _, newShortAlias := range opt.shortAliases {
				if newShortAlias == oldOpt.shortName {
					return fmt.Errorf("option.AddOption: attempt to add shortname alias %c of option %s which is also the shortname of another option %s", newShortAlias, opt.name, oldOpt.name)
				}
				for _, oldShortAlias := range oldOpt.shortAliases {
					if oldShortAlias == newShortAlias {
						return fmt.Errorf("option.AddOption: attempt to add option %s with shortname alias %c which is also a shortname alias for option %s", opt.name, newShortAlias, oldOpt.name)
					}
				}
			}
		}
	}
	// no conflicts, add the option
	*opts = append(*opts, opt)
	return nil
}

// AddOptionMust adds an option to a set of options and panics if there is an error.
func (opts *Options) AddOptionMust(opt *Option) {
	if err := opts.AddOption(opt); err != nil {
		panic(err)
	}
}

// GetOption gets a option by name, returning nil if the option does not exist.
// The name is case insensitive and whitespace is trimmed.
// It can match either the name or any alias of the option.
func GetOption(f Options, name string) *Option {
	trimmedName := strings.ToLower(strings.TrimSpace(name))
	for _, opt := range f {
		if opt.name == trimmedName {
			return opt
		}
		for _, alias := range opt.aliases {
			if alias == trimmedName {
				return opt
			}
		}
	}
	return nil
}

// GetValue is a generic function to get the value of a option.
// returns false if the type of the value is not what was expected.
func GetValue[V OptionTypes](f *Option) (V, bool) {
	v, ok := f.value.(V)
	return v, ok
}

// GetValueOK is a generic function to get the value of a option.
// returns false if the type of the value is not what was expected.
// Same as GetValue.
func GetValueOK[V OptionTypes](f *Option) (V, bool) {
	v, ok := f.value.(V)
	return v, ok
}

// GetValueAny gets the value of a option as an interface{}.
func (opt *Option) GetValueAny() any {
	return opt.value
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
		return fmt.Errorf("option.ParseValue: unknown type %T", v)
	}
	return nil
}

// String returns a string representation of a Option, useful for debugging.
func (fs Options) String() string {
	s := strings.Builder{}
	s.WriteString("Options:\n")
	w := tabwriter.NewWriter(&s, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Name\tAliases\tShort\tShortAliases\tDefault\tType\tDescription\tLongDescription")
	for _, f := range fs {
		sn := string(f.shortName)
		if f.shortName == 0 {
			sn = "" // space if no short name
		}
		sa := make([]string, 0)
		for _, r := range f.shortAliases {
			sas := string(r)
			if r == 0 {
				sas = "" // space if no short name alias
			}
			sa = append(sa, sas)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\t%T\t%s\t%s\n", f.name, strings.Join(f.aliases, ","), sn, strings.Join(sa, ","), f.value, f.value, f.description, f.longDescription)
	}
	w.Flush()
	return s.String()
}

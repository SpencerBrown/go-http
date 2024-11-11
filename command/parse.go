package command

import (
	"fmt"
	"slices"
	"strings"
)

// Parse parses the command line args and sets args and flags accordingly
// Flag parsing stops just before the first non-flag argument ("-" is a non-flag argument) or after the terminator "--",
// and the Args slice is set to the remaining command line arguments.
// The flag can be --name or -shortname, the value can have an = or not.
// You must use the --flag=false form to turn off a boolean flag.
// Integer flags accept 1234, 0664, 0x1234 and may be negative.
// Boolean flags may be 1, 0, t, f, T, F, true, false, TRUE, FALSE, True, False.
// Duration flags accept any input valid for time.ParseDuration.
// []string flags accept a list of comma-separated strings.
// --help automatically prints out the flags.
func Parse(cmds *Commands, cargs []string) error {
	if cmds == nil {
		return fmt.Errorf("nil commands")
	}
	if cmds.root == nil {
		return fmt.Errorf("nil root command")
	}
	cmd := cmds.root // Current command being parsed
	if cargs == nil {
		return fmt.Errorf("nil command line arguments")
	}
	if len(cargs) == 0 {
		return fmt.Errorf("empty command line arguments")
	}
	for i, carguntrimmed := range cargs {
		carg := strings.TrimSpace(carguntrimmed)
		if carg == "" {
			return fmt.Errorf("empty command line argument")
		}
		if i == 0 { // First element of the command line must be the root command name or alias
			if carg != cmd.name && !slices.Contains(cmd.alias, carg) {
				return fmt.Errorf("Command %s not recognized, should be one of %s", carg, strings.Join(append([]string{cmd.name}, cmd.alias...), ", "))
			}
			continue
		}
		if carg == "--" { // -- is signal for end of commands and flags, the rest is just args
			cmds.SetArgs(cargs[i+1:])
			return nil
		}
		if carg[0] == '-' {
			// This is a flag
			var flagName string
			if len(carg) == 1 {
				return fmt.Errorf("Blank flag %s", carg)
			}
			if carg[1] == '-' {
				// this is a long flag
				flagName = carg[2:]
			} else {
				// this is a short flag
				flagName = carg[1:]
			}
		} else {
			// this is a subcommand or the start of the arguments
			for _, sub := range cmd.sub {
				if carg == sub.name || slices.Contains(sub.alias, carg) {
					cmd = sub
					continue
				}
				cmds.SetArgs(cargs[i:])
				return nil
			}
		}
	}
	return nil
}

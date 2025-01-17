package command

import (
	"fmt"
	"strings"

	"github.com/SpencerBrown/go-http/flag"
)

// ParsedCommand represents the parsed command line arguments and flags
// args is the command line arguments
// flags is the consolidated set of flags for the command
type ParsedCommand struct {
	names []string   // command/subcommand names
	args  []string   // command line arguments
	flags flag.Flags // flags for the command
}

// Parse parses the raw args and sets the flags and args accordingly
// The general structure of the command line is:
// command [flags] [subcommand] [flags] [subcommand] [flags] [args]
// Flags cannot be duplicated among the command and its subcommands on one path of the command tree,
// but flags can be duplicated among different paths of the command tree.
// Parse identifies the subcommands being used and returns a ParseCommand struct with the command line arguments and consolidated flags.
// For each argument:
// - If the argument is a flag, it is parsed and set in the flag.Flags struct for the current command.
// - If the argument is not a flag, we check if it is a subcommand.  If it is, we set the subcommand as the current command.
// - We start with the root command, which doesn't have an actual name, but its subcommands are the top-level commands.
// - If the argument is not a flag or a subcommand, it is added to the args slice, and we stop parsing for flags and subcommands.
// Flag parsing stops just before the first non-flag argument ("-" is a non-flag argument) or after the terminator "--",
// and the args slice in the ParsedCommand is set to the remaining command line arguments.
// The flag can be --name or -shortname, the value can have an = or not.
// You must use the --flag=false form to turn off a boolean flag.
// -- is used to separate the flags from the arguments.
// Integer flags accept 1234, 0664, 0x1234 and may be negative.
// Boolean flags may be 1, 0, t, f, T, F, true, false, TRUE, FALSE, True, False.
// TODO Duration flags accept any input valid for time.ParseDuration.
// TODO []string flags accept a list of comma-separated strings.
// --help automatically prints out the flags and parent commands for that subcommand branch of the tree.
func Parse(cmd *Command, args []string) (*ParsedCommand, error) {
	if cmd == nil {
		return nil, fmt.Errorf("command.Parse called with nil Command")
	}
	currentCmd := cmd
	parsedCmd := ParsedCommand{}
	iArg := 0 // Tracks the first argument after parsing the commands and flags
	for ; iArg < len(args); iArg++ {
		arg := args[iArg]
		if arg == "--" {
			// stop parsing flags when you see a bare "--", the rest is args
			iArg++
			break
		}
		if arg == "-" {
			// stop parsing flags when you see a bare "-" it and following are args
			break
		}
		if strings.HasPrefix(arg, "--") {
			// find the flag with the double dash prefix
			// note we have already checked for a bare "--" so we know there's more in the arg string
			flagNameValue := strings.Split(arg[2:], "=") // it might have a value after the "="
			var flagName, flagValue string
			flagOK := false
			switch len(flagNameValue) {
			case 1: // value is in the next arg
				flagName = flagNameValue[0]
				if iArg < len(args)-1 {
					iArg++
					flagValue = args[iArg]
					flagOK = true
				}
			case 2: // value is after the equals sign
				flagName = flagNameValue[0]
				flagValue = flagNameValue[1]
				flagOK = true
			default:
				// leave flagOK as false
			}
			if !flagOK {
				return nil, fmt.Errorf("invalid flag %s", arg)
			}
			theFlag := findFlag(currentCmd, flagName)
			if theFlag == nil {
				return nil, fmt.Errorf("unknown flag %s", arg)
			}
			// parse the flag value according to the type of the default value in the flag
			if err := theFlag.ParseValue(flagValue); err != nil {
				return nil, err
			}
		} else {
			// check if the arg is a command at the current point in the command tree
			foundCmd := false
			for _, subCmd := range currentCmd.sub {
				if subCmd.name == arg {
					currentCmd = subCmd
					foundCmd = true
				}
			}
			if !foundCmd {
				// stop parsing flags and subcommands when you see a non-flag non-command argument
				break
			}
		}
	}
	// set the names of the command and its subcommands in order from the root to the current command
	for cmd := currentCmd; cmd != nil; cmd = cmd.parent {
		parsedCmd.names = append([]string{cmd.name}, parsedCmd.names...)
	}
	// set the remaining args
	if iArg < len(args)-1 {
		parsedCmd.args = args[iArg:]
	} else {
		parsedCmd.args = nil
	}
	// set the consolidated flags for the command
	parsedCmd.flags = flag.NewFlags()
	for cmd := currentCmd; cmd != nil; cmd = cmd.parent {
		flag.MergeFlags(parsedCmd.flags, cmd.flags)
	}
	return &parsedCmd, nil
}

// FindFlag finds a flag in the command/subcommand tree
// It searches the current command and its parent commands for the flag
func findFlag(cmd *Command, flagName string) *flag.Flag {
	if cmd == nil {
		return nil
	}
	// find the flag at this level of the command tree
	f := cmd.flags.FindFlag(flagName)
	if f != nil {
		return f
	}
	// find the flag in a parent command with a recursive call
	return findFlag(cmd.parent, flagName)
}

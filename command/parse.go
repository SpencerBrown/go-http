package command

import (
	"fmt"
	"strings"

	"github.com/SpencerBrown/go-http/option"
)

// ParsedCommands represents the parsed command line arguments and the parsed commands themselves with options
// args is the command line arguments
type ParsedCommands struct {
	commands []ParsedCommand // parsed commands in order
	args     []string        // command line arguments
}

// ParsedCommand represents a parsed command with its options as specified and defaulted
// The command name is the actual name not an alias
// The options are the actual name and value of the option, not an alias or short name
type ParsedCommand struct {
	name    string               // actual command name, not an alias
	options option.ParsedOptions // actual options for this command, not aliases or short names
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
// --help automatically prints out the flags and parent commands for that subcommand branch of the tree.
func Parse(cmds Commands, cmdArgs []string) (*ParsedCommands, error) {
	if cmds == nil || len(cmds) == 0 {
		return nil, fmt.Errorf("command.Parse called with nil or empty Commands")
	}
	// parsedCmds is what we will return, we build this as we parse the command line
	parsedCmds := ParsedCommands{
		commands: make([]ParsedCommand, 0),
		args:     make([]string, 0),
	}
	// nextCommand tracks where we are in the list of commands and parsed commands
	nextCommand := 0
	// iArg is the index of cmdArgs where we are at the moment, when we exit the loop it will point to the first arg after options and commands
	iArg := 0 
	// We loop through the command line arguments, parsing commands and options
	for ; iArg < len(cmdArgs); iArg++ {
		cmdArg := cmdArgs[iArg]
		if cmdArg == "--" {
			// stop parsing flags when you see a bare "--", the rest is args
			iArg++
			break
		}
		if cmdArg == "-" {
			// stop parsing flags when you see a bare "-" it and following are args
			break
		}
		if strings.HasPrefix(cmdArg, "--") {
			// find the option with the double dash prefix
			// note we have already checked for a bare "--" so we know there's more in the arg string
			// it could be --option=value or --option value or just --option for a boolean flag
			// so we need to parse it accordingly
			optName, optValue, hasEquals := strings.Cut(cmdArg[2:], "=") 
			if hasEquals {
				// it's of the form --option=value
			} else {
				// it's of the form --option value or just --option
				// see if it's a boolean option
				theOption := option.GetOption(findFlag(cmds[nextCommand], optName))
				if iArg < len(cmdArgs)-1 { // make sure there's another arg for the value
					iArg++
					optValue = cmdArgs[iArg]
				}
			var flagName, flagValue string
			flagOK := false
			switch len(flagNameValue) {
			case 1: // value is in the next arg because it's of the form --option value
				flagName = flagNameValue[0]
				if iArg < len(cmdArgs)-1 { // make sure there's another arg for the value
					iArg++
					flagValue = cmdArgs[iArg]
					flagOK = true
				}
			case 2: // value is after the equals sign, it's of the form --option=value
				flagName = flagNameValue[0]
				flagValue = flagNameValue[1]
				flagOK = true
			default:
				// there's more than one equals sign, it is possibly valid if the value has equals signs in it
				flagName = flagNameValue[0]
				flagValue = strings.Join(flagNameValue[1:], "=")
				flagOK = true
			}
			if !flagOK {
				return nil, fmt.Errorf("invalid flag %s", cmdArg)
			}
			theOption := option.GetOption(findFlag(currentCmd, flagName))
			if theOption == nil {
				return nil, fmt.Errorf("unknown flag %s", cmdArg)
			}
			// parse the flag value according to the type of the default value in the flag
			if err := theOption.ParseValue(flagValue); err != nil {
				return nil, err
			
		} else {
			// check if the arg is a command at the current point in the command tree
			foundCmd := false
			for _, subCmd := range currentCmd.sub {
				if subCmd.name == cmdArg {
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
	if iArg < len(cmdArgs)-1 {
		parsedCmd.args = cmdArgs[iArg:]
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

// String returns a string representation of the ParsedCommand
func (pc *ParsedCommand) String() string {
	if pc == nil {
		return "ParsedCommand: nil\n"
	}
	var builder strings.Builder
	builder.WriteString("ParsedCommand:\n")
	builder.WriteString(pc.commands.String())
	builder.WriteString(fmt.Sprintf("Args: %s\n", strings.Join(pc.args, " ")))
	return builder.String()
}

// return the commands in the ParsedCommand
func (pc *ParsedCommand) Commands() Commands {
	return pc.commands
}

// return the args in the ParsedCommand
func (pc *ParsedCommand) Args() []string {
	return pc.args
}

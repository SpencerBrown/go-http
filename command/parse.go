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
	name        string               // actual command name, not an alias
	invokedName string               // the name or alias used to invoke this command
	options     option.ParsedOptions // actual options for this command, not aliases or short names
}

// Parse parses the raw args and sets the options and args accordingly
// Parse identifies the subcommands being used and returns a ParseCommands struct with the command line arguments and consolidated options.
// --help automatically prints out the options and parent commands for that subcommand branch of the tree.
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
	// note that command names are case insensitive
	for ; iArg < len(cmdArgs); iArg++ {
		cmdArg := strings.ToLower(strings.TrimSpace(cmdArgs[iArg]))	
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
			if !hasEquals {
				// it's of the form --option value or just --option
				// see if it's a boolean option
				if iArg < len(cmdArgs)-1 { // make sure there's another arg for the value
					iArg++
					optValue = cmdArgs[iArg]
				}
			}
			// TODO: process optName and optValue
			continue
		}
		// check if the arg is a command at the current point in the command tree
		foundCmd := false
		if nextCommand < len(cmds) {
			if cmds[nextCommand].name == cmdArg {
				foundCmd = true
				nextCommand++
			}
		}
		if !foundCmd {
			// stop parsing flags and subcommands when you see a non-flag non-command argument
			break
		}
	}
	// set the remaining args
	if iArg < len(cmdArgs) {
		parsedCmds.args = cmdArgs[iArg:]
	}
	return &parsedCmds, nil
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

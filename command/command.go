package command

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/SpencerBrown/go-http/flag"
)

// Command represents a command or subcommand in a command/subcommand tree
// flags represents the flags for this command at this level of the tree
// parent is the parent command, nil if this is the root command
// sub is a slice of subcommands at this level of the tree
// name is the name of the command
// alias is a slice of aliases for the command
// description is a description of the command
// The command name and aliases must be unique among the subtree starting at this command
//
// So the syntax is something like:
//
//	 command flags subcommand flags subcommand flags args
//	 where flags start with "-" for short form and "--" for long form followed by the flag name and value, separated by "=" or " "
//	 and args are the remaining command line arguments after the last flags or after "--" if needed to separate commands from args
//		and subcommands can be nested arbitrarily deep
//	 for boolean flags, use --flag=false or -f=false to turn off the flag
//
// Arguments start at the first unrecognized token, or after the terminator "--"
// The --help flag automatically prints out the command syntax and flags
// The func associated with a Command is what is called when the provided command line maps to this Command
// The func is given a ParsedCommand with the command line arguments and flags
type Command struct {
	name        string                     // Name of command
	alias       []string                   // Aliases for command
	description string                     // Description of command
	flags       flag.Flags                 // Flags for this command
	// handler     func(*ParsedCommand) error // Handler for this command
	parent      *Command                   // Parent command
	sub         []*Command                 // Subcommands
}

// Name returns the command name
func (cmd *Command) Name() string {
	return cmd.name
}

// Alias returns the command aliases
func (cmd *Command) Alias() []string {
	return cmd.alias
}

// Description returns the command description
func (cmd *Command) Description() string {
	return cmd.description
}

// Flags returns the Flags for the command
func (cmd *Command) Flags() flag.Flags {
	return cmd.flags
}

// NewCommand creates a new command with the given name, aliases, description, and flags
// command name and any aliases cannot be blank, and cannot duplicate each other
// command name and aliases are case insensitive and can include unicode characters
func NewCommand(nm string, al []string, desc string, flgs flag.Flags) *Command {
	// do basic checks of parameters
	name := strings.ToLower(strings.TrimSpace(nm))
	nameLength := utf8.RuneCountInString(name)
	if nameLength == 0 {
		panic("command.NewCommand called with blank command name")
	}
	aliases := make([]string, 0)
	for _, aliasuntrimmed := range al {
		alias := strings.ToLower(strings.TrimSpace(aliasuntrimmed))
		aliasLength := utf8.RuneCountInString(alias)
		if aliasLength == 0 {
			panic("command.NewCommand called with a blank alias")
		}
		aliases = append(aliases, alias)
	}
	// ensure no duplicates among the name and aliases
	checker := make([]string, 0)
	checker = append(checker, name)
	checker = append(checker, aliases...)
	chk := make(map[string]struct{})
	for _, str := range checker {
		_, ok := chk[str]
		if ok {
			panic(fmt.Sprintf("command.NewCommand called with duplicate name or alias %s", str))
		}
		chk[str] = struct{}{}
	}
	return &Command{
		name:        name,
		alias:       aliases,
		description: desc,
		flags:       flgs,
		sub:         nil,
	}
}

// SetSub adds a subcommand to the command/subcommand tree
// We check to ensure that none of the subcommands at this level of the tree or above
// duplicate each others' names or aliases or flags
// We return the subcommand so that we can chain SetSub calls
func (parentcmd *Command) SetSub(subcmd *Command) *Command {
	if parentcmd == nil {
		panic("command.SetSub called with nil parent command")
	}
	if subcmd == nil {
		panic("command.SetSub called with nil subcommand")
	}
	if subcmd.parent != nil {
		panic("command.SetSub called with subcommand that already has a parent")
	}
	// ensure no duplicates among the names and aliases
	nameList := make([]string, 0)
	nameList = append(nameList, subcmd.name)
	nameList = append(nameList, subcmd.alias...)
	for _, cmd := range parentcmd.sub {
		nameList = append(nameList, cmd.name)
		nameList = append(nameList, cmd.alias...)
	}
	for cmd := parentcmd; cmd != nil; cmd = cmd.parent {
		nameList = append(nameList, cmd.name)
		nameList = append(nameList, cmd.alias...)
	}
	dupNameCheck := make(map[string]struct{})
	for _, str := range nameList {
		_, ok := dupNameCheck[str]
		if ok {
			panic(fmt.Sprintf("command.SetSub called with duplicate name or alias %s", str))
		}
		dupNameCheck[str] = struct{}{}
	}
	// ensure no duplicates among the flags for this command and its parent commands
	thisFlags := subcmd.flags
	for cmd := parentcmd; cmd != nil; cmd = cmd.parent {
		if !flag.CheckFlagsForDuplicates(thisFlags, cmd.flags) {
			panic("command.SetSub called with subcommand that has duplicate flags with a parent command")
		}
	}
	// Add subcommand
	subcmd.parent = parentcmd
	parentcmd.sub = append(parentcmd.sub, subcmd)
	return subcmd
}

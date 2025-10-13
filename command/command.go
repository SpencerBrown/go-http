package command

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/SpencerBrown/go-http/option"
)

// Command represents a command or subcommand in a command/subcommand tree
// options represents the options for this command at this level of the tree
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
	name            string         // Name of command
	alias           []string       // Aliases for command
	description     string         // Description of command
	longDescription string         // Long description of command
	options         option.Options // Flags for this command
	subcommands     Commands       // Subcommands that can follow this command
}

type Commands []*Command

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

// LongDescription returns the command long description
func (cmd *Command) LongDescription() string {
	return cmd.longDescription
}

// Options returns the Options for the command
func (cmd *Command) Options() option.Options {
	return cmd.options
}

// Subcommands returns the subcommands for the command
func (cmd *Command) Subcommands() Commands {
	return cmd.subcommands
}

// NewCommand creates a new command with the given name, aliases, descriptions, and options
// command name and any aliases cannot be blank, and cannot duplicate each other
// command name and aliases are case insensitive and can include unicode characters
func NewCommand(nm string, al []string, desc string, longDesc string, opts option.Options) (*Command, error) {
	// do basic checks of parameters
	name := strings.ToLower(strings.TrimSpace(nm))
	if utf8.RuneCountInString(name) == 0 {
		panic("command.NewCommand called with blank command name")
	}
	aliases := make([]string, 0)
	for _, aliasuntrimmed := range al {
		alias := strings.ToLower(strings.TrimSpace(aliasuntrimmed))
		if utf8.RuneCountInString(alias) == 0 {
			return nil, fmt.Errorf("command.NewCommand called with a blank alias")
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
			return nil, fmt.Errorf("command.NewCommand called with duplicate name or alias %s", str)
		}
		chk[str] = struct{}{}
	}
	return &Command{
		name:            name,
		alias:           aliases,
		description:     desc,
		longDescription: longDesc,
		options:         opts,
		subcommands:     nil,
	}, nil
}

// NewCommandMust is like NewCommand but panics if there is an error.
func NewCommandMust(nm string, al []string, desc string, longDesc string, opts option.Options) *Command {
	cmd, err := NewCommand(nm, al, desc, longDesc, opts)
	if err != nil {
		panic(err)
	}
	return cmd
}

// NewCommands creates a new Commands slice
func NewCommands() *Commands {
	newcmds := make(Commands, 0)
	return &newcmds
}

// AddCommand adds a subcommand to the given Commands slice
// does nothing if cmds or cmd is nil, or if the cmd name or any alias duplicates an existing name or alias in cmds
func (cmds *Commands) AddCommand(cmd *Command) error {
	if cmds == nil {
		return fmt.Errorf("command.AddCommand called with nil command list")
	}
	if cmd == nil {
		return fmt.Errorf("command.AddCommand called with nil command")
	}
	// ensure no duplicates among the names and aliases
	nameList := make([]string, 0)
	nameList = append(nameList, cmd.name)
	nameList = append(nameList, cmd.alias...)
	for _, existingCmd := range *cmds {
		nameList = append(nameList, existingCmd.name)
		nameList = append(nameList, existingCmd.alias...)
	}
	dupNameCheck := make(map[string]struct{})
	for _, str := range nameList {
		_, ok := dupNameCheck[str]
		if ok {
			return fmt.Errorf("command.AddCommand called with duplicate name or alias %s", str)
		}
		dupNameCheck[str] = struct{}{}
	}
	*cmds = append(*cmds, cmd)
	return nil
}

// AddCommandMust adds an option to a set of options and panics if there is an error.
func (cmds *Commands) AddCommandMust(cmd *Command) {
	if err := cmds.AddCommand(cmd); err != nil {
		panic(err)
	}
}

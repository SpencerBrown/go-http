package command

import (
	"fmt"
	"strings"

	"github.com/SpencerBrown/go-http/flag"
	"github.com/SpencerBrown/go-http/util"
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
	handler     func(*ParsedCommand) error // Handler for this command
	parent      *Command                   // Parent command
	sub         []*Command                 // Subcommands
}

// ParsedCommand represents the parsed command line arguments and flags
// args is the command line arguments
// flags is the flags for the command

type ParsedCommand struct {
	names []string   // command/subcommand names
	args  []string   // command line arguments
	flags flag.Flags // flags for the command
}

// Commands represents a command/subcommand tree and the command line arguments
type Commands struct {
	args []string
	root *Command
}

// Args returns the command line arguments
func (cmds *Commands) Args() []string {
	return cmds.args
}

// SetArgs sets the command line arguments
func (cmds *Commands) SetArgs(args []string) {
	cmds.args = args
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

// NewCommands creates a new empty command/subcommand tree
func NewCommands() *Commands {
	return &Commands{
		args: make([]string, 0),
		root: nil,
	}
}

// NewCommand creates a new command with the given name, aliases, description, and flags
// command name and any aliases cannot be blank, and cannot duplicate each other
// and cannot duplicate any other command name or alias in the entire tree
// flags must be unique among the entire tree
func NewCommand(nm string, al []string, desc string, flgs flag.Flags) *Command {
	// do basic checks of parameters
	name := strings.TrimSpace(nm)
	if len(name) == 0 {
		panic("command.NewCommand called with blank command name")
	}
	aliases := make([]string, 0)
	for _, aliasuntrimmed := range al {
		alias := strings.TrimSpace(aliasuntrimmed)
		if len(alias) == 0 {
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
	// create the new command and return it
	return &Command{
		name:        name,
		alias:       aliases,
		description: desc,
		flags:       flgs,
		sub:         make([]*Command, 0),
	}
}

// SetRoot makes a command the root command (main command name) for a command/subcommand tree
func (cmds *Commands) SetRoot(cmd *Command) {
	if cmds == nil {
		panic("command.SetRoot called with nil Commands")
	}
	if cmd == nil {
		panic("command.SetRoot called with nil Command")
	}
	if cmds.root == nil {
		cmds.root = cmd
	} else {
		panic("command.SetRoot Internal Error: attempt to set root command twice")
	}
}

// SetSub adds a subcommand to the command/subcommand tree
// We check to ensure that none of the subcommands at this level of the tree
// duplicate each others' names or aliases
func (parentcmd *Command) SetSub(subcmd *Command) {
	// basic checks
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
	checker := make([]string, 0)
	checker = append(checker, subcmd.name)
	checker = append(checker, subcmd.alias...)
	for _, cmd := range parentcmd.sub {
		checker = append(checker, cmd.name)
		checker = append(checker, cmd.alias...)
	}
	chk := make(map[string]struct{})
	for _, str := range checker {
		_, ok := chk[str]
		if ok {
			panic(fmt.Sprintf("command.SetSub called with duplicate name or alias %s", str))
		}
		chk[str] = struct{}{}
	}
	// Add subcommand
	subcmd.parent = parentcmd
	parentcmd.sub = append(parentcmd.sub, subcmd)
}

func (cmds *Commands) String() string {
	if cmds == nil {
		return "Commands: nil\n"
	}
	var builder strings.Builder
	builder.WriteString("Commands:\n")
	builder.WriteString(fmt.Sprintf("Args: %s\n", strings.Join(cmds.args, ", ")))
	cmds.root.commandTree(&builder, 0)
	return builder.String()
}

func (cmd *Command) commandTree(builder *strings.Builder, indent int) string {
	if cmd == nil {
		return "Command: nil\n"
	}
	// stringify this command at the indent level
	builder.WriteString(util.Indent(cmd.String(), indent))
	// output the subcommands indented recursively
	for _, subcmd := range cmd.sub {
		subcmd.commandTree(builder, indent+2)
	}
	return builder.String()
}

func (cmd *Command) String() string {
	if cmd == nil {
		return "Command: nil\n"
	}
	var builder strings.Builder
	builder.WriteString("Command: " + cmd.name + "\n")
	builder.WriteString(util.Indent(fmt.Sprintf("Aliases: %s\n", strings.Join(cmd.alias, "'")), 1))
	builder.WriteString(util.Indent(fmt.Sprintf("Description: %s\n", cmd.description), 1))
	builder.WriteString(util.Indent(cmd.flags.String(), 1))
	return builder.String()
}

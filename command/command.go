package command

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/SpencerBrown/go-http/option"
	"github.com/SpencerBrown/go-http/util"
)

// Command represents a command or subcommand in a command/subcommand tree. The root of the tree is a Commands, representing a set of Command. 
// options represents the options for this command at this level of the tree
// subcommands is a Commands representing the subcommands that can follow this command
// name is the name of the command
// alias is a slice of aliases for the command
// description is a description of the command, longDescription is a longer description of the command
// The command name and aliases must be unique at this level of the tree. Command names and aliases are case insensitive.
// The command line is parsed left to right, matching command names and aliases at each level of the tree.
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

// Commands is a set of Command representing a set of commands at this level of the command tree
type Commands map[string](*Command)

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

// NewCommands creates a new Commands
func NewCommands() Commands {
	return make(Commands, 0)
}

// AddCommand adds a Command to a Commands
// checks for duplicate names and aliases among existing commands
func (cmds *Commands) AddCommand(cmd *Command) error {
	if cmds == nil {
		return fmt.Errorf("command.AddCommand called with nil Commands")
	}
	if cmd == nil {
		return fmt.Errorf("command.AddSubcommand called with nil Command")
	}
	// ensure no duplicates among the names and aliases at this level of the tree
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
			return fmt.Errorf("command.AddSubcommand called with duplicate name or alias %s", str)
		}
		dupNameCheck[str] = struct{}{}
	}
	(*cmds)[cmd.name] = cmd // add the Command to the Commands
	return nil
}

// AddCommandMust adds a Command to aCommands and panics if there is an error.
func (cmds *Commands) AddCommandMust(cmd *Command) {
	if err := cmds.AddCommand(cmd); err != nil {
		panic(err)
	}
}

// AddSubcommand adds a Command as a subcommand to this Command
func (cmd *Command) AddSubcommand(subcmd *Command) error {
	if cmd == nil {
		return fmt.Errorf("command.AddSubcommand called with nil Command")
	}
	if subcmd == nil {
		return fmt.Errorf("command.AddSubcommand called with nil subcommand")
	}
	if cmd.subcommands == nil {
		cmd.subcommands = NewCommands()
	}
	return cmd.subcommands.AddCommand(subcmd)
}

// AddSubcommandMust adds a Command as a subcommand to this Command and panics if there is an error.
func (cmd *Command) AddSubcommandMust(subcmd *Command) {
	if err := cmd.AddSubcommand(subcmd); err != nil {
		panic(err)
	}
}

// String returns a string representation of the Commands and all subcommands
// we loop through each command and its subcommands recursively,
// indenting each level for readability
func (cmds *Commands) String() string {
	if cmds == nil {
		return "Commands: nil\n"
	}
	var builder strings.Builder
	builder.WriteString("Commands:\n")
	showCommands(&builder, cmds, 0)
	return builder.String()
}

func showCommands(builder *strings.Builder, cmds *Commands, indent int) {
	for _, cmd := range *cmds {
		cmd.showCommand(builder, indent)
		for _, subcmd := range cmd.subcommands {
			subcmd.showCommand(builder, indent+1)
			showCommands(builder, &subcmd.subcommands, indent+1)
		}
	}
}

// showCommand appends the string representation of a Command to the builder
func (cmd *Command) showCommand(builder *strings.Builder, indent int) {
	if cmd == nil {
		builder.WriteString("Command: nil\n")
		return
	}
	builder.WriteString("Command: " + cmd.name + "\n")
	builder.WriteString(util.Indent(fmt.Sprintf("Aliases: %s\n", strings.Join(cmd.alias, "'")), indent+1))
	builder.WriteString(util.Indent(fmt.Sprintf("Description: %s\n", cmd.description), indent+1))
	builder.WriteString(util.Indent(fmt.Sprintf("Long Description: %s\n", cmd.longDescription), indent+1))
	builder.WriteString(util.Indent(cmd.options.String(), indent+1))
}

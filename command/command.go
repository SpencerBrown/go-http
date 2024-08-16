package command

import (
	"fmt"
	"strings"

	"github.com/SpencerBrown/go-http/flag"
)

type Command struct {
	name        string
	alias       []string
	description string
	flags       flag.Flags
	parent      *Command
	sub         []*Command
}

type Commands struct {
	args []string
	root *Command
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

// NewCommand creates a new command to be added to the tree
// command name and any aliases cannot be blank, and cannot duplicate each other
func NewCommand(nm string, al []string, desc string, flgs flag.Flags) *Command {
	// do basic checks of parameters
	if len(strings.TrimSpace(nm)) == 0 {
		panic("command.NewCommand called with blank command name")
	}
	for _, alias := range al {
		if len(strings.TrimSpace(alias)) == 0 {
			panic("command.NewCommand called with a blankj alias")
		}
	}
	// ensure no duplicates among the name and aliases
	checker := make([]string, 0)
	checker = append(checker, nm)
	checker = append(checker, al...)
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
		name:        nm,
		alias:       al,
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
	if subcmd.parent == nil {
		subcmd.parent = parentcmd
		parentcmd.sub = append(parentcmd.sub, subcmd)
	} else {
		panic("command.SetSub Internal Error: attempt to set subcommand twice")
	}
}

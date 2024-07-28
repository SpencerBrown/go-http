package command

import "github.com/SpencerBrown/go-http/flag"

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

func NewCommands() *Commands {
	return &Commands{
		args: make([]string, 0),
		root: nil,
	}
}

func NewCommand(nm string, al []string, desc string, flgs flag.Flags) *Command {
	return &Command{
		name:        nm,
		alias:       al,
		description: desc,
		flags:       flgs,
		sub:         make([]*Command, 0),
	}
}

func (cmds *Commands) SetRoot(cmd *Command) {
	if cmds.root == nil {
		cmds.root = cmd
	} else {
		panic("command.SetRoot Internal Error: attempt to set root command twice")
	}
}

func (cmd *Command) SetSub(subcmd *Command) {
	if subcmd.parent == nil {
		subcmd.parent = cmd
		cmd.sub = append(cmd.sub, subcmd)
	} else {
		panic("command.SetSub Internal Error: attempt to set subcommand twice")
	}
}

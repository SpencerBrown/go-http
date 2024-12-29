package command

import (
	"fmt"
	"strings"

	"github.com/SpencerBrown/go-http/util"
)

func (cmd *Command) String() string {
	if cmd == nil {
		return "Command: nil\n"
	}
	var builder strings.Builder
	cmd.commandTree(&builder, 0)
	return builder.String()
}

func (cmd *Command) commandTree(builder *strings.Builder, indent int) string {
	if cmd == nil {
		return ""
	}
	// stringify this command at the indent level
	builder.WriteString(util.Indent(cmd.OneString(), indent))
	// output the subcommands indented recursively
	for _, subcmd := range cmd.sub {
		subcmd.commandTree(builder, indent+2)
	}
	return builder.String()
}

func (cmd *Command) OneString() string {
	if cmd == nil {
		return ""
	}
	var builder strings.Builder
	builder.WriteString("Command: " + cmd.name + "\n")
	builder.WriteString(util.Indent(fmt.Sprintf("Aliases: %s\n", strings.Join(cmd.alias, "'")), 1))
	builder.WriteString(util.Indent(fmt.Sprintf("Description: %s\n", cmd.description), 1))
	builder.WriteString(util.Indent(cmd.flags.String(), 1))
	return builder.String()
}

package command

import (
	"fmt"
	"strings"

	"github.com/SpencerBrown/go-http/util"
)

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

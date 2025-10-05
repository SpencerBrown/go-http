Go foundation package

Support any application that starts as a CLI, and  that can optionally act as a Web server

At each level there can be a set of commands with corresponding options. Each command can be the root of another list of subcommands and options. 

Command —option <value> cmd2 —option2 args….   

Options for a command must follow the command and precede the next subcommand. Two subcommands at the same level can have the same options. An unrecognized command/option starts the args. A lone dash also signals start of args. 

Options have types for their values. Type can be bool which means they are true if asserted with no value or value if true, and set to false with a value of false. Options can have multiple alias names. 

A single dash indicates a single letter option. There can be multiple aliases of a single letter as well. 

Data should look like:

type struct Level []Command
type struct Command {
	name string
	description string
	longDescription string
	aliases []string
	options []Option
	handler commandHandler
}
type struct Option {
	name string
	description string
	longDescription string
	aliases []string
	shortname rune
	shortaliases []rune
	value any
}
type commandHandler func(cmd *Command) error 
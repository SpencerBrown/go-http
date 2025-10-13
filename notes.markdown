Go foundation package

Support any application that starts as a CLI, and  that can optionally act as a Web server

At each level there can be a set of commands with corresponding options. Each command can be the root of another list of subcommands and options. 

Command —option <value> cmd2 —option2 args….   

Options for a command must follow the command and precede the next subcommand. Two subcommands at the same level can have the same options. An unrecognized command/option starts the args. A lone dash also signals start of args. 

Options have types for their values. Type can be bool which means they are true if asserted with no value or value if true, and set to false with a value of false. Options can have multiple alias names. 

A single dash indicates a single letter option. There can be multiple aliases of a single letter as well. 

Data should look like:

```go
type Commands []Command
type Args []string

type struct Command {
	name string
	subcommands Commands
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
	handler optionHandler
}

type commandHandler func(cmd *Command) error 

type optionHandler func(opt *Option) error
```

after setting up the commands, we will have:

Commands (list of Command)
	Command
		options (list of Option)
		subcommands (Commands - list of Command)

at each level, Commands must be unique both in name and aliases

parsing of the command string results in:

a Commands which is a list of the actual commands and subcommands from the string with all the options for each command set to the specified values (or the default value)
an Args which is a list of the arguments from the command string
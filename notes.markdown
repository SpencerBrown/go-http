Go foundation package

Support any application that starts as a CLI, and  that can optionally act as a Web server

At each level there can be a set of commands with corresponding options. Each command can be the root of another list of subcommands and options. 

Command —-option <value> cmd2 —-option2 args….   

Options for a command must follow the command and precede the next subcommand. Two subcommands at the same level can have the same options. This implies that for a single command without subcommands, the options must precede any command line arguments. 

An unrecognized command starts the command line arguments. A lone double dash also signals the start of the arguments. A lone single dash is the first argument. 

Options have types for their values. Type can be bool which means they are true if asserted with no value or value if true, and set to false with a value of false. Options can have multiple alias names. 

A single dash indicates a single letter option. There can be multiple aliases of a single letter as well. 

An option's names must be two or more alphanumeric characters. Non-ASCII characters are fine but cannot be white space. Unicode is supported. THe name is folded to lowercase, so options cannot be the same except for case.

An option can also have a single character name and single character aliases. These are NOT folded to lower case.

A single character option is preceded by a single dash. Single letter boolean options can be stacked together after a single dash. The last stacked single letter option can be any type and have a value. 

Options have default values, and a specified type of accepted value. The default value is applied when an option is immediately followed by another option. 

The user can also specify a handler function to do their own manipulation of the value.

Command names ans aliases are also folded to lowercase and cannot contain whitespace. Unicode is supported. Command names cannot start with a dash, because that makes them look like options.

Data should look like:

```go
type Commands map[string]Command

type ParsedCommands  struct {
	commands []ParsedCommand // list of commands in order that they were invoked, each with its set of options
	args []string // list of remaining arguments in the command line
}

type ParsedCommand struct { // command and options as processed by parser
	name string // name of command
	invokedName string // name by which the command was invoked
	options map[string]ParsedOption // options for this command/subcommand
}

type ParsedOption struct {
	name string
	invokedName string
	isDefault bool
	isSet bool
	value any
}

type struct Command { // command information as provided by the user
	name string // command name
	description string
	longDescription string
	aliases []string
	options []Option
	subcommands Commands //subcommands it can have
	handler commandHandler
}

type struct Option {
	name string
	description string
	longDescription string
	aliases []string
	shortname rune
	hasDefault bool // option has a default value
	shortaliases []rune
	value any // set to the default value, or some value of the type if no default
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

a ParsedCommands which is a list of the actual commands and subcommands from the string with all the options for each command set to the specified values (or the default value)
including the args which is a list of the arguments from the command string after parsing

## Outline of logic for parsing a command line

Initialize empty ParsedCommands. 
Step through the args slice. 
Look for a lone dash, if found, leave to handle the args.
Look for a lone double dash, if found, skip it and leave to handle the args.
First look for a command, if none found, leave to handle the args.
Create a ParsedCommand and add to the list in ParsedCommands.
Now look for and handle options.
Loop. 

Handle Options for a command/subcommand:
given current Command, ParsedCommand, ParsedCommands, current args pointer
loop through options list
	for each, create a ParsedOption 
loop through args
	canonArg := current arg stripped whitespace before/after, lowercased
	if a lone dash, leave 
	if a lone double dash, increment args pointer and leave
	if doesn't start with a dash, leave
	if it starts with a triple dash, leave
	if starts with double dash, handle long option
	if starts with a single dash, handle short option
  increment args pointer
end loop

handle long option: 
given current Command, ParsedCommand, args index, canonical arg
split across equal sign if there is one, assign second part to value and first part to optionName, value cannot be empty
look up optionName after canonicalizing it (whitespace, lowerCase)
if not found, return error
// we want to allow these cases:
// 	--option=value (value cannot be empty)
//	--option value (not allowed if option is boolean)
//  --option (if it is boolean)
if value is not empty, there must have been an equals sign, do setValue and return
if value is empty, if option is boolean, do setValue to "the default true or false value" and return
the next arg must be a value, even if there is a default! grab it and setValue and return

handle short option:
given current Command, ParsedCommand, args index, we know it starts with a single dash and has something
// handle cases:
// -o (only if option is boolean)
// -ovalue (not allowed if option is boolean)
// -o=value
// -o value (not allowed if option is boolean)
// -abo any of the above "o" and only if a and b are boolean
set shortOption index to 1
loop over boolean options:
	if current shortOption is a valid shortOption, and it's boolean:
		if followed by nothing, set value to True and return
		if followed by =, assign remaining to value, setValue, and return
		otherwise, continue, next rune must be another option
	if we get here, we have processed all the boolean option prefixes and there must be one and only one non-boolean option remaining
split across equal sign if there is one, assign second part to value and first part to optionRune, value cannot be empty, setValue and return
if no equal sign, and there is more in the token, assign that to value, setValue and return
final case: -o value and non-boolean option. get next token, setValue and return.

setValue:
given current Option and ParsedOption and parsed option name and token string:
set invokedName
switch on type of value
parse value appropriately and assign to value, returning error as appropriate
set isSet flag and return

Handle args: 
Get slice of args starting with current arg and assign to args.
Now complete the ParsedCommands and return it. 

## Better rewrite of specs for command line options

### POSIX specs

see [the official specs](https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html)

briefly:

* options are single dash and single letter
* lone single dash is an argument used as a filename indicating stdin or stdout
* options may have an option argument, or not
* options without an option argument are just boolean flags
* options with an option argument come in two flavors:
	* mandatory argument follows the option with a space in between
	* optional argument immediately follows the option letter with no space
* multiple options without an argument, except for the last one which can have an argument, may be grouped in a single string after the single dash -
* the first -- is a delimiter for the "operands" of the command

Numbers are parsed as 0 to 2**31-1 and can have a minus sign. 

### GNU specs

extends POSIX, as [they describe POSIX](https://sourceware.org/glibc/manual/latest/html_node/Argument-Syntax.html):

* Arguments are options if they begin with a hyphen delimiter (‘-’).
* Multiple options may follow a hyphen delimiter in a single token if the options do not take arguments. Thus, ‘-abc’ is equivalent to ‘-a -b -c’.
* Option names are single alphanumeric characters (as for isalnum; see Classification of Characters).
* Certain options require an argument. For example, the -o option of the ld command requires an argument—an output file name.
* An option and its argument may or may not appear as separate tokens. (In other words, the whitespace separating them is optional.) Thus, -o foo and -ofoo are equivalent.
* Options typically precede other non-option arguments.

The implementations of getopt and argp_parse in the GNU C Library normally make it appear as if all the option arguments were specified before all the non-option arguments for the purposes of parsing, even if the user of your program intermixed option and non-option arguments. They do this by reordering the elements of the argv array. This behavior is nonstandard; if you want to suppress it, define the _POSIX_OPTION_ORDER environment variable. See Standard Environment Variables.

* The argument -- terminates all options; any following arguments are treated as non-option arguments, even if they begin with a hyphen.
* A token consisting of a single hyphen character is interpreted as an ordinary non-option argument. By convention, it is used to specify input from or output to the standard input and output streams.
* Options may be supplied in any order, or appear multiple times. The interpretation is left up to the particular application program.

and they add:
* Long options consist of -- followed by a name made of alphanumeric characters and dashes. Option names are typically one to three words long, with hyphens to separate words. Users can abbreviate the option names as long as the abbreviations are unique.

### Go flags package specs

* One or two dashes, they are equivalent
* Nothing special about single-character flag names
* Flag value types are Bool, Duration, Float64, Int, Int64, Uint, Uint64, String
	* Can specify an unmarshaler for String, they call this type Text
* Also allow for user-defined types with a callback function to create the value of the user-defined type
* Boolean flags can be `-flag` (sets true) or `-flag=true/false` (sets true or false) but not `-flag true/false`.
* Integer flags accept 1234, 0664, 0x1234 and may be negative. 
* Boolean flags may be: 1, 0, t, f, T, F, true, false, TRUE, FALSE, True, False
* Duration flags accept any input valid for time.ParseDuration.
	* A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

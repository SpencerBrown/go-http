package main

import (
	"context"
	"fmt"
	"os"

	"github.com/SpencerBrown/go-http/command"
	"github.com/SpencerBrown/go-http/flag"
	"github.com/SpencerBrown/go-http/run"
)

func main() {
	ctx := context.Background()
	r := run.Runner{
		Commands:    command.NewCommands(),
		GetEnvVar:   os.Getenv,
		GetWorkDir:  os.Getwd,
		Input:       os.Stdin,
		Output:      os.Stdout,
		ErrorOutput: os.Stderr,
	}
	flgs := flag.NewFlags()
	flag.NewFlag(flgs, "foo", nil, "f", "first part of foobar", "one")
	flag.NewFlag(flgs, "bar", nil, "b", "second part of foobar", 2)
	flag.NewFlag(flgs, "foobar", []string{"fb"}, "", "is it?", true)
	cmds := r.Commands
	cmd := command.NewCommand("foobarfoo", []string{"fbf"}, "foobar command", flgs)
	cmds.SetRoot(cmd)
	cmds.SetArgs(os.Args[1:])
	subcmd := command.NewCommand("subfoobar", []string{"sfb"}, "sub foobar command", nil)
	cmd.SetSub(subcmd)
	if err := r.Run(ctx, true); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

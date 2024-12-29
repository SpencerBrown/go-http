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
		Command:    nil,
		Args:       os.Args,
		GetEnvVar:   os.Getenv,
		GetWorkDir:  os.Getwd,
		Input:       os.Stdin,
		Output:      os.Stdout,
		ErrorOutput: os.Stderr,
	}
	flgs := flag.NewFlags()
	flgs.AddFlag(flag.NewFlag("foo", nil, "f", "first part of foobar", "one"))
	flgs.AddFlag(flag.NewFlag("bar", nil, "b", "second part of foobar", 2))
	flgs.AddFlag(flag.NewFlag("foobar", []string{"fb"}, "", "is it?", true))
	cmd := command.NewCommand("foobarfoo", []string{"fbf"}, "foobar command", flgs)
	subcmd := command.NewCommand("subfoobar", []string{"sfb"}, "sub foobar command", nil)
	cmd.SetSub(subcmd)
	if err := r.Run(ctx, true); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

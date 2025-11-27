package main

import (
	"context"
	"fmt"
	"os"

	"github.com/SpencerBrown/go-http/command"
	"github.com/SpencerBrown/go-http/option"
	"github.com/SpencerBrown/go-http/run"
)

func main() {
	ctx := context.Background()
	r := run.Runner{
		Commands:    nil,
		Args:        os.Args,
		GetEnvVar:   os.Getenv,
		GetWorkDir:  os.Getwd,
		Input:       os.Stdin,
		Output:      os.Stdout,
		ErrorOutput: os.Stderr,
	}

	opts := option.NewOptions()
	opts.AddOptionMust(option.NewOptionMust("foo", nil, 'f', []rune{'g'}, "first", "first part of foobar", false, "one", nil))
	opts.AddOptionMust(option.NewOptionMust("bar", nil, 'b', nil, "second", "second part of foobar", false, 2, nil))
	opts.AddOptionMust(option.NewOptionMust("foobar", []string{"fb"}, 0, nil, "is?", "is it?", false, true, nil))

	cmds := command.Commands{}
	cmds.AddCommandMust(command.NewCommandMust("foobarfoo", []string{"fbf"}, "foobar", "foobar command", opts))
	cmds.AddCommandMust(command.NewCommandMust("subfoobar", []string{"sfb"}, "sub foobar", "sub foobar command", nil))
	r.Commands = &cmds
	if err := r.Run(ctx, true); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

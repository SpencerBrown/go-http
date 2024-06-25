package main

import (
	"context"
	"fmt"
	"os"

	"github.com/SpencerBrown/go-http/flag"
	"github.com/SpencerBrown/go-http/run"
)

func main() {
	ctx := context.Background()
	r := run.Runner{
		Args:        os.Args,
		Flags:       flag.NewFlags(),
		GetEnvVar:   os.Getenv,
		GetWorkDir:  os.Getwd,
		Input:       os.Stdin,
		Output:      os.Stdout,
		ErrorOutput: os.Stderr,
	}
	flag.NewFlag(r.Flags, "foo", "f", "first part of foobar", "one")
	flag.NewFlag(r.Flags, "bar", "b", "second part of foobar", 2)
	flag.NewFlag(r.Flags, "foobar", "", "is it?", true)
	if err := flag.GetFlags(r.Flags); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	if err := r.Run(ctx, true); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

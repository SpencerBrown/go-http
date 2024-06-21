package main

import (
	"context"
	"fmt"
	"os"

	"github.com/SpencerBrown/go-http/run"
)

func main() {
	ctx := context.Background()
	r := run.Runner{
		Args:        os.Args,
		GetEnvVar:   os.Getenv,
		GetWorkDir:  os.Getwd,
		Input:       os.Stdin,
		Output:      os.Stdout,
		ErrorOutput: os.Stderr,
	}
	// flags := []run.RunFlag{
	// 	{
	// 		Name:        "foo",
	// 		ShortName:   "f",
	// 		Description: "first part of foobar",
	// 		Value:       "one",
	// 	},
	// 	{
	// 		Name:        "bar",
	// 		ShortName:   "b",
	// 		Description: "second part of foobar",
	// 		Value:       2,
	// 	},
	// 	{
	// 		Name:        "foobar",
	// 		ShortName:   "",
	// 		Description: "is it?",
	// 		Value:       false,
	// 	},
	// }
	// r.Flags = flags
	run.NewFlag(&r, "foo", "f", "first part of foobar", "one")
	run.NewFlag(&r, "bar", "b", "second part of foobar", 2)
	run.NewFlag(&r, "foobar", "", "is it?", true)
	if err := r.GetFlags(true); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	if err := r.Run(ctx, true); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/SpencerBrown/go-http/run"
)

func main() {
	ctx := context.Background()
	r := run.RunArgs{
		PrintArgs:   false,
		Args:        os.Args,
		EnvVars:     nil,
		GetEnvVar:   os.Getenv,
		GetWorkDir:  os.Getwd,
		Input:       os.Stdin,
		Output:      os.Stdout,
		ErrorOutput: os.Stderr,
	}
	if err := r.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

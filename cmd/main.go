package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"text/template"
	"time"
)

type RunArgs struct {
	Args        []string // command line arguments; args[0] is the command name
	Input       io.Reader
	Output      io.Writer
	ErrorOutput io.Writer
	EnvVars     map[string]string // caller fills in keys with "" values, they get filled in at run time
	GetEnvVar   func(string) string
	GetWorkDir  func() (string, error)
}

// the following copied from Mat Ryer's blog post "How I write HTTP servrices in Go after 13 years"
// https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/

func run(ctx context.Context, r *RunArgs) error {
	// Gather required and optional environment variables
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	// the actual program, which watches for the ctx being signaled
	// there's probably a fancy way of putting some error code in the context

	fmt.Println(r.String())

	i := 0
	le := len(r.Args) - 1
	for {
		select {
		case <-ctx.Done():
			return errors.New("Interrupted!")
		default:
			time.Sleep(time.Second)
			if i++; i > le {
				i = 1
			}
			if le <= 0 {
				fmt.Fprintln(r.Output, 0)
			} else {
				fmt.Fprintln(r.Output, i, r.Args[i])
			}
		}
	}
	//
}

func main() {
	ctx := context.Background()
	r := RunArgs{
		Args:        os.Args,
		EnvVars:     nil,
		GetEnvVar:   os.Getenv,
		GetWorkDir:  os.Getwd,
		Input:       os.Stdin,
		Output:      os.Stdout,
		ErrorOutput: os.Stderr,
	}
	if err := run(ctx, &r); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func (r *RunArgs) String() string {
	tmpl, err := template.New("runargs").Parse(runArgsTemplate)
	if err != nil {
		panic(fmt.Sprintf("Internal error parsing template for RunArgs: %v", err))
	}
	s := strings.Builder{}
	err = tmpl.Execute(&s, r)
	if err != nil {
		panic(fmt.Sprintf("Internal error executing template for RunArgs: %v", err))
	}
	return s.String()
}

const runArgsTemplate = `Run arguments{{if .Args}} for {{index .Args 0}}:
Args: {{if gt (len .Args) 1}}{{slice .Args 1}}{{else}}None{{end}}{{end}}
WOrking directory: {{call .WorkingDirectory}}`

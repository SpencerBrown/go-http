package run

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

type Runnable interface {
	Run(context.Context) error
}

type RunArgs struct {
	PrintArgs   bool
	Args        []string          // command line arguments; args[0] is the command name
	EnvVars     map[string]string // caller fills in keys with "" values, they get filled in at run time
	GetEnvVar   func(string) string
	GetWorkDir  func() (string, error)
	Input       io.Reader
	Output      io.Writer
	ErrorOutput io.Writer
}

type RunFlag struct {
}

// the following copied from Mat Ryer's blog post "How I write HTTP services in Go after 13 years"
// https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/

func (r *RunArgs) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	if r.PrintArgs {
		fmt.Println(r.String())
	}

	// the actual program, which watches for the ctx being signaled
	// there's probably a fancy way of putting some error code in the context

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
WOrking directory: {{call .GetWorkDir}}`

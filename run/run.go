package run

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/SpencerBrown/go-http/flag"
)

type Runnable interface {
	Run(context.Context, bool) error
}

type Runner struct {
	Args        []string // command line arguments; args[0] is the command name
	Flags       flag.Flags
	GetEnvVar   func(string) string
	GetWorkDir  func() (string, error)
	Input       io.Reader
	Output      io.Writer
	ErrorOutput io.Writer
}

// the following copied from Mat Ryer's blog post "How I write HTTP services in Go after 13 years"
// https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/

func (r *Runner) Run(ctx context.Context, debug bool) error {
	if debug {
		fmt.Println(r.String())
	}
	// run the program
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	// the actual program, which watches for the ctx being signaled
	// there's probably a fancy way of putting some error code in the context

	i := 0
	le := len(r.Args) - 1
	for {
		select {
		case <-ctx.Done():
			return errors.New("Interrupted!")
		default:
			time.Sleep(time.Second / 2)
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
}

func (r *Runner) String() string {
	s := strings.Builder{}
	// 	tmpl, err := template.New("runargs").Parse(runArgsTemplate)
	// 	if err != nil {
	// 		panic(fmt.Sprintf("Internal error parsing template for RunArgs: %v", err))
	// 	}
	// 	err = tmpl.Execute(&s, r)
	// 	if err != nil {
	// 		panic(fmt.Sprintf("Internal error executing template for RunArgs: %v", err))
	// 	}
	s.WriteString("---Runner struct---\n")
	if r.Args == nil || len(r.Args) == 0 {
		s.WriteString("No command or arguments\n")
	} else {
		fmt.Fprintf(&s, "Command: %s\n", r.Args[0])
		if len(r.Args) == 1 {
			s.WriteString("No arguments\n")
		} else {
			fmt.Fprintf(&s, "Arguments: %v\n", r.Args[1:])
		}
	}
	workdir, err := r.GetWorkDir()
	if err == nil {
		fmt.Fprintf(&s, "Working directory: %s\n", workdir)
	}
	s.WriteString("Flags:\n")
	s.WriteString(r.Flags.String())
	s.WriteString("---Runner struct---\n")
	return s.String()
}

// const runArgsTemplate = `Run arguments{{if .Args}} for {{index .Args 0}}:
// Args: {{if gt (len .Args) 1}}{{slice .Args 1}}{{else}}None{{end}}{{end}}
// Working directory: {{call .GetWorkDir}}
// Flags:
// Name Short Default Type Description{{range .Flags}}
// {{.Name}} {{if .ShortName}}{{.ShortName}}{{else}}-{{end}} {{.Value}} {{printf "%T" .Value}} "{{.Description}}"{{end}}`

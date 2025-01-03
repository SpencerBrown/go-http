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

	"github.com/SpencerBrown/go-http/command"
)

type Runnable interface {
	Run(context.Context, bool) error
}

type Runner struct {
	Command     *command.Command       // The template for the expected command line with command, subcommands, flags, and args
	Args        []string               // The actual command line
	GetEnvVar   func(string) string    // A function to get an environment variable
	GetWorkDir  func() (string, error) // A function to get the working directory
	Input       io.Reader              // The input stream
	Output      io.Writer              // The output stream
	ErrorOutput io.Writer              // The error output stream
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
	s.WriteString(r.Command.String())
	workdir, err := r.GetWorkDir()
	if err == nil {
		fmt.Fprintf(&s, "Working directory: %s\n", workdir)
	}
	s.WriteString("---Runner struct---\n")
	return s.String()
}

// const runArgsTemplate = `Run arguments{{if .Args}} for {{index .Args 0}}:
// Args: {{if gt (len .Args) 1}}{{slice .Args 1}}{{else}}None{{end}}{{end}}
// Working directory: {{call .GetWorkDir}}
// Flags:
// Name Short Default Type Description{{range .Flags}}
// {{.Name}} {{if .ShortName}}{{.ShortName}}{{else}}-{{end}} {{.Value}} {{printf "%T" .Value}} "{{.Description}}"{{end}}`

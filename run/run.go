package run

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"text/tabwriter"
	"time"
)

type Runnable interface {
	Run(context.Context, bool) error
}

type Runner struct {
	Args        []string           // command line arguments; args[0] is the command name
	flags       map[string]RunFlag // working data for flag processing
	GetEnvVar   func(string) string
	GetWorkDir  func() (string, error)
	Input       io.Reader
	Output      io.Writer
	ErrorOutput io.Writer
}

type RunFlag struct {
	ShortName   string // short flag name
	Description string // description of flag
	Value       any    // default value and type of flag
}

type RunFlagTypes interface {
	int | int64 | string | bool
}

// NewFlag creates a new flag. It is a generic function that sets the default value
func NewFlag[V RunFlagTypes](r *Runner, name string, shortName string, description string, value V) {
	if r.flags == nil {
		r.flags = make(map[string]RunFlag, 0)
	}
	r.flags[name] = RunFlag{
		ShortName:   shortName,
		Description: description,
		Value:       value,
	}
}

// GetFlag is a generic type to get the value of a flag.
// ok is false if the type of the value is not what was expected.
func GetFlag[V RunFlagTypes](r *Runner, name string) (V, bool) {
	v := r.flags[name]
	vv, ok := v.Value.(V)
	return vv, ok
}

// GetFlagMust is a generic type to get the value of a flag.
// It panics if the type of the flag value is not what was expected.
func GetFlagMust[V RunFlagTypes](r *Runner, name string) V {
	v := r.flags[name].Value
	vv, ok := v.(V)
	if !ok {
		var wantV V
		panic(fmt.Sprintf("Internal error: flag %s is type %T, tried to get as type %T", name, v, wantV))
	}
	return vv
}

// GetFlags parses the command line args and sets flags accordingly
// Flag parsing stops just before the first non-flag argument ("-" is a non-flag argument) or after the terminator "--",
// and the Args slice is set to the remaining command line arguments.
// The flag can be --name or -shortname, the value can have an = or not.
// You must use the --flag=false form to turn off a boolean flag.
// Integer flags accept 1234, 0664, 0x1234 and may be negative.
// Boolean flags may be 1, 0, t, f, T, F, true, false, TRUE, FALSE, True, False.
// Duration flags accept any input valid for time.ParseDuration.
// []string flags accept a list of comma-separated strings.
// --help automatically prints out the flags.
func (r *Runner) GetFlags(debug bool) error {
	if debug {
		fmt.Println(r.String())
	}
	return nil
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

	// i := 0
	// le := len(r.Args) - 1
	for {
		select {
		case <-ctx.Done():
			return errors.New("Interrupted!")
		default:
			time.Sleep(time.Second / 2)
			f1 := GetFlagMust[string](r, "foo")
			fmt.Printf("foo value: %v (%T)\n", f1, f1)
			f2, ok := GetFlag[int](r, "bar")
			if ok {
				fmt.Printf("bar value: %v (%T)\n", f2, f2)
			}
			f3 := GetFlagMust[bool](r, "foobar")
			fmt.Printf("foobar value: %v (%T)\n", f3, f3)
			f4 := GetFlagMust[int64](r, "foobar")
			fmt.Printf("foobar value: %v (%T)\n", f4, f4)
			// if i++; i > le {
			// 	i = 1
			// }
			// if le <= 0 {
			// 	fmt.Fprintln(r.Output, 0)
			// } else {
			// 	fmt.Fprintln(r.Output, i, r.Args[i])
			// }
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
	fmt.Fprintf(&s, "Args: %v\n", r.Args)
	workdir, err := r.GetWorkDir()
	if err == nil {
		fmt.Fprintf(&s, "Working directory: %s\n", workdir)
	}
	fmt.Fprintln(&s, "Flags:")
	w := tabwriter.NewWriter(&s, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Name\tShort\tDefault\tType\tDescription")
	for n, f := range r.flags {
		fmt.Fprintf(w, "%s\t%s\t%v\t%T\t%s\n", n, f.ShortName, f.Value, f.Value, f.Description)
	}
	w.Flush()
	return s.String()
}

// const runArgsTemplate = `Run arguments{{if .Args}} for {{index .Args 0}}:
// Args: {{if gt (len .Args) 1}}{{slice .Args 1}}{{else}}None{{end}}{{end}}
// Working directory: {{call .GetWorkDir}}
// Flags:
// Name Short Default Type Description{{range .Flags}}
// {{.Name}} {{if .ShortName}}{{.ShortName}}{{else}}-{{end}} {{.Value}} {{printf "%T" .Value}} "{{.Description}}"{{end}}`

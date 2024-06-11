package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"
)

// the following copied from Mat Ryer's blog post "How I write HTTP servrices in Go after 13 years"
// https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/

func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	// the actual program, which watches for the ctx being signaled
	// there's probably a fancy way of putting some error code in the context
	i := 0
	for {
		select {
		case <-ctx.Done():
			return errors.New("Interrupted!")
		default:
			time.Sleep(time.Second)
			if i++; i >= len(args) {
				i = 1
			}
			fmt.Fprintln(w, i, args[i])
		}
	}
	//
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

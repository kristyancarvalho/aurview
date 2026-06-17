package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kristyancarvalho/aurview/internal/app"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	args := os.Args[1:]
	if len(args) == 1 && (args[0] == "--version" || args[0] == "version") {
		fmt.Fprintf(os.Stdout, "aurview %s\ncommit: %s\ndate: %s\n", version, commit, date)
		return
	}

	if err := app.Run(context.Background(), app.Options{InitialQuery: strings.Join(args, " ")}); err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

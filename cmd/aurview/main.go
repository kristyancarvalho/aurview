package main

import (
	"context"
	"os"
	"strings"

	"github.com/kristyancarvalho/aurview/internal/app"
	"github.com/kristyancarvalho/aurview/internal/version"
)

func main() {
	args := os.Args[1:]
	if len(args) == 1 && (args[0] == "--version" || args[0] == "version") {
		_, _ = os.Stdout.WriteString(version.Current().String("aurview"))
		return
	}

	if err := app.Run(context.Background(), app.Options{InitialQuery: strings.Join(args, " ")}); err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

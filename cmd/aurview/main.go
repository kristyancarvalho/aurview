package main

import (
	"context"
	"os"
	"strings"

	"github.com/kristyancarvalho/aurview/internal/app"
)

func main() {
	if err := app.Run(context.Background(), app.Options{InitialQuery: strings.Join(os.Args[1:], " ")}); err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

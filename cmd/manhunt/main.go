package main

import (
	"os"

	"github.com/bashnko/manhunt/internal/app"
)

func main() {
	if err := app.Run(os.Args[1:]); err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

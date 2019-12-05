package main

import (
	"fmt"
	"os"

	"github.com/sawadashota/orb-update/cmd"
)

func main() {
	if err := cmd.RootCmd().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

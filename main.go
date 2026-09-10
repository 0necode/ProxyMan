package main

import (
	"fmt"
	"github.com/proxy-cli/proxy-cli/internal/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

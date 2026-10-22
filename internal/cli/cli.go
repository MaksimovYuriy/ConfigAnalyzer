// Package cli provides the command-line interface and process exit codes.
package cli

import (
	"fmt"
	"io"
)

const (
	ExitOK       = 0
	ExitFindings = 1
	ExitError    = 2
)

// Run executes the command and returns its process exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
		_, err := fmt.Fprintln(stdout, "Usage: config-analyzer <file> [-s|--silent]\n       config-analyzer --stdin [-s|--silent]\n\nConfiguration analysis is not implemented yet.")
		if err != nil {
			return ExitError
		}
		return ExitOK
	}

	fmt.Fprintln(stderr, "config-analyzer: configuration analysis is not implemented yet; use --help for usage")
	return ExitError
}

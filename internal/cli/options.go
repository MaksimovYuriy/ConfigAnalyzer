package cli

import (
	"fmt"
	"strings"
)

const usage = `Usage: config-analyzer <file> [-s|--silent]
       config-analyzer --stdin [-s|--silent]

Checks: debug logging, plaintext passwords, TLS, wildcard listeners, insecure algorithms.`

type options struct {
	path      string
	silent    bool
	fromStdin bool
	help      bool
}

func parseArgs(args []string) (options, error) {
	var opts options
	if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
		opts.help = true
		return opts, nil
	}
	positional := false
	fileCount := 0
	for _, arg := range args {
		switch {
		case !positional && arg == "--":
			positional = true
		case !positional && (arg == "-s" || arg == "--silent"):
			opts.silent = true
		case !positional && arg == "--stdin":
			opts.fromStdin = true
		case !positional && strings.HasPrefix(arg, "-"):
			return options{}, fmt.Errorf("unknown option %q", arg)
		default:
			fileCount++
			if fileCount > 1 {
				return options{}, fmt.Errorf("expected one configuration file")
			}
			opts.path = arg
		}
	}
	if opts.fromStdin && fileCount > 0 || !opts.fromStdin && opts.path == "" {
		return options{}, fmt.Errorf("provide either a configuration file or --stdin")
	}
	return opts, nil
}

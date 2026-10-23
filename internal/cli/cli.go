package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/analyzer"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules"
	"github.com/MaksimovYuriy/ConfigAnalyzer/pkg/parser"
)

const (
	ExitOK       = 0
	ExitFindings = 1
	ExitError    = 2
)

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	opts, err := parseArgs(args)
	if err != nil {
		return fail(stderr, err)
	}
	if opts.help {
		if _, err := fmt.Fprintln(stdout, usage); err != nil {
			return fail(stderr, err)
		}
		return ExitOK
	}
	input := stdin
	if !opts.fromStdin {
		file, err := os.Open(opts.path)
		if err != nil {
			return fail(stderr, err)
		}
		defer file.Close()
		input = file
	}
	doc, err := parser.Parse(input)
	if err != nil {
		return fail(stderr, err)
	}
	problems := analyzer.Analyze(doc.Root, rules.Default()...)
	if len(problems) == 0 {
		if _, err := fmt.Fprintln(stdout, "Анализ завершён. Проблем не найдено."); err != nil {
			return fail(stderr, err)
		}
	}
	for _, problem := range problems {
		if _, err := fmt.Fprintf(stdout, "%s: %s: %s %s\n", problem.Severity, problem.Path, problem.Explanation, problem.Recommendation); err != nil {
			return fail(stderr, err)
		}
	}
	if len(problems) > 0 && !opts.silent {
		return ExitFindings
	}
	return ExitOK
}

func fail(stderr io.Writer, err error) int {
	fmt.Fprintln(stderr, "config-analyzer:", err)
	return ExitError
}

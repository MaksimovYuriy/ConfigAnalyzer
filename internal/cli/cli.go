package cli

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase/configuration"
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
	var metadata *entity.FileMetadata
	if !opts.fromStdin {
		file, err := os.Open(opts.path)
		if err != nil {
			return fail(stderr, err)
		}
		defer file.Close()
		info, err := os.Stat(opts.path)
		if err != nil {
			return fail(stderr, err)
		}
		metadata = &entity.FileMetadata{Path: opts.path, Mode: info.Mode()}
		input = file
	}
	content, err := io.ReadAll(input)
	if err != nil {
		return fail(stderr, err)
	}
	service := configuration.New(rules.Default())
	problems, err := service.Analyze(context.Background(), entity.Configuration{Content: string(content), FileMetadata: metadata})
	if err != nil {
		return fail(stderr, err)
	}
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

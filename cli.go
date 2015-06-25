package main

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/b4b4r07/pygments"
)

const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
	ExitCodeParseFlagError
	ExitCodeInitializeError
	ExitCodePygmentizeError
	ExitCodeTooFewArguments
)

type CLI struct {
	outStream, errStream io.Writer
}

var (
	errStack []int
)

func (cli *CLI) Run(args []string) (status int, err error) {
	var style string
	flags := flag.NewFlagSet("pygmentize", flag.ContinueOnError)
	flags.SetOutput(cli.errStream)
	flags.StringVar(&style, "style", "default", "")
	flags.StringVar(&style, "s", "default", "")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagError, err
	}

	p, err := pygments.New(
		pygments.Executable(),
		pygments.Formatter("console256"),
		pygments.Style(style),
	)
	if err != nil {
		return ExitCodeInitializeError, err
	}

	if flags.NArg() == 0 {
		return ExitCodeTooFewArguments, errors.New("too few arguments")
	}

	for _, f := range flags.Args() {
		out, err := p.Pygmentize(f)
		if err != nil {
			errStack = append(errStack, 1)
		}
		fmt.Fprintln(cli.outStream, out)
	}

	if len(errStack) > 0 {
		return ExitCodePygmentizeError, err
	}
	return ExitCodeOK, nil
}

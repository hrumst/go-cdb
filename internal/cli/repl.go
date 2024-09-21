package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/hrumst/go-cdb/internal/database/compute"
)

const (
	exitCommand = ".exit"
)

type database interface {
	Execute(ctx context.Context, input string) (*compute.CommandExecResult, error)
}

type repl struct {
	cliName    string
	inputLimit int
	db         database
	input      io.Reader
	output     io.Writer
}

func NewREPL(
	cliName string,
	inputLimit int,
	db database,
	input io.Reader,
	output io.Writer,
) *repl {
	return &repl{
		cliName:    cliName,
		inputLimit: inputLimit,
		db:         db,
		input:      input,
		output:     output,
	}
}

func (rc *repl) isExitCommand(input string) bool {
	return strings.TrimSpace(input) == exitCommand
}

func (rc *repl) printPrompt() error {
	_, err := fmt.Fprintf(rc.output, "%s > ", rc.cliName)
	return err
}

func (rc *repl) printErrResult(message string) error {
	_, err := fmt.Fprintf(rc.output, "%s Error: %s \n", rc.cliName, message)
	return err
}

func (rc *repl) printOkResult(message string) error {
	_, err := fmt.Fprintf(rc.output, "%s Ok: %s \n", rc.cliName, message)
	return err
}

func (rc *repl) Run() error {
	reader := bufio.NewScanner(rc.input)
	reader.Buffer(make([]byte, rc.inputLimit), rc.inputLimit)

	if err := rc.printPrompt(); err != nil {
		return err
	}
	for reader.Scan() {
		input := reader.Text()
		if rc.isExitCommand(input) {
			break
		}

		result, err := rc.db.Execute(context.Background(), input)
		if err != nil {
			if err := rc.printErrResult(err.Error()); err != nil {
				return err
			}
		} else {
			if err := rc.printOkResult(result.Result); err != nil {
				return err
			}
		}

		if err := rc.printPrompt(); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintln(rc.output, "Bye!"); err != nil {
		return err
	}
	return reader.Err()

}

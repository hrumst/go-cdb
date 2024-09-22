package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
)

const (
	exitCommand = ".exit"
)

type database interface {
	Execute(ctx context.Context, input string) string
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

func (rc *repl) printResult(message string) error {
	_, err := fmt.Fprintf(rc.output, "%s: %s \n", rc.cliName, message)
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

		result := rc.db.Execute(context.Background(), input)
		if err := rc.printResult(result); err != nil {
			return err
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

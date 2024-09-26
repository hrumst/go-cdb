package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	exitCommand = ".exit"
)

type repl struct {
	cliName    string
	inputLimit int
	input      io.Reader
	output     io.Writer
	handler    func(string) (string, error)
}

func NewREPL(
	cliName string,
	inputLimit int,
	input io.Reader,
	output io.Writer,
	handler func(string) (string, error),
) *repl {
	return &repl{
		cliName:    cliName,
		inputLimit: inputLimit,
		input:      input,
		output:     output,
		handler:    handler,
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

		resp, err := rc.handler(input)
		if err != nil {
			return err
		}
		if err := rc.printResult(resp); err != nil {
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

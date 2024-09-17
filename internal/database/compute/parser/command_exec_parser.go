package parser

import "github.com/hrumst/go-cdb/internal/database/compute"

type inputParserFactory func(string) inputParser

type inputParser interface {
	ParseStart()
	ParseCommandToken() (compute.CommandType, error)
	ParseWhitespace() error
	ParseKey() (string, error)
	ParseValue() (string, error)
	ParseEnd() error
}

type commandExecParser struct {
	ipf inputParserFactory
}

func NewCommandExecParserPlain() *commandExecParser {
	ipf := func(s string) inputParser {
		return NewInputParserPlain(s)
	}
	return newCommandExecParser(ipf)
}

func newCommandExecParser(ipf inputParserFactory) *commandExecParser {
	return &commandExecParser{
		ipf: ipf,
	}
}

func (cp *commandExecParser) Parse(input string) (*compute.CommandExec, error) {
	ip := cp.ipf(input)

	newCommandExec := &compute.CommandExec{}

	ip.ParseStart()

	commandType, err := ip.ParseCommandToken()
	if err != nil {
		return nil, err
	}
	newCommandExec.Command = commandType

	if err := ip.ParseWhitespace(); err != nil {
		return nil, err
	}

	key, err := ip.ParseKey()
	if err != nil {
		return nil, err
	}
	newCommandExec.Key = key

	if commandType == compute.CommandTypeSet {
		if err := ip.ParseWhitespace(); err != nil {
			return nil, err
		}
		val, err := ip.ParseValue()
		if err != nil {
			return nil, err
		}
		newCommandExec.Val = val
	}

	if err := ip.ParseEnd(); err != nil {
		return nil, err
	}
	return newCommandExec, nil
}

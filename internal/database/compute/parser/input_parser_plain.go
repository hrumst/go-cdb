package parser

import (
	"fmt"
	"unicode"

	"github.com/hrumst/go-cdb/internal/database/compute"
)

type InputParserPlain struct {
	posi     int
	rawInput []rune
}

func NewInputParserPlain(rawInput string) *InputParserPlain {
	return &InputParserPlain{
		rawInput: []rune(rawInput),
	}
}

func (ip *InputParserPlain) isEnd() bool {
	return ip.posi >= len(ip.rawInput)
}

func (ip *InputParserPlain) isSpace() bool {
	return unicode.IsSpace(ip.rawInput[ip.posi])
}

func (ip *InputParserPlain) parseError(parseErrType error, startPosi int) error {
	return fmt.Errorf(
		"%w: %w, at position: %d, invalid input: '%s'",
		ParseError,
		parseErrType,
		startPosi,
		string(ip.rawInput),
	)
}

func (ip *InputParserPlain) movePosParse() bool {
	if ip.posi+1 >= len(ip.rawInput) {
		return false
	}
	ip.posi += 1
	return true
}

func (ip *InputParserPlain) ParseCommandToken() (compute.CommandType, error) {
	startPosi := ip.posi
	for !ip.isEnd() {
		if ip.isSpace() || ip.posi-ip.posi > 3 {
			return 0, ip.parseError(invalidCommandTokenError, startPosi)
		}

		if ip.posi-startPosi == 2 {
			var ct compute.CommandType
			switch string(ip.rawInput[startPosi : ip.posi+1]) {
			case "SET":
				ct = compute.CommandTypeSet
			case "GET":
				ct = compute.CommandTypeGet
			case "DEL":
				ct = compute.CommandTypeDel
			default:
				return 0, ip.parseError(invalidCommandTokenError, startPosi)
			}
			ip.movePosParse()
			return ct, nil
		}

		if !ip.movePosParse() {
			break
		}
	}
	return 0, ip.parseError(invalidCommandTokenError, startPosi)
}

func (ip *InputParserPlain) ParseWhitespace() error {
	startPosi := ip.posi
	spaceLen := 0
	for !ip.isEnd() {
		if !ip.isSpace() {
			break
		}
		spaceLen += 1
		if !ip.movePosParse() {
			break
		}
	}
	if spaceLen == 0 {
		return ip.parseError(noWhitespaceError, startPosi)
	}
	return nil
}

func (ip *InputParserPlain) ParseKey() (string, error) {
	startPosi := ip.posi
	var key string
	for !ip.isEnd() {
		if ip.isSpace() {
			key = string(ip.rawInput[startPosi:ip.posi])
			break
		}
		if !ip.movePosParse() {
			key = string(ip.rawInput[startPosi : ip.posi+1])
			break
		}
	}

	if len(key) == 0 {
		return "", ip.parseError(invalidKeyError, startPosi)
	}
	return key, nil
}

func (ip *InputParserPlain) ParseValue() (string, error) {
	startPosi := ip.posi
	var val string
	for !ip.isEnd() {
		if ip.isSpace() {
			val = string(ip.rawInput[startPosi:ip.posi])
			break
		}
		if !ip.movePosParse() {
			val = string(ip.rawInput[startPosi : ip.posi+1])
			break
		}
	}

	if len(val) == 0 {
		return "", ip.parseError(invalidKeyValueError, startPosi)
	}
	return val, nil
}

func (ip *InputParserPlain) ParseEnd() error {
	startPosi := ip.posi
	for !ip.isEnd() {
		if !ip.isSpace() {
			if ip.posi-startPosi > 0 {
				return ip.parseError(unexpectedInputEndError, startPosi)
			}
		}
		if !ip.movePosParse() {
			break
		}
	}
	return nil
}

func (ip *InputParserPlain) ParseStart() {
	for {
		if ip.isEnd() || !ip.isSpace() || !ip.movePosParse() {
			break
		}
	}
}

package parser

import "errors"

var (
	ParseError = errors.New("parse error")

	invalidCommandTokenError = errors.New("invalid token")
	invalidKeyError          = errors.New("invalid key")
	invalidKeyValueError     = errors.New("invalid value")
	noWhitespaceError        = errors.New("no whitespace")
	unexpectedInputEndError  = errors.New("unexpected input end")
)

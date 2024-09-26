package config

import (
	"strconv"
	"strings"
)

func ParserSize(input string) int64 {
	input = strings.TrimSpace(strings.ToLower(input))
	var (
		numPart string
		unit    string
	)
	for i, r := range input {
		if r < '0' || r > '9' {
			numPart = input[:i]
			unit = input[i:]
			break
		}
	}

	num, err := strconv.ParseInt(numPart, 10, 64)
	if err != nil {
		return 0
	}

	switch unit {
	case "kb":
		return num << 10
	case "mb":
		return num << 20
	case "gb":
		return num << 30
	case "b":
		return num
	default:
		return 0
	}
}

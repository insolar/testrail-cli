//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package internal

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

func skipLeadingWhitespaces(leftovers string) string {
	for i, r := range leftovers {
		if !unicode.IsSpace(r) {
			return leftovers[i:]
		}
	}
	return ""
}

func parseMessage(leftovers string, skipEqualSign int) (string, []byte, error) {
	lastWordPos := 0

	for i, r := range leftovers {
		switch {
		case r == '=':
			if skipEqualSign >= 0 {
				skipEqualSign--
				continue
			} else if lastWordPos == 0 {
				return leftovers, nil, nil
			}
			return leftovers[lastWordPos-1:], []byte(leftovers[:lastWordPos-1]), nil
		case unicode.IsSpace(r):
			lastWordPos = i+1
		}
	}

	return "", []byte(leftovers), nil
}

var (
	MalformedKeyHaveSpace = errors.New("malformed key: got space")
)

func parseKeyString(leftovers string) (string, []byte, error) {
	for i, r := range leftovers {
		switch {
		case unicode.IsSpace(r):
			return leftovers, nil, MalformedKeyHaveSpace
		case r == '=':
			return leftovers[i+1:], []byte(leftovers[:i]), nil
		}
	}
	return leftovers, nil, errors.New("malformed key 2")
}


func parseReverseKeyString(leftovers string) (string, []byte, error) {
	for i, r := range leftovers {
		switch {
		case unicode.IsSpace(r):
			return leftovers[i+1:], []byte(leftovers[1:i]), nil
		}
	}
	return leftovers, nil, errors.New("malformed key 2")
}

func parseJsonValueString(leftovers string) (string, []byte, error) {
	if !(leftovers[0] == '"') {
		for i, r := range leftovers {
			if unicode.IsSpace(r) {
				return leftovers[i:], []byte(leftovers[:i]), nil
			}
		}
		return "", []byte(leftovers), nil
	}

	var result []byte

	i := 1
	for i < len(leftovers) {
		c := leftovers[i]
		switch leftovers[i] {
		case '\\':
			if len(leftovers) < i+2 {
				return leftovers, nil, errors.New("malformed string 1")
			}

			newPosition := i+2
			switch leftovers[i+1] {
			case '"':
				result = append(result, '"')
			case '\\':
				result = append(result, '\\')
			case 'b':
				result = append(result, '\b')
			case 'f':
				result = append(result, '\f')
			case 'n':
				result = append(result, '\n')
			case 'r':
				result = append(result, '\r')
			case 't':
				result = append(result, '\t')
			case 'u':
				if len(leftovers) < i+6 {
					return leftovers, nil, errors.New("malformed string 2")
				}

				unquotedNumber, err := strconv.ParseInt(leftovers[i+2:i+6], 16, 32)
				if err != nil {
					return leftovers, nil, fmt.Errorf("malformed string 3: %w", err)
				}

				var b []byte
				for unquotedNumber > 0 {
					b = append(b, byte(unquotedNumber) & 255)
					unquotedNumber >>= 8
				}

				for i := 0; i < len(b); i++ {
					result = append(result, b[len(b) - (i+1)])
				}

				newPosition = i+6
			default:
				return leftovers, nil, errors.New("malformed string 4")
			}

			i = newPosition
		case '"':
			return leftovers[:i], result, nil
		default:
			result = append(result, c)
			i++
		}
	}

	return leftovers, nil, errors.New("malformed string 5")
}

func parseReverseJsonValueString(leftovers string) (string, []byte, error) {
	if !(leftovers[0] == '"') {
		for i, r := range leftovers {
			if r == '=' {
				return leftovers[i:], []byte(leftovers[:i]), nil
			} else if unicode.IsSpace(r) {
				return leftovers, nil, errors.New("malformed string 0")
			}
		}
		return "", []byte(leftovers), nil
	}

	var result []byte

	i := 1
	for i < len(leftovers) {
		c := leftovers[i]
		switch leftovers[i] {
		case '\\':
			if len(leftovers) < i+2 {
				return leftovers, nil, errors.New("malformed string 1")
			}

			newPosition := i+2
			switch leftovers[i+1] {
			case '"':
				result = append(result, '"')
			case '\\':
				result = append(result, '\\')
			case 'b':
				result = append(result, '\b')
			case 'f':
				result = append(result, '\f')
			case 'n':
				result = append(result, '\n')
			case 'r':
				result = append(result, '\r')
			case 't':
				result = append(result, '\t')
			case 'u':
				if len(leftovers) < i+6 {
					return leftovers, nil, errors.New("malformed string 2")
				}

				unquotedNumber, err := strconv.ParseInt(leftovers[i+2:i+6], 16, 32)
				if err != nil {
					return leftovers, nil, fmt.Errorf("malformed string 3: %w", err)
				}

				var b []byte
				for unquotedNumber > 0 {
					b = append(b, byte(unquotedNumber) & 255)
					unquotedNumber >>= 8
				}

				for i := 0; i < len(b); i++ {
					result = append(result, b[len(b) - (i+1)])
				}

				newPosition = i+6
			default:
				return leftovers, nil, errors.New("malformed string 4")
			}

			i = newPosition
		case '"':
			if leftovers[i+1] != '=' {
				return leftovers, nil, errors.New("malformed string 6")
			}
			return leftovers[:i], result, nil
		default:
			result = append(result, c)
			i++
		}
	}

	return leftovers, nil, errors.New("malformed string 5")
}

package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const (
	stateStart = iota
	stateBackslash
	stateDigit
	stateCharacter
)

type state struct {
	lastChar  rune
	prevState int
}

func (s *state) next(char rune, state int) {
	s.lastChar = char
	s.prevState = state
}

const backSlash = rune(92)

func Unpack(s string) (string, error) {
	builder := strings.Builder{}
	st := state{
		lastChar:  0,
		prevState: stateStart,
	}

	for _, c := range s {
		isDigit := unicode.IsDigit(c)
		isBackslash := c == backSlash

		if st.prevState == stateStart {
			if isDigit {
				return "", ErrInvalidString
			}

			if isBackslash {
				st.next(0, stateBackslash)
				continue
			}

			st.next(c, stateCharacter)

			continue
		}

		if st.prevState == stateBackslash {
			st.next(c, stateCharacter)
			continue
		}

		if st.prevState == stateCharacter {
			repeat := 1
			if isDigit {
				repeat, _ = strconv.Atoi(string(c))
			}
			for repeat > 0 {
				builder.WriteRune(st.lastChar)
				repeat--
			}
			switch {
			case isDigit:
				st.next(0, stateDigit)
			case isBackslash:
				st.next(c, stateBackslash)
			default:
				st.next(c, stateCharacter)
			}
			continue
		}

		if st.prevState == stateDigit {
			if isDigit {
				return "", ErrInvalidString
			}
			st.next(c, stateCharacter)
			continue
		}
	}

	if st.lastChar != 0 {
		builder.WriteRune(st.lastChar)
	}

	return builder.String(), nil
}

package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

var ErrInvalidState = errors.New("invalid state")

type parsingState int

const (
	stateStart parsingState = iota + 1
	stateBackslash
	stateDigit
	stateCharacter
)

type state struct {
	lastChar  rune
	prevState parsingState
}

func newState() state {
	return state{
		lastChar:  0,
		prevState: stateStart,
	}
}

func (st *state) rememberLastChar(c rune, s parsingState) {
	st.lastChar = c
	st.prevState = s
}

const backSlash = '\\'

func Unpack(s string) (string, error) {
	builder := strings.Builder{}
	st := newState()

	for _, c := range s {
		isDigit := unicode.IsDigit(c)
		isBackslash := c == backSlash

		switch st.prevState {
		case stateStart:
			if isDigit {
				return "", ErrInvalidString
			}

			if isBackslash {
				st.rememberLastChar(0, stateBackslash)
				continue
			}

			st.rememberLastChar(c, stateCharacter)
		case stateBackslash:
			st.rememberLastChar(c, stateCharacter)
		case stateCharacter:
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
				st.rememberLastChar(0, stateDigit)
			case isBackslash:
				st.rememberLastChar(c, stateBackslash)
			default:
				st.rememberLastChar(c, stateCharacter)
			}
		case stateDigit:
			if isDigit {
				return "", ErrInvalidString
			}
			st.rememberLastChar(c, stateCharacter)
		default:
			return "", ErrInvalidState
		}
	}

	if st.lastChar != 0 {
		builder.WriteRune(st.lastChar)
	}

	return builder.String(), nil
}

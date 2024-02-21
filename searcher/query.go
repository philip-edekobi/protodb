package searcher

import (
	"fmt"
	"strings"
)

type CompareFilter struct {
	key   []string
	value string
	op    string
}

type Query struct {
	ands []CompareFilter
}

func parseFilter(q string) (*Query, error) {
	if q == "" {
		return &Query{}, nil
	}

	i := 0
	var parsed Query
	qRune := []rune(q)

	for i < len(qRune) {
		i = eatWhiteSpace(qRune, i)

		key, nextIndex, err := LexString(qRune, i)
		if err != nil {
			return nil, fmt.Errorf("expected a valid key, got [%s]: `%s`", err, q[nextIndex:])
		}

		// expecting comparison operator, ":"
		if q[nextIndex] != ':' {
			return nil, fmt.Errorf("expected colon at %d, got `%s`", nextIndex, q[nextIndex:])
		}
		i = nextIndex + 1

		op := "="
		i = eatWhiteSpace(qRune, i)
		if q[i] == '>' || q[i] == '<' {
			op = string(q[i])
			i++
		}

		value, nextIndex, err := LexString(qRune, i)
		if err != nil {
			return nil, fmt.Errorf("expected a valid value, got [%s]: `%s`", err, q[nextIndex:])
		}
		i = nextIndex

		argument := CompareFilter{
			key:   strings.Split(key, "."),
			value: value,
			op:    op,
		}

		parsed.ands = append(parsed.ands, argument)
	}

	return &parsed, nil
}

package searcher

import (
	"fmt"
	"strconv"
	"strings"
)

type CompareFilter struct {
	Key   []string
	Value string
	Op    string
}

type Query struct {
	Ands []CompareFilter
}

func (q Query) Match(doc map[string]any) bool {
	for _, argument := range q.Ands {
		Value, ok := getPath(doc, argument.Key)
		if !ok {
			return false
		}

		// equality
		if argument.Op == "=" {
			match := fmt.Sprintf("%v", Value) == argument.Value
			if !match {
				return false
			}

			continue
		}

		// < and >
		right, err := strconv.ParseFloat(argument.Value, 64)
		if err != nil {
			return false
		}

		var left float64

		switch t := Value.(type) {
		case float64:
			left = t
		case float32:
			left = float64(t)
		case uint:
			left = float64(t)
		case uint8:
			left = float64(t)
		case uint16:
			left = float64(t)
		case uint32:
			left = float64(t)
		case uint64:
			left = float64(t)
		case int:
			left = float64(t)
		case int8:
			left = float64(t)
		case int16:
			left = float64(t)
		case int32:
			left = float64(t)
		case int64:
			left = float64(t)
		case string:
			left, err = strconv.ParseFloat(t, 64)
			if err != nil {
				return false
			}
		default:
			return false
		}

		if argument.Op == ">" {
			if left <= right {
				return false
			}

			continue
		}

		if left >= right {
			return false
		}
	}

	return true
}

func ParseFilter(q string) (*Query, error) {
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
			return nil, fmt.Errorf("expected a valid Key, got [%s]: `%s`", err, q[nextIndex:])
		}

		if nextIndex >= len(q) {
			break
		}

		// expecting comparison Operator, ":"
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
			return nil, fmt.Errorf("expected a valid Value, got [%s]: `%s`", err, q[nextIndex:])
		}
		i = nextIndex

		argument := CompareFilter{
			Key:   strings.Split(key, "."),
			Value: value,
			Op:    op,
		}

		parsed.Ands = append(parsed.Ands, argument)
	}

	return &parsed, nil
}

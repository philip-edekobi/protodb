package searcher

import (
	"fmt"
	"unicode"
)

func LexString(input []rune, index int) (string, int, error) {
	if index >= len(input) {
		return "", index, nil
	}

	index = eatWhiteSpace(input, index)

	if input[index] == '"' {
		index++
		foundEnd := false

		var s []rune

		// TODO: fix nested quotes
		for index < len(input) {
			if input[index] == '"' {
				foundEnd = true
				break
			}

			s = append(s, input[index])
			index++
		}

		if !foundEnd {
			return "", index, fmt.Errorf("quoted string does not terminate")
		}

		return string(s), index + 1, nil
	}

	// string is unquoted
	var s []rune
	var c rune

	// TODO: handle the case of ... in input
	for index < len(input) {
		c = input[index]
		if !(unicode.IsLetter(c) || unicode.IsDigit(c) || c == '.') {
			break
		}

		s = append(s, c)
		index++
	}

	if len(s) == 0 {
		fmt.Println("currchar is:", string(c))
		return "", index, fmt.Errorf("no string found")
	}

	return string(s), index, nil
}

func eatWhiteSpace(input []rune, index int) int {
	if index >= len(input) {
		return index
	}

	idx := index

	for unicode.IsSpace(input[idx]) {
		idx++
	}

	return idx
}

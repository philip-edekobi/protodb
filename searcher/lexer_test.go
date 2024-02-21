package searcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLex(t *testing.T) {
	str := "age: 17"
	idx := 0

	age, idx, err := LexString([]rune(str), 0)
	require.Nil(t, err)
	require.Equal(t, "age", age)

	num, idx, err := LexString([]rune(str), idx+1)
	require.Nil(t, err)
	require.Equal(t, "17", num)
}

func TestEatWhiteSpace(t *testing.T) {
	str := " hola  dee"

	index := eatWhiteSpace([]rune(str), 0)
	require.Equal(t, 1, index)

	index = eatWhiteSpace([]rune(str), 5)
	require.Equal(t, 7, index)
}

package searcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFilter(t *testing.T) {
	q := "age:< 17 name.b:andy rating: > 3"

	expected := &Query{
		ands: []CompareFilter{
			{
				key:   []string{"age"},
				value: "17",
				op:    "<",
			},
			{
				key:   []string{"name", "b"},
				value: "andy",
				op:    "=",
			},
			{
				key:   []string{"rating"},
				value: "3",
				op:    ">",
			},
		},
	}

	ans, err := parseFilter(q)
	require.Nil(t, err)

	require.Equal(t, expected, ans)
}

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

func TestMatch(t *testing.T) {
	docs := []map[string]any{
		{
			"name": "andy",
			"age":  "17",
		},
		{
			"class": map[string]any{
				"upper": "12",
			},
		},
	}

	q, err := parseFilter("age:>16")
	require.Nil(t, err)

	q2, err := parseFilter("class.upper:>9")
	require.Nil(t, err)

	require.Equal(t, true, q.Match(docs[0]))
	require.Equal(t, true, q2.Match(docs[1]))
	require.Equal(t, false, q2.Match(docs[0]))
}

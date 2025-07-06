package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{ // Whitespace case
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{ // Word Case case
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{ // Mixed Case
			input:    "tHis iS a tEsT        of tHe funCtion",
			expected: []string{"this", "is", "a", "test", "of", "the", "function"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf(`Slice is not expected length.`)
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf(`Unexpected word.\nExpected: %s\nGot: %s`, expectedWord, word)
			}
		}
	}
}

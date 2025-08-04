package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "   ",
			expected: []string{},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "singleword",
			expected: []string{"singleword"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		// Check the length of the actual slice against the expected slice
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q) returned slice of length %d, expected %d",
				c.input, len(actual), len(c.expected))
			continue
		}

		// Check each word in the slice
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("cleanInput(%q) returned %q at index %d, expected %q",
					c.input, word, i, expectedWord)
			}
		}
	}
}

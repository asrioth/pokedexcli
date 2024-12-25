package main

import (
	"testing"
)

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
			input:    "  hElLO  woRld  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  hel    lo  wor    ld  ",
			expected: []string{"hel", "lo", "wor", "ld"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "helloworld",
			expected: []string{"helloworld"},
		},
		{
			input:    "  helloworld  ",
			expected: []string{"helloworld"},
		},

		// add more casees here
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Not as many words as expected.\nInput: %v\n Actual: %v, Expected: %v", c.input, len(actual), len(c.expected))
			t.Fail()
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Actual word didn't match expected word.\nInput: %v\n Actual: %v, Expected: %v", c.input, word, expectedWord)
				t.Fail()
			}
		}
	}
}

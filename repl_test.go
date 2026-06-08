package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input : "hello world",
			expected : []string{"hello", "world"},
		},
		{
			input : "   leading and trailing spaces   ",
			expected : []string{"leading", "and", "trailing", "spaces"},
		},
		{
			input : "multiple   spaces",
			expected : []string{"multiple", "spaces"},
		},
		{
			input : "",
			expected : []string{},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Expected length %d, got %d", len(c.expected), len(actual))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Expected word '%s', got '%s'", expectedWord, word)
			}
		}
	}
}

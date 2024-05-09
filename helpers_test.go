package main

import (
	"testing"
)

func TestReplaceWords(t *testing.T) {
	cases := []struct {
		input    []string
		expected string
	}{
		{
			input:    []string{"This is a kerfuffle opinion I need to share with the world", "kerfuffle", "****"},
			expected: "This is a **** opinion I need to share with the world",
		},
		{
			input:    []string{"I hear Mastodon is better than Chirpy. sharbert I need to migrate", "sharbert", "****"},
			expected: "I hear Mastodon is better than Chirpy. **** I need to migrate",
		},
		{
			input:    []string{"I really need a foRnAX to go to bed sooner, Fornax !", "fornax", "****"},
			expected: "I really need a **** to go to bed sooner, **** !",
		},
	}

	for _, c := range cases {
		actual := replaceWord(c.input[0], c.input[1], c.input[2])
		if len(actual) != len(c.expected) {
			t.Errorf("lengths don't match: '%v' vs '%v'", actual, c.expected)
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("cleanInput(%v) == %v, expected %v", c.input, actual, c.expected)
			}
		}
	}
}

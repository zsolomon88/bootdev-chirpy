package database

import (
	"testing"
)

func TestReplaceWords(t *testing.T) {
	cases := []struct {
		input    string
		expected Chirp
	}{
		{
			input: "This is a test chirp 1",
			expected: Chirp{
				Id:   1,
				Body: "This is a test chirp 1",
			},
		},
		{
			input: "This is a test chirp 2",
			expected: Chirp{
				Id:   2,
				Body: "This is a test chirp 2",
			},
		},
		{
			input: "This is a test chirp 3",
			expected: Chirp{
				Id:   3,
				Body: "This is a test chirp 3",
			},
		},
	}

	db, err := NewDB("./testdb.json")
	if err != nil {
		t.Errorf("unable to create db: %s", err)
	}

	for _, c := range cases {
		actual, createErr := db.CreateChirp(c.input)
		if createErr != nil {
			t.Errorf("unable to create chirp: %v", err)
		}
		if len(actual.Body) != len(c.expected.Body) {
			t.Errorf("lengths don't match: '%v' vs '%v'", actual.Body, c.expected.Body)
			continue
		}
		if actual.Id != c.expected.Id {
			t.Errorf("chirp id's dont match: '%v' vs '%v'", actual.Id, c.expected.Id)
		}
		for i := range actual.Body {
			word := actual.Body[i]
			expectedWord := c.expected.Body[i]
			if word != expectedWord {
				t.Errorf("cleanInput(%v) == %v, expected %v", c.input, actual.Body, c.expected.Body)
			}
		}
	}

}

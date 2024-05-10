package database

import (
	"fmt"
	"testing"
)

func TestCreateUser(t *testing.T) {
	cases := []struct {
		input    []string
		expected User
		updated  User
	}{
		{
			input: []string{"usr1@boot.dev", "pwd1"},
			expected: User{
				Id:       1,
				Email:    "usr1@boot.dev",
				Password: "pwd1",
			},
			updated: User{
				Id:       1,
				Email:    "updated1@boot.dev",
				Password: "updated",
			},
		},
		{
			input: []string{"usr2@boot.dev", "pwd2"},
			expected: User{
				Id:       2,
				Email:    "usr2@boot.dev",
				Password: "pwd2",
			},
			updated: User{
				Id:       2,
				Email:    "updated2@boot.dev",
				Password: "updated",
			},
		},
		{
			input: []string{"usr3@boot.dev", "pwd3"},
			expected: User{
				Id:       3,
				Email:    "usr3@boot.dev",
				Password: "pwd3",
			},
			updated: User{
				Id:       3,
				Email:    "updated3@boot.dev",
				Password: "updated",
			},
		},
	}

	db, err := NewDB("./testdb.json")
	if err != nil {
		t.Errorf("unable to create db: %s", err)
	}

	for _, c := range cases {
		actual, createErr := db.CreateUser(c.input[0], c.input[1])
		if createErr != nil {
			t.Errorf("unable to create user: %v", err)
		}
		if len(actual.Email) != len(c.expected.Email) {
			t.Errorf("lengths don't match: '%v' vs '%v'", actual.Email, c.expected.Email)
			continue
		}
		if actual.Id != c.expected.Id {
			t.Errorf("chirp id's dont match: '%v' vs '%v'", actual.Id, c.expected.Id)
		}
		for i := range actual.Email {
			word := actual.Email[i]
			expectedWord := c.expected.Email[i]
			if word != expectedWord {
				t.Errorf("cleanInput(%v) == %v, expected %v", c.input, actual.Email, c.expected.Email)
			}
		}
		for i := range actual.Password {
			word := actual.Password[i]
			expectedWord := c.expected.Password[i]
			if word != expectedWord {
				t.Errorf("cleanInput(%v) == %v, expected %v", c.input, actual.Password, c.expected.Password)
			}
		}
		actualUpdate, updateErr := db.UpdateUser(actual.Id, User{Password: "updated", Email: fmt.Sprintf("updated%d@boot.dev", actual.Id)})
		if updateErr != nil {
			t.Errorf("unable to update password: %v", err)
		}
		if len(actualUpdate.Password) != len(c.updated.Password) {
			t.Errorf("updated lengths don't match: '%v' vs '%v'", actual.Password, c.updated.Password)
		}
		for i := range actualUpdate.Email {
			word := actualUpdate.Email[i]
			expectedWord := c.updated.Email[i]
			if word != expectedWord {
				t.Errorf("cleanInput(%v) == %v, expected %v", "updated", actualUpdate.Password, c.updated.Password)
			}
		}
		for i := range actualUpdate.Password {
			word := actualUpdate.Password[i]
			expectedWord := c.updated.Password[i]
			if word != expectedWord {
				t.Errorf("cleanInput(%v) == %v, expected %v", "updated", actualUpdate.Password, c.updated.Password)
			}
		}
	}

}

func TestCreateTweet(t *testing.T) {
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

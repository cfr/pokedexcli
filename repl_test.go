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
            input:    " ",
            expected: []string{},
        },
        {
            input:    "  hello  world  ",
            expected: []string{"hello", "world"},
        },
        {
            input:    "  hEllO  wORld  ",
            expected: []string{"hello", "world"},
        },
        {
            input:    "world  ",
            expected: []string{"world"},
        },
        {
            input:    "≥ •",
            expected: []string{"≥", "•"},
        },
    }

    for _, c := range cases {
        actual := cleanInput(c.input)
        if len(actual) != len(c.expected) {
            t.Errorf("length does not match: %v vs %v in %v", len(actual), len(c.expected), c.input);
        }
        for i := range actual {
            word := actual[i]
            expectedWord := c.expected[i]
            if (word != expectedWord) {
                t.Errorf("word does not match: %v vs %v", word, expectedWord);
            }
        }
    }
}




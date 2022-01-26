package main

import (
	"fmt"
	"testing"
)

func TestVerifyData(t *testing.T) {

}

func TestTableVerifyData(t *testing.T) {
	var tests = []struct {
		input    string
		expected bool
	}{
		{"gggbbb", true},
		{"bbbggg", true},

		{"bbgggb", true},
		{"bbggbg", true},
		{"bbgbgg", true},

		{"ggbbbg", true},
		{"ggbbgb", true},
		{"ggbgbb", true},

		{"bgbbgg", true},
		{"bgbgbg", true},
		{"bgbggb", true},

		{"bgggbb", true},
		{"bggbgb", true},
		{"bggbbg", true},

		{"gbbbgg", true},
		{"gbbgbg", true},
		{"gbbggb", true},

		{"gbggbb", true},
		{"gbgbgb", true},
		{"gbgbbg", true},

		{"", false},
		{"b", false},
		{"g", false},
		{"ggggbb", false},
		{"ggggbbb", false},
		{"gggbbbb", false},
		{"gggbb", false},
		{"gggbbh", false},
		{"gg7bbb", false},
		{"@*^&", false},
		{"gggggggggggggggggggggbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", false},
	}
	for _, test := range tests {
		actual := verifyDataIsValid(test.input)
		if actual != test.expected {
			t.Error(fmt.Sprintf("Test Failed: Input: %s, Expected: %t, Actual: %t", test.input, test.expected, actual))
		}
	}
}

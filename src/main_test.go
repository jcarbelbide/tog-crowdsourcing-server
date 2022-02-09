package main

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestVerifyData(t *testing.T) {
	var tests = []struct {
		inputWorldNumber int
		inputStreamOrder string
		expectedOutput   bool
	}{
		{301, "gggbbb", true},
		{500, "gggbbb", true},
		{1000, "gggbbb", false},
		{10000, "gggbbb", false},
		{100, "gggbbb", false},
		{300, "gggbbb", false},
		{-1, "gggbbb", false},
	}
	var worldInformation WorldInformation
	for _, test := range tests {
		worldInformation.WorldNumber = test.inputWorldNumber
		worldInformation.StreamOrder = test.inputStreamOrder
		actual := verifyDataIsValid(worldInformation)
		if actual != test.expectedOutput {
			t.Error(fmt.Sprintf("Test Failed: Input: %d, %s, Expected: %t, Actual: %t", test.inputWorldNumber, test.inputStreamOrder, test.expectedOutput, actual))
		}
	}
}

func TestVerifyStreamOrder(t *testing.T) {
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
		actual := verifyStreamOrderIsValid(test.input)
		if actual != test.expected {
			t.Error(fmt.Sprintf("Test Failed: Input: %s, Expected: %t, Actual: %t", test.input, test.expected, actual))
		}
	}
}

func TestHasIPAlreadySubmittedDataForWorld(t *testing.T) {
	// Database
	testDB, _ := sql.Open("sqlite3", "./TestDB.db")
	statement, _ := testDB.Prepare("CREATE TABLE IF NOT EXISTS IP_WORLD_BLACKLIST (ip_world_hash INTEGER PRIMARY KEY)")
	statement.Exec()

	var ipWorldList = []struct {
		worldNumber int
		ipAddress   string
	}{
		{301, "123"},
		{301, "12345"},
		{302, "123"},
		{305, "1235"},
		{456, "123"},
	}

	statement, _ = testDB.Prepare("INSERT INTO IP_WORLD_BLACKLIST (ip_world_hash) VALUES ((?))")
	for _, v := range ipWorldList {
		ipWorldHash := hashIPAndWorldInfo(v.ipAddress, v.worldNumber)
		statement.Exec(ipWorldHash)
	}

	var tests = []struct {
		inputWorld     int
		inputIPAddress string
		expectedOutput bool
	}{
		{301, "343", false},
		{301, "35243", false},
		{302, "343235", false},
		{305, "1235", true},
		{456, "123", true},
		{301, "123", true},
		{301, "12345", true},
	}
	for _, test := range tests {
		actual, _ := hasIPAlreadySubmittedDataForWorld(test.inputIPAddress, test.inputWorld, testDB)
		if actual != test.expectedOutput {
			t.Error(fmt.Sprintf("Test Failed: Input: %d, %s, Expected: %t, Actual: %t", test.inputWorld, test.inputIPAddress, test.expectedOutput, actual))
		}
	}

	statement, _ = testDB.Prepare("DROP TABLE IP_WORLD_BLACKLIST")
	statement.Exec()

}

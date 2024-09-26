package rest

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sf "github.com/wissance/stringFormatter"
)

var cases = []struct {
	name     string
	value    string
	expected bool
}{
	{name: "correct", value: "GsFGdfgdfgdfg", expected: true},
	{name: "cyrillic test", value: "вапфрЫВАхиа", expected: true},
	{name: "number and test", value: "sdfg12515", expected: true},
	{name: "with underscore", value: "dfglq_sfdg_as", expected: true},
	{name: "with dash", value: "dfgdf-qwlwer-qwel", expected: true},
	{name: "with dash start", value: "-dfgdfqwlwerqwel", expected: true},

	{name: "empty string", value: "", expected: false},
	{name: "with space", value: "sdf sdf", expected: false},
	{name: "with double dash", value: "fdsdf--qlwqk", expected: false},
	{name: "with double dash start", value: "--fdsdfqlwqk", expected: false},
	{name: "with double dash finish", value: "fdsdfqlwqk--", expected: false},
	{name: "with slash", value: "kddfg/asd", expected: false},
	{name: "with backslash", value: `dfgdfg \n a\sd`, expected: false},
	{name: "QUOTATION MARK", value: "0\"", expected: false},
	{name: "QUOTATION MARK", value: `"`, expected: false},
	{name: "backslash", value: "\\", expected: false},
	{name: "slash", value: "/", expected: false},
	{name: "with semicolon;", value: "sdfsdf;", expected: false},
}

func TestValidate(t *testing.T) {
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, Validate(tc.value), sf.Format("INPUT: {0}", tc.value))
		})
	}
}

func FuzzValidate(f *testing.F) {
	for _, tc := range cases {
		f.Add(tc.value)
	}
	f.Fuzz(func(t *testing.T, input string) {
		isValid := Validate(input)
		forbiddenSymbols := []rune{'"', '\'', '%', '/', '\\'}
		if isContainsOne(input, forbiddenSymbols...) {
			require.False(t, isValid, sf.Format("INPUT: {0}", input))
		}
		if strings.Contains(input, "--") {
			require.False(t, isValid, sf.Format("INPUT: {0}", input))
		}
	})
}

// Returns true if at least one rune is contained
func isContainsOne(input string, args ...rune) bool {
	for _, r := range args {
		isContains := strings.ContainsRune(input, r)
		if isContains {
			return true
		}
	}
	return false
}
